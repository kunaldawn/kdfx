package blur

import (
	"kimg/pkg/context"
	"kimg/pkg/core"
	"kimg/pkg/node"
)

const radialBlurFS = `
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

// RadialBlurNode applies a radial (zoom) blur to the input texture.
type RadialBlurNode interface {
	node.Node
	// SetCenter sets the center of the blur (0.0 to 1.0).
	SetCenter(x, y float32)
	// SetStrength sets the strength of the blur.
	SetStrength(s float32)
}

type radialBlurNode struct {
	node.Node
}

// NewRadialBlurNode creates a new radial blur node.
func NewRadialBlurNode(ctx context.Context, width, height int) (RadialBlurNode, error) {
	base, err := node.NewBaseNode(ctx, width, height)
	if err != nil {
		return nil, err
	}

	program, err := core.NewShaderProgram(core.SimpleVS, radialBlurFS)
	if err != nil {
		base.Release()
		return nil, err
	}

	base.SetShaderProgram(program)

	n := &radialBlurNode{
		Node: base,
	}
	n.SetCenter(0.5, 0.5)
	n.SetStrength(0.1)

	return n, nil
}

func (n *radialBlurNode) SetCenter(x, y float32) {
	n.SetUniform("u_center", []float32{x, y})
}

func (n *radialBlurNode) SetStrength(s float32) {
	n.SetUniform("u_strength", s)
}
