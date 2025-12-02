package node

import (
	"kimg/context"
	"kimg/core"
)

// BaseNode implements common logic for Nodes.
// Embed this in specific Node implementations.
type BaseNode struct {
	Inputs   map[string]Input
	Uniforms map[string]interface{}
	Output   *core.Framebuffer
	Program  *core.ShaderProgram
	Quad     *core.Quad
	Dirty    bool
	Context  context.Context
}

// NewBaseNode initializes a BaseNode.
func NewBaseNode(ctx context.Context, width, height int) (*BaseNode, error) {
	fbo, err := core.NewFramebuffer(width, height)
	if err != nil {
		return nil, err
	}
	return &BaseNode{
		Inputs:   make(map[string]Input),
		Uniforms: make(map[string]interface{}),
		Output:   fbo,
		Quad:     core.NewQuad(), // BaseNode now manages a default Quad
		Dirty:    true,
		Context:  ctx,
	}, nil
}

func (n *BaseNode) SetInput(name string, input Input) {
	n.Inputs[name] = input
	n.Dirty = true
}

func (n *BaseNode) SetUniform(name string, value interface{}) {
	n.Uniforms[name] = value
	n.Dirty = true
}

func (n *BaseNode) GetInput(name string) Input {
	return n.Inputs[name]
}

func (n *BaseNode) GetTexture() *core.Texture {
	return n.Output.Texture
}

func (n *BaseNode) IsDirty() bool {
	if n.Dirty {
		return true
	}
	for _, input := range n.Inputs {
		if input.IsDirty() {
			return true
		}
	}
	return false
}

func (n *BaseNode) Release() {
	if n.Output != nil {
		n.Output.Release()
	}
	if n.Quad != nil {
		n.Quad.Release()
	}
	// Program is usually shared or managed by the specific node implementation,
	// but if we move it here, we should release it if we own it.
	// For now, let's assume the specific node sets n.Program and might share it?
	// Actually, usually each node instance has its own program if it has unique state,
	// but programs are usually stateless and shared.
	// However, in this design, we are putting Program in BaseNode.
	// If we create a new Program for each node, we should release it.
	if n.Program != nil {
		n.Program.Release()
	}
}

// CheckDirty checks if processing is needed.
// Returns true if dirty, and resets the dirty flag of the node itself (but not inputs).
// Note: This logic might need refinement depending on how we want to propagate "consumed" dirty state.
// For now, we assume Process() calls this and if it returns true, we render.
func (n *BaseNode) CheckDirty() bool {
	isDirty := n.IsDirty()
	if isDirty {
		n.Dirty = false
	}
	return isDirty
}

// Process executes the node's operation.
// It handles dirty checking, input processing, binding, and rendering.
func (n *BaseNode) Process(ctx context.Context) error {
	// 1. Process Inputs
	if err := n.ProcessInputs(ctx); err != nil {
		return err
	}

	// 2. Check Dirty
	if !n.CheckDirty() {
		return nil
	}

	// 3. Setup Render
	n.Output.Bind()
	if n.Program != nil {
		n.Program.Use()

		// 4. Bind Inputs
		textureUnit := 0
		for name, input := range n.Inputs {
			tex := input.GetTexture()
			if tex != nil {
				tex.BindToUnit(textureUnit) // We need BindToUnit in Texture
				n.Program.SetUniform1i(name, int32(textureUnit))
				textureUnit++
			}
		}

		// 5. Set Uniforms
		for name, value := range n.Uniforms {
			switch v := value.(type) {
			case float32:
				n.Program.SetUniform1f(name, v)
			case int:
				n.Program.SetUniform1i(name, int32(v))
			case int32:
				n.Program.SetUniform1i(name, v)
			case []float32:
				if len(v) == 2 {
					n.Program.SetUniform2f(name, v[0], v[1])
				} else if len(v) == 3 {
					n.Program.SetUniform3f(name, v[0], v[1], v[2])
				}
				// Add more cases as needed
			}
		}

		// 6. Draw
		if n.Quad != nil {
			posLoc := n.Program.GetAttribLocation("a_position")
			texLoc := n.Program.GetAttribLocation("a_texCoord")
			n.Quad.Draw(posLoc, texLoc)
		}

		// Unbind textures? Not strictly necessary if we overwrite next time, but good practice.
	}

	n.Output.Unbind()
	return nil
}

// ProcessInputs ensures that all input nodes are processed.
func (n *BaseNode) ProcessInputs(ctx context.Context) error {
	for _, input := range n.Inputs {
		if node, ok := input.(Node); ok {
			if err := node.Process(ctx); err != nil {
				return err
			}
		}
	}
	return nil
}
