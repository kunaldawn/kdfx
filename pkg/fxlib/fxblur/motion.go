package fxblur

import (
	"kdfx/pkg/fxcontext"
	"kdfx/pkg/fxcore"
	"kdfx/pkg/fxnode"
	"math"
)

const FXMotionBlurFS = `
precision mediump float;
varying vec2 v_texCoord;
uniform sampler2D u_texture;
uniform vec2 u_velocity; // Direction * Strength

void main() {
	vec4 color = vec4(0.0);
	float total = 0.0;
	
	// Sample 10 times along the velocity vector
	for (float t = 0.0; t <= 1.0; t += 0.1) {
		vec2 offset = u_velocity * (t - 0.5); // Center the blur? Or trail?
		// Let's do trail (0 to 1)
		offset = u_velocity * t;
		
		color += texture2D(u_texture, v_texCoord - offset); // Sample backwards for trail
		total += 1.0;
	}
	
	gl_FragColor = color / total;
}
`

// FXMotionBlurNode applies a directional motion blur to the input texture.
type FXMotionBlurNode interface {
	fxnode.FXNode
	// SetAngle sets the angle of the blur in degrees.
	SetAngle(degrees float32)
	// SetStrength sets the strength/length of the fxblur.
	SetStrength(s float32)
}

type fxMotionBlurNode struct {
	fxnode.FXNode
	angle    float32
	strength float32
}

// NewFXMotionBlurNode creates a new motion blur fxnode.
func NewFXMotionBlurNode(ctx fxcontext.FXContext, width, height int) (FXMotionBlurNode, error) {
	base, err := fxnode.NewFXBaseNode(ctx, width, height)
	if err != nil {
		return nil, err
	}

	program, err := fxcore.NewFXShaderProgram(fxcore.FXSimpleVS, FXMotionBlurFS)
	if err != nil {
		base.Release()
		return nil, err
	}

	base.SetShaderProgram(program)

	n := &fxMotionBlurNode{
		FXNode:   base,
		angle:    0.0,
		strength: 0.01,
	}
	n.updateVelocity()

	return n, nil
}

func (n *fxMotionBlurNode) SetAngle(degrees float32) {
	n.angle = degrees
	n.updateVelocity()
}

func (n *fxMotionBlurNode) SetStrength(s float32) {
	n.strength = s
	n.updateVelocity()
}

func (n *fxMotionBlurNode) updateVelocity() {
	rad := float64(n.angle) * math.Pi / 180.0
	vx := float32(math.Cos(rad)) * n.strength
	vy := float32(math.Sin(rad)) * n.strength
	n.SetUniform("u_velocity", []float32{vx, vy})
}
