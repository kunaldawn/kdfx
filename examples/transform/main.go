package main

import (
	"image"
	"image/color"
	"image/png"
	"math"
	"os"

	"kdfx/pkg/fxcontext"
	"kdfx/pkg/fxcore"
	"kdfx/pkg/fxlib/fxcolor"
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
	width, height := 800, 600
	ctx, err := fxcontext.NewFXOffscreenContext(width, height)
	if err != nil {
		panic(err)
	}
	defer ctx.Destroy()

	// 1. Create a test image (Grid pattern)
	img := image.NewRGBA(image.Rect(0, 0, 256, 256))
	for y := 0; y < 256; y++ {
		for x := 0; x < 256; x++ {
			if x%32 == 0 || y%32 == 0 {
				img.Set(x, y, color.RGBA{255, 255, 255, 255})
			} else {
				img.Set(x, y, color.RGBA{0, 0, 255, 255}) // Blue
			}
		}
	}
	saveImage("input.png", img)

	// 2. Upload to Texture
	inputTex, err := fxcore.FXLoadTextureFromFile("input.png")
	if err != nil {
		panic(err)
	}

	// 3. Create a Node (Color Adjustment)
	node, err := fxcolor.NewFXColorAdjustmentNode(ctx, width, height)
	if err != nil {
		panic(err)
	}

	inputNode := &InputNode{Texture: inputTex}
	node.SetInput("u_texture", inputNode)

	// 4. Apply Transformations
	// Move to the right and up slightly
	node.SetPosition(0.5, 0.2)

	// Scale down to half size
	node.SetSize(0.5, 0.5)

	// Rotate 45 degrees
	node.SetRotation(math.Pi / 4)

	// 5. Process
	if err := node.Process(ctx); err != nil {
		panic(err)
	}

	// 6. Save Output
	outImg, err := node.GetFramebuffer().GetTexture().Download()
	if err != nil {
		panic(err)
	}
	saveImage("output.png", outImg)
}

func saveImage(filename string, img image.Image) {
	f, err := os.Create(filename)
	if err != nil {
		panic(err)
	}
	defer f.Close()
	png.Encode(f, img)
}
