package fxnode

import (
	"kdfx/pkg/fxcontext"
	"kdfx/pkg/fxcore"
)

// FXInput represents a source of a texture.
// It can be a FXNode or a standalone texture provider.
type FXInput interface {
	// GetTexture returns the output texture of this input.
	GetTexture() fxcore.FXTexture
	// IsDirty returns true if the input has changed and needs reprocessing.
	IsDirty() bool
}

// FXNode represents a processing unit in the fxPipeline.
type FXNode interface {
	FXInput // FXNode is also an FXInput

	// SetInput connects an input to a named slot.
	SetInput(name string, input FXInput)
	// GetInput returns the input connected to a named slot.
	GetInput(name string) FXInput
	// GetFramebuffer returns the node's output framebuffer.
	GetFramebuffer() fxcore.FXFramebuffer

	// SetUniform sets a uniform value for the node's shader.
	SetUniform(name string, value interface{})

	// SetPosition sets the position of the node in normalized coordinates (-1 to 1).
	SetPosition(x, y float32)
	// SetSize sets the size of the node (scale factor).
	SetSize(w, h float32)
	// SetRotation sets the rotation of the node in radians.
	SetRotation(angle float32)

	// SetShaderProgram sets the shader program for the fxnode.
	SetShaderProgram(program fxcore.FXShaderProgram)

	// UpdateTransformationUniforms sets the transformation uniforms on the given program.
	UpdateTransformationUniforms(program fxcore.FXShaderProgram)

	// Process executes the node's operation if necessary.
	// It should check IsDirty() and its inputs' IsDirty().
	Process(ctx fxcontext.FXContext) error

	// Release frees resources held by the fxnode.
	Release()
}

// FXGraph represents a collection of nodes and their connections.
type FXGraph interface {
	AddNode(name string, node FXNode)
	Connect(sourceNodeName, targetNodeName, inputSlot string) error
	// GetNode returns a node by name.
	GetNode(name string) FXNode
	// Release frees resources held by all nodes in the fxGraph.
	Release()
}

// FXPipeline manages the execution of a fxGraph.
type FXPipeline interface {
	Execute(outputNodeName string) error
	Release()
}
