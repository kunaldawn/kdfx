package core

import (
	"fmt"

	"github.com/go-gl/gl/v3.1/gles2"
)

type Framebuffer interface {
	Bind()
	Unbind()
	Release()
	GetTexture() Texture
}

type framebuffer struct {
	id      uint32
	texture Texture
}

func NewFramebuffer(width, height int) (Framebuffer, error) {
	var id uint32
	gles2.GenFramebuffers(1, &id)

	tex := NewTexture(width, height)

	gles2.BindFramebuffer(gles2.FRAMEBUFFER, id)
	gles2.FramebufferTexture2D(gles2.FRAMEBUFFER, gles2.COLOR_ATTACHMENT0, gles2.TEXTURE_2D, tex.GetID(), 0)

	status := gles2.CheckFramebufferStatus(gles2.FRAMEBUFFER)
	if status != gles2.FRAMEBUFFER_COMPLETE {
		gles2.BindFramebuffer(gles2.FRAMEBUFFER, 0)
		tex.Release()
		gles2.DeleteFramebuffers(1, &id)
		return nil, fmt.Errorf("framebuffer incomplete: status %x", status)
	}

	gles2.BindFramebuffer(gles2.FRAMEBUFFER, 0)

	return &framebuffer{id: id, texture: tex}, nil
}

func (fb *framebuffer) Bind() {
	gles2.BindFramebuffer(gles2.FRAMEBUFFER, fb.id)
	w, h := fb.texture.GetSize()
	gles2.Viewport(0, 0, int32(w), int32(h))
}

func (fb *framebuffer) Unbind() {
	gles2.BindFramebuffer(gles2.FRAMEBUFFER, 0)
}

func (fb *framebuffer) Release() {
	if fb.texture != nil {
		fb.texture.Release()
	}
	gles2.DeleteFramebuffers(1, &fb.id)
}

func (fb *framebuffer) GetTexture() Texture {
	return fb.texture
}
