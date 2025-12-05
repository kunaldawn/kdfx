package fxartistic

import (
	"kdfx/pkg/fxcontext"
	"kdfx/pkg/fxcore"
	"kdfx/pkg/fxnode"
)

// FXBloomFS is the fragment shader for bloom effect.
// This is a simplified single-pass bloom.
const FXBloomFS = `
precision mediump float;
varying vec2 v_texCoord;
uniform sampler2D u_texture;
uniform vec2 u_resolution;
uniform float u_threshold;
uniform float u_intensity;
uniform float u_blurSize;

void main() {
	vec4 color = texture2D(u_texture, v_texCoord);
	vec3 rgb = color.rgb;
	
	// Extract bright areas
	float brightness = dot(rgb, vec3(0.2126, 0.7152, 0.0722));
	vec3 brightColor = vec3(0.0);
	if(brightness > u_threshold) {
		brightColor = rgb;
	}

	// Simple box blur for the bloom glow
	vec2 onePixel = vec2(1.0, 1.0) / u_resolution;
	vec3 blur = vec3(0.0);
	float total = 0.0;
	float size = u_blurSize;
	
	// Small kernel for performance in this single pass example
	for (float x = -2.0; x <= 2.0; x++) {
		for (float y = -2.0; y <= 2.0; y++) {
			vec2 offset = vec2(x, y) * size * onePixel;
			vec3 c = texture2D(u_texture, v_texCoord + offset).rgb;
			float b = dot(c, vec3(0.2126, 0.7152, 0.0722));
			if(b > u_threshold) {
				blur += c;
			}
			total += 1.0;
		}
	}
	blur /= total;

	// Combine original with bloom
	rgb += blur * u_intensity;

	gl_FragColor = vec4(rgb, color.a);
}
`

// FXBloomNode applies a bloom effect to the input texture.
type FXBloomNode interface {
	fxnode.FXNode
	// SetThreshold sets the brightness threshold (0.0 to 1.0).
	SetThreshold(threshold float32)
	// SetIntensity sets the intensity of the bloom (0.0 to >1.0).
	SetIntensity(intensity float32)
	// SetBlurSize sets the size of the blur (0.0 to 10.0).
	SetBlurSize(size float32)
}

// fxBloomNode implements FXBloomNode.
type fxBloomNode struct {
	fxnode.FXNode
}

// NewFXBloomNode creates a new bloom fxnode.
func NewFXBloomNode(ctx fxcontext.FXContext, width, height int) (FXBloomNode, error) {
	base, err := fxnode.NewFXBaseNode(ctx, width, height)
	if err != nil {
		return nil, err
	}

	program, err := fxcore.NewFXShaderProgram(fxcore.FXSimpleVS, FXBloomFS)
	if err != nil {
		base.Release()
		return nil, err
	}

	base.SetShaderProgram(program)

	n := &fxBloomNode{
		FXNode: base,
	}

	n.SetUniform("u_resolution", []float32{float32(width), float32(height)})
	n.SetThreshold(0.7)
	n.SetIntensity(1.0)
	n.SetBlurSize(2.0)

	return n, nil
}

func (n *fxBloomNode) SetThreshold(threshold float32) {
	n.SetUniform("u_threshold", threshold)
}

func (n *fxBloomNode) SetIntensity(intensity float32) {
	n.SetUniform("u_intensity", intensity)
}

func (n *fxBloomNode) SetBlurSize(size float32) {
	n.SetUniform("u_blurSize", size)
}
