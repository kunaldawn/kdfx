package fxartistic

import (
	"kdfx/pkg/fxcontext"
	"kdfx/pkg/fxcore"
	"kdfx/pkg/fxnode"
)

// FXEdgeDetectionFS is the fragment shader for edge detection (Sobel).
const FXEdgeDetectionFS = `
precision mediump float;
varying vec2 v_texCoord;
uniform sampler2D u_texture;
uniform vec2 u_resolution;
uniform float u_threshold;

float getLuminance(vec4 color) {
	return dot(color.rgb, vec3(0.299, 0.587, 0.114));
}

void main() {
	vec2 onePixel = vec2(1.0, 1.0) / u_resolution;
	
	float tl = getLuminance(texture2D(u_texture, v_texCoord + vec2(-onePixel.x, -onePixel.y)));
	float t  = getLuminance(texture2D(u_texture, v_texCoord + vec2(0.0, -onePixel.y)));
	float tr = getLuminance(texture2D(u_texture, v_texCoord + vec2(onePixel.x, -onePixel.y)));
	float l  = getLuminance(texture2D(u_texture, v_texCoord + vec2(-onePixel.x, 0.0)));
	float r  = getLuminance(texture2D(u_texture, v_texCoord + vec2(onePixel.x, 0.0)));
	float bl = getLuminance(texture2D(u_texture, v_texCoord + vec2(-onePixel.x, onePixel.y)));
	float b  = getLuminance(texture2D(u_texture, v_texCoord + vec2(0.0, onePixel.y)));
	float br = getLuminance(texture2D(u_texture, v_texCoord + vec2(onePixel.x, onePixel.y)));

	float gx = tl + 2.0*l + bl - tr - 2.0*r - br;
	float gy = tl + 2.0*t + tr - bl - 2.0*b - br;
	
	float g = sqrt(gx*gx + gy*gy);
	
	if (g < u_threshold) {
		g = 0.0;
	}

	gl_FragColor = vec4(vec3(g), 1.0);
}
`

// FXEdgeDetectionNode applies edge detection to the input texture.
type FXEdgeDetectionNode interface {
	fxnode.FXNode
	// SetThreshold sets the edge detection threshold (0.0 to 1.0).
	SetThreshold(threshold float32)
}

// fxEdgeDetectionNode implements FXEdgeDetectionNode.
type fxEdgeDetectionNode struct {
	fxnode.FXNode
}

// NewFXEdgeDetectionNode creates a new edge detection fxnode.
func NewFXEdgeDetectionNode(ctx fxcontext.FXContext, width, height int) (FXEdgeDetectionNode, error) {
	base, err := fxnode.NewFXBaseNode(ctx, width, height)
	if err != nil {
		return nil, err
	}

	program, err := fxcore.NewFXShaderProgram(fxcore.FXSimpleVS, FXEdgeDetectionFS)
	if err != nil {
		base.Release()
		return nil, err
	}

	base.SetShaderProgram(program)

	n := &fxEdgeDetectionNode{
		FXNode: base,
	}

	n.SetUniform("u_resolution", []float32{float32(width), float32(height)})
	n.SetThreshold(0.0)

	return n, nil
}

func (n *fxEdgeDetectionNode) SetThreshold(threshold float32) {
	n.SetUniform("u_threshold", threshold)
}
