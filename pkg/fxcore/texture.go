package fxcore

import (
	"fmt"
	"image"
	"image/draw"
	_ "image/jpeg"
	_ "image/png"
	"os"

	"github.com/go-gl/gl/v3.1/gles2"
)

// FXTexture represents an OpenGL fxTexture.
type FXTexture interface {
	// Bind binds the fxTexture to the current fxcontext.
	Bind()
	// BindToUnit binds the fxTexture to a specific fxTexture unit.
	BindToUnit(unit int)
	// Unbind unbinds the fxTexture.
	Unbind()
	// Release frees the OpenGL resources associated with the fxTexture.
	Release()
	// Download reads the fxTexture data back to an image.RGBA.
	Download() (*image.RGBA, error)
	// GetID returns the OpenGL fxTexture ID.
	GetID() uint32
	// GetSize returns the width and height of the fxTexture.
	// GetSize returns the width and height of the fxTexture.
	GetSize() (int, int)
	// Upload updates the fxTexture content from an image.RGBA.
	Upload(img *image.RGBA)
}

// fxTexture implements FXTexture.
type fxTexture struct {
	// id is the OpenGL texture ID.
	id uint32
	// width is the width of the texture.
	width int
	// height is the height of the texture.
	height int
}

// NewFXTexture creates a new empty fxTexture.
func NewFXTexture(width, height int) FXTexture {
	var id uint32
	// Generate a new texture ID.
	gles2.GenTextures(1, &id)
	t := &fxTexture{id: id, width: width, height: height}
	// Bind the texture to configure it.
	t.Bind()

	// Set default parameters
	// Linear filtering for minification and magnification.
	gles2.TexParameteri(gles2.TEXTURE_2D, gles2.TEXTURE_MIN_FILTER, gles2.LINEAR)
	gles2.TexParameteri(gles2.TEXTURE_2D, gles2.TEXTURE_MAG_FILTER, gles2.LINEAR)
	// Clamp to edge to avoid artifacts at the borders.
	gles2.TexParameteri(gles2.TEXTURE_2D, gles2.TEXTURE_WRAP_S, gles2.CLAMP_TO_EDGE)
	gles2.TexParameteri(gles2.TEXTURE_2D, gles2.TEXTURE_WRAP_T, gles2.CLAMP_TO_EDGE)

	// Allocate storage (empty)
	// Initialize the texture with null data, allocating memory on the GPU.
	gles2.TexImage2D(gles2.TEXTURE_2D, 0, gles2.RGBA, int32(width), int32(height), 0, gles2.RGBA, gles2.UNSIGNED_BYTE, nil)

	// Unbind the texture.
	t.Unbind()

	return t
}

// FXLoadTextureFromFile loads a fxTexture from an image file.
func FXLoadTextureFromFile(path string) (FXTexture, error) {
	// Open the file.
	file, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	// Decode the image.
	img, _, err := image.Decode(file)
	if err != nil {
		return nil, err
	}

	// Convert to RGBA. This ensures we have a consistent format for OpenGL.
	rgba := image.NewRGBA(img.Bounds())
	draw.Draw(rgba, rgba.Bounds(), img, image.Point{0, 0}, draw.Src)

	// Create a new texture and upload the data.
	t := NewFXTexture(rgba.Rect.Size().X, rgba.Rect.Size().Y)
	t.Bind()
	w, h := t.GetSize()
	gles2.TexImage2D(gles2.TEXTURE_2D, 0, gles2.RGBA, int32(w), int32(h), 0, gles2.RGBA, gles2.UNSIGNED_BYTE, gles2.Ptr(rgba.Pix))
	t.Unbind()

	return t, nil
}

func (t *fxTexture) Bind() {
	gles2.BindTexture(gles2.TEXTURE_2D, t.id)
}

func (t *fxTexture) BindToUnit(unit int) {
	// Activate the specified texture unit.
	gles2.ActiveTexture(gles2.TEXTURE0 + uint32(unit))
	// Bind the texture to that unit.
	gles2.BindTexture(gles2.TEXTURE_2D, t.id)
}

func (t *fxTexture) Unbind() {
	gles2.BindTexture(gles2.TEXTURE_2D, 0)
}

func (t *fxTexture) Release() {
	// Delete the texture to free GPU memory.
	gles2.DeleteTextures(1, &t.id)
}

func (t *fxTexture) GetID() uint32 {
	return t.id
}

func (t *fxTexture) GetSize() (int, int) {
	return t.width, t.height
}

// Download reads the fxTexture data back to an image.RGBA.
func (t *fxTexture) Download() (*image.RGBA, error) {
	// Create a temporary FBO to read from
	// We can't read directly from a texture, so we attach it to an FBO.
	var fbo uint32
	gles2.GenFramebuffers(1, &fbo)
	gles2.BindFramebuffer(gles2.FRAMEBUFFER, fbo)
	gles2.FramebufferTexture2D(gles2.FRAMEBUFFER, gles2.COLOR_ATTACHMENT0, gles2.TEXTURE_2D, t.id, 0)

	// Check if the FBO is complete.
	status := gles2.CheckFramebufferStatus(gles2.FRAMEBUFFER)
	if status != gles2.FRAMEBUFFER_COMPLETE {
		gles2.BindFramebuffer(gles2.FRAMEBUFFER, 0)
		gles2.DeleteFramebuffers(1, &fbo)
		return nil, fmt.Errorf("fxFramebuffer incomplete: status %x", status)
	}

	// Read pixels from the FBO.
	pixels := make([]uint8, t.width*t.height*4)
	gles2.ReadPixels(0, 0, int32(t.width), int32(t.height), gles2.RGBA, gles2.UNSIGNED_BYTE, gles2.Ptr(pixels))

	// Unbind the FBO.
	gles2.BindFramebuffer(gles2.FRAMEBUFFER, 0)
	gles2.DeleteFramebuffers(1, &fbo)

	// Create an image from the pixels.
	rect := image.Rect(0, 0, t.width, t.height)
	img := image.NewRGBA(rect)

	// Copy and flip Y
	// OpenGL coordinates have (0,0) at bottom-left, while Go images have it at top-left.
	stride := t.width * 4
	for y := 0; y < t.height; y++ {
		srcRow := pixels[y*stride : (y+1)*stride]
		dstRow := img.Pix[(t.height-1-y)*stride : (t.height-y)*stride]
		copy(dstRow, srcRow)
	}

	return img, nil
}

func (t *fxTexture) Upload(img *image.RGBA) {
	// Bind the texture to upload new data.
	t.Bind()
	// Upload new data to the existing texture storage.
	gles2.TexSubImage2D(gles2.TEXTURE_2D, 0, 0, 0, int32(t.width), int32(t.height), gles2.RGBA, gles2.UNSIGNED_BYTE, gles2.Ptr(img.Pix))
	// Unbind the texture.
	t.Unbind()
}
