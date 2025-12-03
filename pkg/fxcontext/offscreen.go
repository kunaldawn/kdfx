package fxcontext

import (
	"fmt"
	"runtime"

	"github.com/go-gl/gl/v3.1/gles2"
	"github.com/go-gl/glfw/v3.3/glfw"
)

func init() {
	// GLFW event handling must run on the main OS thread
	runtime.LockOSThread()
}

// fxOffscreenContext implements FXContext using a hidden GLFW window.
type fxOffscreenContext struct {
	window *glfw.Window
	width  int
	height int
}

// NewFXOffscreenContext creates a new offscreen context with the specified dimensions.
func NewFXOffscreenContext(width, height int) (FXContext, error) {
	if err := glfw.Init(); err != nil {
		return nil, fmt.Errorf("failed to initialize glfw: %v", err)
	}

	glfw.WindowHint(glfw.Visible, glfw.False)
	glfw.WindowHint(glfw.ContextVersionMajor, 2)
	glfw.WindowHint(glfw.ContextVersionMinor, 0)
	glfw.WindowHint(glfw.ClientAPI, glfw.OpenGLESAPI)

	window, err := glfw.CreateWindow(width, height, "kimg-offscreen", nil, nil)
	if err != nil {
		glfw.Terminate()
		return nil, fmt.Errorf("failed to create glfw window: %v", err)
	}

	window.MakeContextCurrent()

	if err := gles2.Init(); err != nil {
		window.Destroy()
		glfw.Terminate()
		return nil, fmt.Errorf("failed to initialize gles2: %v", err)
	}

	return &fxOffscreenContext{
		window: window,
		width:  width,
		height: height,
	}, nil
}

func (c *fxOffscreenContext) MakeCurrent() {
	c.window.MakeContextCurrent()
}

func (c *fxOffscreenContext) SwapBuffers() {
	c.window.SwapBuffers()
}

func (c *fxOffscreenContext) Destroy() {
	c.window.Destroy()
	glfw.Terminate()
}

func (c *fxOffscreenContext) GetSize() (int, int) {
	return c.width, c.height
}

func (c *fxOffscreenContext) Viewport(x, y, width, height int) {
	gles2.Viewport(int32(x), int32(y), int32(width), int32(height))
}
