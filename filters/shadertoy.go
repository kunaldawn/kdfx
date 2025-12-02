package filters

import (
	"kimg/context"
	"kimg/core"
	"kimg/node"
)

// ShaderToyNode allows running arbitrary GLSL code compatible with ShaderToy.
// It provides iResolution, iTime, iMouse (partial), iChannel0..3.
type ShaderToyNode struct {
	*node.BaseNode
	Program *core.ShaderProgram
	Time    float32
}

const shaderToyVS = `
attribute vec2 a_position;
varying vec2 fragCoord;
uniform vec2 iResolution;

void main() {
	gl_Position = vec4(a_position, 0.0, 1.0);
	// Convert -1..1 to 0..Resolution
	fragCoord = (a_position * 0.5 + 0.5) * iResolution;
}
`

// Wrap user code with ShaderToy boilerplate
func wrapShaderToyCode(userCode string) string {
	return `
precision mediump float;
uniform vec3 iResolution;
uniform float iTime;
uniform sampler2D iChannel0;
// uniform sampler2D iChannel1; // TODO: Support more channels
// uniform vec4 iMouse; // TODO: Support mouse

` + userCode + `

void main() {
	mainImage(gl_FragColor, gl_FragCoord.xy);
}
`
}

func NewShaderToyNode(ctx context.Context, width, height int, code string) (*ShaderToyNode, error) {
	base, err := node.NewBaseNode(ctx, width, height)
	if err != nil {
		return nil, err
	}

	fullSource := wrapShaderToyCode(code)
	program, err := core.NewShaderProgram(simpleVS, fullSource) // Use simpleVS for standard UVs if user code uses texture coordinates, OR use shaderToyVS if they use fragCoord?
	// ShaderToy uses mainImage(out vec4 fragColor, in vec2 fragCoord).
	// fragCoord is in pixels (0.5 to resolution-0.5).
	// gl_FragCoord provides this automatically in Fragment Shader!
	// So we can use simpleVS (which just draws a quad) and rely on gl_FragCoord.

	if err != nil {
		base.Release()
		return nil, err
	}

	return &ShaderToyNode{
		BaseNode: base,
		Program:  program,
		Time:     0.0,
	}, nil
}

func (n *ShaderToyNode) SetTime(t float32) {
	if n.Time != t {
		n.Time = t
		n.Dirty = true
	}
}

func (n *ShaderToyNode) Process(ctx context.Context) error {
	if err := n.ProcessInputs(ctx); err != nil {
		return err
	}

	if !n.CheckDirty() {
		return nil
	}

	// Bind inputs to channels
	// For now only iChannel0
	input := n.GetInput("iChannel0")
	if input != nil {
		tex := input.GetTexture()
		if tex != nil {
			tex.Bind()
			n.Program.SetUniform1i("iChannel0", 0)
		}
	}

	n.Output.Bind()
	n.Program.Use()

	w, h := n.Context.GetSize()
	n.Program.SetUniform3f("iResolution", float32(w), float32(h), 1.0)
	n.Program.SetUniform1f("iTime", n.Time)

	// Draw Quad
	// We need a quad. BaseNode doesn't have one.
	// We should probably move Quad to BaseNode or create one here.
	// For now create one here.
	q := core.NewQuad()
	defer q.Release() // Inefficient to create/destroy every frame, but ok for now.

	posLoc := n.Program.GetAttribLocation("a_position")
	// texLoc := n.Program.GetAttribLocation("a_texCoord") // Not used in ShaderToy usually, but we might need to satisfy simpleVS

	// simpleVS expects a_texCoord. If we don't enable it, it might crash or warn.
	// Let's use a custom VS for ShaderToy that doesn't need texCoord if we want,
	// or just pass it.
	texLoc := n.Program.GetAttribLocation("a_texCoord")

	q.Draw(posLoc, texLoc)

	if input != nil {
		input.GetTexture().Unbind()
	}
	n.Output.Unbind()

	return nil
}

func (n *ShaderToyNode) Release() {
	n.BaseNode.Release()
	n.Program.Release()
}
