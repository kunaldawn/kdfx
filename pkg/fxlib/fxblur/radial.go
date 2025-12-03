package fxblur

import (
	"kdfx/pkg/fxcontext"
	"kdfx/pkg/fxcore"
	"kdfx/pkg/fxnode"
)

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
	// SetCenter sets the center of the blur (0.0 to 1.0).
	SetCenter(x, y float32)
	// SetStrength sets the strength of the fxblur.
	SetStrength(s float32)
}

type fxRadialBlurNode struct {
	fxnode.FXNode
}

// NewFXRadialBlurNode creates a new radial blur fxnode.
func NewFXRadialBlurNode(ctx fxcontext.FXContext, width, height int) (FXRadialBlurNode, error) {
	base, err := fxnode.NewFXBaseNode(ctx, width, height)
	if err != nil {
		return nil, err
	}

	program, err := fxcore.NewFXShaderProgram(fxcore.FXSimpleVS, FXRadialBlurFS)
	if err != nil {
		base.Release()
		return nil, err
	}

	base.SetShaderProgram(program)

	n := &fxRadialBlurNode{
		FXNode: base,
	}
	n.SetCenter(0.5, 0.5)
	n.SetStrength(0.1)

	return n, nil
}

func (n *fxRadialBlurNode) SetCenter(x, y float32) {
	n.SetUniform("u_center", []float32{x, y})
}

func (n *fxRadialBlurNode) SetStrength(s float32) {
	n.SetUniform("u_strength", s)
}
