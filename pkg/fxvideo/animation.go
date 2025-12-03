package fxvideo

import (
	"fmt"
	"io"
	"time"

	"kdfx/pkg/fxcontext"
	"kdfx/pkg/fxnode"
)

// FXAnimation defines the interface for an fxAnimation.
type FXAnimation interface {
	// Render renders the fxAnimation to the provided writer using the specified node as output.
	Render(ctx fxcontext.FXContext, node fxnode.FXNode, writer io.Writer) error
}

type fxAnimation struct {
	duration time.Duration
	fps      int
	update   func(t time.Duration)
}

// NewFXAnimation creates a new fxAnimation.
func NewFXAnimation(duration time.Duration, fps int, update func(t time.Duration)) FXAnimation {
	return &fxAnimation{
		duration: duration,
		fps:      fps,
		update:   update,
	}
}

// Render renders the fxAnimation to the provided writer using the specified node as output.
func (a *fxAnimation) Render(ctx fxcontext.FXContext, node fxnode.FXNode, writer io.Writer) error {
	width, height := ctx.GetSize()

	encoder, err := NewFXMP4StreamEncoder(writer, width, height, a.fps)
	if err != nil {
		return fmt.Errorf("failed to create encoder: %w", err)
	}
	defer encoder.Close()

	frameCount := int(a.duration.Seconds() * float64(a.fps))
	dt := time.Second / time.Duration(a.fps)

	for i := 0; i < frameCount; i++ {
		currentTime := time.Duration(i) * dt

		// Update scene state
		if a.update != nil {
			a.update(currentTime)
		}

		// Process the graph
		if err := node.Process(ctx); err != nil {
			return fmt.Errorf("failed to process frame %d: %w", i, err)
		}

		// Download the result
		tex := node.GetTexture()
		if tex == nil {
			return fmt.Errorf("node returned nil texture at frame %d", i)
		}

		img, err := tex.Download()
		if err != nil {
			return fmt.Errorf("failed to download texture at frame %d: %w", i, err)
		}

		// Add frame to encoder
		if err := encoder.AddFrame(img); err != nil {
			return fmt.Errorf("failed to add frame %d to encoder: %w", i, err)
		}
	}

	return nil
}
