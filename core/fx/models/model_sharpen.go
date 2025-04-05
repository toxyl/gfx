package models

import (
	"image"
	stdcolor "image/color"
	"image/draw"

	"github.com/toxyl/gfx/core/color"
	"github.com/toxyl/gfx/core/fx"
)

// SharpenFunction represents a function that applies a sharpening effect to an image
type SharpenFunction struct {
	*fx.BaseFunction
	amount float64
}

// NewSharpen creates a new sharpen function
func NewSharpen(amount float64) *SharpenFunction {
	return &SharpenFunction{
		BaseFunction: fx.NewBaseFunction("sharpen", "Sharpen", color.New(0, 0, 0, 1), nil),
		amount:       amount,
	}
}

// Apply implements the Function interface
func (f *SharpenFunction) Apply(img *image.Image) (*image.Image, error) {
	bounds := (*img).Bounds()
	dst := image.NewRGBA(bounds)
	draw.Draw(dst, bounds, *img, bounds.Min, draw.Src)

	// Create temporary image for the unsharp mask
	temp := image.NewRGBA(bounds)
	draw.Draw(temp, bounds, *img, bounds.Min, draw.Src)

	// Apply unsharp mask
	for y := bounds.Min.Y; y < bounds.Max.Y; y++ {
		for x := bounds.Min.X; x < bounds.Max.X; x++ {
			col := f.applyUnsharpMask(x, y, img)
			dst.Set(x, y, col.ToUint8())
		}
	}

	*img = dst
	return img, nil
}

// applyUnsharpMask applies the unsharp mask algorithm
func (f *SharpenFunction) applyUnsharpMask(x, y int, img *image.Image) *color.Color64 {
	bounds := (*img).Bounds()
	original := color.New(0, 0, 0, 1)
	original.SetUint8((*img).At(x, y).(stdcolor.RGBA))

	// Calculate average of surrounding pixels
	var sumR, sumG, sumB float64
	count := 0

	for dy := -1; dy <= 1; dy++ {
		for dx := -1; dx <= 1; dx++ {
			px := x + dx
			py := y + dy

			if px < bounds.Min.X || px >= bounds.Max.X || py < bounds.Min.Y || py >= bounds.Max.Y {
				continue
			}

			col := color.New(0, 0, 0, 1)
			col.SetUint8((*img).At(px, py).(stdcolor.RGBA))

			sumR += col.R
			sumG += col.G
			sumB += col.B
			count++
		}
	}

	if count == 0 {
		return original
	}

	// Calculate average
	avgR := sumR / float64(count)
	avgG := sumG / float64(count)
	avgB := sumB / float64(count)

	// Apply unsharp mask formula
	r := original.R + f.amount*(original.R-avgR)
	g := original.G + f.amount*(original.G-avgG)
	b := original.B + f.amount*(original.B-avgB)

	// Clamp values
	r = clamp(r)
	g = clamp(g)
	b = clamp(b)

	return color.New(r, g, b, original.A)
}

// clamp ensures a value is between 0 and 1
func clamp(v float64) float64 {
	if v < 0 {
		return 0
	}
	if v > 1 {
		return 1
	}
	return v
}
