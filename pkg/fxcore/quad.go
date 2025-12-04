// Package fxcore provides core graphics primitives and resource management for the kdfx library.
package fxcore

import (
	"github.com/go-gl/gl/v3.1/gles2"
)

// FXQuad represents a full-screen fxQuad for rendering.
type FXQuad interface {
	// Draw renders the fxQuad using the specified attribute locations.
	Draw(positionAttrib, texCoordAttrib int32)
	// Release frees the OpenGL resources associated with the fxQuad.
	Release()
}

// fxQuad implements FXQuad.
type fxQuad struct {
	// vbo is the Vertex Buffer Object ID.
	vbo uint32
}

// fxQuadVertices contains the vertices for a full-screen quad.
var fxQuadVertices = []float32{
	// Pos      // Tex
	-1.0, 1.0, 0.0, 1.0,
	-1.0, -1.0, 0.0, 0.0,
	1.0, 1.0, 1.0, 1.0,
	1.0, -1.0, 1.0, 0.0,
}

// NewFXQuad creates a new full-screen fxQuad.
func NewFXQuad() FXQuad {
	var vbo uint32
	// Generate a Vertex Buffer Object (VBO) to store vertex data.
	gles2.GenBuffers(1, &vbo)
	// Bind the VBO to the GL_ARRAY_BUFFER target.
	gles2.BindBuffer(gles2.ARRAY_BUFFER, vbo)
	// Upload the vertex data (position and texture coordinates) to the GPU.
	// The data is static, meaning it won't change often.
	gles2.BufferData(gles2.ARRAY_BUFFER, len(fxQuadVertices)*4, gles2.Ptr(fxQuadVertices), gles2.STATIC_DRAW)
	// Unbind the VBO to avoid accidental modification.
	gles2.BindBuffer(gles2.ARRAY_BUFFER, 0)
	return &fxQuad{vbo: vbo}
}

func (q *fxQuad) Draw(positionAttrib, texCoordAttrib int32) {
	// Bind the VBO containing the quad vertices.
	gles2.BindBuffer(gles2.ARRAY_BUFFER, q.vbo)

	// Enable the position attribute and define its layout.
	// Stride is 4 floats (2 for position, 2 for texCoord).
	// Offset is 0 for position.
	gles2.EnableVertexAttribArray(uint32(positionAttrib))
	gles2.VertexAttribPointer(uint32(positionAttrib), 2, gles2.FLOAT, false, 4*4, gles2.PtrOffset(0))

	// Enable the texture coordinate attribute and define its layout.
	// Offset is 2 floats (8 bytes) for texture coordinates.
	gles2.EnableVertexAttribArray(uint32(texCoordAttrib))
	gles2.VertexAttribPointer(uint32(texCoordAttrib), 2, gles2.FLOAT, false, 4*4, gles2.PtrOffset(2*4))

	// Draw the quad as a triangle strip.
	gles2.DrawArrays(gles2.TRIANGLE_STRIP, 0, 4)

	// Disable attributes and unbind the buffer to clean up state.
	gles2.DisableVertexAttribArray(uint32(positionAttrib))
	gles2.DisableVertexAttribArray(uint32(texCoordAttrib))

	gles2.BindBuffer(gles2.ARRAY_BUFFER, 0)
}

func (q *fxQuad) Release() {
	// Delete the VBO to free GPU memory.
	gles2.DeleteBuffers(1, &q.vbo)
}
