package core

import (
	"github.com/go-gl/gl/v3.1/gles2"
)

// Quad represents a full-screen quad for rendering.
type Quad struct {
	VBO uint32
}

var quadVertices = []float32{
	// Pos      // Tex
	-1.0, 1.0, 0.0, 1.0,
	-1.0, -1.0, 0.0, 0.0,
	1.0, 1.0, 1.0, 1.0,
	1.0, -1.0, 1.0, 0.0,
}

func NewQuad() *Quad {
	var vbo uint32
	gles2.GenBuffers(1, &vbo)
	gles2.BindBuffer(gles2.ARRAY_BUFFER, vbo)
	gles2.BufferData(gles2.ARRAY_BUFFER, len(quadVertices)*4, gles2.Ptr(quadVertices), gles2.STATIC_DRAW)
	gles2.BindBuffer(gles2.ARRAY_BUFFER, 0)
	return &Quad{VBO: vbo}
}

func (q *Quad) Draw(positionAttrib, texCoordAttrib int32) {
	gles2.BindBuffer(gles2.ARRAY_BUFFER, q.VBO)

	gles2.EnableVertexAttribArray(uint32(positionAttrib))
	gles2.VertexAttribPointer(uint32(positionAttrib), 2, gles2.FLOAT, false, 4*4, gles2.PtrOffset(0))

	gles2.EnableVertexAttribArray(uint32(texCoordAttrib))
	gles2.VertexAttribPointer(uint32(texCoordAttrib), 2, gles2.FLOAT, false, 4*4, gles2.PtrOffset(2*4))

	gles2.DrawArrays(gles2.TRIANGLE_STRIP, 0, 4)

	gles2.DisableVertexAttribArray(uint32(positionAttrib))
	gles2.DisableVertexAttribArray(uint32(texCoordAttrib))

	gles2.BindBuffer(gles2.ARRAY_BUFFER, 0)
}

func (q *Quad) Release() {
	gles2.DeleteBuffers(1, &q.VBO)
}
