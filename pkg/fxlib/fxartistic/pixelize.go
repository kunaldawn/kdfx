package fxartistic

import (
	"kdfx/pkg/fxcontext"
	"kdfx/pkg/fxcore"
	"kdfx/pkg/fxnode"
)

// FXPixelizeFS is the fragment shader for pixelize effect.
const FXPixelizeFS = `
precision mediump float;
varying vec2 v_texCoord;
uniform sampler2D u_texture;
uniform vec2 u_resolution;
uniform float u_pixelSize;

void main() {
	vec2 d = vec2(u_pixelSize, u_pixelSize) / u_resolution;
	vec2 coord = d * floor(v_texCoord / d) + 0.5 * d;
	gl_FragColor = texture2D(u_texture, coord);
}
`

// FXPixelizeNode applies a pixelize effect to the input texture.
type FXPixelizeNode interface {
	fxnode.FXNode
	// SetPixelSize sets the size of the pixels (1.0 to >100.0).
	SetPixelSize(size float32)
}

// fxPixelizeNode implements FXPixelizeNode.
type fxPixelizeNode struct {
	fxnode.FXNode
}

// NewFXPixelizeNode creates a new pixelize fxnode.
func NewFXPixelizeNode(ctx fxcontext.FXContext, width, height int) (FXPixelizeNode, error) {
	base, err := fxnode.NewFXBaseNode(ctx, width, height)
	if err != nil {
		return nil, err
	}

	program, err := fxcore.NewFXShaderProgram(fxcore.FXSimpleVS, FXPixelizeFS)
	if err != nil {
		base.Release()
		return nil, err
	}

	base.SetShaderProgram(program)

	n := &fxPixelizeNode{
		FXNode: base,
	}

	n.SetUniform("u_resolution", []float32{float32(width), float32(height)})
	n.SetPixelSize(10.0)

	return n, nil
}

func (n *fxPixelizeNode) SetPixelSize(size float32) {
	n.SetUniform("u_pixelSize", size)
}
