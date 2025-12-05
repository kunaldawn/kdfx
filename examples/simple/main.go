package main

import (
	"fmt"
	"image"
	"image/color"
	"image/png"
	"os"

	"kdfx/pkg/fxcontext"
	"kdfx/pkg/fximage"
	"kdfx/pkg/fxlib/fxblur"
	"kdfx/pkg/fxlib/fxcolor"
)

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

	// 2. Create Input Node
	// Load input from file using the new FXImageInput node
	inputNode, err := fximage.NewFXImageInputFromFile("input.png")
	if err != nil {
		panic(err)
	}

	// 3. Build Graph
	// Input -> Brightness -> Blur -> Output

	bcNode, err := fxcolor.NewFXColorAdjustmentNode(ctx, width, height)
	if err != nil {
		panic(err)
	}
	bcNode.SetInput("u_texture", inputNode)
	bcNode.SetBrightness(0.1) // Increase brightness
	bcNode.SetContrast(1.2)   // Increase contrast

	blurNode, err := fxblur.NewFXGaussianBlurNode(ctx, width, height)
	if err != nil {
		panic(err)
	}
	blurNode.SetInput("u_texture", bcNode)
	blurNode.SetRadius(10.0)

	// 4. Create Output Node
	outputNode := fximage.NewFXImageOutput()
	outputNode.SetInput(blurNode)

	// 5. Execute Pipeline
	// We call Process on the output node (which delegates to input) or directly on the last processing node.
	// Since FXImageOutput.Process delegates, we can use it.
	err = outputNode.Process(ctx)
	if err != nil {
		panic(err)
	}

	// 6. Save Result
	err = outputNode.Save("output.png")
	if err != nil {
		panic(err)
	}

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
