package blur

import (
	"fmt"
	"kimg/pkg/context"
	"kimg/pkg/core"
	"kimg/pkg/node"
)

const gaussianFS = `
precision mediump float;
varying vec2 v_texCoord;
uniform sampler2D u_texture;
uniform vec2 u_resolution;
uniform float u_radius;
uniform vec2 u_direction; // (1, 0) for horizontal, (0, 1) for vertical

// 9-tap Gaussian kernel weights
// 0.227027, 0.1945946, 0.1216216, 0.054054, 0.016216
void main() {
	vec2 off = vec2(1.0) / u_resolution * u_direction;
	vec4 color = texture2D(u_texture, v_texCoord) * 0.227027;
	
	color += texture2D(u_texture, v_texCoord + off * 1.0 * u_radius) * 0.1945946;
	color += texture2D(u_texture, v_texCoord - off * 1.0 * u_radius) * 0.1945946;
	
	color += texture2D(u_texture, v_texCoord + off * 2.0 * u_radius) * 0.1216216;
	color += texture2D(u_texture, v_texCoord - off * 2.0 * u_radius) * 0.1216216;
	
	color += texture2D(u_texture, v_texCoord + off * 3.0 * u_radius) * 0.054054;
	color += texture2D(u_texture, v_texCoord - off * 3.0 * u_radius) * 0.054054;
	
	color += texture2D(u_texture, v_texCoord + off * 4.0 * u_radius) * 0.016216;
	color += texture2D(u_texture, v_texCoord - off * 4.0 * u_radius) * 0.016216;
	
	gl_FragColor = color;
}
`

// GaussianBlurNode applies a Gaussian blur to the input texture.
type GaussianBlurNode interface {
	node.Node
	// SetRadius sets the blur radius.
	SetRadius(r float32)
}

type gaussianBlurNode struct {
	node.Node
	tempFB  core.Framebuffer
	ctx     context.Context
	program core.ShaderProgram
	radius  float32
}

// NewGaussianBlurNode creates a new Gaussian blur node.
func NewGaussianBlurNode(ctx context.Context, width, height int) (GaussianBlurNode, error) {
	base, err := node.NewBaseNode(ctx, width, height)
	if err != nil {
		return nil, err
	}

	program, err := core.NewShaderProgram(core.SimpleVS, gaussianFS)
	if err != nil {
		base.Release()
		return nil, err
	}
	base.SetShaderProgram(program)

	// Create temporary framebuffer for two-pass blur
	tempFB, err := core.NewFramebuffer(width, height)
	if err != nil {
		base.Release()
		return nil, err
	}

	return &gaussianBlurNode{
		Node:    base,
		tempFB:  tempFB,
		ctx:     ctx,
		program: program,
	}, nil
}

func (n *gaussianBlurNode) SetRadius(r float32) {
	n.SetUniform("u_radius", r)
	n.radius = r
}

// Process overrides the default process to implement two-pass blur
func (n *gaussianBlurNode) Process(ctx context.Context) error {
	if !n.IsDirty() {
		return nil
	}

	// 1. Get Input
	input := n.GetInput("u_texture")
	if input == nil {
		return fmt.Errorf("missing input 'u_texture'")
	}

	// 2. Process Input if it's a Node
	if inputNode, ok := input.(node.Node); ok {
		if inputNode.IsDirty() {
			if err := inputNode.Process(ctx); err != nil {
				return err
			}
		}
	}
	inputTex := input.GetTexture()

	// 3. Setup Quad
	quad := core.NewQuad()
	defer quad.Release()

	// Get Attrib Locations
	posLoc := n.program.GetAttribLocation("a_position")
	texLoc := n.program.GetAttribLocation("a_texCoord")

	// 4. Pass 1: Horizontal Blur (Input -> TempFB)
	n.tempFB.Bind()
	w, h := n.tempFB.GetTexture().GetSize()
	n.ctx.Viewport(0, 0, w, h)
	n.program.Use()

	// Set Uniforms for Pass 1
	n.program.SetUniform2f("u_direction", 1.0, 0.0)
	n.program.SetUniform1f("u_radius", n.radius)
	n.program.SetUniform2f("u_resolution", float32(w), float32(h))

	// Bind Input Texture
	inputTex.BindToUnit(0)
	n.program.SetUniform1i("u_texture", 0)

	// Draw
	quad.Draw(posLoc, texLoc)

	// 5. Pass 2: Vertical Blur (TempFB -> OutputFB)
	outputFB := n.GetFramebuffer()
	outputFB.Bind()
	w, h = outputFB.GetTexture().GetSize()
	n.ctx.Viewport(0, 0, w, h)
	n.program.Use() // Ensure program is used (though it should be)

	// Set Uniforms for Pass 2
	n.program.SetUniform2f("u_direction", 0.0, 1.0)
	n.program.SetUniform1f("u_radius", n.radius)
	n.program.SetUniform2f("u_resolution", float32(w), float32(h))

	// Bind Temp Texture
	n.tempFB.GetTexture().BindToUnit(0)
	n.program.SetUniform1i("u_texture", 0)

	// Draw
	quad.Draw(posLoc, texLoc)

	return nil
}
