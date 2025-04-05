package models

import (
	"image"
	stdcolor "image/color"
	"image/draw"

	"github.com/toxyl/gfx/core/color"
	"github.com/toxyl/gfx/core/fx"
)

// SaturationFunction represents a function that adjusts the saturation of an image
type SaturationFunction struct {
	*fx.BaseFunction
	amount float64
}

// NewSaturation creates a new saturation adjustment function
func NewSaturation(amount float64) *SaturationFunction {
	return &SaturationFunction{
		BaseFunction: fx.NewBaseFunction("saturation", "Adjusts image saturation", color.New(0, 0, 0, 1), nil),
		amount:       amount,
	}
}

// Apply implements the Function interface
func (f *SaturationFunction) Apply(img *image.Image) (*image.Image, error) {
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
func (f *SaturationFunction) ProcessPixel(x, y int, img *image.Image) *color.Color64 {
	col := color.New(0, 0, 0, 1)
	col.SetUint8((*img).At(x, y).(stdcolor.RGBA))

	// Convert to HSL for saturation adjustment
	stdRgb := col.ToUint16()
	customRgb := &color.RGBA64{}
	customRgb.SetUint16(stdRgb)
	hsl := color.HSLFromRGB(customRgb)

	// Adjust saturation
	saturationMeta := color.NewChannel("saturation", 0, 1, "", "Saturation value")
	hsl.S = color.ClampChannelValue(hsl.S+f.amount, saturationMeta)

	// Convert back to RGB
	rgb := hsl.ToRGB()
	return color.New(rgb.R, rgb.G, rgb.B, rgb.A)
}
