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
)

func main() {
	width, height := 512, 512
	ctx, err := fxcontext.NewFXOffscreenContext(width, height)
	if err != nil {
		panic(err)
	}
	defer ctx.Destroy()

	// 1. Create Test Image (Grid with Circle)
	img := image.NewRGBA(image.Rect(0, 0, width, height))
	for y := 0; y < height; y++ {
		for x := 0; x < width; x++ {
			// Grid
			if x%64 == 0 || y%64 == 0 {
				img.Set(x, y, color.RGBA{255, 255, 255, 255})
				continue
			}
			// Circle
			dx, dy := float64(x-width/2), float64(y-height/2)
			if dx*dx+dy*dy < 100*100 {
				img.Set(x, y, color.RGBA{255, 0, 0, 255})
			} else {
				img.Set(x, y, color.RGBA{0, 0, 0, 255})
			}
		}
	}
	saveImage("input.png", img)
	inputNode, err := fximage.NewFXImageInputFromFile("input.png")
	if err != nil {
		panic(err)
	}

	// Output Node
	outputNode := fximage.NewFXImageOutput()

	// 2. Test Gaussian Blur
	fmt.Println("Testing Gaussian Blur...")
	gNode, err := fxblur.NewFXGaussianBlurNode(ctx, width, height)
	if err != nil {
		panic(err)
	}
	gNode.SetInput("u_texture", inputNode)
	gNode.SetRadius(10.0)
	outputNode.SetInput(gNode)
	if err := outputNode.Process(ctx); err != nil {
		panic(err)
	}
	if err := outputNode.Save("output_gaussian.png"); err != nil {
		panic(err)
	}

	// 3. Test Box Blur
	fmt.Println("Testing Box Blur...")
	bNode, err := fxblur.NewFXBoxBlurNode(ctx, width, height)
	if err != nil {
		panic(err)
	}
	bNode.SetInput("u_texture", inputNode)
	bNode.SetRadius(5.0)
	outputNode.SetInput(bNode)
	if err := outputNode.Process(ctx); err != nil {
		panic(err)
	}
	if err := outputNode.Save("output_box.png"); err != nil {
		panic(err)
	}

	// 4. Test Radial Blur
	fmt.Println("Testing Radial Blur...")
	rNode, err := fxblur.NewFXRadialBlurNode(ctx, width, height)
	if err != nil {
		panic(err)
	}
	rNode.SetInput("u_texture", inputNode)
	rNode.SetStrength(0.05)
	outputNode.SetInput(rNode)
	if err := outputNode.Process(ctx); err != nil {
		panic(err)
	}
	if err := outputNode.Save("output_radial.png"); err != nil {
		panic(err)
	}

	// 5. Test Motion Blur
	fmt.Println("Testing Motion Blur...")
	mNode, err := fxblur.NewFXMotionBlurNode(ctx, width, height)
	if err != nil {
		panic(err)
	}
	mNode.SetInput("u_texture", inputNode)
	mNode.SetAngle(45.0)
	mNode.SetStrength(0.05)
	outputNode.SetInput(mNode)
	if err := outputNode.Process(ctx); err != nil {
		panic(err)
	}
	if err := outputNode.Save("output_motion.png"); err != nil {
		panic(err)
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
