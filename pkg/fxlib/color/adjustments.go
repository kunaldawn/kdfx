package color

import (
	"kimg/pkg/context"
	"kimg/pkg/core"
	"kimg/pkg/node"
)

const adjustmentsFS = `
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

type ColorAdjustmentNode interface {
	node.Node
	SetBrightness(b float32) // -1.0 to 1.0, default 0.0
	SetContrast(c float32)   // 0.0 to >1.0, default 1.0
	SetHue(h float32)        // -0.5 to 0.5 (rotates hue), default 0.0
	SetSaturation(s float32) // 0.0 to >1.0, default 1.0
	SetGamma(g float32)      // 0.1 to >1.0, default 1.0
	SetExposure(e float32)   // 0.0 to >1.0, default 1.0
}

type colorAdjustmentNode struct {
	node.Node
}

func NewColorAdjustmentNode(ctx context.Context, width, height int) (ColorAdjustmentNode, error) {
	base, err := node.NewBaseNode(ctx, width, height)
	if err != nil {
		return nil, err
	}

	program, err := core.NewShaderProgram(core.SimpleVS, adjustmentsFS)
	if err != nil {
		base.Release()
		return nil, err
	}

	base.SetShaderProgram(program)

	n := &colorAdjustmentNode{
		Node: base,
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

func (n *colorAdjustmentNode) SetBrightness(b float32) {
	n.SetUniform("u_brightness", b)
}

func (n *colorAdjustmentNode) SetContrast(c float32) {
	n.SetUniform("u_contrast", c)
}

func (n *colorAdjustmentNode) SetHue(h float32) {
	n.SetUniform("u_hue", h)
}

func (n *colorAdjustmentNode) SetSaturation(s float32) {
	n.SetUniform("u_saturation", s)
}

func (n *colorAdjustmentNode) SetGamma(g float32) {
	n.SetUniform("u_gamma", g)
}

func (n *colorAdjustmentNode) SetExposure(e float32) {
	n.SetUniform("u_exposure", e)
}
