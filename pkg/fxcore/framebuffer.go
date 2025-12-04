package fxcore

import (
	"fmt"

	"github.com/go-gl/gl/v3.1/gles2"
)

// FXFramebuffer represents an OpenGL fxFramebuffer object.
type FXFramebuffer interface {
	// Bind binds the fxFramebuffer for rendering.
	Bind()
	// Unbind unbinds the fxFramebuffer.
	Unbind()
	// Release frees the OpenGL resources associated with the fxFramebuffer.
	Release()
	// GetTexture returns the fxTexture attached to the fxFramebuffer.
	GetTexture() FXTexture
}

// fxFramebuffer implements FXFramebuffer.
type fxFramebuffer struct {
	// id is the OpenGL framebuffer ID.
	id uint32
	// fxTexture is the texture attached to the framebuffer.
	fxTexture FXTexture
}

// NewFXFramebuffer creates a new fxFramebuffer with a fxTexture attachment of the specified size.
func NewFXFramebuffer(width, height int) (FXFramebuffer, error) {
	var id uint32
	gles2.GenFramebuffers(1, &id)

	tex := NewFXTexture(width, height)

	gles2.BindFramebuffer(gles2.FRAMEBUFFER, id)
	gles2.FramebufferTexture2D(gles2.FRAMEBUFFER, gles2.COLOR_ATTACHMENT0, gles2.TEXTURE_2D, tex.GetID(), 0)

	status := gles2.CheckFramebufferStatus(gles2.FRAMEBUFFER)
	if status != gles2.FRAMEBUFFER_COMPLETE {
		gles2.BindFramebuffer(gles2.FRAMEBUFFER, 0)
		tex.Release()
		gles2.DeleteFramebuffers(1, &id)
		return nil, fmt.Errorf("fxFramebuffer incomplete: status %x", status)
	}

	gles2.BindFramebuffer(gles2.FRAMEBUFFER, 0)

	return &fxFramebuffer{id: id, fxTexture: tex}, nil
}

func (fb *fxFramebuffer) Bind() {
	gles2.BindFramebuffer(gles2.FRAMEBUFFER, fb.id)
	w, h := fb.fxTexture.GetSize()
	gles2.Viewport(0, 0, int32(w), int32(h))
}

func (fb *fxFramebuffer) Unbind() {
	gles2.BindFramebuffer(gles2.FRAMEBUFFER, 0)
}

func (fb *fxFramebuffer) Release() {
	if fb.fxTexture != nil {
		fb.fxTexture.Release()
	}
	gles2.DeleteFramebuffers(1, &fb.id)
}

func (fb *fxFramebuffer) GetTexture() FXTexture {
	return fb.fxTexture
}
