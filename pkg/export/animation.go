package export

import (
	"fmt"
	"io"
	"time"

	"kdfx/pkg/context"
	"kdfx/pkg/node"
)

// Animation defines the interface for an animation.
type Animation interface {
	// Render renders the animation to the provided writer using the specified node as output.
	Render(ctx context.Context, node node.Node, writer io.Writer) error
}

type animation struct {
	duration time.Duration
	fps      int
	update   func(t time.Duration)
}

// NewAnimation creates a new animation.
func NewAnimation(duration time.Duration, fps int, update func(t time.Duration)) Animation {
	return &animation{
		duration: duration,
		fps:      fps,
		update:   update,
	}
}

// Render renders the animation to the provided writer using the specified node as output.
func (a *animation) Render(ctx context.Context, node node.Node, writer io.Writer) error {
	width, height := ctx.GetSize()

	encoder, err := NewMP4StreamEncoder(writer, width, height, a.fps)
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
