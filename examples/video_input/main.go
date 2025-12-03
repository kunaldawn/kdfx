package main

import (
	"fmt"
	"os"
	"time"

	"kdfx/pkg/fxcontext"
	"kdfx/pkg/fxlib/fxcolor"
	"kdfx/pkg/fxvideo"
)

func main() {
	// Use the output from animation example as input
	inputPath := "output.mp4"
	if _, err := os.Stat(inputPath); os.IsNotExist(err) {
		fmt.Println("Error: output.mp4 not found. Please run examples/animation/main.go first.")
		return
	}

	// Probe input to get dimensions
	info, err := fxvideo.FXProbeVideo(inputPath)
	if err != nil {
		panic(err)
	}
	fmt.Printf("Input Video: %dx%d @ %d fps, Duration: %v\n", info.Width, info.Height, info.FPS, info.Duration)

	width, height := info.Width, info.Height
	ctx, err := fxcontext.NewFXOffscreenContext(width, height)
	if err != nil {
		panic(err)
	}
	defer ctx.Destroy()

	// 1. Create Video Node
	videoNode, err := fxvideo.NewFXVideoInputNode(ctx, inputPath)
	if err != nil {
		panic(err)
	}
	defer videoNode.Release()

	// Set Loop Mode
	videoNode.SetMode(fxvideo.FXModeLoop)

	// 2. Apply Filter (Sepia)
	filterNode, err := fxcolor.NewFXColorFilterNode(ctx, width, height)
	if err != nil {
		panic(err)
	}
	filterNode.SetInput("u_texture", videoNode)
	filterNode.SetMode(fxcolor.FXFilterSepia)

	// 3. Render Output
	// Render for 2x the input duration to verify looping
	outDuration := info.Duration * 2
	outFile, err := os.Create("video_out.mp4")
	if err != nil {
		panic(err)
	}
	defer outFile.Close()

	anim := fxvideo.NewFXAnimation(outDuration, info.FPS, func(t time.Duration) {
		videoNode.SetTime(t)
	})

	fmt.Println("Rendering video_out.mp4...")
	startTime := time.Now()

	if err := anim.Render(ctx, filterNode, outFile); err != nil {
		panic(err)
	}

	fmt.Printf("Done! Rendered in %v\n", time.Since(startTime))
}
