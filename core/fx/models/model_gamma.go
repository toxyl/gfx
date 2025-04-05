package models

import (
	"image"
	stdcolor "image/color"
	"image/draw"
	"math"

	"github.com/toxyl/gfx/core/color"
	"github.com/toxyl/gfx/core/fx"
)

// GammaFunction represents a function that applies gamma correction to an image
type GammaFunction struct {
	*fx.BaseFunction
	gamma float64
}

// NewGamma creates a new gamma function
func NewGamma(gamma float64) *GammaFunction {
	return &GammaFunction{
		BaseFunction: fx.NewBaseFunction("gamma", "Applies gamma correction to an image", color.New(0, 0, 0, 1), nil),
		gamma:        gamma,
	}
}

// Apply implements the Function interface
func (f *GammaFunction) Apply(img *image.Image) (*image.Image, error) {
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
func (f *GammaFunction) ProcessPixel(x, y int, img *image.Image) *color.Color64 {
	col := color.New(0, 0, 0, 1)
	col.SetUint8((*img).At(x, y).(stdcolor.RGBA))

	// Apply gamma correction
	channelMeta := color.NewChannel("gamma", 0.1, 5.0, "", "Gamma correction value")

	// Apply gamma correction to each channel
	col.R = color.ClampChannelValue(math.Pow(col.R, 1.0/f.gamma), channelMeta)
	col.G = color.ClampChannelValue(math.Pow(col.G, 1.0/f.gamma), channelMeta)
	col.B = color.ClampChannelValue(math.Pow(col.B, 1.0/f.gamma), channelMeta)

	return col
}
