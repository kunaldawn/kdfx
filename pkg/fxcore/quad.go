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
	gles2.GenBuffers(1, &vbo)
	gles2.BindBuffer(gles2.ARRAY_BUFFER, vbo)
	gles2.BufferData(gles2.ARRAY_BUFFER, len(fxQuadVertices)*4, gles2.Ptr(fxQuadVertices), gles2.STATIC_DRAW)
	gles2.BindBuffer(gles2.ARRAY_BUFFER, 0)
	return &fxQuad{vbo: vbo}
}

func (q *fxQuad) Draw(positionAttrib, texCoordAttrib int32) {
	gles2.BindBuffer(gles2.ARRAY_BUFFER, q.vbo)

	gles2.EnableVertexAttribArray(uint32(positionAttrib))
	gles2.VertexAttribPointer(uint32(positionAttrib), 2, gles2.FLOAT, false, 4*4, gles2.PtrOffset(0))

	gles2.EnableVertexAttribArray(uint32(texCoordAttrib))
	gles2.VertexAttribPointer(uint32(texCoordAttrib), 2, gles2.FLOAT, false, 4*4, gles2.PtrOffset(2*4))

	gles2.DrawArrays(gles2.TRIANGLE_STRIP, 0, 4)

	gles2.DisableVertexAttribArray(uint32(positionAttrib))
	gles2.DisableVertexAttribArray(uint32(texCoordAttrib))

	gles2.BindBuffer(gles2.ARRAY_BUFFER, 0)
}

func (q *fxQuad) Release() {
	gles2.DeleteBuffers(1, &q.vbo)
}
