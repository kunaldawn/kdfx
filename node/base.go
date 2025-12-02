package node

import (
	"kimg/context"
	"kimg/core"
)

// BaseNode implements common logic for Nodes.
// Embed this in specific Node implementations.
type BaseNode struct {
	Inputs  map[string]Input
	Output  *core.Framebuffer
	Dirty   bool
	Context context.Context
}

// NewBaseNode initializes a BaseNode.
func NewBaseNode(ctx context.Context, width, height int) (*BaseNode, error) {
	fbo, err := core.NewFramebuffer(width, height)
	if err != nil {
		return nil, err
	}
	return &BaseNode{
		Inputs:  make(map[string]Input),
		Output:  fbo,
		Dirty:   true, // Start dirty to ensure first process
		Context: ctx,
	}, nil
}

func (n *BaseNode) SetInput(name string, input Input) {
	n.Inputs[name] = input
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
