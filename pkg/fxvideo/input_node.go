package fxvideo

import (
	"image"
	"time"

	"kdfx/pkg/fxcontext"
	"kdfx/pkg/fxcore"
	"kdfx/pkg/fxnode"
)

// FXVideoPlaybackMode defines how the video behaves when the requested time is outside its duration.
type FXVideoPlaybackMode int

const (
	FXModeLoop    FXVideoPlaybackMode = iota // Repeat video
	FXModeStretch                            // Stretch to target duration
	FXModeClamp                              // Hold last frame
	FXModeNone                               // Black/Transparent after end
)

// FXVideoInputNode represents a node that outputs video frames.
type FXVideoInputNode interface {
	fxnode.FXNode
	// SetMode sets the playback mode (Loop, Stretch, Clamp, None).
	SetMode(mode FXVideoPlaybackMode)
	// SetTargetDuration sets the target duration for Stretch mode.
	// The video will be sped up or slowed down to match this duration.
	SetTargetDuration(d time.Duration)
	// SetTime sets the current playback time.
	// This is typically called by the animation loop.
	SetTime(t time.Duration)
}

// fxVideoInputNode implements FXVideoInputNode.
type fxVideoInputNode struct {
	fxnode.FXNode
	// decoder is the video stream decoder.
	decoder FXStreamDecoder
	// texture is the texture where the video frame is uploaded.
	texture fxcore.FXTexture
	// img is the temporary image buffer.
	img *image.RGBA
	// mode is the playback mode (Loop, Stretch, Clamp, None).
	mode FXVideoPlaybackMode
	// targetDuration is the target duration for Stretch mode.
	targetDuration time.Duration
	// currentTime is the current playback time.
	currentTime time.Duration
}

// NewFXVideoInputNode creates a new video fxnode.
func NewFXVideoInputNode(ctx fxcontext.FXContext, path string) (FXVideoInputNode, error) {
	// Initialize the video decoder.
	decoder, err := NewFXStreamDecoder(path)
	if err != nil {
		return nil, err
	}

	info := decoder.Info()
	// Create a texture to store video frames.
	tex := fxcore.NewFXTexture(info.Width, info.Height)
	// if err != nil { ... } // NewTexture doesn't return error currently

	// Create a base node.
	base, err := fxnode.NewFXBaseNode(ctx, info.Width, info.Height)
	if err != nil {
		decoder.Close()
		return nil, err
	}

	return &fxVideoInputNode{
		FXNode:  base,
		decoder: decoder,
		texture: tex,
		img:     image.NewRGBA(image.Rect(0, 0, info.Width, info.Height)),
		mode:    FXModeLoop, // Default to loop
	}, nil
}

func (n *fxVideoInputNode) SetMode(mode FXVideoPlaybackMode) {
	n.mode = mode
}

func (n *fxVideoInputNode) SetTargetDuration(d time.Duration) {
	n.targetDuration = d
}

func (n *fxVideoInputNode) SetTime(t time.Duration) {
	n.currentTime = t
}

func (n *fxVideoInputNode) Process(ctx fxcontext.FXContext) error {
	info := n.decoder.Info()
	videoDuration := info.Duration
	var videoTime time.Duration

	// Calculate the video time based on the playback mode.
	switch n.mode {
	case FXModeLoop:
		// Loop the video if the current time exceeds duration.
		if videoDuration > 0 {
			videoTime = time.Duration(int64(n.currentTime) % int64(videoDuration))
		}
	case FXModeStretch:
		// Stretch the video to fit the target duration.
		if n.targetDuration > 0 {
			videoTime = time.Duration(float64(n.currentTime) * float64(videoDuration) / float64(n.targetDuration))
		} else {
			videoTime = n.currentTime
		}
	case FXModeClamp:
		// Clamp the video to the last frame if current time exceeds duration.
		if n.currentTime > videoDuration {
			videoTime = videoDuration
		} else {
			videoTime = n.currentTime
		}
	case FXModeNone:
		// Play normally, potentially going past the end (handling EOF later).
		videoTime = n.currentTime
	}

	// Seek decoder to the calculated time
	// This is efficient because the decoder handles seeking internally.
	if err := n.decoder.Seek(videoTime); err != nil {
		return err
	}

	// Read frame
	// Decode the frame into the image buffer.
	if err := n.decoder.ReadFrame(n.img); err != nil {
		if err.Error() == "EOF" && n.mode == FXModeLoop {
			// If we hit EOF in loop mode, try seeking to 0 and reading again
			if err := n.decoder.Seek(0); err != nil {
				return err
			}
			if err := n.decoder.ReadFrame(n.img); err != nil {
				return err
			}
		} else if n.mode == FXModeNone && err.Error() == "EOF" {
			// Ignore EOF for FXModeNone (keep last frame or clear?)
			// Keeping last frame is safer for now.
			return nil
		} else {
			return err
		}
	}

	// Upload to texture
	// Upload the decoded frame to the GPU texture.
	n.texture.Upload(n.img)

	return nil
}

func (n *fxVideoInputNode) GetTexture() fxcore.FXTexture {
	return n.texture
}

func (n *fxVideoInputNode) Release() {
	n.decoder.Close()
	n.FXNode.Release()
}
