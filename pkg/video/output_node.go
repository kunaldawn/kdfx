package video

import (
	"fmt"

	"kdfx/pkg/context"
	"kdfx/pkg/core"
	"kdfx/pkg/node"
)

// VideoOutputNode represents a node that writes video frames to an encoder.
type VideoOutputNode interface {
	node.Node
	Close() error
}

type videoOutputNode struct {
	node.Node
	encoder StreamEncoder
	input   node.Input
}

// NewVideoOutputNode creates a new video output node.
func NewVideoOutputNode(ctx context.Context, encoder StreamEncoder, width, height int) (VideoOutputNode, error) {
	base, err := node.NewBaseNode(ctx, width, height)
	if err != nil {
		return nil, err
	}

	return &videoOutputNode{
		Node:    base,
		encoder: encoder,
	}, nil
}

func (n *videoOutputNode) SetInput(name string, input node.Input) {
	if name == "u_texture" {
		n.input = input
	}
}

func (n *videoOutputNode) Process(ctx context.Context) error {
	if n.input != nil {
		if inputNode, ok := n.input.(node.Node); ok {
			if err := inputNode.Process(ctx); err != nil {
				return err
			}
		}
	}

	// Get texture from input
	var tex core.Texture
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

func (n *videoOutputNode) GetTexture() core.Texture {
	if n.input != nil {
		return n.input.GetTexture()
	}
	return nil
}

func (n *videoOutputNode) Close() error {
	return n.encoder.Close()
}
