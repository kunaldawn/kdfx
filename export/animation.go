package export

import (
	"fmt"
	"io"
	"time"

	"kimg/context"
	"kimg/core"
)

// Node represents a node in the processing graph.
// We need an interface that provides the output texture and a way to process it.
// This mirrors the structure seen in examples, where nodes have Process(ctx) and GetTexture().
type Node interface {
	Process(ctx context.Context) error
	GetTexture() *core.Texture
}

// Animation defines the parameters for an animation.
type Animation struct {
	Duration time.Duration
	FPS      int
	Update   func(t time.Duration)
}

// Render renders the animation to the provided writer using the specified node as output.
func (a *Animation) Render(ctx context.Context, node Node, writer io.Writer) error {
	width, height := ctx.GetSize()

	encoder, err := NewMP4StreamEncoder(writer, width, height, a.FPS)
	if err != nil {
		return fmt.Errorf("failed to create encoder: %w", err)
	}
	defer encoder.Close()

	frameCount := int(a.Duration.Seconds() * float64(a.FPS))
	dt := time.Second / time.Duration(a.FPS)

	for i := 0; i < frameCount; i++ {
		currentTime := time.Duration(i) * dt

		// Update scene state
		if a.Update != nil {
			a.Update(currentTime)
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
