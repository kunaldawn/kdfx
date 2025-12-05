package fxartistic

import (
	"kdfx/pkg/fxcontext"
	"kdfx/pkg/fxcore"
	"kdfx/pkg/fxnode"
)

// FXOilPaintFS is the fragment shader for oil paint effect.
// Simplified Kuwahara filter.
const FXOilPaintFS = `
precision mediump float;
varying vec2 v_texCoord;
uniform sampler2D u_texture;
uniform vec2 u_resolution;
uniform int u_radius;

void main() {
	vec2 src_size = u_resolution;
	vec2 uv = v_texCoord;
	float n = float((u_radius + 1) * (u_radius + 1));

	vec3 m[4];
	vec3 s[4];
	for (int k = 0; k < 4; ++k) {
		m[k] = vec3(0.0);
		s[k] = vec3(0.0);
	}

	for (int j = -u_radius; j <= 0; ++j)  {
		for (int i = -u_radius; i <= 0; ++i)  {
			vec3 c = texture2D(u_texture, uv + vec2(i,j) / src_size).rgb;
			m[0] += c;
			s[0] += c * c;
		}
	}

	for (int j = -u_radius; j <= 0; ++j)  {
		for (int i = 0; i <= u_radius; ++i)  {
			vec3 c = texture2D(u_texture, uv + vec2(i,j) / src_size).rgb;
			m[1] += c;
			s[1] += c * c;
		}
	}

	for (int j = 0; j <= u_radius; ++j)  {
		for (int i = 0; i <= u_radius; ++i)  {
			vec3 c = texture2D(u_texture, uv + vec2(i,j) / src_size).rgb;
			m[2] += c;
			s[2] += c * c;
		}
	}

	for (int j = 0; j <= u_radius; ++j)  {
		for (int i = -u_radius; i <= 0; ++i)  {
			vec3 c = texture2D(u_texture, uv + vec2(i,j) / src_size).rgb;
			m[3] += c;
			s[3] += c * c;
		}
	}

	float min_sigma2 = 1e+2;
	vec3 final_color = vec3(0.0);
	for (int k = 0; k < 4; ++k) {
		m[k] /= n;
		s[k] = abs(s[k] / n - m[k] * m[k]);

		float sigma2 = s[k].r + s[k].g + s[k].b;
		if (sigma2 < min_sigma2) {
			min_sigma2 = sigma2;
			final_color = m[k];
		}
	}

	gl_FragColor = vec4(final_color, 1.0);
}
`

// FXOilPaintNode applies an oil paint effect to the input texture.
type FXOilPaintNode interface {
	fxnode.FXNode
	// SetRadius sets the radius of the effect (1 to 10).
	SetRadius(radius int)
}

// fxOilPaintNode implements FXOilPaintNode.
type fxOilPaintNode struct {
	fxnode.FXNode
}

// NewFXOilPaintNode creates a new oil paint fxnode.
func NewFXOilPaintNode(ctx fxcontext.FXContext, width, height int) (FXOilPaintNode, error) {
	base, err := fxnode.NewFXBaseNode(ctx, width, height)
	if err != nil {
		return nil, err
	}

	program, err := fxcore.NewFXShaderProgram(fxcore.FXSimpleVS, FXOilPaintFS)
	if err != nil {
		base.Release()
		return nil, err
	}

	base.SetShaderProgram(program)

	n := &fxOilPaintNode{
		FXNode: base,
	}

	n.SetUniform("u_resolution", []float32{float32(width), float32(height)})
	n.SetRadius(4)

	return n, nil
}

func (n *fxOilPaintNode) SetRadius(radius int) {
	n.SetUniform("u_radius", radius)
}
