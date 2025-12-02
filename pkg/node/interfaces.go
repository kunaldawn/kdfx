package node

import (
	"kimg/pkg/context"
	"kimg/pkg/core"
)

// Input represents a source of a texture.
// It can be a Node or a standalone texture provider.
type Input interface {
	// GetTexture returns the output texture of this input.
	GetTexture() core.Texture
	// IsDirty returns true if the input has changed and needs reprocessing.
	IsDirty() bool
}

// Node represents a processing unit in the pipeline.
type Node interface {
	Input // Node is also an Input

	// SetInput connects an input to a named slot.
	SetInput(name string, input Input)
	// GetInput returns the input connected to a named slot.
	GetInput(name string) Input

	// SetUniform sets a uniform value for the node's shader.
	SetUniform(name string, value interface{})

	// SetShaderProgram sets the shader program for the node.
	SetShaderProgram(program core.ShaderProgram)

	// Process executes the node's operation if necessary.
	// It should check IsDirty() and its inputs' IsDirty().
	Process(ctx context.Context) error

	// Release frees resources held by the node.
	Release()
}

// Graph represents a collection of nodes and their connections.
type Graph interface {
	AddNode(name string, node Node)
	Connect(sourceNodeName, targetNodeName, inputSlot string) error
	GetNode(name string) Node
	Release()
}

// Pipeline manages the execution of a graph.
type Pipeline interface {
	Execute(outputNodeName string) error
	Release()
}
