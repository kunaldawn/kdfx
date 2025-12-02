package main

import (
	"fmt"
	"image"
	"image/color"
	"image/png"
	"os"

	"kdfx/pkg/context"
	"kdfx/pkg/core"
	"kdfx/pkg/fxlib/blur"
)

// InputNode is a simple node that just provides a texture.
type InputNode struct {
	Texture core.Texture
}

func (n *InputNode) GetTexture() core.Texture                { return n.Texture }
func (n *InputNode) IsDirty() bool                           { return false }
func (n *InputNode) Process(ctx context.Context) error       { return nil }
func (n *InputNode) SetInput(name string, input interface{}) {} // Dummy

func main() {
	width, height := 512, 512
	ctx, err := context.NewOffscreenContext(width, height)
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
	inputTex, _ := core.LoadTextureFromFile("input.png")
	inputNode := &InputNode{Texture: inputTex}

	// 2. Test Gaussian Blur
	fmt.Println("Testing Gaussian Blur...")
	gNode, err := blur.NewGaussianBlurNode(ctx, width, height)
	if err != nil {
		panic(err)
	}
	gNode.SetInput("u_texture", inputNode)
	gNode.SetRadius(10.0)
	if err := gNode.Process(ctx); err != nil {
		panic(err)
	}
	saveImage("output_gaussian.png", mustDownload(gNode.GetTexture()))

	// 3. Test Box Blur
	fmt.Println("Testing Box Blur...")
	bNode, err := blur.NewBoxBlurNode(ctx, width, height)
	if err != nil {
		panic(err)
	}
	bNode.SetInput("u_texture", inputNode)
	bNode.SetRadius(5.0)
	if err := bNode.Process(ctx); err != nil {
		panic(err)
	}
	saveImage("output_box.png", mustDownload(bNode.GetTexture()))

	// 4. Test Radial Blur
	fmt.Println("Testing Radial Blur...")
	rNode, err := blur.NewRadialBlurNode(ctx, width, height)
	if err != nil {
		panic(err)
	}
	rNode.SetInput("u_texture", inputNode)
	rNode.SetStrength(0.05)
	if err := rNode.Process(ctx); err != nil {
		panic(err)
	}
	saveImage("output_radial.png", mustDownload(rNode.GetTexture()))

	// 5. Test Motion Blur
	fmt.Println("Testing Motion Blur...")
	mNode, err := blur.NewMotionBlurNode(ctx, width, height)
	if err != nil {
		panic(err)
	}
	mNode.SetInput("u_texture", inputNode)
	mNode.SetAngle(45.0)
	mNode.SetStrength(0.05)
	if err := mNode.Process(ctx); err != nil {
		panic(err)
	}
	saveImage("output_motion.png", mustDownload(mNode.GetTexture()))

	fmt.Println("Done! Check output_*.png files.")
}

func mustDownload(t core.Texture) image.Image {
	img, err := t.Download()
	if err != nil {
		panic(err)
	}
	return img
}

func saveImage(filename string, img image.Image) {
	f, err := os.Create(filename)
	if err != nil {
		panic(err)
	}
	defer f.Close()
	png.Encode(f, img)
}
