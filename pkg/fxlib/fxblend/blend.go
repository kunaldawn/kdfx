package fxblend

import (
	"kdfx/pkg/fxcontext"
	"kdfx/pkg/fxcore"
	"kdfx/pkg/fxnode"
)

// FXBlendMode represents the blending mode used to combine textures.
type FXBlendMode int

const (
	FXBlendNormal FXBlendMode = iota
	FXBlendAdd
	FXBlendMultiply
	FXBlendScreen
	FXBlendOverlay
	FXBlendDarken
	FXBlendLighten
	FXBlendColorDodge
	FXBlendColorBurn
	FXBlendHardLight
	FXBlendSoftLight
	FXBlendDifference
	FXBlendExclusion
)

const FXBlendFS = `
precision mediump float;
varying vec2 v_texCoord;
uniform sampler2D u_texture1; // Base
uniform sampler2D u_texture2; // Blend
uniform float u_factor;       // Opacity
uniform int u_mode;

float blendAdd(float base, float blend) {
	return min(base + blend, 1.0);
}

float blendMultiply(float base, float blend) {
	return base * blend;
}

float blendScreen(float base, float blend) {
	return 1.0 - ((1.0 - base) * (1.0 - blend));
}

float blendOverlay(float base, float blend) {
	return base < 0.5 ? (2.0 * base * blend) : (1.0 - 2.0 * (1.0 - base) * (1.0 - blend));
}

float blendDarken(float base, float blend) {
	return min(base, blend);
}

float blendLighten(float base, float blend) {
	return max(base, blend);
}

float blendColorDodge(float base, float blend) {
	return (blend == 1.0) ? blend : min(base / (1.0 - blend), 1.0);
}

float blendColorBurn(float base, float blend) {
	return (blend == 0.0) ? blend : max((1.0 - ((1.0 - base) / blend)), 0.0);
}

float blendHardLight(float base, float blend) {
	return blendOverlay(blend, base);
}

float blendSoftLight(float base, float blend) {
	return (blend < 0.5) ? (2.0 * base * blend + base * base * (1.0 - 2.0 * blend)) : (sqrt(base) * (2.0 * blend - 1.0) + 2.0 * base * (1.0 - blend));
}

float blendDifference(float base, float blend) {
	return abs(base - blend);
}

float blendExclusion(float base, float blend) {
	return base + blend - 2.0 * base * blend;
}

void main() {
	vec4 c1 = texture2D(u_texture1, v_texCoord);
	vec4 c2 = texture2D(u_texture2, v_texCoord);
	
	vec3 base = c1.rgb;
	vec3 blend = c2.rgb;
	vec3 result = base;

	if (u_mode == 1) { // Add
		result = vec3(blendAdd(base.r, blend.r), blendAdd(base.g, blend.g), blendAdd(base.b, blend.b));
	} else if (u_mode == 2) { // Multiply
		result = vec3(blendMultiply(base.r, blend.r), blendMultiply(base.g, blend.g), blendMultiply(base.b, blend.b));
	} else if (u_mode == 3) { // Screen
		result = vec3(blendScreen(base.r, blend.r), blendScreen(base.g, blend.g), blendScreen(base.b, blend.b));
	} else if (u_mode == 4) { // Overlay
		result = vec3(blendOverlay(base.r, blend.r), blendOverlay(base.g, blend.g), blendOverlay(base.b, blend.b));
	} else if (u_mode == 5) { // Darken
		result = vec3(blendDarken(base.r, blend.r), blendDarken(base.g, blend.g), blendDarken(base.b, blend.b));
	} else if (u_mode == 6) { // Lighten
		result = vec3(blendLighten(base.r, blend.r), blendLighten(base.g, blend.g), blendLighten(base.b, blend.b));
	} else if (u_mode == 7) { // ColorDodge
		result = vec3(blendColorDodge(base.r, blend.r), blendColorDodge(base.g, blend.g), blendColorDodge(base.b, blend.b));
	} else if (u_mode == 8) { // ColorBurn
		result = vec3(blendColorBurn(base.r, blend.r), blendColorBurn(base.g, blend.g), blendColorBurn(base.b, blend.b));
	} else if (u_mode == 9) { // HardLight
		result = vec3(blendHardLight(base.r, blend.r), blendHardLight(base.g, blend.g), blendHardLight(base.b, blend.b));
	} else if (u_mode == 10) { // SoftLight
		result = vec3(blendSoftLight(base.r, blend.r), blendSoftLight(base.g, blend.g), blendSoftLight(base.b, blend.b));
	} else if (u_mode == 11) { // Difference
		result = vec3(blendDifference(base.r, blend.r), blendDifference(base.g, blend.g), blendDifference(base.b, blend.b));
	} else if (u_mode == 12) { // Exclusion
		result = vec3(blendExclusion(base.r, blend.r), blendExclusion(base.g, blend.g), blendExclusion(base.b, blend.b));
	} else { // Normal (0)
		result = blend;
	}

	// Apply opacity (factor)
	// Interpolate between base color and blended result
	gl_FragColor = vec4(mix(base, result, u_factor), c1.a);
}
`

// FXBlendNode blends two input textures.
type FXBlendNode interface {
	fxnode.FXNode
	// SetFactor sets the opacity of the blend (0.0 to 1.0).
	SetFactor(f float32)
	// SetMode sets the blending mode.
	SetMode(mode FXBlendMode)
	// SetInput1 sets the base texture input.
	SetInput1(input fxnode.FXInput)
	// SetInput2 sets the blend texture input.
	SetInput2(input fxnode.FXInput)
}

type fxBlendNode struct {
	fxnode.FXNode
}

// NewFXBlendNode creates a new blend fxnode.
func NewFXBlendNode(ctx fxcontext.FXContext, width, height int) (FXBlendNode, error) {
	base, err := fxnode.NewFXBaseNode(ctx, width, height)
	if err != nil {
		return nil, err
	}

	program, err := fxcore.NewFXShaderProgram(fxcore.FXSimpleVS, FXBlendFS)
	if err != nil {
		base.Release()
		return nil, err
	}

	base.SetShaderProgram(program)

	n := &fxBlendNode{
		FXNode: base,
	}
	n.SetFactor(1.0)
	n.SetMode(FXBlendNormal)

	return n, nil
}

func (n *fxBlendNode) SetFactor(f float32) {
	n.SetUniform("u_factor", f)
}

func (n *fxBlendNode) SetMode(mode FXBlendMode) {
	n.SetUniform("u_mode", int(mode))
}

func (n *fxBlendNode) SetInput1(input fxnode.FXInput) {
	n.SetInput("u_texture1", input)
}

func (n *fxBlendNode) SetInput2(input fxnode.FXInput) {
	n.SetInput("u_texture2", input)
}
