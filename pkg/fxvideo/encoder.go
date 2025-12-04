package fxvideo

import (
	"fmt"
	"image"
	"io"
	"os/exec"
)

// FXStreamEncoder defines an interface for streaming video encoding.
type FXStreamEncoder interface {
	// AddFrame adds an image frame to the video stream.
	AddFrame(img *image.RGBA) error
	// Close finishes the video stream and releases resources.
	Close() error
}

// fxMp4StreamEncoder implements StreamEncoder for MP4 output using ffmpeg.
// fxMp4StreamEncoder implements FXStreamEncoder for MP4 output using ffmpeg.
type fxMp4StreamEncoder struct {
	// cmd is the ffmpeg command.
	cmd *exec.Cmd
	// stdin is the stdin pipe to ffmpeg.
	stdin io.WriteCloser
	// width is the width of the video.
	width int
	// height is the height of the video.
	height int
}

// NewFXMP4StreamEncoder creates a new MP4StreamEncoder that writes to the provided writer.
func NewFXMP4StreamEncoder(writer io.Writer, width, height, fps int) (FXStreamEncoder, error) {
	// ffmpeg command to read raw rgba video from stdin and output mp4 to stdout
	args := []string{
		"-y", // Overwrite output files without asking
		"-f", "rawvideo",
		"-pix_fmt", "rgba",
		"-s", fmt.Sprintf("%dx%d", width, height),
		"-r", fmt.Sprintf("%d", fps),
		"-i", "-",
		"-c:v", "libx264",
		"-preset", "ultrafast",
		"-f", "mp4",
		"-movflags", "frag_keyframe+empty_moov",
		"-",
	}

	cmd := exec.Command("ffmpeg", args...)
	cmd.Stdout = writer

	stdin, err := cmd.StdinPipe()
	if err != nil {
		return nil, fmt.Errorf("failed to create stdin pipe: %w", err)
	}

	if err := cmd.Start(); err != nil {
		return nil, fmt.Errorf("failed to start ffmpeg: %w", err)
	}

	return &fxMp4StreamEncoder{
		cmd:    cmd,
		stdin:  stdin,
		width:  width,
		height: height,
	}, nil
}

// AddFrame writes a single frame to the encoder.
func (e *fxMp4StreamEncoder) AddFrame(img *image.RGBA) error {
	if img.Rect.Dx() != e.width || img.Rect.Dy() != e.height {
		return fmt.Errorf("frame dimension mismatch: expected %dx%d, got %dx%d", e.width, e.height, img.Rect.Dx(), img.Rect.Dy())
	}

	// Write raw pixels to ffmpeg stdin
	_, err := e.stdin.Write(img.Pix)
	if err != nil {
		return fmt.Errorf("failed to write frame to ffmpeg: %w", err)
	}
	return nil
}

// Close closes the input stream and waits for the encoding to finish.
func (e *fxMp4StreamEncoder) Close() error {
	if err := e.stdin.Close(); err != nil {
		return fmt.Errorf("failed to close stdin: %w", err)
	}

	if err := e.cmd.Wait(); err != nil {
		return fmt.Errorf("ffmpeg process failed: %w", err)
	}

	return nil
}
