package models

import (
	"image"
	stdcolor "image/color"
	"image/draw"
	"math"

	"github.com/toxyl/gfx/core/color"
	"github.com/toxyl/gfx/core/fx"
)

// EdgeFunction represents a function that detects edges in an image
type EdgeFunction struct {
	*fx.BaseFunction
	threshold float64
}

// NewEdge creates a new edge detection function
func NewEdge(threshold float64) *EdgeFunction {
	return &EdgeFunction{
		BaseFunction: fx.NewBaseFunction("edge", "Applies edge detection to an image", color.New(0, 0, 0, 1), nil),
		threshold:    threshold,
	}
}

// Apply implements the Function interface
func (f *EdgeFunction) Apply(img *image.Image) (*image.Image, error) {
	bounds := (*img).Bounds()
	dst := image.NewRGBA(bounds)
	draw.Draw(dst, bounds, *img, bounds.Min, draw.Src)

	// Apply Sobel operator
	for y := bounds.Min.Y; y < bounds.Max.Y; y++ {
		for x := bounds.Min.X; x < bounds.Max.X; x++ {
			col := f.applySobel(x, y, img)
			dst.Set(x, y, col.ToUint8())
		}
	}

	*img = dst
	return img, nil
}

// applySobel applies the Sobel operator for edge detection
func (f *EdgeFunction) applySobel(x, y int, img *image.Image) *color.Color64 {
	bounds := (*img).Bounds()

	// Sobel kernels
	kernelX := [3][3]float64{
		{-1, 0, 1},
		{-2, 0, 2},
		{-1, 0, 1},
	}
	kernelY := [3][3]float64{
		{-1, -2, -1},
		{0, 0, 0},
		{1, 2, 1},
	}

	var gxR, gxG, gxB float64
	var gyR, gyG, gyB float64

	// Apply Sobel kernels
	for dy := -1; dy <= 1; dy++ {
		for dx := -1; dx <= 1; dx++ {
			px := x + dx
			py := y + dy

			if px < bounds.Min.X || px >= bounds.Max.X || py < bounds.Min.Y || py >= bounds.Max.Y {
				continue
			}

			col := color.New(0, 0, 0, 1)
			col.SetUint8((*img).At(px, py).(stdcolor.RGBA))

			// Apply X kernel
			gxR += col.R * kernelX[dy+1][dx+1]
			gxG += col.G * kernelX[dy+1][dx+1]
			gxB += col.B * kernelX[dy+1][dx+1]

			// Apply Y kernel
			gyR += col.R * kernelY[dy+1][dx+1]
			gyG += col.G * kernelY[dy+1][dx+1]
			gyB += col.B * kernelY[dy+1][dx+1]
		}
	}

	// Calculate gradient magnitude
	magR := math.Sqrt(gxR*gxR + gyR*gyR)
	magG := math.Sqrt(gxG*gxG + gyG*gyG)
	magB := math.Sqrt(gxB*gxB + gyB*gyB)

	// Normalize and apply threshold
	magR = edgeClamp(magR / f.threshold)
	magG = edgeClamp(magG / f.threshold)
	magB = edgeClamp(magB / f.threshold)

	return color.New(magR, magG, magB, 1.0)
}

// edgeClamp ensures a value is between 0 and 1
func edgeClamp(v float64) float64 {
	if v < 0 {
		return 0
	}
	if v > 1 {
		return 1
	}
	return v
}
