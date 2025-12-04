// Package fxcontext provides interfaces and implementations for OpenGL context management.
package fxcontext

// FXContext defines the interface for an OpenGL context wrapper.
// It abstracts the underlying windowing system (GLFW, EGL, etc.).
type FXContext interface {
	// MakeCurrent makes the context current on the calling thread.
	// This must be called before issuing any OpenGL commands on the thread.
	MakeCurrent()
	// SwapBuffers swaps the front and back buffers (if applicable).
	// For offscreen contexts, this might be a no-op or used to synchronize with the window system.
	SwapBuffers()
	// Destroy destroys the context and releases all associated resources.
	// It should be called when the context is no longer needed to prevent leaks.
	Destroy()
	// GetSize returns the width and height of the context/surface in pixels.
	GetSize() (int, int)
	// Viewport sets the viewport for the fxcontext.
	// This maps the normalized device coordinates to window coordinates.
	Viewport(x, y, width, height int)
}
