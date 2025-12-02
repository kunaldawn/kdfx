package blend

import (
	"kimg/pkg/context"
	"kimg/pkg/core"
	"kimg/pkg/node"
)

const blendFS = `
precision mediump float;
varying vec2 v_texCoord;
uniform sampler2D u_texture1;
uniform sampler2D u_texture2;
uniform float u_factor;

void main() {
	vec4 c1 = texture2D(u_texture1, v_texCoord);
	vec4 c2 = texture2D(u_texture2, v_texCoord);
	
	// Simple mix
	gl_FragColor = mix(c1, c2, u_factor);
}
`

type BlendNode interface {
	node.Node
	SetFactor(f float32)
	SetInput1(input node.Input)
	SetInput2(input node.Input)
}

type blendNode struct {
	node.Node
}

func NewBlendNode(ctx context.Context, width, height int) (BlendNode, error) {
	base, err := node.NewBaseNode(ctx, width, height)
	if err != nil {
		return nil, err
	}

	program, err := core.NewShaderProgram(core.SimpleVS, blendFS)
	if err != nil {
		base.Release()
		return nil, err
	}

	base.SetShaderProgram(program)

	n := &blendNode{
		Node: base,
	}
	n.SetFactor(0.5)

	return n, nil
}

func (n *blendNode) SetFactor(f float32) {
	n.SetUniform("u_factor", f)
}

func (n *blendNode) SetInput1(input node.Input) {
	n.SetInput("u_texture1", input)
}

func (n *blendNode) SetInput2(input node.Input) {
	n.SetInput("u_texture2", input)
}
