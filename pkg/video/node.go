package video

import (
	"image"
	"time"

	"kdfx/pkg/context"
	"kdfx/pkg/core"
	"kdfx/pkg/node"
)

// VideoPlaybackMode defines how the video behaves when the requested time is outside its duration.
type VideoPlaybackMode int

const (
	ModeLoop    VideoPlaybackMode = iota // Repeat video
	ModeStretch                          // Stretch to target duration
	ModeClamp                            // Hold last frame
	ModeNone                             // Black/Transparent after end
)

// VideoNode represents a node that outputs video frames.
type VideoNode interface {
	node.Node
	SetMode(mode VideoPlaybackMode)
	SetTargetDuration(d time.Duration) // For Stretch mode
	SetTime(t time.Duration)           // Called by animation loop
}

type videoNode struct {
	node.Node
	decoder        StreamDecoder
	texture        core.Texture
	img            *image.RGBA
	mode           VideoPlaybackMode
	targetDuration time.Duration
	currentTime    time.Duration
}

// NewVideoNode creates a new video node.
func NewVideoNode(ctx context.Context, path string) (VideoNode, error) {
	decoder, err := NewStreamDecoder(path)
	if err != nil {
		return nil, err
	}

	info := decoder.Info()
	tex := core.NewTexture(info.Width, info.Height)
	// if err != nil { ... } // NewTexture doesn't return error currently

	base, err := node.NewBaseNode(ctx, info.Width, info.Height)
	if err != nil {
		decoder.Close()
		return nil, err
	}

	return &videoNode{
		Node:    base,
		decoder: decoder,
		texture: tex,
		img:     image.NewRGBA(image.Rect(0, 0, info.Width, info.Height)),
		mode:    ModeLoop, // Default to loop
	}, nil
}

func (n *videoNode) SetMode(mode VideoPlaybackMode) {
	n.mode = mode
}

func (n *videoNode) SetTargetDuration(d time.Duration) {
	n.targetDuration = d
}

func (n *videoNode) SetTime(t time.Duration) {
	n.currentTime = t
}

func (n *videoNode) Process(ctx context.Context) error {
	info := n.decoder.Info()
	videoDuration := info.Duration
	var videoTime time.Duration

	switch n.mode {
	case ModeLoop:
		if videoDuration > 0 {
			videoTime = time.Duration(int64(n.currentTime) % int64(videoDuration))
		}
	case ModeStretch:
		if n.targetDuration > 0 {
			videoTime = time.Duration(float64(n.currentTime) * float64(videoDuration) / float64(n.targetDuration))
		} else {
			videoTime = n.currentTime
		}
	case ModeClamp:
		if n.currentTime > videoDuration {
			videoTime = videoDuration
		} else {
			videoTime = n.currentTime
		}
	case ModeNone:
		videoTime = n.currentTime
	}

	// Seek decoder to the calculated time
	if err := n.decoder.Seek(videoTime); err != nil {
		return err
	}

	// Read frame
	if err := n.decoder.ReadFrame(n.img); err != nil {
		if err.Error() == "EOF" && n.mode == ModeLoop {
			// If we hit EOF in loop mode, try seeking to 0 and reading again
			if err := n.decoder.Seek(0); err != nil {
				return err
			}
			if err := n.decoder.ReadFrame(n.img); err != nil {
				return err
			}
		} else if n.mode == ModeNone && err.Error() == "EOF" {
			// Ignore EOF for ModeNone (keep last frame or clear?)
			// Keeping last frame is safer for now.
			return nil
		} else {
			return err
		}
	}

	// Upload to texture
	n.texture.Upload(n.img)

	return nil
}

func (n *videoNode) GetTexture() core.Texture {
	return n.texture
}

func (n *videoNode) Release() {
	n.decoder.Close()
	n.Node.Release()
}
