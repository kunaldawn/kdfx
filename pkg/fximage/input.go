package fximage

import (
	"kdfx/pkg/fxcore"
)

// FXImageInput is a simple node that provides a texture from an image or existing texture.
type FXImageInput struct {
	Texture fxcore.FXTexture
}

// NewFXImageInput creates a new FXImageInput node.
func NewFXImageInput(texture fxcore.FXTexture) *FXImageInput {
	return &FXImageInput{Texture: texture}
}

// NewFXImageInputFromFile creates a new FXImageInput node from an image file.
func NewFXImageInputFromFile(path string) (*FXImageInput, error) {
	tex, err := fxcore.FXLoadTextureFromFile(path)
	if err != nil {
		return nil, err
	}
	return &FXImageInput{Texture: tex}, nil
}

func (n *FXImageInput) GetTexture() fxcore.FXTexture { return n.Texture }
func (n *FXImageInput) IsDirty() bool                { return false } // Static input
