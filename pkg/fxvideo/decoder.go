package fxvideo

import (
	"fmt"
	"image"
	"io"
	"os/exec"
	"time"
)

// FXStreamDecoder decodes a video stream.
type FXStreamDecoder interface {
	// Seek seeks to the specified timestamp.
	// If t < current, restarts ffmpeg.
	// If t > current, reads and discards frames.
	Seek(t time.Duration) error
	// ReadFrame reads the next frame into the provided image buffer.
	ReadFrame(img *image.RGBA) error
	Close() error
	Info() FXVideoInfo
}

// fxFfmpegStreamDecoder implements FXStreamDecoder using ffmpeg.
type fxFfmpegStreamDecoder struct {
	// path is the path to the video file.
	path string
	// info contains metadata about the video.
	info FXVideoInfo
	// cmd is the ffmpeg command.
	cmd *exec.Cmd
	// stdout is the stdout pipe from ffmpeg.
	stdout io.ReadCloser
	// currentTime is the current playback time.
	currentTime time.Duration
}

// NewFXStreamDecoder creates a new StreamDecoder for the given video file.
func NewFXStreamDecoder(path string) (FXStreamDecoder, error) {
	info, err := FXProbeVideo(path)
	if err != nil {
		return nil, err
	}

	decoder := &fxFfmpegStreamDecoder{
		path: path,
		info: *info,
	}

	if err := decoder.startFFmpeg(0); err != nil {
		return nil, err
	}

	return decoder, nil
}

func (d *fxFfmpegStreamDecoder) startFFmpeg(startTime time.Duration) error {
	if d.cmd != nil {
		d.Close()
	}

	// Start ffmpeg at the specified time
	// -ss before -i is faster (keyframe seeking) but less accurate.
	// -ss after -i is slower (decoding) but accurate.
	// For exact frame, we might need accurate seeking.
	// Let's try -ss before -i for now, it's usually good enough for playback.
	args := []string{
		"-ss", fmt.Sprintf("%f", startTime.Seconds()),
		"-i", d.path,
		"-f", "rawvideo",
		"-pix_fmt", "rgba",
		"-",
	}

	cmd := exec.Command("ffmpeg", args...)
	stdout, err := cmd.StdoutPipe()
	if err != nil {
		return fmt.Errorf("failed to create stdout pipe: %w", err)
	}

	if err := cmd.Start(); err != nil {
		return fmt.Errorf("failed to start ffmpeg: %w", err)
	}

	d.cmd = cmd
	d.stdout = stdout
	d.currentTime = startTime
	return nil
}

func (d *fxFfmpegStreamDecoder) Seek(t time.Duration) error {
	// If seeking backwards, we MUST restart.
	// If seeking forwards, we could skip frames, but restarting with -ss is often more efficient for large jumps.
	// For small jumps, skipping might be better.
	// For simplicity and robustness, let's always restart for now if the difference is significant (> 1 sec).
	// If difference is small and positive, we can read and discard.

	delta := t - d.currentTime
	if delta < 0 || delta > 1*time.Second {
		return d.startFFmpeg(t)
	}

	// Skip frames
	framesToSkip := int(delta.Seconds() * float64(d.info.FPS))
	if framesToSkip > 0 {
		frameSize := d.info.Width * d.info.Height * 4
		discard := make([]byte, frameSize)

		for i := 0; i < framesToSkip; i++ {
			if _, err := io.ReadFull(d.stdout, discard); err != nil {
				return err
			}
			d.currentTime += time.Second / time.Duration(d.info.FPS)
		}
	}

	return nil
}

func (d *fxFfmpegStreamDecoder) ReadFrame(img *image.RGBA) error {
	if img.Rect.Dx() != d.info.Width || img.Rect.Dy() != d.info.Height {
		return fmt.Errorf("image dimension mismatch: expected %dx%d, got %dx%d", d.info.Width, d.info.Height, img.Rect.Dx(), img.Rect.Dy())
	}

	_, err := io.ReadFull(d.stdout, img.Pix)
	if err != nil {
		return err
	}

	d.currentTime += time.Second / time.Duration(d.info.FPS)
	return nil
}

func (d *fxFfmpegStreamDecoder) Close() error {
	if d.stdout != nil {
		d.stdout.Close()
	}
	if d.cmd != nil {
		d.cmd.Process.Kill() // Force kill if still running
		d.cmd.Wait()
	}
	return nil
}

func (d *fxFfmpegStreamDecoder) Info() FXVideoInfo {
	return d.info
}
