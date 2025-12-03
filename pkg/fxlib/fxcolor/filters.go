package fxcolor

import (
	"kdfx/pkg/fxcontext"
	"kdfx/pkg/fxcore"
	"kdfx/pkg/fxnode"
)

// FXFilterMode represents the type of color filter.
type FXFilterMode int

const (
	FXFilterNone FXFilterMode = iota
	FXFilterInvert
	FXFilterSepia
	FXFilterGrayscale
	FXFilterThreshold
	FXFilterPosterize
)

const FXFiltersFS = `
precision mediump float;
varying vec2 v_texCoord;
uniform sampler2D u_texture;
uniform int u_mode;
uniform float u_param; // Threshold value or Posterize levels

void main() {
	vec4 color = texture2D(u_texture, v_texCoord);
	vec3 rgb = color.rgb;

	if (u_mode == 1) { // Invert
		rgb = 1.0 - rgb;
	} else if (u_mode == 2) { // Sepia
		vec3 sepia;
		sepia.r = dot(rgb, vec3(0.393, 0.769, 0.189));
		sepia.g = dot(rgb, vec3(0.349, 0.686, 0.168));
		sepia.b = dot(rgb, vec3(0.272, 0.534, 0.131));
		rgb = sepia;
	} else if (u_mode == 3) { // Grayscale
		float gray = dot(rgb, vec3(0.299, 0.587, 0.114));
		rgb = vec3(gray);
	} else if (u_mode == 4) { // Threshold
		float gray = dot(rgb, vec3(0.299, 0.587, 0.114));
		rgb = vec3(step(u_param, gray));
	} else if (u_mode == 5) { // Posterize
		float levels = max(2.0, u_param);
		rgb = floor(rgb * levels) / (levels - 1.0);
	}

	gl_FragColor = vec4(rgb, color.a);
}
`

// FXColorFilterNode applies artistic color filters to the input texture.
type FXColorFilterNode interface {
	fxnode.FXNode
	// SetMode sets the filter mode.
	SetMode(mode FXFilterMode)
	// SetParam sets the parameter for the filter (e.g., threshold or levels).
	SetParam(p float32)
}

type fxColorFilterNode struct {
	fxnode.FXNode
}

// NewFXColorFilterNode creates a new color filter fxnode.
func NewFXColorFilterNode(ctx fxcontext.FXContext, width, height int) (FXColorFilterNode, error) {
	base, err := fxnode.NewFXBaseNode(ctx, width, height)
	if err != nil {
		return nil, err
	}

	program, err := fxcore.NewFXShaderProgram(fxcore.FXSimpleVS, FXFiltersFS)
	if err != nil {
		base.Release()
		return nil, err
	}

	base.SetShaderProgram(program)

	n := &fxColorFilterNode{
		FXNode: base,
	}

	n.SetMode(FXFilterNone)
	n.SetParam(0.5) // Default param

	return n, nil
}

func (n *fxColorFilterNode) SetMode(mode FXFilterMode) {
	n.SetUniform("u_mode", int(mode))
}

func (n *fxColorFilterNode) SetParam(p float32) {
	n.SetUniform("u_param", p)
}
