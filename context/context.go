package context

// Context defines the interface for an OpenGL context wrapper.
// It abstracts the underlying windowing system (GLFW, EGL, etc.).
type Context interface {
	// MakeCurrent makes the context current on the calling thread.
	MakeCurrent()
	// SwapBuffers swaps the front and back buffers (if applicable).
	SwapBuffers()
	// Destroy destroys the context and releases resources.
	Destroy()
	// GetSize returns the width and height of the context/surface.
	GetSize() (int, int)
}
