package node

import (
	"kimg/pkg/context"
	"kimg/pkg/core"
)

// baseNode implements common logic for Nodes.
type baseNode struct {
	inputs   map[string]Input
	uniforms map[string]interface{}
	output   core.Framebuffer
	program  core.ShaderProgram
	quad     core.Quad
	dirty    bool
	context  context.Context
}

// NewBaseNode initializes a BaseNode.
func NewBaseNode(ctx context.Context, width, height int) (Node, error) {
	fbo, err := core.NewFramebuffer(width, height)
	if err != nil {
		return nil, err
	}
	return &baseNode{
		inputs:   make(map[string]Input),
		uniforms: make(map[string]interface{}),
		output:   fbo,
		quad:     core.NewQuad(),
		dirty:    true,
		context:  ctx,
	}, nil
}

func (n *baseNode) SetInput(name string, input Input) {
	n.inputs[name] = input
	n.dirty = true
}

func (n *baseNode) SetUniform(name string, value interface{}) {
	n.uniforms[name] = value
	n.dirty = true
}

func (n *baseNode) SetShaderProgram(program core.ShaderProgram) {
	n.program = program
}

func (n *baseNode) GetInput(name string) Input {
	return n.inputs[name]
}

func (n *baseNode) GetTexture() core.Texture {
	return n.output.GetTexture()
}

func (n *baseNode) IsDirty() bool {
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

func (n *baseNode) Release() {
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
func (n *baseNode) CheckDirty() bool {
	isDirty := n.IsDirty()
	if isDirty {
		n.dirty = false
	}
	return isDirty
}

// Process executes the node's operation.
func (n *baseNode) Process(ctx context.Context) error {
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

		// 6. Draw
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
func (n *baseNode) ProcessInputs(ctx context.Context) error {
	for _, input := range n.inputs {
		if node, ok := input.(Node); ok {
			if err := node.Process(ctx); err != nil {
				return err
			}
		}
	}
	return nil
}
