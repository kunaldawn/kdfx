package fxcolor

import (
	"kdfx/pkg/fxcontext"
	"kdfx/pkg/fxcore"
	"kdfx/pkg/fxnode"
)

const FXAdjustmentsFS = `
precision mediump float;
varying vec2 v_texCoord;
uniform sampler2D u_texture;

uniform float u_brightness;
uniform float u_contrast;
uniform float u_hue;
uniform float u_saturation;
uniform float u_gamma;
uniform float u_exposure;

// RGB to HSV conversion
vec3 rgb2hsv(vec3 c) {
	vec4 K = vec4(0.0, -1.0 / 3.0, 2.0 / 3.0, -1.0);
	vec4 p = mix(vec4(c.bg, K.wz), vec4(c.gb, K.xy), step(c.b, c.g));
	vec4 q = mix(vec4(p.xyw, c.r), vec4(c.r, p.yzx), step(p.x, c.r));

	float d = q.x - min(q.w, q.y);
	float e = 1.0e-10;
	return vec3(abs(q.z + (q.w - q.y) / (6.0 * d + e)), d / (q.x + e), q.x);
}

// HSV to RGB conversion
vec3 hsv2rgb(vec3 c) {
	vec4 K = vec4(1.0, 2.0 / 3.0, 1.0 / 3.0, 3.0);
	vec3 p = abs(fract(c.xxx + K.xyz) * 6.0 - K.www);
	return c.z * mix(K.xxx, clamp(p - K.xxx, 0.0, 1.0), c.y);
}

void main() {
	vec4 color = texture2D(u_texture, v_texCoord);
	vec3 rgb = color.rgb;

	// 1. Exposure
	rgb *= u_exposure;

	// 2. Brightness
	rgb += u_brightness;

	// 3. Contrast
	rgb = (rgb - 0.5) * u_contrast + 0.5;

	// 4. Hue & Saturation
	vec3 hsv = rgb2hsv(rgb);
	hsv.x += u_hue;
	hsv.y *= u_saturation;
	rgb = hsv2rgb(hsv);

	// 5. Gamma
	rgb = pow(rgb, vec3(1.0 / u_gamma));

	gl_FragColor = vec4(rgb, color.a);
}
`

// FXColorAdjustmentNode applies various color adjustments to the input texture.
type FXColorAdjustmentNode interface {
	fxnode.FXNode
	// SetBrightness sets the brightness (-1.0 to 1.0, default 0.0).
	SetBrightness(b float32)
	// SetContrast sets the contrast (0.0 to >1.0, default 1.0).
	SetContrast(c float32)
	// SetHue sets the hue rotation (-0.5 to 0.5, default 0.0).
	SetHue(h float32)
	// SetSaturation sets the saturation (0.0 to >1.0, default 1.0).
	SetSaturation(s float32)
	// SetGamma sets the gamma correction (0.1 to >1.0, default 1.0).
	SetGamma(g float32)
	// SetExposure sets the exposure (0.0 to >1.0, default 1.0).
	SetExposure(e float32)
}

// fxColorAdjustmentNode implements FXColorAdjustmentNode.
type fxColorAdjustmentNode struct {
	fxnode.FXNode
}

// NewFXColorAdjustmentNode creates a new color adjustment fxnode.
func NewFXColorAdjustmentNode(ctx fxcontext.FXContext, width, height int) (FXColorAdjustmentNode, error) {
	base, err := fxnode.NewFXBaseNode(ctx, width, height)
	if err != nil {
		return nil, err
	}

	program, err := fxcore.NewFXShaderProgram(fxcore.FXSimpleVS, FXAdjustmentsFS)
	if err != nil {
		base.Release()
		return nil, err
	}

	base.SetShaderProgram(program)

	n := &fxColorAdjustmentNode{
		FXNode: base,
	}

	// Set defaults
	n.SetBrightness(0.0)
	n.SetContrast(1.0)
	n.SetHue(0.0)
	n.SetSaturation(1.0)
	n.SetGamma(1.0)
	n.SetExposure(1.0)

	return n, nil
}

func (n *fxColorAdjustmentNode) SetBrightness(b float32) {
	n.SetUniform("u_brightness", b)
}

func (n *fxColorAdjustmentNode) SetContrast(c float32) {
	n.SetUniform("u_contrast", c)
}

func (n *fxColorAdjustmentNode) SetHue(h float32) {
	n.SetUniform("u_hue", h)
}

func (n *fxColorAdjustmentNode) SetSaturation(s float32) {
	n.SetUniform("u_saturation", s)
}

func (n *fxColorAdjustmentNode) SetGamma(g float32) {
	n.SetUniform("u_gamma", g)
}

func (n *fxColorAdjustmentNode) SetExposure(e float32) {
	n.SetUniform("u_exposure", e)
}
