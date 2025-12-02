package main

import (
	"fmt"
	"image"
	"image/color"
	"image/png"
	"os"
	"time"

	"kimg/context"
	"kimg/core"
	"kimg/export"
	"kimg/filters"
)

// InputNode is a simple node that just provides a texture.
type InputNode struct {
	Texture core.Texture
}

func (n *InputNode) GetTexture() core.Texture          { return n.Texture }
func (n *InputNode) IsDirty() bool                     { return false }
func (n *InputNode) Process(ctx context.Context) error { return nil }

func main() {
	width, height := 512, 512
	ctx, err := context.NewOffscreenContext(width, height)
	if err != nil {
		panic(err)
	}
	defer ctx.Destroy()

	// 1. Create a test image (Checkered pattern)
	img := image.NewRGBA(image.Rect(0, 0, width, height))
	for y := 0; y < height; y++ {
		for x := 0; x < width; x++ {
			if (x/32+y/32)%2 == 0 {
				img.Set(x, y, color.RGBA{255, 255, 255, 255})
			} else {
				img.Set(x, y, color.RGBA{0, 0, 0, 255})
			}
		}
	}

	// Save input for reference
	saveImage("input.png", img)

	// 2. Upload to Texture
	inputTex, err := core.LoadTextureFromFile("input.png")
	if err != nil {
		panic(err)
	}

	// 3. Build Graph
	inputNode := &InputNode{Texture: inputTex}

	bcNode, err := filters.NewBrightnessContrastNode(ctx, width, height)
	if err != nil {
		panic(err)
	}
	bcNode.SetInput("u_texture", inputNode)

	// 4. Setup Animation
	anim := export.NewAnimation(2*time.Second, 30, func(t time.Duration) {
		// Animate brightness from 0.0 to 2.0
		progress := float64(t) / float64(2*time.Second)
		brightness := 2.0 * progress
		bcNode.SetBrightness(float32(brightness))
		bcNode.SetContrast(1.0)
	})

	// 5. Render to MP4
	outFile, err := os.Create("output.mp4")
	if err != nil {
		panic(err)
	}
	defer outFile.Close()

	fmt.Println("Rendering animation to output.mp4...")
	startTime := time.Now()

	if err := anim.Render(ctx, bcNode, outFile); err != nil {
		panic(err)
	}

	fmt.Printf("Done! Rendered in %v\n", time.Since(startTime))
}

func saveImage(filename string, img image.Image) {
	f, err := os.Create(filename)
	if err != nil {
		panic(err)
	}
	defer f.Close()
	png.Encode(f, img)
}
