package fximage

import (
	"fmt"
	"image/png"
	"os"

	"kdfx/pkg/fxcontext"
	"kdfx/pkg/fxcore"
	"kdfx/pkg/fxnode"
)

// FXImageOutput is a node that captures the output of another node and can save it to a file.
type FXImageOutput struct {
	Input fxnode.FXInput
}

// NewFXImageOutput creates a new FXImageOutput node.
func NewFXImageOutput() *FXImageOutput {
	return &FXImageOutput{}
}

func (n *FXImageOutput) SetInput(input fxnode.FXInput) {
	n.Input = input
}

func (n *FXImageOutput) GetTexture() fxcore.FXTexture {
	if n.Input != nil {
		return n.Input.GetTexture()
	}
	return nil
}

func (n *FXImageOutput) IsDirty() bool {
	if n.Input != nil {
		return n.Input.IsDirty()
	}
	return false
}

// Save saves the current texture to a PNG file.
func (n *FXImageOutput) Save(filename string) error {
	tex := n.GetTexture()
	if tex == nil {
		return fmt.Errorf("no input texture to save")
	}

	img, err := tex.Download()
	if err != nil {
		return err
	}

	f, err := os.Create(filename)
	if err != nil {
		return err
	}
	defer f.Close()

	return png.Encode(f, img)
}

// Process is a no-op for FXImageOutput as it just passes through the texture,
// but it might be needed if added to a graph that calls Process on it.
// However, usually we call Process on the node we want to view.
// If FXImageOutput is the end of the chain, we might want to call Process on its input.
func (n *FXImageOutput) Process(ctx fxcontext.FXContext) error {
	// If the input is a node, we should probably process it?
	// But usually the pipeline handles processing.
	// For now, let's assume the user drives processing on the node they want to render.
	// If FXImageOutput is used as a sink, maybe it should trigger processing of its input?
	if node, ok := n.Input.(fxnode.FXNode); ok {
		return node.Process(ctx)
	}
	return nil
}
