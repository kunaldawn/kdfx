package fxdistortion

import (
	"kdfx/pkg/fxcontext"
	"kdfx/pkg/fxcore"
	"kdfx/pkg/fxnode"
)

// FXTwirlFS is the fragment shader for twirl distortion.
const FXTwirlFS = `
precision mediump float;
varying vec2 v_texCoord;
uniform sampler2D u_texture;
uniform float u_angle;
uniform float u_radius;
uniform vec2 u_center;

void main() {
	vec2 uv = v_texCoord;
	vec2 delta = uv - u_center;
	float dist = length(delta);
	
	if (dist < u_radius) {
		float percent = (u_radius - dist) / u_radius;
		float theta = percent * percent * u_angle;
		float s = sin(theta);
		float c = cos(theta);
		uv = vec2(dot(delta, vec2(c, -s)), dot(delta, vec2(s, c))) + u_center;
	}

	gl_FragColor = texture2D(u_texture, uv);
}
`

// FXTwirlNode applies a twirl distortion to the input texture.
type FXTwirlNode interface {
	fxnode.FXNode
	// SetAngle sets the rotation angle in radians.
	SetAngle(angle float32)
	// SetRadius sets the radius of the effect (0.0 to 1.0).
	SetRadius(radius float32)
	// SetCenter sets the center point of the twirl (0.0 to 1.0).
	SetCenter(x, y float32)
}

// fxTwirlNode implements FXTwirlNode.
type fxTwirlNode struct {
	fxnode.FXNode
}

// NewFXTwirlNode creates a new twirl fxnode.
func NewFXTwirlNode(ctx fxcontext.FXContext, width, height int) (FXTwirlNode, error) {
	base, err := fxnode.NewFXBaseNode(ctx, width, height)
	if err != nil {
		return nil, err
	}

	program, err := fxcore.NewFXShaderProgram(fxcore.FXSimpleVS, FXTwirlFS)
	if err != nil {
		base.Release()
		return nil, err
	}

	base.SetShaderProgram(program)

	n := &fxTwirlNode{
		FXNode: base,
	}

	n.SetAngle(3.14) // 180 degrees
	n.SetRadius(0.5)
	n.SetCenter(0.5, 0.5)

	return n, nil
}

func (n *fxTwirlNode) SetAngle(angle float32) {
	n.SetUniform("u_angle", angle)
}

func (n *fxTwirlNode) SetRadius(radius float32) {
	n.SetUniform("u_radius", radius)
}

func (n *fxTwirlNode) SetCenter(x, y float32) {
	n.SetUniform("u_center", []float32{x, y})
}
