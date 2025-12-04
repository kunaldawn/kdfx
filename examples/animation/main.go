package main

import (
	"fmt"
	"image"
	"image/color"
	"image/png"
	"os"
	"time"

	"kdfx/pkg/fxcontext"
	"kdfx/pkg/fxcore"
	"kdfx/pkg/fxlib/fxblur"
	"kdfx/pkg/fxlib/fxcolor"
	"kdfx/pkg/fxvideo"
	"math"
)

// InputNode is a simple node that just provides a texture.
type InputNode struct {
	Texture fxcore.FXTexture
}

func (n *InputNode) GetTexture() fxcore.FXTexture            { return n.Texture }
func (n *InputNode) IsDirty() bool                           { return false }
func (n *InputNode) Process(ctx fxcontext.FXContext) error   { return nil }
func (n *InputNode) SetInput(name string, input interface{}) {} // Dummy

func main() {
	width, height := 512, 512
	ctx, err := fxcontext.NewFXOffscreenContext(width, height)
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
	inputTex, err := fxcore.FXLoadTextureFromFile("input.png")
	if err != nil {
		panic(err)
	}

	// 3. Build Graph
	inputNode := &InputNode{Texture: inputTex}

	// Color Adjustment Node
	bcNode, err := fxcolor.NewFXColorAdjustmentNode(ctx, width, height)
	if err != nil {
		panic(err)
	}
	bcNode.SetInput("u_texture", inputNode)

	// Motion Blur Node
	mbNode, err := fxblur.NewFXMotionBlurNode(ctx, width, height)
	if err != nil {
		panic(err)
	}
	mbNode.SetInput("u_texture", bcNode)

	// 4. Setup Animation
	duration := 60 * time.Second
	anim := fxvideo.NewFXAnimation(duration, 30, func(t time.Duration) {
		progress := float64(t) / float64(duration)

		// Animate brightness (pulse)
		brightness := float32(math.Sin(progress * math.Pi * 2))
		bcNode.SetBrightness(brightness)

		// Animate Transform (Rotate and Scale)
		rot := float32(progress * math.Pi * 2)
		bcNode.SetRotation(rot)

		scale := float32(0.8 + 0.2*math.Sin(progress*math.Pi*4))
		bcNode.SetSize(scale, scale)

		// Animate Position (Orbit)
		posX := float32(0.2 * math.Cos(progress*math.Pi*2))
		posY := float32(0.2 * math.Sin(progress*math.Pi*2))
		bcNode.SetPosition(posX, posY)

		// Animate Motion Blur (rotate angle, pulse strength)
		angle := float32(progress * 360.0)
		strength := float32(0.05 * (1.0 + math.Sin(progress*math.Pi*4)))

		mbNode.SetAngle(angle)
		mbNode.SetStrength(strength)
	})

	// 5. Render to MP4
	outFile, err := os.Create("output.mp4")
	if err != nil {
		panic(err)
	}
	defer outFile.Close()

	fmt.Println("Rendering animation to output.mp4...")
	startTime := time.Now()

	if err := anim.Render(ctx, mbNode, outFile); err != nil {
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
