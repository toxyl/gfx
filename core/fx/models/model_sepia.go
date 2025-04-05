package models

import (
	"image"
	"image/draw"

	"github.com/toxyl/gfx/core/color"
	"github.com/toxyl/gfx/core/fx"
)

// SepiaFunction represents a function that applies a sepia tone to an image
type SepiaFunction struct {
	*fx.BaseFunction
}

// NewSepia creates a new sepia tone function
func NewSepia() *SepiaFunction {
	return &SepiaFunction{
		BaseFunction: fx.NewBaseFunction("sepia", "Applies sepia tone to image", color.New(0, 0, 0, 1), nil),
	}
}

// Apply implements the Function interface
func (f *SepiaFunction) Apply(img *image.Image) (*image.Image, error) {
	bounds := (*img).Bounds()
	dst := image.NewRGBA(bounds)
	draw.Draw(dst, bounds, *img, bounds.Min, draw.Src)

	for y := bounds.Min.Y; y < bounds.Max.Y; y++ {
		for x := bounds.Min.X; x < bounds.Max.X; x++ {
			col := f.ProcessRGBA(x, y, dst)
			dst.Set(x, y, col.ToUint8())
		}
	}
	*img = dst
	return img, nil
}

// ProcessRGBA implements the Function interface
func (f *SepiaFunction) ProcessRGBA(x, y int, img *image.RGBA) *color.Color64 {
	col := color.New(0, 0, 0, 1)
	col.SetUint8(img.RGBAAt(x, y))

	// Convert to grayscale first
	gray := 0.299*col.R + 0.587*col.G + 0.114*col.B

	// Apply sepia tone
	col.R = min(gray*1.2, 1.0) // Red channel
	col.G = min(gray*0.9, 1.0) // Green channel
	col.B = min(gray*0.6, 1.0) // Blue channel

	return col
}

func min(a, b float64) float64 {
	if a < b {
		return a
	}
	return b
}
