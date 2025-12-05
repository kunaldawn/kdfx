package fxartistic

import (
	"kdfx/pkg/fxcontext"
	"kdfx/pkg/fxcore"
	"kdfx/pkg/fxnode"
)

// FXVignetteFS is the fragment shader for vignette effect.
const FXVignetteFS = `
precision mediump float;
varying vec2 v_texCoord;
uniform sampler2D u_texture;
uniform float u_radius;
uniform float u_softness;
uniform float u_opacity;

void main() {
	vec4 color = texture2D(u_texture, v_texCoord);
	vec2 position = v_texCoord - 0.5;
	float dist = length(position);
	
	float vignette = smoothstep(u_radius, u_radius - u_softness, dist);
	
	// Apply opacity to the vignette effect
	vignette = mix(1.0, vignette, u_opacity);

	gl_FragColor = vec4(color.rgb * vignette, color.a);
}
`

// FXVignetteNode applies a vignette effect to the input texture.
type FXVignetteNode interface {
	fxnode.FXNode
	// SetRadius sets the radius of the vignette (0.0 to 1.0).
	SetRadius(radius float32)
	// SetSoftness sets the softness of the vignette edge (0.0 to 1.0).
	SetSoftness(softness float32)
	// SetOpacity sets the opacity of the vignette (0.0 to 1.0).
	SetOpacity(opacity float32)
}

// fxVignetteNode implements FXVignetteNode.
type fxVignetteNode struct {
	fxnode.FXNode
}

// NewFXVignetteNode creates a new vignette fxnode.
func NewFXVignetteNode(ctx fxcontext.FXContext, width, height int) (FXVignetteNode, error) {
	base, err := fxnode.NewFXBaseNode(ctx, width, height)
	if err != nil {
		return nil, err
	}

	program, err := fxcore.NewFXShaderProgram(fxcore.FXSimpleVS, FXVignetteFS)
	if err != nil {
		base.Release()
		return nil, err
	}

	base.SetShaderProgram(program)

	n := &fxVignetteNode{
		FXNode: base,
	}

	n.SetRadius(0.75)
	n.SetSoftness(0.45)
	n.SetOpacity(0.5)

	return n, nil
}

func (n *fxVignetteNode) SetRadius(radius float32) {
	n.SetUniform("u_radius", radius)
}

func (n *fxVignetteNode) SetSoftness(softness float32) {
	n.SetUniform("u_softness", softness)
}

func (n *fxVignetteNode) SetOpacity(opacity float32) {
	n.SetUniform("u_opacity", opacity)
}
