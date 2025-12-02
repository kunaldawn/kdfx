package filters

import (
	"kimg/context"
	"kimg/core"
	"kimg/node"
)

// Simple 9-tap gaussian blur for demonstration.
// For production, a two-pass blur is better.
const gaussianBlurFS = `
precision mediump float;
varying vec2 v_texCoord;
uniform sampler2D u_texture;
uniform vec2 u_resolution;
uniform float u_radius;

void main() {
	vec4 color = vec4(0.0);
	vec2 off = vec2(u_radius) / u_resolution;
	
	// 3x3 kernel approximation (simplified)
	color += texture2D(u_texture, v_texCoord + vec2(-off.x, -off.y)) * 0.0625;
	color += texture2D(u_texture, v_texCoord + vec2(0.0,    -off.y)) * 0.125;
	color += texture2D(u_texture, v_texCoord + vec2(off.x,  -off.y)) * 0.0625;
	
	color += texture2D(u_texture, v_texCoord + vec2(-off.x, 0.0))    * 0.125;
	color += texture2D(u_texture, v_texCoord + vec2(0.0,    0.0))    * 0.25;
	color += texture2D(u_texture, v_texCoord + vec2(off.x,  0.0))    * 0.125;
	
	color += texture2D(u_texture, v_texCoord + vec2(-off.x, off.y))  * 0.0625;
	color += texture2D(u_texture, v_texCoord + vec2(0.0,    off.y))  * 0.125;
	color += texture2D(u_texture, v_texCoord + vec2(off.x,  off.y))  * 0.0625;
	
	gl_FragColor = color;
}
`

type GaussianBlurNode struct {
	*node.BaseNode
	Radius float32
}

func NewGaussianBlurNode(ctx context.Context, width, height int) (*GaussianBlurNode, error) {
	base, err := node.NewBaseNode(ctx, width, height)
	if err != nil {
		return nil, err
	}

	program, err := core.NewShaderProgram(simpleVS, gaussianBlurFS)
	if err != nil {
		base.Release()
		return nil, err
	}

	base.Program = program

	return &GaussianBlurNode{
		BaseNode: base,
		Radius:   1.0,
	}, nil
}

func (n *GaussianBlurNode) SetRadius(r float32) {
	n.SetUniform("u_radius", r)
	// We also need resolution.
	w, h := n.Context.GetSize()
	n.SetUniform("u_resolution", []float32{float32(w), float32(h)})
}

func (n *GaussianBlurNode) Release() {
	n.BaseNode.Release()
}
