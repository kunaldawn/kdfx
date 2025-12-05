package fxblur

import (
	"kdfx/pkg/fxcontext"
	"kdfx/pkg/fxcore"
	"kdfx/pkg/fxnode"
)

// FXSharpenFS is the fragment shader for sharpening.
const FXSharpenFS = `
precision mediump float;
varying vec2 v_texCoord;
uniform sampler2D u_texture;
uniform vec2 u_resolution;
uniform float u_amount;

void main() {
	vec2 onePixel = vec2(1.0, 1.0) / u_resolution;
	vec4 color = texture2D(u_texture, v_texCoord);
	
	vec4 n = texture2D(u_texture, v_texCoord + vec2(0.0, -onePixel.y));
	vec4 s = texture2D(u_texture, v_texCoord + vec2(0.0, onePixel.y));
	vec4 e = texture2D(u_texture, v_texCoord + vec2(onePixel.x, 0.0));
	vec4 w = texture2D(u_texture, v_texCoord + vec2(-onePixel.x, 0.0));

	vec4 result = color + (color * 4.0 - n - s - e - w) * u_amount;
	
	gl_FragColor = vec4(result.rgb, color.a);
}
`

// FXSharpenNode applies a sharpen filter to the input texture.
type FXSharpenNode interface {
	fxnode.FXNode
	// SetAmount sets the sharpening amount (0.0 to >1.0).
	SetAmount(amount float32)
}

// fxSharpenNode implements FXSharpenNode.
type fxSharpenNode struct {
	fxnode.FXNode
}

// NewFXSharpenNode creates a new sharpen fxnode.
func NewFXSharpenNode(ctx fxcontext.FXContext, width, height int) (FXSharpenNode, error) {
	base, err := fxnode.NewFXBaseNode(ctx, width, height)
	if err != nil {
		return nil, err
	}

	program, err := fxcore.NewFXShaderProgram(fxcore.FXSimpleVS, FXSharpenFS)
	if err != nil {
		base.Release()
		return nil, err
	}

	base.SetShaderProgram(program)

	n := &fxSharpenNode{
		FXNode: base,
	}

	n.SetUniform("u_resolution", []float32{float32(width), float32(height)})
	n.SetAmount(0.5)

	return n, nil
}

func (n *fxSharpenNode) SetAmount(amount float32) {
	n.SetUniform("u_amount", amount)
}
