package shadertoy

import (
	"kdfx/pkg/context"
	"kdfx/pkg/core"
	"kdfx/pkg/node"
)

// ShaderToyNode allows running arbitrary GLSL code compatible with ShaderToy.
// It provides iResolution, iTime, iMouse (partial), iChannel0..3.
type ShaderToyNode interface {
	node.Node
	// SetTime sets the current time for the shader (iTime).
	SetTime(t float32)
}

type shaderToyNode struct {
	node.Node
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

// NewShaderToyNode creates a new ShaderToy node with the provided GLSL code.
func NewShaderToyNode(ctx context.Context, width, height int, code string) (ShaderToyNode, error) {
	base, err := node.NewBaseNode(ctx, width, height)
	if err != nil {
		return nil, err
	}

	fullSource := wrapShaderToyCode(code)
	program, err := core.NewShaderProgram(core.SimpleVS, fullSource)
	// ShaderToy uses mainImage(out vec4 fragColor, in vec2 fragCoord).
	// fragCoord is in pixels (0.5 to resolution-0.5).
	// gl_FragCoord provides this automatically in Fragment Shader!
	// So we can use simpleVS (which just draws a quad) and rely on gl_FragCoord.

	if err != nil {
		base.Release()
		return nil, err
	}

	base.SetShaderProgram(program)

	// Set initial resolution
	base.SetUniform("iResolution", []float32{float32(width), float32(height), 1.0})

	return &shaderToyNode{
		Node: base,
	}, nil
}

func (n *shaderToyNode) SetTime(t float32) {
	n.SetUniform("iTime", t)
}
