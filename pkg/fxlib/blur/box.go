package blur

import (
	"kimg/pkg/context"
	"kimg/pkg/core"
	"kimg/pkg/node"
)

const boxBlurFS = `
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

type BoxBlurNode interface {
	node.Node
	SetRadius(r float32)
}

type boxBlurNode struct {
	node.Node
	ctx context.Context
}

func NewBoxBlurNode(ctx context.Context, width, height int) (BoxBlurNode, error) {
	base, err := node.NewBaseNode(ctx, width, height)
	if err != nil {
		return nil, err
	}

	program, err := core.NewShaderProgram(core.SimpleVS, boxBlurFS)
	if err != nil {
		base.Release()
		return nil, err
	}

	base.SetShaderProgram(program)

	return &boxBlurNode{
		Node: base,
		ctx:  ctx,
	}, nil
}

func (n *boxBlurNode) SetRadius(r float32) {
	n.SetUniform("u_radius", r)
	w, h := n.ctx.GetSize()
	n.SetUniform("u_resolution", []float32{float32(w), float32(h)})
}
