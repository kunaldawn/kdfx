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
	SetMode(mode FXVideoPlaybackMode)
	SetTargetDuration(d time.Duration) // For Stretch mode
	SetTime(t time.Duration)           // Called by fxAnimation loop
}

type fxVideoInputNode struct {
	fxnode.FXNode
	decoder        FXStreamDecoder
	texture        fxcore.FXTexture
	img            *image.RGBA
	mode           FXVideoPlaybackMode
	targetDuration time.Duration
	currentTime    time.Duration
}

// NewFXVideoInputNode creates a new video fxnode.
func NewFXVideoInputNode(ctx fxcontext.FXContext, path string) (FXVideoInputNode, error) {
	decoder, err := NewFXStreamDecoder(path)
	if err != nil {
		return nil, err
	}

	info := decoder.Info()
	tex := fxcore.NewFXTexture(info.Width, info.Height)
	// if err != nil { ... } // NewTexture doesn't return error currently

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

	switch n.mode {
	case FXModeLoop:
		if videoDuration > 0 {
			videoTime = time.Duration(int64(n.currentTime) % int64(videoDuration))
		}
	case FXModeStretch:
		if n.targetDuration > 0 {
			videoTime = time.Duration(float64(n.currentTime) * float64(videoDuration) / float64(n.targetDuration))
		} else {
			videoTime = n.currentTime
		}
	case FXModeClamp:
		if n.currentTime > videoDuration {
			videoTime = videoDuration
		} else {
			videoTime = n.currentTime
		}
	case FXModeNone:
		videoTime = n.currentTime
	}

	// Seek decoder to the calculated time
	if err := n.decoder.Seek(videoTime); err != nil {
		return err
	}

	// Read frame
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
