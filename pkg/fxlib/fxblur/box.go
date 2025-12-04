// Package fxblur provides blur effects for the kdfx library.
package fxblur

import (
	"kdfx/pkg/fxcontext"
	"kdfx/pkg/fxcore"
	"kdfx/pkg/fxnode"
)

// FXBoxBlurFS is the fragment fxShader for box blur.
const FXBoxBlurFS = `
precision mediump float;
varying vec2 v_texCoord;
uniform sampler2D u_texture;
uniform vec2 u_resolution;
uniform float u_radius;

void main() {
	vec4 color = vec4(0.0);
	vec2 off = vec2(1.0) / u_resolution;
	float total = 0.0;
	
	// Simple 3x3 box blur for now, or loop?
	// Loops are expensive in WebGL 1.0 / GLES 2.0 if not constant.
	// Let's do a fixed 5x5 kernel.
	
	for (float x = -2.0; x <= 2.0; x++) {
		for (float y = -2.0; y <= 2.0; y++) {
			vec2 offset = vec2(x, y) * u_radius * off;
			color += texture2D(u_texture, v_texCoord + offset);
			total += 1.0;
		}
	}
	
	gl_FragColor = color / total;
}
`

// FXBoxBlurNode applies a box blur to the input texture.
type FXBoxBlurNode interface {
	fxnode.FXNode
	// SetRadius sets the blur radius in pixels.
	// A larger radius results in a stronger blur.
	SetRadius(r float32)
}

// fxBoxBlurNode implements FXBoxBlurNode.
type fxBoxBlurNode struct {
	fxnode.FXNode
	// ctx is the context used for rendering.
	ctx fxcontext.FXContext
}

// NewFXBoxBlurNode creates a new box blur fxnode.
func NewFXBoxBlurNode(ctx fxcontext.FXContext, width, height int) (FXBoxBlurNode, error) {
	base, err := fxnode.NewFXBaseNode(ctx, width, height)
	if err != nil {
		return nil, err
	}

	// Compile the shader program with the simple vertex shader and box blur fragment shader.
	program, err := fxcore.NewFXShaderProgram(fxcore.FXSimpleVS, FXBoxBlurFS)
	if err != nil {
		base.Release()
		return nil, err
	}

	// Set the shader program.
	base.SetShaderProgram(program)

	return &fxBoxBlurNode{
		FXNode: base,
		ctx:    ctx,
	}, nil
}

func (n *fxBoxBlurNode) SetRadius(r float32) {
	// Set the blur radius uniform.
	n.SetUniform("u_radius", r)
	// Set the resolution uniform, needed to calculate pixel offsets.
	w, h := n.ctx.GetSize()
	n.SetUniform("u_resolution", []float32{float32(w), float32(h)})
}
