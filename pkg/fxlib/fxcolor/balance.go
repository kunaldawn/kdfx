package fxcolor

import (
	"kdfx/pkg/fxcontext"
	"kdfx/pkg/fxcore"
	"kdfx/pkg/fxnode"
)

// FXColorBalanceFS is the fragment shader for color balance.
const FXColorBalanceFS = `
precision mediump float;
varying vec2 v_texCoord;
uniform sampler2D u_texture;

uniform vec3 u_shadows;
uniform vec3 u_midtones;
uniform vec3 u_highlights;
uniform int u_preserveLuminosity;

// RGB to HSL conversion
vec3 rgb2hsl(vec3 c) {
	float h = 0.0;
	float s = 0.0;
	float l = 0.0;
	float r = c.r;
	float g = c.g;
	float b = c.b;
	float cMin = min(r, min(g, b));
	float cMax = max(r, max(g, b));

	l = (cMax + cMin) / 2.0;
	if (cMax > cMin) {
		float d = cMax - cMin;
		s = l > 0.5 ? d / (2.0 - cMax - cMin) : d / (cMax + cMin);
		if (cMax == r) {
			h = (g - b) / d + (g < b ? 6.0 : 0.0);
		} else if (cMax == g) {
			h = (b - r) / d + 2.0;
		} else {
			h = (r - g) / d + 4.0;
		}
		h /= 6.0;
	}
	return vec3(h, s, l);
}

// HSL to RGB conversion
float hue2rgb(float p, float q, float t) {
	if (t < 0.0) t += 1.0;
	if (t > 1.0) t -= 1.0;
	if (t < 1.0/6.0) return p + (q - p) * 6.0 * t;
	if (t < 1.0/2.0) return q;
	if (t < 2.0/3.0) return p + (q - p) * (2.0/3.0 - t) * 6.0;
	return p;
}

vec3 hsl2rgb(vec3 c) {
	vec3 rgb;
	if (c.y == 0.0) {
		rgb = vec3(c.z); // Achromatic
	} else {
		float q = c.z < 0.5 ? c.z * (1.0 + c.y) : c.z + c.y - c.z * c.y;
		float p = 2.0 * c.z - q;
		rgb.r = hue2rgb(p, q, c.x + 1.0/3.0);
		rgb.g = hue2rgb(p, q, c.x);
		rgb.b = hue2rgb(p, q, c.x - 1.0/3.0);
	}
	return rgb;
}

void main() {
	vec4 color = texture2D(u_texture, v_texCoord);
	vec3 rgb = color.rgb;

	// Calculate lightness for tonal ranges
	float lightness = rgb.r * 0.299 + rgb.g * 0.587 + rgb.b * 0.114;

	// Shadows (0.0 - 0.33), Midtones (0.33 - 0.66), Highlights (0.66 - 1.0)
	// Smooth interpolation
	float s = 1.0 - smoothstep(0.0, 0.5, lightness);
	float h = smoothstep(0.5, 1.0, lightness);
	float m = 1.0 - s - h;

	vec3 adjustment = u_shadows * s + u_midtones * m + u_highlights * h;
	vec3 newRGB = clamp(rgb + adjustment, 0.0, 1.0);

	if (u_preserveLuminosity == 1) {
		vec3 hslOriginal = rgb2hsl(rgb);
		vec3 hslNew = rgb2hsl(newRGB);
		newRGB = hsl2rgb(vec3(hslNew.x, hslNew.y, hslOriginal.z));
	}

	gl_FragColor = vec4(newRGB, color.a);
}
`

// FXColorBalanceNode applies color balance adjustment.
type FXColorBalanceNode interface {
	fxnode.FXNode
	// SetShadows sets the cyan/red, magenta/green, yellow/blue balance for shadows (-1.0 to 1.0).
	SetShadows(r, g, b float32)
	// SetMidtones sets the cyan/red, magenta/green, yellow/blue balance for midtones (-1.0 to 1.0).
	SetMidtones(r, g, b float32)
	// SetHighlights sets the cyan/red, magenta/green, yellow/blue balance for highlights (-1.0 to 1.0).
	SetHighlights(r, g, b float32)
	// SetPreserveLuminosity sets whether to preserve the original image luminosity.
	SetPreserveLuminosity(preserve bool)
}

// fxColorBalanceNode implements FXColorBalanceNode.
type fxColorBalanceNode struct {
	fxnode.FXNode
}

// NewFXColorBalanceNode creates a new color balance fxnode.
func NewFXColorBalanceNode(ctx fxcontext.FXContext, width, height int) (FXColorBalanceNode, error) {
	base, err := fxnode.NewFXBaseNode(ctx, width, height)
	if err != nil {
		return nil, err
	}

	program, err := fxcore.NewFXShaderProgram(fxcore.FXSimpleVS, FXColorBalanceFS)
	if err != nil {
		base.Release()
		return nil, err
	}

	base.SetShaderProgram(program)

	n := &fxColorBalanceNode{
		FXNode: base,
	}

	// Set defaults
	n.SetShadows(0, 0, 0)
	n.SetMidtones(0, 0, 0)
	n.SetHighlights(0, 0, 0)
	n.SetPreserveLuminosity(true)

	return n, nil
}

func (n *fxColorBalanceNode) SetShadows(r, g, b float32) {
	n.SetUniform("u_shadows", []float32{r, g, b})
}

func (n *fxColorBalanceNode) SetMidtones(r, g, b float32) {
	n.SetUniform("u_midtones", []float32{r, g, b})
}

func (n *fxColorBalanceNode) SetHighlights(r, g, b float32) {
	n.SetUniform("u_highlights", []float32{r, g, b})
}

func (n *fxColorBalanceNode) SetPreserveLuminosity(preserve bool) {
	val := 0
	if preserve {
		val = 1
	}
	n.SetUniform("u_preserveLuminosity", val)
}
