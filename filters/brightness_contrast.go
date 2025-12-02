package filters

import (
	"kimg/context"
	"kimg/core"
	"kimg/node"
)

const brightnessContrastFS = `
precision mediump float;
varying vec2 v_texCoord;
uniform sampler2D u_texture;
uniform float u_brightness;
uniform float u_contrast;

void main() {
	vec4 color = texture2D(u_texture, v_texCoord);
	
	// Apply brightness
	color.rgb += u_brightness;
	
	// Apply contrast
	color.rgb = (color.rgb - 0.5) * u_contrast + 0.5;
	
	gl_FragColor = color;
}
`

const simpleVS = `
attribute vec2 a_position;
attribute vec2 a_texCoord;
varying vec2 v_texCoord;
void main() {
	gl_Position = vec4(a_position, 0.0, 1.0);
	v_texCoord = a_texCoord;
}
`

type BrightnessContrastNode interface {
	node.Node
	SetBrightness(b float32)
	SetContrast(c float32)
}

type brightnessContrastNode struct {
	node.Node
}

func NewBrightnessContrastNode(ctx context.Context, width, height int) (BrightnessContrastNode, error) {
	base, err := node.NewBaseNode(ctx, width, height)
	if err != nil {
		return nil, err
	}

	program, err := core.NewShaderProgram(simpleVS, brightnessContrastFS)
	if err != nil {
		base.Release()
		return nil, err
	}

	base.SetShaderProgram(program)

	n := &brightnessContrastNode{
		Node: base,
	}
	n.SetUniform("u_brightness", float32(0.0))
	n.SetUniform("u_contrast", float32(1.0))

	return n, nil
}

func (n *brightnessContrastNode) SetBrightness(b float32) {
	n.SetUniform("u_brightness", b)
}

func (n *brightnessContrastNode) SetContrast(c float32) {
	n.SetUniform("u_contrast", c)
}
