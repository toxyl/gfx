package models

import (
	"image"
	stdcolor "image/color"
	"image/draw"

	"github.com/toxyl/gfx/core/color"
	"github.com/toxyl/gfx/core/fx"
)

// HueFunction represents a function that adjusts the hue of an image
type HueFunction struct {
	*fx.BaseFunction
	amount float64
}

// NewHue creates a new hue adjustment function
func NewHue(amount float64) *HueFunction {
	return &HueFunction{
		BaseFunction: fx.NewBaseFunction("hue", "Adjusts image hue", color.New(0, 0, 0, 1), nil),
		amount:       amount,
	}
}

// Apply implements the Function interface
func (f *HueFunction) Apply(img *image.Image) (*image.Image, error) {
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
func (f *HueFunction) ProcessPixel(x, y int, img *image.Image) *color.Color64 {
	col := color.New(0, 0, 0, 1)
	col.SetUint8((*img).At(x, y).(stdcolor.RGBA))

	// Convert to HSL for hue adjustment
	stdRgb := col.ToUint16()
	rgb, _ := color.NewRGBA64(float64(stdRgb.R), float64(stdRgb.G), float64(stdRgb.B), float64(stdRgb.A))
	hsl := color.HSLFromRGB(rgb)

	// Adjust hue
	hueMeta := color.NewChannel("hue", 0, 360, "degrees", "Hue angle")
	hsl.H = color.ClampChannelValue(hsl.H+f.amount, hueMeta)

	// Convert back to RGB
	rgb = hsl.ToRGB()
	return color.New(rgb.R, rgb.G, rgb.B, rgb.A)
}
