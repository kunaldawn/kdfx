package core

import (
	"fmt"

	"github.com/go-gl/gl/v3.1/gles2"
)

type Framebuffer struct {
	ID      uint32
	Texture *Texture
}

func NewFramebuffer(width, height int) (*Framebuffer, error) {
	var id uint32
	gles2.GenFramebuffers(1, &id)

	tex := NewTexture(width, height)

	gles2.BindFramebuffer(gles2.FRAMEBUFFER, id)
	gles2.FramebufferTexture2D(gles2.FRAMEBUFFER, gles2.COLOR_ATTACHMENT0, gles2.TEXTURE_2D, tex.ID, 0)

	status := gles2.CheckFramebufferStatus(gles2.FRAMEBUFFER)
	if status != gles2.FRAMEBUFFER_COMPLETE {
		gles2.BindFramebuffer(gles2.FRAMEBUFFER, 0)
		tex.Release()
		gles2.DeleteFramebuffers(1, &id)
		return nil, fmt.Errorf("framebuffer incomplete: status %x", status)
	}

	gles2.BindFramebuffer(gles2.FRAMEBUFFER, 0)

	return &Framebuffer{ID: id, Texture: tex}, nil
}

func (fb *Framebuffer) Bind() {
	gles2.BindFramebuffer(gles2.FRAMEBUFFER, fb.ID)
	gles2.Viewport(0, 0, int32(fb.Texture.Width), int32(fb.Texture.Height))
}

func (fb *Framebuffer) Unbind() {
	gles2.BindFramebuffer(gles2.FRAMEBUFFER, 0)
}

func (fb *Framebuffer) Release() {
	if fb.Texture != nil {
		fb.Texture.Release()
	}
	gles2.DeleteFramebuffers(1, &fb.ID)
}
