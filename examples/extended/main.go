package main

import (
	"fmt"
	"os"

	"kdfx/pkg/fxcontext"
	"kdfx/pkg/fximage"
	"kdfx/pkg/fxlib/fxartistic"
	"kdfx/pkg/fxlib/fxblur"
	"kdfx/pkg/fxlib/fxcolor"
	"kdfx/pkg/fxlib/fxdistortion"
)

func main() {
	width, height := 512, 512
	// Initialize context
	ctx, err := fxcontext.NewFXOffscreenContext(width, height)
	if err != nil {
		panic(err)
	}
	defer ctx.Destroy()

	// Check if base.png exists, if not panic
	if _, err := os.Stat("base.png"); os.IsNotExist(err) {
		panic("base.png not found in current directory")
	}

	// Create input node
	inputNode, err := fximage.NewFXImageInputFromFile("base.png")
	if err != nil {
		panic(err)
	}

	// 1. Levels
	levelsNode, _ := fxcolor.NewFXLevelsNode(ctx, width, height)
	defer levelsNode.Release()
	levelsNode.SetInput("u_texture", inputNode)
	levelsNode.SetInputLevels(0.1, 0.9) // Increase contrast
	levelsNode.SetGamma(1.2)

	// 2. Color Balance
	balanceNode, _ := fxcolor.NewFXColorBalanceNode(ctx, width, height)
	defer balanceNode.Release()
	balanceNode.SetInput("u_texture", levelsNode)
	balanceNode.SetMidtones(0.1, 0.0, -0.1) // Warmer midtones

	// 3. Sharpen
	sharpenNode, _ := fxblur.NewFXSharpenNode(ctx, width, height)
	defer sharpenNode.Release()
	sharpenNode.SetInput("u_texture", balanceNode)
	sharpenNode.SetAmount(0.8)

	// 4. Ripple
	rippleNode, _ := fxdistortion.NewFXRippleNode(ctx, width, height)
	defer rippleNode.Release()
	rippleNode.SetInput("u_texture", sharpenNode)
	rippleNode.SetAmplitude(0.005)
	rippleNode.SetFrequency(30.0)

	// 5. Twirl
	twirlNode, _ := fxdistortion.NewFXTwirlNode(ctx, width, height)
	defer twirlNode.Release()
	twirlNode.SetInput("u_texture", rippleNode)
	twirlNode.SetRadius(0.3)
	twirlNode.SetAngle(1.0)

	// 6. Pixelize
	pixelizeNode, _ := fxartistic.NewFXPixelizeNode(ctx, width, height)
	defer pixelizeNode.Release()
	pixelizeNode.SetInput("u_texture", twirlNode)
	pixelizeNode.SetPixelSize(4.0)

	// 7. Vignette
	vignetteNode, _ := fxartistic.NewFXVignetteNode(ctx, width, height)
	defer vignetteNode.Release()
	vignetteNode.SetInput("u_texture", pixelizeNode)
	vignetteNode.SetRadius(0.8)
	vignetteNode.SetOpacity(0.6)

	// 8. Oil Paint
	oilPaintNode, _ := fxartistic.NewFXOilPaintNode(ctx, width, height)
	defer oilPaintNode.Release()
	oilPaintNode.SetInput("u_texture", vignetteNode)
	oilPaintNode.SetRadius(3)

	// 9. Bloom
	bloomNode, _ := fxartistic.NewFXBloomNode(ctx, width, height)
	defer bloomNode.Release()
	bloomNode.SetInput("u_texture", oilPaintNode)
	bloomNode.SetThreshold(0.6)
	bloomNode.SetIntensity(0.8)

	// Output main chain
	outputNode := fximage.NewFXImageOutput()
	outputNode.SetInput(bloomNode)

	err = outputNode.Process(ctx)
	if err != nil {
		panic(err)
	}

	err = outputNode.Save("output_extended.png")
	if err != nil {
		panic(err)
	}

	// Edge Detection Demo
	edgeNode, _ := fxartistic.NewFXEdgeDetectionNode(ctx, width, height)
	defer edgeNode.Release()
	edgeNode.SetInput("u_texture", inputNode) // Use original input
	edgeNode.SetThreshold(0.05)

	edgeOutput := fximage.NewFXImageOutput()
	edgeOutput.SetInput(edgeNode)

	err = edgeOutput.Process(ctx)
	if err != nil {
		panic(err)
	}

	err = edgeOutput.Save("output_edge.png")
	if err != nil {
		panic(err)
	}

	fmt.Println("Extended effects applied and saved to output_extended.png and output_edge.png")
}
