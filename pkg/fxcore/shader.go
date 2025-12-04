package fxcore

import (
	"fmt"
	"strings"

	"github.com/go-gl/gl/v3.1/gles2"
)

// FXShaderType represents the type of fxShader (Vertex or Fragment).
type FXShaderType uint32

const (
	FXVertexShader   FXShaderType = gles2.VERTEX_SHADER
	FXFragmentShader FXShaderType = gles2.FRAGMENT_SHADER
)

// FXSimpleVS is a basic vertex fxShader that passes through position and fxTexture coordinates.
const FXSimpleVS = `
attribute vec2 a_position;
attribute vec2 a_texCoord;
varying vec2 v_texCoord;

uniform vec2 u_translation;
uniform vec2 u_scale;
uniform float u_rotation;

void main() {
	// Apply scaling
	vec2 scaledPos = a_position * u_scale;

	// Apply rotation
	float c = cos(u_rotation);
	float s = sin(u_rotation);
	vec2 rotatedPos = vec2(
		scaledPos.x * c - scaledPos.y * s,
		scaledPos.x * s + scaledPos.y * c
	);

	// Apply translation
	vec2 finalPos = rotatedPos + u_translation;

	gl_Position = vec4(finalPos, 0.0, 1.0);
	v_texCoord = a_texCoord;
}
`

// FXShader represents a compiled OpenGL fxShader.
type FXShader interface {
	// Release frees the OpenGL resources associated with the fxShader.
	Release()
	// GetID returns the OpenGL fxShader ID.
	GetID() uint32
}

type fxShader struct {
	id   uint32
	kind FXShaderType
}

// NewFXShader compiles a new fxShader from source code.
func NewFXShader(source string, shaderType FXShaderType) (FXShader, error) {
	id := gles2.CreateShader(uint32(shaderType))

	// Add default precision for ES 2.0 if not present
	if shaderType == FXFragmentShader && !strings.Contains(source, "precision") {
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
		return nil, fmt.Errorf("failed to compile fxShader: %v", log)
	}

	return &fxShader{id: id, kind: shaderType}, nil
}

func (s *fxShader) Release() {
	gles2.DeleteShader(s.id)
}

func (s *fxShader) GetID() uint32 {
	return s.id
}

// FXShaderProgram represents a linked OpenGL fxShader program.
type FXShaderProgram interface {
	// Use activates the fxShader program.
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

type fxShaderProgram struct {
	id uint32
}

// NewFXShaderProgram links a vertex and fragment fxShader into a program.
func NewFXShaderProgram(vertexSource, fragmentSource string) (FXShaderProgram, error) {
	vs, err := NewFXShader(vertexSource, FXVertexShader)
	if err != nil {
		return nil, fmt.Errorf("vertex fxShader error: %v", err)
	}
	defer vs.Release()

	fs, err := NewFXShader(fragmentSource, FXFragmentShader)
	if err != nil {
		return nil, fmt.Errorf("fragment fxShader error: %v", err)
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

	return &fxShaderProgram{id: id}, nil
}

func (p *fxShaderProgram) Use() {
	gles2.UseProgram(p.id)
}

func (p *fxShaderProgram) Release() {
	gles2.DeleteProgram(p.id)
}

func (p *fxShaderProgram) GetUniformLocation(name string) int32 {
	cstrs, free := gles2.Strs(name + "\x00")
	defer free()
	return gles2.GetUniformLocation(p.id, *cstrs)
}

func (p *fxShaderProgram) SetUniform1i(name string, value int32) {
	loc := p.GetUniformLocation(name)
	if loc != -1 {
		gles2.Uniform1i(loc, value)
	}
}

func (p *fxShaderProgram) SetUniform1f(name string, value float32) {
	loc := p.GetUniformLocation(name)
	if loc != -1 {
		gles2.Uniform1f(loc, value)
	}
}

func (p *fxShaderProgram) SetUniform2f(name string, v0, v1 float32) {
	loc := p.GetUniformLocation(name)
	if loc != -1 {
		gles2.Uniform2f(loc, v0, v1)
	}
}

func (p *fxShaderProgram) SetUniform3f(name string, v0, v1, v2 float32) {
	loc := p.GetUniformLocation(name)
	if loc != -1 {
		gles2.Uniform3f(loc, v0, v1, v2)
	}
}

func (p *fxShaderProgram) GetAttribLocation(name string) int32 {
	cstrs, free := gles2.Strs(name + "\x00")
	defer free()
	return gles2.GetAttribLocation(p.id, *cstrs)
}
