package filters

import (
	"fmt"
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

type BrightnessContrastNode struct {
	*node.BaseNode
	Program    *core.ShaderProgram
	Brightness float32
	Contrast   float32
	Quad       *core.Quad
}

func NewBrightnessContrastNode(ctx context.Context, width, height int) (*BrightnessContrastNode, error) {
	base, err := node.NewBaseNode(ctx, width, height)
	if err != nil {
		return nil, err
	}

	program, err := core.NewShaderProgram(simpleVS, brightnessContrastFS)
	if err != nil {
		base.Release()
		return nil, err
	}

	return &BrightnessContrastNode{
		BaseNode:   base,
		Program:    program,
		Brightness: 0.0,
		Contrast:   1.0,
		Quad:       core.NewQuad(),
	}, nil
}

func (n *BrightnessContrastNode) SetBrightness(b float32) {
	if n.Brightness != b {
		n.Brightness = b
		n.Dirty = true
	}
}

func (n *BrightnessContrastNode) SetContrast(c float32) {
	if n.Contrast != c {
		n.Contrast = c
		n.Dirty = true
	}
}

func (n *BrightnessContrastNode) Process(ctx context.Context) error {
	if err := n.ProcessInputs(ctx); err != nil {
		return err
	}

	if !n.CheckDirty() {
		return nil
	}

	input := n.GetInput("image")
	if input == nil {
		return fmt.Errorf("missing input 'image'")
	}

	tex := input.GetTexture()
	if tex == nil {
		return fmt.Errorf("input 'image' has no texture")
	}

	n.Output.Bind()
	n.Program.Use()

	tex.Bind()
	n.Program.SetUniform1i("u_texture", 0)
	n.Program.SetUniform1f("u_brightness", n.Brightness)
	n.Program.SetUniform1f("u_contrast", n.Contrast)

	posLoc := n.Program.GetAttribLocation("a_position")
	texLoc := n.Program.GetAttribLocation("a_texCoord")

	n.Quad.Draw(posLoc, texLoc)

	tex.Unbind()
	n.Output.Unbind()

	return nil
}

func (n *BrightnessContrastNode) Release() {
	n.BaseNode.Release()
	n.Program.Release()
	n.Quad.Release()
}
