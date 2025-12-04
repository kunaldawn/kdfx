package fxvideo

import (
	"fmt"

	"kdfx/pkg/fxcontext"
	"kdfx/pkg/fxcore"
	"kdfx/pkg/fxnode"
)

// FXVideoOutputNode represents a node that writes video frames to an encoder.
type FXVideoOutputNode interface {
	fxnode.FXNode
	Close() error
}

// fxVideoOutputNode implements FXVideoOutputNode.
type fxVideoOutputNode struct {
	fxnode.FXNode
	// encoder is the video stream encoder.
	encoder FXStreamEncoder
	// input is the input connection.
	input fxnode.FXInput
	// ctx is the context used for rendering.
	ctx fxcontext.FXContext
}

// NewFXVideoOutputNode creates a new video output fxnode.
func NewFXVideoOutputNode(ctx fxcontext.FXContext, encoder FXStreamEncoder) (FXVideoOutputNode, error) {
	// Create a base node. Dimensions are 0 because this node doesn't render to its own texture,
	// but rather consumes an input texture.
	base, err := fxnode.NewFXBaseNode(ctx, 0, 0) // width and height are not directly used by this node, but base node requires them.
	if err != nil {
		return nil, err
	}

	return &fxVideoOutputNode{
		FXNode:  base,
		encoder: encoder,
	}, nil
}

func (n *fxVideoOutputNode) SetInput(name string, input fxnode.FXInput) {
	if name == "u_texture" {
		n.input = input
	}
}

func (n *fxVideoOutputNode) Process(ctx fxcontext.FXContext) error {
	// Process the input node first to ensure the texture is ready.
	if n.input != nil {
		if inputNode, ok := n.input.(fxnode.FXNode); ok {
			if err := inputNode.Process(ctx); err != nil {
				return err
			}
		}
	}

	// Get texture from input
	var tex fxcore.FXTexture
	if n.input != nil {
		tex = n.input.GetTexture()
	}

	if tex == nil {
		return fmt.Errorf("no input texture")
	}

	// Download texture
	// Read the texture data from GPU memory to CPU memory.
	img, err := tex.Download()
	if err != nil {
		return fmt.Errorf("failed to download texture: %w", err)
	}

	// Add frame to encoder
	// Send the image to the video encoder.
	if err := n.encoder.AddFrame(img); err != nil {
		return fmt.Errorf("failed to add frame: %w", err)
	}

	return nil
}

func (n *fxVideoOutputNode) GetTexture() fxcore.FXTexture {
	if n.input != nil {
		return n.input.GetTexture()
	}
	return nil
}

func (n *fxVideoOutputNode) Close() error {
	return n.encoder.Close()
}
