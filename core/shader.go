package core

import (
	"fmt"
	"strings"

	"github.com/go-gl/gl/v3.1/gles2"
)

type ShaderType uint32

const (
	VertexShader   ShaderType = gles2.VERTEX_SHADER
	FragmentShader ShaderType = gles2.FRAGMENT_SHADER
)

type Shader struct {
	ID   uint32
	Type ShaderType
}

func NewShader(source string, shaderType ShaderType) (*Shader, error) {
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

	return &Shader{ID: id, Type: shaderType}, nil
}

func (s *Shader) Release() {
	gles2.DeleteShader(s.ID)
}

type ShaderProgram struct {
	ID uint32
}

func NewShaderProgram(vertexSource, fragmentSource string) (*ShaderProgram, error) {
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
	gles2.AttachShader(id, vs.ID)
	gles2.AttachShader(id, fs.ID)
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

	return &ShaderProgram{ID: id}, nil
}

func (p *ShaderProgram) Use() {
	gles2.UseProgram(p.ID)
}

func (p *ShaderProgram) Release() {
	gles2.DeleteProgram(p.ID)
}

func (p *ShaderProgram) GetUniformLocation(name string) int32 {
	cstrs, free := gles2.Strs(name + "\x00")
	defer free()
	return gles2.GetUniformLocation(p.ID, *cstrs)
}

func (p *ShaderProgram) SetUniform1i(name string, value int32) {
	loc := p.GetUniformLocation(name)
	if loc != -1 {
		gles2.Uniform1i(loc, value)
	}
}

func (p *ShaderProgram) SetUniform1f(name string, value float32) {
	loc := p.GetUniformLocation(name)
	if loc != -1 {
		gles2.Uniform1f(loc, value)
	}
}

func (p *ShaderProgram) SetUniform2f(name string, v0, v1 float32) {
	loc := p.GetUniformLocation(name)
	if loc != -1 {
		gles2.Uniform2f(loc, v0, v1)
	}
}

func (p *ShaderProgram) SetUniform3f(name string, v0, v1, v2 float32) {
	loc := p.GetUniformLocation(name)
	if loc != -1 {
		gles2.Uniform3f(loc, v0, v1, v2)
	}
}

func (p *ShaderProgram) GetAttribLocation(name string) int32 {
	cstrs, free := gles2.Strs(name + "\x00")
	defer free()
	return gles2.GetAttribLocation(p.ID, *cstrs)
}
