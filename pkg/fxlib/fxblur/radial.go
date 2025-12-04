package fxblur

import (
	"kdfx/pkg/fxcontext"
	"kdfx/pkg/fxcore"
	"kdfx/pkg/fxnode"
)

// FXRadialBlurFS is the fragment fxShader for radial blur.
const FXRadialBlurFS = `
precision mediump float;
varying vec2 v_texCoord;
uniform sampler2D u_texture;
uniform vec2 u_center;
uniform float u_strength;

void main() {
	vec4 color = vec4(0.0);
	float total = 0.0;
	vec2 toCenter = u_center - v_texCoord;
	
	// Sample 10 times along the vector to center
	for (float t = 0.0; t <= 1.0; t += 0.1) {
		float percent = (t + 0.0) * u_strength; // Randomize? No.
		float weight = 1.0 - t; // Weight drops off?
		vec4 sample = texture2D(u_texture, v_texCoord + toCenter * percent);
		
		// Simple average
		color += sample;
		total += 1.0;
	}
	
	gl_FragColor = color / total;
}
`

// FXRadialBlurNode applies a radial (zoom) blur to the input texture.
type FXRadialBlurNode interface {
	fxnode.FXNode
	// SetCenter sets the center of the blur in normalized coordinates (0.0 to 1.0).
	// (0.5, 0.5) is the center of the image.
	SetCenter(x, y float32)
	// SetStrength sets the strength of the fxblur.
	// This controls how much the pixels are pulled towards the center.
	SetStrength(s float32)
}

// fxRadialBlurNode implements FXRadialBlurNode.
type fxRadialBlurNode struct {
	fxnode.FXNode
}

// NewFXRadialBlurNode creates a new radial blur fxnode.
func NewFXRadialBlurNode(ctx fxcontext.FXContext, width, height int) (FXRadialBlurNode, error) {
	base, err := fxnode.NewFXBaseNode(ctx, width, height)
	if err != nil {
		return nil, err
	}

	// Compile the shader program with the simple vertex shader and radial blur fragment shader.
	program, err := fxcore.NewFXShaderProgram(fxcore.FXSimpleVS, FXRadialBlurFS)
	if err != nil {
		base.Release()
		return nil, err
	}

	base.SetShaderProgram(program)

	n := &fxRadialBlurNode{
		FXNode: base,
	}
	// Set default center to the middle of the texture.
	n.SetCenter(0.5, 0.5)
	// Set default strength.
	n.SetStrength(0.1)

	return n, nil
}

func (n *fxRadialBlurNode) SetCenter(x, y float32) {
	// Set the center point of the radial blur.
	n.SetUniform("u_center", []float32{x, y})
}

func (n *fxRadialBlurNode) SetStrength(s float32) {
	// Set the strength of the blur.
	n.SetUniform("u_strength", s)
}
