package fxcolor

import (
	"kdfx/pkg/fxcontext"
	"kdfx/pkg/fxcore"
	"kdfx/pkg/fxnode"
)

// FXLevelsFS is the fragment shader for levels adjustment.
const FXLevelsFS = `
precision mediump float;
varying vec2 v_texCoord;
uniform sampler2D u_texture;

uniform float u_inBlack;
uniform float u_inWhite;
uniform float u_outBlack;
uniform float u_outWhite;
uniform float u_gamma;

void main() {
	vec4 color = texture2D(u_texture, v_texCoord);
	vec3 rgb = color.rgb;

	// Input levels
	rgb = (rgb - u_inBlack) / (u_inWhite - u_inBlack);

	// Gamma
	rgb = pow(max(rgb, 0.0), vec3(1.0 / u_gamma));

	// Output levels
	rgb = rgb * (u_outWhite - u_outBlack) + u_outBlack;

	gl_FragColor = vec4(clamp(rgb, 0.0, 1.0), color.a);
}
`

// FXLevelsNode applies levels adjustment to the input texture.
type FXLevelsNode interface {
	fxnode.FXNode
	// SetInputLevels sets the input black and white points (0.0 to 1.0).
	SetInputLevels(black, white float32)
	// SetOutputLevels sets the output black and white points (0.0 to 1.0).
	SetOutputLevels(black, white float32)
	// SetGamma sets the gamma correction (0.1 to >1.0, default 1.0).
	SetGamma(gamma float32)
}

// fxLevelsNode implements FXLevelsNode.
type fxLevelsNode struct {
	fxnode.FXNode
}

// NewFXLevelsNode creates a new levels fxnode.
func NewFXLevelsNode(ctx fxcontext.FXContext, width, height int) (FXLevelsNode, error) {
	base, err := fxnode.NewFXBaseNode(ctx, width, height)
	if err != nil {
		return nil, err
	}

	program, err := fxcore.NewFXShaderProgram(fxcore.FXSimpleVS, FXLevelsFS)
	if err != nil {
		base.Release()
		return nil, err
	}

	base.SetShaderProgram(program)

	n := &fxLevelsNode{
		FXNode: base,
	}

	// Set defaults
	n.SetInputLevels(0.0, 1.0)
	n.SetOutputLevels(0.0, 1.0)
	n.SetGamma(1.0)

	return n, nil
}

func (n *fxLevelsNode) SetInputLevels(black, white float32) {
	n.SetUniform("u_inBlack", black)
	n.SetUniform("u_inWhite", white)
}

func (n *fxLevelsNode) SetOutputLevels(black, white float32) {
	n.SetUniform("u_outBlack", black)
	n.SetUniform("u_outWhite", white)
}

func (n *fxLevelsNode) SetGamma(gamma float32) {
	n.SetUniform("u_gamma", gamma)
}
