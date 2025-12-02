package main

import (
	"fmt"
	"image"
	"image/color"
	"image/png"
	"os"

	"kimg/context"
	"kimg/core"
	"kimg/filters"
)

// InputNode is a simple node that just provides a texture.
type InputNode struct {
	Texture *core.Texture
}

func (n *InputNode) GetTexture() *core.Texture { return n.Texture }
func (n *InputNode) IsDirty() bool             { return false } // Static input

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
	// Load input from file
	inputTex, err := core.LoadTextureFromFile("input.png")
	if err != nil {
		panic(err)
	}

	// 3. Build Graph
	// Input -> Brightness -> Blur -> Output

	inputNode := &InputNode{Texture: inputTex}

	bcNode, err := filters.NewBrightnessContrastNode(ctx, width, height)
	if err != nil {
		panic(err)
	}
	bcNode.SetInput("u_texture", inputNode)
	bcNode.SetBrightness(0.1) // Increase brightness
	bcNode.SetContrast(1.2)   // Increase contrast

	blurNode, err := filters.NewGaussianBlurNode(ctx, width, height)
	if err != nil {
		panic(err)
	}
	blurNode.SetInput("u_texture", bcNode)
	blurNode.SetRadius(10.0)

	// 4. Execute Pipeline
	// We can just call Process on the last node
	err = blurNode.Process(ctx)
	if err != nil {
		panic(err)
	}

	// 5. Download Result
	outTex := blurNode.GetTexture() // BaseNode.GetTexture() returns Output.Texture
	// Wait, BaseNode.GetOutput() returns *core.Texture?
	// In BaseNode: func (n *BaseNode) GetTexture() *core.Texture { return n.Output.Texture }
	// So yes.

	outImg, err := outTex.Download()
	if err != nil {
		panic(err)
	}

	saveImage("output.png", outImg)
	fmt.Println("Done! Saved output.png")
}

func saveImage(filename string, img image.Image) {
	f, err := os.Create(filename)
	if err != nil {
		panic(err)
	}
	defer f.Close()
	png.Encode(f, img)
}
