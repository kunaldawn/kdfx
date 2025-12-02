package export

import (
	"fmt"
	"image"
	"io"
	"os/exec"
)

// StreamEncoder defines an interface for streaming video encoding.
type StreamEncoder interface {
	AddFrame(img *image.RGBA) error
	Close() error
}

// MP4StreamEncoder implements StreamEncoder for MP4 output using ffmpeg.
type MP4StreamEncoder struct {
	cmd    *exec.Cmd
	stdin  io.WriteCloser
	width  int
	height int
}

// NewMP4StreamEncoder creates a new MP4StreamEncoder that writes to the provided writer.
func NewMP4StreamEncoder(writer io.Writer, width, height, fps int) (*MP4StreamEncoder, error) {
	// ffmpeg command to read raw rgba video from stdin and output mp4 to stdout
	// -f rawvideo: input format
	// -pix_fmt rgba: input pixel format
	// -s WxH: input resolution
	// -r FPS: input frame rate
	// -i -: read from stdin
	// -c:v libx264: video codec
	// -preset ultrafast: encoding speed (faster = larger file, less cpu)
	// -f mp4: output format
	// -movflags frag_keyframe+empty_moov: required for streaming mp4 (fragmented mp4)
	// -: write to stdout

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
	// We can capture stderr for debugging if needed, or pipe it to os.Stderr
	// cmd.Stderr = os.Stderr

	stdin, err := cmd.StdinPipe()
	if err != nil {
		return nil, fmt.Errorf("failed to create stdin pipe: %w", err)
	}

	if err := cmd.Start(); err != nil {
		return nil, fmt.Errorf("failed to start ffmpeg: %w", err)
	}

	return &MP4StreamEncoder{
		cmd:    cmd,
		stdin:  stdin,
		width:  width,
		height: height,
	}, nil
}

// AddFrame writes a single frame to the encoder.
func (e *MP4StreamEncoder) AddFrame(img *image.RGBA) error {
	if img.Rect.Dx() != e.width || img.Rect.Dy() != e.height {
		return fmt.Errorf("frame dimension mismatch: expected %dx%d, got %dx%d", e.width, e.height, img.Rect.Dx(), img.Rect.Dy())
	}

	// Write raw pixels to ffmpeg stdin
	// image.RGBA Pix is already in the right format (R, G, B, A packed)
	_, err := e.stdin.Write(img.Pix)
	if err != nil {
		return fmt.Errorf("failed to write frame to ffmpeg: %w", err)
	}
	return nil
}

// Close closes the input stream and waits for the encoding to finish.
func (e *MP4StreamEncoder) Close() error {
	if err := e.stdin.Close(); err != nil {
		return fmt.Errorf("failed to close stdin: %w", err)
	}

	if err := e.cmd.Wait(); err != nil {
		return fmt.Errorf("ffmpeg process failed: %w", err)
	}

	return nil
}
