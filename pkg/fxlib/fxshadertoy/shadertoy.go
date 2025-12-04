// Package fxshadertoy provides a node for running ShaderToy-compatible GLSL code.
package fxshadertoy

import (
	"kdfx/pkg/fxcontext"
	"kdfx/pkg/fxcore"
	"kdfx/pkg/fxnode"
)

// FXShaderToyNode allows running arbitrary GLSL code compatible with ShaderToy.
// It provides standard ShaderToy uniforms:
// - iResolution (vec3): viewport resolution (in pixels)
// - iTime (float): shader playback time (in seconds)
// - iChannel0 (sampler2D): input texture channel 0
type FXShaderToyNode interface {
	fxnode.FXNode
	// SetTime sets the current time for the shader (iTime).
	// This should be updated every frame to animate the shader.
	SetTime(t float32)
}

// fxShadertoyNode implements FXShaderToyNode.
type fxShadertoyNode struct {
	fxnode.FXNode
}

// FXShaderToyVS is the vertex fxShader for ShaderToy.
const FXShaderToyVS = `
attribute vec2 a_position;
varying vec2 fragCoord;
uniform vec2 iResolution;

uniform vec2 u_translation;
uniform vec2 u_scale;
uniform float u_rotation;

void main() {
	// Apply scaling
	vec2 scaledPos = a_position * u_scale;

	// Apply rotation
	float c = cos(u_rotation);
	float s = sin(u_rotation);
	vec2 rotatedPos = vec2(
		scaledPos.x * c - scaledPos.y * s,
		scaledPos.x * s + scaledPos.y * c
	);

	// Apply translation
	vec2 finalPos = rotatedPos + u_translation;

	gl_Position = vec4(finalPos, 0.0, 1.0);
	
	// Convert -1..1 to 0..Resolution
	// We use the original a_position for texture coordinates to map the full image onto the transformed quad
	fragCoord = (a_position * 0.5 + 0.5) * iResolution;
}
`

// Wrap user code with ShaderToy boilerplate
func wrapShaderToyCode(userCode string) string {
	// Inject standard ShaderToy uniforms and the main entry point.
	// The user code is expected to define mainImage(out vec4 fragColor, in vec2 fragCoord).
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

// NewFXShaderToyNode creates a new ShaderToy node with the provided GLSL code.
func NewFXShaderToyNode(ctx fxcontext.FXContext, width, height int, code string) (FXShaderToyNode, error) {
	base, err := fxnode.NewFXBaseNode(ctx, width, height)
	if err != nil {
		return nil, err
	}

	// Wrap the user code to make it compatible with our shader system.
	fullSource := wrapShaderToyCode(code)
	program, err := fxcore.NewFXShaderProgram(fxcore.FXSimpleVS, fullSource)
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
	// iResolution is a vec3 in ShaderToy (width, height, pixel_aspect_ratio).
	base.SetUniform("iResolution", []float32{float32(width), float32(height), 1.0})

	return &fxShadertoyNode{
		FXNode: base,
	}, nil
}

func (n *fxShadertoyNode) SetTime(t float32) {
	// Set the iTime uniform.
	n.SetUniform("iTime", t)
}
