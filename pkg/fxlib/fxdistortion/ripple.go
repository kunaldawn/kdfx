package fxdistortion

import (
	"kdfx/pkg/fxcontext"
	"kdfx/pkg/fxcore"
	"kdfx/pkg/fxnode"
)

// FXRippleFS is the fragment shader for ripple distortion.
const FXRippleFS = `
precision mediump float;
varying vec2 v_texCoord;
uniform sampler2D u_texture;
uniform float u_time;
uniform float u_amplitude;
uniform float u_frequency;
uniform float u_speed;

void main() {
	vec2 uv = v_texCoord;
	
	// Apply sine wave displacement
	uv.x += sin(uv.y * u_frequency + u_time * u_speed) * u_amplitude;
	uv.y += cos(uv.x * u_frequency + u_time * u_speed) * u_amplitude;

	gl_FragColor = texture2D(u_texture, uv);
}
`

// FXRippleNode applies a ripple distortion to the input texture.
type FXRippleNode interface {
	fxnode.FXNode
	// SetAmplitude sets the strength of the ripple (0.0 to 0.1).
	SetAmplitude(amplitude float32)
	// SetFrequency sets the frequency of the ripple (0.0 to 100.0).
	SetFrequency(frequency float32)
	// SetSpeed sets the animation speed (0.0 to 10.0).
	SetSpeed(speed float32)
	// SetTime sets the current time for animation.
	SetTime(time float32)
}

// fxRippleNode implements FXRippleNode.
type fxRippleNode struct {
	fxnode.FXNode
}

// NewFXRippleNode creates a new ripple fxnode.
func NewFXRippleNode(ctx fxcontext.FXContext, width, height int) (FXRippleNode, error) {
	base, err := fxnode.NewFXBaseNode(ctx, width, height)
	if err != nil {
		return nil, err
	}

	program, err := fxcore.NewFXShaderProgram(fxcore.FXSimpleVS, FXRippleFS)
	if err != nil {
		base.Release()
		return nil, err
	}

	base.SetShaderProgram(program)

	n := &fxRippleNode{
		FXNode: base,
	}

	n.SetAmplitude(0.01)
	n.SetFrequency(20.0)
	n.SetSpeed(1.0)
	n.SetTime(0.0)

	return n, nil
}

func (n *fxRippleNode) SetAmplitude(amplitude float32) {
	n.SetUniform("u_amplitude", amplitude)
}

func (n *fxRippleNode) SetFrequency(frequency float32) {
	n.SetUniform("u_frequency", frequency)
}

func (n *fxRippleNode) SetSpeed(speed float32) {
	n.SetUniform("u_speed", speed)
}

func (n *fxRippleNode) SetTime(time float32) {
	n.SetUniform("u_time", time)
}
