// Package fxnode defines the node graph system, including node interfaces, graph management, and pipeline execution.
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
// It can receive inputs, process them using a shader, and produce an output texture.
type FXNode interface {
	FXInput // FXNode is also an FXInput, allowing it to be connected to other nodes.

	// SetInput connects an input to a named slot.
	// The name usually corresponds to a sampler2D uniform in the shader.
	SetInput(name string, input FXInput)
	// GetInput returns the input connected to a named slot.
	GetInput(name string) FXInput
	// GetFramebuffer returns the node's output framebuffer.
	// This contains the result of the node's processing.
	GetFramebuffer() fxcore.FXFramebuffer

	// SetUniform sets a uniform value for the node's shader.
	// Supported types: float32, int, int32, []float32 (vec2, vec3).
	SetUniform(name string, value interface{})

	// SetPosition sets the position of the node in normalized coordinates (-1 to 1).
	// This affects the rendering of the node's quad.
	SetPosition(x, y float32)
	// SetSize sets the size of the node (scale factor).
	// This scales the node's quad.
	SetSize(w, h float32)
	// SetRotation sets the rotation of the node in radians.
	// This rotates the node's quad.
	SetRotation(angle float32)

	// SetShaderProgram sets the shader program for the fxnode.
	// This program defines how the node processes its inputs.
	SetShaderProgram(program fxcore.FXShaderProgram)

	// UpdateTransformationUniforms sets the transformation uniforms on the given program.
	// Uniforms: u_translation (vec2), u_scale (vec2), u_rotation (float).
	UpdateTransformationUniforms(program fxcore.FXShaderProgram)

	// Process executes the node's operation if necessary.
	// It checks if the node or any of its inputs are dirty.
	// If so, it renders the result to the output framebuffer.
	Process(ctx fxcontext.FXContext) error

	// Release frees resources held by the fxnode.
	// This includes the output framebuffer, quad, and shader program.
	Release()
}

// FXGraph represents a collection of nodes and their connections.
type FXGraph interface {
	// AddNode adds a node to the graph.
	AddNode(name string, node FXNode)
	// Connect connects two nodes.
	Connect(sourceNodeName, targetNodeName, inputSlot string) error
	// GetNode returns a node by name.
	GetNode(name string) FXNode
	// Release frees resources held by all nodes in the fxGraph.
	Release()
}

// FXPipeline manages the execution of a fxGraph.
type FXPipeline interface {
	// Execute executes the pipeline.
	Execute(outputNodeName string) error
	// Release frees resources held by the pipeline.
	Release()
}
