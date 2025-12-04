package fxnode

import (
	"kdfx/pkg/fxcontext"
	"kdfx/pkg/fxcore"
)

// fxBaseNode implements common logic for Nodes.
type fxBaseNode struct {
	inputs   map[string]FXInput
	uniforms map[string]interface{}
	output   fxcore.FXFramebuffer
	program  fxcore.FXShaderProgram
	quad     fxcore.FXQuad
	dirty    bool
	context  fxcontext.FXContext

	// Transformations
	posX, posY     float32
	scaleX, scaleY float32
	rotation       float32
}

// NewFXBaseNode initializes a FXBaseNode.
func NewFXBaseNode(ctx fxcontext.FXContext, width, height int) (FXNode, error) {
	fbo, err := fxcore.NewFXFramebuffer(width, height)
	if err != nil {
		return nil, err
	}
	return &fxBaseNode{
		inputs:   make(map[string]FXInput),
		uniforms: make(map[string]interface{}),
		output:   fbo,
		quad:     fxcore.NewFXQuad(),
		dirty:    true,
		context:  ctx,
		scaleX:   1.0,
		scaleY:   1.0,
	}, nil
}

func (n *fxBaseNode) SetInput(name string, input FXInput) {
	n.inputs[name] = input
	n.dirty = true
}

func (n *fxBaseNode) GetInput(name string) FXInput {
	return n.inputs[name]
}

func (n *fxBaseNode) GetFramebuffer() fxcore.FXFramebuffer {
	return n.output
}

func (n *fxBaseNode) SetUniform(name string, value interface{}) {
	n.uniforms[name] = value
	n.dirty = true
}

func (n *fxBaseNode) SetPosition(x, y float32) {
	n.posX = x
	n.posY = y
	n.dirty = true
}

func (n *fxBaseNode) SetSize(w, h float32) {
	n.scaleX = w
	n.scaleY = h
	n.dirty = true
}

func (n *fxBaseNode) SetRotation(angle float32) {
	n.rotation = angle
	n.dirty = true
}

func (n *fxBaseNode) SetShaderProgram(program fxcore.FXShaderProgram) {
	n.program = program
}

func (n *fxBaseNode) UpdateTransformationUniforms(program fxcore.FXShaderProgram) {
	program.SetUniform2f("u_translation", n.posX, n.posY)
	program.SetUniform2f("u_scale", n.scaleX, n.scaleY)
	program.SetUniform1f("u_rotation", n.rotation)
}

func (n *fxBaseNode) GetTexture() fxcore.FXTexture {
	return n.output.GetTexture()
}

func (n *fxBaseNode) IsDirty() bool {
	if n.dirty {
		return true
	}
	for _, input := range n.inputs {
		if input.IsDirty() {
			return true
		}
	}
	return false
}

func (n *fxBaseNode) Release() {
	if n.output != nil {
		n.output.Release()
	}
	if n.quad != nil {
		n.quad.Release()
	}
	if n.program != nil {
		n.program.Release()
	}
}

// CheckDirty checks if processing is needed.
func (n *fxBaseNode) CheckDirty() bool {
	isDirty := n.IsDirty()
	if isDirty {
		n.dirty = false
	}
	return isDirty
}

// Process executes the node's operation.
func (n *fxBaseNode) Process(ctx fxcontext.FXContext) error {
	// 1. Process Inputs
	if err := n.ProcessInputs(ctx); err != nil {
		return err
	}

	// 2. Check Dirty
	if !n.CheckDirty() {
		return nil
	}

	// 3. Setup Render
	n.output.Bind()
	if n.program != nil {
		n.program.Use()

		// 4. Bind Inputs
		textureUnit := 0
		for name, input := range n.inputs {
			tex := input.GetTexture()
			if tex != nil {
				tex.BindToUnit(textureUnit)
				n.program.SetUniform1i(name, int32(textureUnit))
				textureUnit++
			}
		}

		// 5. Set Uniforms
		for name, value := range n.uniforms {
			switch v := value.(type) {
			case float32:
				n.program.SetUniform1f(name, v)
			case int:
				n.program.SetUniform1i(name, int32(v))
			case int32:
				n.program.SetUniform1i(name, v)
			case []float32:
				if len(v) == 2 {
					n.program.SetUniform2f(name, v[0], v[1])
				} else if len(v) == 3 {
					n.program.SetUniform3f(name, v[0], v[1], v[2])
				}
			}
		}

		// 6. Set Transformation Uniforms
		n.UpdateTransformationUniforms(n.program)

		// 7. Draw
		if n.quad != nil {
			posLoc := n.program.GetAttribLocation("a_position")
			texLoc := n.program.GetAttribLocation("a_texCoord")
			n.quad.Draw(posLoc, texLoc)
		}
	}

	n.output.Unbind()
	return nil
}

// ProcessInputs ensures that all input nodes are processed.
func (n *fxBaseNode) ProcessInputs(ctx fxcontext.FXContext) error {
	for _, input := range n.inputs {
		if node, ok := input.(FXNode); ok {
			if err := node.Process(ctx); err != nil {
				return err
			}
		}
	}
	return nil
}
