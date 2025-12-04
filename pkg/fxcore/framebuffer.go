package fxcore

import (
	"fmt"

	"github.com/go-gl/gl/v3.1/gles2"
)

// FXFramebuffer represents an OpenGL fxFramebuffer object.
type FXFramebuffer interface {
	// Bind binds the fxFramebuffer for rendering.
	// All subsequent drawing operations will target this framebuffer.
	Bind()
	// Unbind unbinds the fxFramebuffer.
	// This restores the default framebuffer (usually the window).
	Unbind()
	// Release frees the OpenGL resources associated with the fxFramebuffer.
	// This includes the FBO itself and the attached texture.
	Release()
	// GetTexture returns the fxTexture attached to the fxFramebuffer.
	// This texture contains the rendered output.
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
// NewFXFramebuffer creates a new fxFramebuffer with a fxTexture attachment of the specified size.
func NewFXFramebuffer(width, height int) (FXFramebuffer, error) {
	var id uint32
	// Generate a new Framebuffer Object (FBO) ID.
	gles2.GenFramebuffers(1, &id)

	// Create a texture to attach to the FBO. This will store the rendered output.
	tex := NewFXTexture(width, height)

	// Bind the FBO to configure it.
	gles2.BindFramebuffer(gles2.FRAMEBUFFER, id)
	// Attach the texture to the color attachment point 0.
	gles2.FramebufferTexture2D(gles2.FRAMEBUFFER, gles2.COLOR_ATTACHMENT0, gles2.TEXTURE_2D, tex.GetID(), 0)

	// Check if the framebuffer is complete and ready for use.
	status := gles2.CheckFramebufferStatus(gles2.FRAMEBUFFER)
	if status != gles2.FRAMEBUFFER_COMPLETE {
		gles2.BindFramebuffer(gles2.FRAMEBUFFER, 0)
		tex.Release()
		gles2.DeleteFramebuffers(1, &id)
		return nil, fmt.Errorf("fxFramebuffer incomplete: status %x", status)
	}

	// Unbind the FBO.
	gles2.BindFramebuffer(gles2.FRAMEBUFFER, 0)

	return &fxFramebuffer{id: id, fxTexture: tex}, nil
}

func (fb *fxFramebuffer) Bind() {
	// Bind the framebuffer for rendering.
	gles2.BindFramebuffer(gles2.FRAMEBUFFER, fb.id)
	// Set the viewport to match the framebuffer size.
	w, h := fb.fxTexture.GetSize()
	gles2.Viewport(0, 0, int32(w), int32(h))
}

func (fb *fxFramebuffer) Unbind() {
	// Unbind the framebuffer, switching back to the default framebuffer (usually the window).
	gles2.BindFramebuffer(gles2.FRAMEBUFFER, 0)
}

func (fb *fxFramebuffer) Release() {
	// Release the attached texture.
	if fb.fxTexture != nil {
		fb.fxTexture.Release()
	}
	// Delete the framebuffer object.
	gles2.DeleteFramebuffers(1, &fb.id)
}

func (fb *fxFramebuffer) GetTexture() FXTexture {
	return fb.fxTexture
}
