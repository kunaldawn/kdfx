package filters

import (
	"fmt"
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
	Program *core.ShaderProgram
	Radius  float32
	Quad    *core.Quad
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

	return &GaussianBlurNode{
		BaseNode: base,
		Program:  program,
		Radius:   1.0,
		Quad:     core.NewQuad(),
	}, nil
}

func (n *GaussianBlurNode) SetRadius(r float32) {
	if n.Radius != r {
		n.Radius = r
		n.Dirty = true
	}
}

func (n *GaussianBlurNode) Process(ctx context.Context) error {
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
	n.Program.SetUniform1f("u_radius", n.Radius)

	w, h := n.Context.GetSize()
	// Pass resolution as vec2. We need SetUniform2f
	// For now, let's just use 1f for width and height separately or add helper.
	// I'll add SetUniform2f helper to shader.go first.
	// But since I can't edit shader.go in this tool call, I will assume it exists or use a workaround?
	// No, I must add it. I will add it in next step.
	// For now I will comment out resolution usage or hardcode it in shader?
	// Better: I will add SetUniform2f in next step and then update this file.
	// Or I can use gl.Uniform2f directly if I import gles2.
	// I will import gles2 here for now to be safe and fast.

	// Actually, I'll just add SetUniform2f to shader.go in a separate call before this one?
	// No, I can't reorder.
	// I will write this file assuming SetUniform2f exists, and then immediately add it to shader.go.
	n.Program.SetUniform2f("u_resolution", float32(w), float32(h))

	posLoc := n.Program.GetAttribLocation("a_position")
	texLoc := n.Program.GetAttribLocation("a_texCoord")

	n.Quad.Draw(posLoc, texLoc)

	tex.Unbind()
	n.Output.Unbind()

	return nil
}

func (n *GaussianBlurNode) Release() {
	n.BaseNode.Release()
	n.Program.Release()
	n.Quad.Release()
}
