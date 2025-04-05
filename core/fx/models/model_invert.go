package models

import (
	"image"
	stdcolor "image/color"
	"image/draw"

	"github.com/toxyl/gfx/core/color"
	"github.com/toxyl/gfx/core/fx"
)

// InvertFunction represents a function that inverts the colors of an image
type InvertFunction struct {
	*fx.BaseFunction
}

// NewInvert creates a new color inversion function
func NewInvert() *InvertFunction {
	return &InvertFunction{
		BaseFunction: fx.NewBaseFunction("invert", "Inverts image colors", color.New(0, 0, 0, 1), nil),
	}
}

// Apply implements the Function interface
func (f *InvertFunction) Apply(img *image.Image) (*image.Image, error) {
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
func (f *InvertFunction) ProcessPixel(x, y int, img *image.Image) *color.Color64 {
	col := color.New(0, 0, 0, 1)
	col.SetUint8((*img).At(x, y).(stdcolor.RGBA))

	// Invert each color channel
	col.R = 1.0 - col.R
	col.G = 1.0 - col.G
	col.B = 1.0 - col.B

	return col
}
