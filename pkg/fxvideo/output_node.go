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

type fxVideoOutputNode struct {
	fxnode.FXNode
	encoder FXStreamEncoder
	input   fxnode.FXInput
	ctx     fxcontext.FXContext
}

// NewFXVideoOutputNode creates a new video output fxnode.
func NewFXVideoOutputNode(ctx fxcontext.FXContext, encoder FXStreamEncoder) (FXVideoOutputNode, error) {
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
	img, err := tex.Download()
	if err != nil {
		return fmt.Errorf("failed to download texture: %w", err)
	}

	// Add frame to encoder
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
