package main

import (
	"fmt"
	"image"
	"image/color"
	"image/png"
	"os"

	"kimg/pkg/context"
	"kimg/pkg/core"
	"kimg/pkg/fxlib/blend"
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

	// 1. Create Base Image (Horizontal Gradient)
	baseImg := image.NewRGBA(image.Rect(0, 0, width, height))
	for y := 0; y < height; y++ {
		for x := 0; x < width; x++ {
			val := uint8(float64(x) / float64(width) * 255)
			baseImg.Set(x, y, color.RGBA{val, val, val, 255})
		}
	}
	saveImage("base.png", baseImg)
	baseTex, _ := core.LoadTextureFromFile("base.png")
	baseNode := &InputNode{Texture: baseTex}

	// 2. Create Blend Image (Vertical Gradient)
	blendImg := image.NewRGBA(image.Rect(0, 0, width, height))
	for y := 0; y < height; y++ {
		for x := 0; x < width; x++ {
			val := uint8(float64(y) / float64(height) * 255)
			blendImg.Set(x, y, color.RGBA{val, 0, 0, 255}) // Red gradient
		}
	}
	saveImage("blend.png", blendImg)
	blendTex, _ := core.LoadTextureFromFile("blend.png")
	blendInputNode := &InputNode{Texture: blendTex}

	// 3. Create Blend Node
	blendNode, err := blend.NewBlendNode(ctx, width, height)
	if err != nil {
		panic(err)
	}
	blendNode.SetInput1(baseNode)
	blendNode.SetInput2(blendInputNode)
	blendNode.SetFactor(1.0)

	// 4. Test Modes
	modes := map[string]blend.BlendMode{
		"normal":     blend.BlendNormal,
		"add":        blend.BlendAdd,
		"multiply":   blend.BlendMultiply,
		"screen":     blend.BlendScreen,
		"overlay":    blend.BlendOverlay,
		"difference": blend.BlendDifference,
	}

	for name, mode := range modes {
		fmt.Printf("Testing mode: %s\n", name)
		blendNode.SetMode(mode)
		if err := blendNode.Process(ctx); err != nil {
			panic(err)
		}

		outTex := blendNode.GetTexture()
		outImg, _ := outTex.Download()
		saveImage(fmt.Sprintf("output_%s.png", name), outImg)
	}

	fmt.Println("Done! Check output_*.png files.")
}

func saveImage(filename string, img image.Image) {
	f, err := os.Create(filename)
	if err != nil {
		panic(err)
	}
	defer f.Close()
	png.Encode(f, img)
}
