package core

import (
	"fmt"
	"strings"

	"github.com/go-gl/gl/v3.1/gles2"
)

// ShaderType represents the type of shader (Vertex or Fragment).
type ShaderType uint32

const (
	VertexShader   ShaderType = gles2.VERTEX_SHADER
	FragmentShader ShaderType = gles2.FRAGMENT_SHADER
)

// SimpleVS is a basic vertex shader that passes through position and texture coordinates.
const SimpleVS = `
attribute vec2 a_position;
attribute vec2 a_texCoord;
varying vec2 v_texCoord;
void main() {
	gl_Position = vec4(a_position, 0.0, 1.0);
	v_texCoord = a_texCoord;
}
`

// Shader represents a compiled OpenGL shader.
type Shader interface {
	// Release frees the OpenGL resources associated with the shader.
	Release()
	// GetID returns the OpenGL shader ID.
	GetID() uint32
}

type shader struct {
	id   uint32
	kind ShaderType
}

// NewShader compiles a new shader from source code.
func NewShader(source string, shaderType ShaderType) (Shader, error) {
	id := gles2.CreateShader(uint32(shaderType))

	// Add default precision for ES 2.0 if not present
	if shaderType == FragmentShader && !strings.Contains(source, "precision") {
		source = "precision mediump float;\n" + source
	}

	cstrs, free := gles2.Strs(source + "\x00")
	defer free()
	gles2.ShaderSource(id, 1, cstrs, nil)
	gles2.CompileShader(id)

	var status int32
	gles2.GetShaderiv(id, gles2.COMPILE_STATUS, &status)
	if status == gles2.FALSE {
		var logLength int32
		gles2.GetShaderiv(id, gles2.INFO_LOG_LENGTH, &logLength)
		log := strings.Repeat("\x00", int(logLength+1))
		gles2.GetShaderInfoLog(id, logLength, nil, gles2.Str(log))
		gles2.DeleteShader(id)
		return nil, fmt.Errorf("failed to compile shader: %v", log)
	}

	return &shader{id: id, kind: shaderType}, nil
}

func (s *shader) Release() {
	gles2.DeleteShader(s.id)
}

func (s *shader) GetID() uint32 {
	return s.id
}

// ShaderProgram represents a linked OpenGL shader program.
type ShaderProgram interface {
	// Use activates the shader program.
	Use()
	// Release frees the OpenGL resources associated with the program.
	Release()
	// GetUniformLocation returns the location of a uniform variable.
	GetUniformLocation(name string) int32
	// SetUniform1i sets a single integer uniform.
	SetUniform1i(name string, value int32)
	// SetUniform1f sets a single float uniform.
	SetUniform1f(name string, value float32)
	// SetUniform2f sets a vec2 uniform.
	SetUniform2f(name string, v0, v1 float32)
	// SetUniform3f sets a vec3 uniform.
	SetUniform3f(name string, v0, v1, v2 float32)
	// GetAttribLocation returns the location of an attribute variable.
	GetAttribLocation(name string) int32
}

type shaderProgram struct {
	id uint32
}

// NewShaderProgram links a vertex and fragment shader into a program.
func NewShaderProgram(vertexSource, fragmentSource string) (ShaderProgram, error) {
	vs, err := NewShader(vertexSource, VertexShader)
	if err != nil {
		return nil, fmt.Errorf("vertex shader error: %v", err)
	}
	defer vs.Release()

	fs, err := NewShader(fragmentSource, FragmentShader)
	if err != nil {
		return nil, fmt.Errorf("fragment shader error: %v", err)
	}
	defer fs.Release()

	id := gles2.CreateProgram()
	gles2.AttachShader(id, vs.GetID())
	gles2.AttachShader(id, fs.GetID())
	gles2.LinkProgram(id)

	var status int32
	gles2.GetProgramiv(id, gles2.LINK_STATUS, &status)
	if status == gles2.FALSE {
		var logLength int32
		gles2.GetProgramiv(id, gles2.INFO_LOG_LENGTH, &logLength)
		log := strings.Repeat("\x00", int(logLength+1))
		gles2.GetProgramInfoLog(id, logLength, nil, gles2.Str(log))
		gles2.DeleteProgram(id)
		return nil, fmt.Errorf("failed to link program: %v", log)
	}

	return &shaderProgram{id: id}, nil
}

func (p *shaderProgram) Use() {
	gles2.UseProgram(p.id)
}

func (p *shaderProgram) Release() {
	gles2.DeleteProgram(p.id)
}

func (p *shaderProgram) GetUniformLocation(name string) int32 {
	cstrs, free := gles2.Strs(name + "\x00")
	defer free()
	return gles2.GetUniformLocation(p.id, *cstrs)
}

func (p *shaderProgram) SetUniform1i(name string, value int32) {
	loc := p.GetUniformLocation(name)
	if loc != -1 {
		gles2.Uniform1i(loc, value)
	}
}

func (p *shaderProgram) SetUniform1f(name string, value float32) {
	loc := p.GetUniformLocation(name)
	if loc != -1 {
		gles2.Uniform1f(loc, value)
	}
}

func (p *shaderProgram) SetUniform2f(name string, v0, v1 float32) {
	loc := p.GetUniformLocation(name)
	if loc != -1 {
		gles2.Uniform2f(loc, v0, v1)
	}
}

func (p *shaderProgram) SetUniform3f(name string, v0, v1, v2 float32) {
	loc := p.GetUniformLocation(name)
	if loc != -1 {
		gles2.Uniform3f(loc, v0, v1, v2)
	}
}

func (p *shaderProgram) GetAttribLocation(name string) int32 {
	cstrs, free := gles2.Strs(name + "\x00")
	defer free()
	return gles2.GetAttribLocation(p.id, *cstrs)
}
