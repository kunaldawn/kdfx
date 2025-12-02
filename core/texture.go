package core

import (
	"fmt"
	"image"
	"image/draw"
	_ "image/jpeg"
	_ "image/png"
	"os"

	"github.com/go-gl/gl/v3.1/gles2"
)

// Texture represents an OpenGL texture.
type Texture struct {
	ID     uint32
	Width  int
	Height int
}

// NewTexture creates a new empty texture.
func NewTexture(width, height int) *Texture {
	var id uint32
	gles2.GenTextures(1, &id)
	t := &Texture{ID: id, Width: width, Height: height}
	t.Bind()

	// Set default parameters
	gles2.TexParameteri(gles2.TEXTURE_2D, gles2.TEXTURE_MIN_FILTER, gles2.LINEAR)
	gles2.TexParameteri(gles2.TEXTURE_2D, gles2.TEXTURE_MAG_FILTER, gles2.LINEAR)
	gles2.TexParameteri(gles2.TEXTURE_2D, gles2.TEXTURE_WRAP_S, gles2.CLAMP_TO_EDGE)
	gles2.TexParameteri(gles2.TEXTURE_2D, gles2.TEXTURE_WRAP_T, gles2.CLAMP_TO_EDGE)

	// Allocate storage (empty)
	gles2.TexImage2D(gles2.TEXTURE_2D, 0, gles2.RGBA, int32(width), int32(height), 0, gles2.RGBA, gles2.UNSIGNED_BYTE, nil)

	t.Unbind()
	return t
}

// LoadTextureFromFile loads a texture from an image file.
func LoadTextureFromFile(path string) (*Texture, error) {
	file, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	img, _, err := image.Decode(file)
	if err != nil {
		return nil, err
	}

	rgba := image.NewRGBA(img.Bounds())
	draw.Draw(rgba, rgba.Bounds(), img, image.Point{0, 0}, draw.Src)

	t := NewTexture(rgba.Rect.Size().X, rgba.Rect.Size().Y)
	t.Bind()
	gles2.TexImage2D(gles2.TEXTURE_2D, 0, gles2.RGBA, int32(t.Width), int32(t.Height), 0, gles2.RGBA, gles2.UNSIGNED_BYTE, gles2.Ptr(rgba.Pix))
	t.Unbind()

	return t, nil
}

func (t *Texture) Bind() {
	gles2.BindTexture(gles2.TEXTURE_2D, t.ID)
}

func (t *Texture) Unbind() {
	gles2.BindTexture(gles2.TEXTURE_2D, 0)
}

func (t *Texture) Release() {
	gles2.DeleteTextures(1, &t.ID)
}

// Download reads the texture data back to an image.RGBA.
// Note: This requires a framebuffer to read from in GLES2.
// We can attach it to a temporary FBO or the current one.
func (t *Texture) Download() (*image.RGBA, error) {
	// Create a temporary FBO to read from
	var fbo uint32
	gles2.GenFramebuffers(1, &fbo)
	gles2.BindFramebuffer(gles2.FRAMEBUFFER, fbo)
	gles2.FramebufferTexture2D(gles2.FRAMEBUFFER, gles2.COLOR_ATTACHMENT0, gles2.TEXTURE_2D, t.ID, 0)

	status := gles2.CheckFramebufferStatus(gles2.FRAMEBUFFER)
	if status != gles2.FRAMEBUFFER_COMPLETE {
		gles2.BindFramebuffer(gles2.FRAMEBUFFER, 0)
		gles2.DeleteFramebuffers(1, &fbo)
		return nil, fmt.Errorf("framebuffer incomplete: status %x", status)
	}

	pixels := make([]uint8, t.Width*t.Height*4)
	gles2.ReadPixels(0, 0, int32(t.Width), int32(t.Height), gles2.RGBA, gles2.UNSIGNED_BYTE, gles2.Ptr(pixels))

	gles2.BindFramebuffer(gles2.FRAMEBUFFER, 0)
	gles2.DeleteFramebuffers(1, &fbo)

	// Flip Y because OpenGL is bottom-left origin
	// Actually, let's just return as is and let caller handle or flip here.
	// Standard image.RGBA is top-left.
	// We should flip it here to match standard image expectations.

	rect := image.Rect(0, 0, t.Width, t.Height)
	img := image.NewRGBA(rect)

	// Copy and flip Y
	stride := t.Width * 4
	for y := 0; y < t.Height; y++ {
		srcRow := pixels[y*stride : (y+1)*stride]
		dstRow := img.Pix[(t.Height-1-y)*stride : (t.Height-y)*stride]
		copy(dstRow, srcRow)
	}

	return img, nil
}
