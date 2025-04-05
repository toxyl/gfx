package models

import (
	"image"
	stdcolor "image/color"
	"image/draw"

	"github.com/toxyl/gfx/core/color"
	"github.com/toxyl/gfx/core/fx"
)

// ContrastFunction represents a function that adjusts the contrast of an image
type ContrastFunction struct {
	*fx.BaseFunction
	amount float64
}

// NewContrast creates a new contrast adjustment function
func NewContrast(amount float64) *ContrastFunction {
	return &ContrastFunction{
		BaseFunction: fx.NewBaseFunction("contrast", "Applies contrast transformation to an image", color.New(0, 0, 0, 1), nil),
		amount:       amount,
	}
}

// Apply implements the Function interface
func (f *ContrastFunction) Apply(img *image.Image) (*image.Image, error) {
	bounds := (*img).Bounds()
	dst := image.NewRGBA(bounds)
	draw.Draw(dst, bounds, *img, bounds.Min, draw.Src)

	for y := bounds.Min.Y; y < bounds.Max.Y; y++ {
		for x := bounds.Min.X; x < bounds.Max.X; x++ {
			col := f.ProcessPixel(x, y, img)
			dst.Set(x, y, col.ToUint8())
		}
	}
	*img = dst
	return img, nil
}

// ProcessPixel implements the ImageFunction interface
func (f *ContrastFunction) ProcessPixel(x, y int, img *image.Image) *color.Color64 {
	col := color.New(0, 0, 0, 1)
	col.SetUint8((*img).At(x, y).(stdcolor.RGBA))

	// Calculate contrast factor
	factor := (259.0 * (f.amount + 255.0)) / (255.0 * (259.0 - f.amount))

	// Adjust contrast for each channel
	contrastMeta := color.NewChannel("contrast", -255.0, 255.0, "", "Contrast adjustment value")
	col.R = color.ClampChannelValue((col.R-0.5)*factor+0.5, contrastMeta)
	col.G = color.ClampChannelValue((col.G-0.5)*factor+0.5, contrastMeta)
	col.B = color.ClampChannelValue((col.B-0.5)*factor+0.5, contrastMeta)

	return col
}
