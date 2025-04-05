package models

import (
	"image"
	stdcolor "image/color"
	"image/draw"

	"github.com/toxyl/gfx/core/color"
	"github.com/toxyl/gfx/core/fx"
)

// GrayFunction represents a function that converts an image to grayscale
type GrayFunction struct {
	*fx.BaseFunction
}

// NewGray creates a new grayscale function
func NewGray() *GrayFunction {
	return &GrayFunction{
		BaseFunction: fx.NewBaseFunction("gray", "Converts image to grayscale", color.New(0, 0, 0, 1), nil),
	}
}

// Apply implements the Function interface
func (f *GrayFunction) Apply(img *image.Image) (*image.Image, error) {
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
func (f *GrayFunction) ProcessPixel(x, y int, img *image.Image) *color.Color64 {
	col := color.New(0, 0, 0, 1)
	col.SetUint8((*img).At(x, y).(stdcolor.RGBA))

	// Convert to grayscale using luminance formula
	gray := 0.299*col.R + 0.587*col.G + 0.114*col.B

	// Set all channels to the same gray value
	col.R = gray
	col.G = gray
	col.B = gray

	return col
}
