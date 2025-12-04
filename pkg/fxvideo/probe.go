package fxvideo

import (
	"fmt"
	"os/exec"
	"strconv"
	"strings"
	"time"
)

// FXVideoInfo contains metadata about a video file.
type FXVideoInfo struct {
	// Width is the width of the video in pixels.
	Width int
	// Height is the height of the video in pixels.
	Height int
	// FPS is the frames per second of the video.
	FPS int
	// Duration is the total duration of the video.
	Duration time.Duration
}

// FXProbeVideo extracts metadata from a video file using ffprobe.
func FXProbeVideo(path string) (*FXVideoInfo, error) {
	// ffprobe -v error -select_streams v:0 -show_entries stream=width,height,r_frame_rate,duration -of csv=p=0 <path>
	cmd := exec.Command("ffprobe",
		"-v", "error",
		"-select_streams", "v:0",
		"-show_entries", "stream=width,height,r_frame_rate,duration",
		"-of", "csv=p=0",
		path,
	)

	output, err := cmd.Output()
	if err != nil {
		return nil, fmt.Errorf("ffprobe failed: %w", err)
	}

	// Output format: width,height,r_frame_rate,duration
	// Example: 1920,1080,30/1,10.5
	parts := strings.Split(strings.TrimSpace(string(output)), ",")
	if len(parts) < 4 {
		return nil, fmt.Errorf("unexpected ffprobe output: %s", string(output))
	}

	width, err := strconv.Atoi(parts[0])
	if err != nil {
		return nil, fmt.Errorf("invalid width: %v", err)
	}

	height, err := strconv.Atoi(parts[1])
	if err != nil {
		return nil, fmt.Errorf("invalid height: %v", err)
	}

	fps, err := parseFPS(parts[2])
	if err != nil {
		return nil, fmt.Errorf("invalid fps: %v", err)
	}

	durationSec, err := strconv.ParseFloat(parts[3], 64)
	if err != nil {
		return nil, fmt.Errorf("invalid duration: %v", err)
	}

	return &FXVideoInfo{
		Width:    width,
		Height:   height,
		FPS:      fps,
		Duration: time.Duration(durationSec * float64(time.Second)),
	}, nil
}

func parseFPS(fpsStr string) (int, error) {
	// fps can be "30/1" or "30"
	parts := strings.Split(fpsStr, "/")
	if len(parts) == 1 {
		return strconv.Atoi(parts[0])
	}
	if len(parts) == 2 {
		num, err := strconv.Atoi(parts[0])
		if err != nil {
			return 0, err
		}
		den, err := strconv.Atoi(parts[1])
		if err != nil {
			return 0, err
		}
		if den == 0 {
			return 0, fmt.Errorf("division by zero in fps")
		}
		return num / den, nil // Integer FPS for now
	}
	return 0, fmt.Errorf("invalid fps format: %s", fpsStr)
}
