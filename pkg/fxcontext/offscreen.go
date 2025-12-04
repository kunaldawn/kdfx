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
	// window is the hidden GLFW window used for the context.
	window *glfw.Window
	// width is the width of the offscreen context.
	width int
	// height is the height of the offscreen context.
	height int
}

// NewFXOffscreenContext creates a new offscreen context with the specified dimensions.
// It initializes GLFW and creates a hidden window to provide an OpenGL ES 2.0 context.
// This is suitable for headless rendering or background processing where no visible window is required.
func NewFXOffscreenContext(width, height int) (FXContext, error) {
	// Initialize GLFW. This is required before creating any window.
	if err := glfw.Init(); err != nil {
		return nil, fmt.Errorf("failed to initialize glfw: %v", err)
	}

	// Set window hints for an invisible window and OpenGL ES 2.0 context.
	// We use OpenGL ES 2.0 for broader compatibility, especially on embedded devices.
	glfw.WindowHint(glfw.Visible, glfw.False)
	glfw.WindowHint(glfw.ContextVersionMajor, 2)
	glfw.WindowHint(glfw.ContextVersionMinor, 0)
	glfw.WindowHint(glfw.ClientAPI, glfw.OpenGLESAPI)

	// Create the window.
	window, err := glfw.CreateWindow(width, height, "kdfx-offscreen", nil, nil)
	if err != nil {
		glfw.Terminate()
		return nil, fmt.Errorf("failed to create glfw window: %v", err)
	}

	// Make the context current immediately so we can initialize GLES.
	window.MakeContextCurrent()

	// Initialize GLES bindings.
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
	// Delegate to GLFW to make the context current on this thread.
	c.window.MakeContextCurrent()
}

func (c *fxOffscreenContext) SwapBuffers() {
	// Swap the front and back buffers. Even for offscreen, this might be needed
	// if we were reading from the front buffer, but typically we render to FBOs.
	// However, it keeps the GLFW state happy.
	c.window.SwapBuffers()
}

func (c *fxOffscreenContext) Destroy() {
	// Clean up the window and terminate GLFW.
	c.window.Destroy()
	glfw.Terminate()
}

func (c *fxOffscreenContext) GetSize() (int, int) {
	return c.width, c.height
}

func (c *fxOffscreenContext) Viewport(x, y, width, height int) {
	// Set the OpenGL viewport.
	gles2.Viewport(int32(x), int32(y), int32(width), int32(height))
}
