package main

import (
	"fmt"
	"os"
	"time"

	"kdfx/pkg/context"
	"kdfx/pkg/export"
	colorfx "kdfx/pkg/fxlib/color"
	"kdfx/pkg/video"
)

func main() {
	// Use the output from animation example as input
	inputPath := "output.mp4"
	if _, err := os.Stat(inputPath); os.IsNotExist(err) {
		fmt.Println("Error: output.mp4 not found. Please run examples/animation/main.go first.")
		return
	}

	// Probe input to get dimensions
	info, err := video.ProbeVideo(inputPath)
	if err != nil {
		panic(err)
	}
	fmt.Printf("Input Video: %dx%d @ %d fps, Duration: %v\n", info.Width, info.Height, info.FPS, info.Duration)

	width, height := info.Width, info.Height
	ctx, err := context.NewOffscreenContext(width, height)
	if err != nil {
		panic(err)
	}
	defer ctx.Destroy()

	// 1. Create Video Node
	vNode, err := video.NewVideoNode(ctx, inputPath)
	if err != nil {
		panic(err)
	}
	defer vNode.Release()

	// Set Loop Mode
	vNode.SetMode(video.ModeLoop)

	// 2. Apply Filter (Sepia)
	filterNode, err := colorfx.NewColorFilterNode(ctx, width, height)
	if err != nil {
		panic(err)
	}
	filterNode.SetInput("u_texture", vNode)
	filterNode.SetMode(colorfx.FilterSepia)

	// 3. Render Output
	// Render for 2x the input duration to verify looping
	outDuration := info.Duration * 2
	outFile, err := os.Create("video_out.mp4")
	if err != nil {
		panic(err)
	}
	defer outFile.Close()

	anim := export.NewAnimation(outDuration, info.FPS, func(t time.Duration) {
		vNode.SetTime(t)
	})

	fmt.Println("Rendering video_out.mp4...")
	startTime := time.Now()

	if err := anim.Render(ctx, filterNode, outFile); err != nil {
		panic(err)
	}

	fmt.Printf("Done! Rendered in %v\n", time.Since(startTime))
}
