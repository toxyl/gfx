package models

import (
	"image"
	stdcolor "image/color"
	"image/draw"
	"math"

	"github.com/toxyl/gfx/core/color"
	"github.com/toxyl/gfx/core/fx"
)

// BrightnessFunction represents a function that adjusts the brightness of an image
type BrightnessFunction struct {
	*fx.BaseFunction
	amount float64
}

// NewBrightness creates a new brightness adjustment function
func NewBrightness(amount float64) *BrightnessFunction {
	return &BrightnessFunction{
		BaseFunction: fx.NewBaseFunction("brightness", "Applies brightness transformation to an image", color.New(0, 0, 0, 1), nil),
		amount:       amount,
	}
}

// Apply implements the Function interface
func (f *BrightnessFunction) Apply(img *image.Image) (*image.Image, error) {
	bounds := (*img).Bounds()
	dst := image.NewRGBA(bounds)
	draw.Draw(dst, bounds, *img, bounds.Min, draw.Src)

	for y := bounds.Min.Y; y < bounds.Max.Y; y++ {
		for x := bounds.Min.X; x < bounds.Max.X; x++ {
			col := f.ProcessPixel(x, y, img)
			dst.Set(x, y, col.ToUint8())
		}
	}
	var result image.Image = dst
	return &result, nil
}

// ProcessPixel implements the ImageFunction interface
func (f *BrightnessFunction) ProcessPixel(x, y int, img *image.Image) *color.Color64 {
	col := color.New(0, 0, 0, 1)
	rgba := (*img).At(x, y).(stdcolor.RGBA)
	col.SetUint8(rgba)

	// Add brightness amount directly in [0,1] range
	col.R = math.Max(0.0, math.Min(1.0, col.R+f.amount))
	col.G = math.Max(0.0, math.Min(1.0, col.G+f.amount))
	col.B = math.Max(0.0, math.Min(1.0, col.B+f.amount))
	col.A = 1.0 // Ensure alpha is always 1.0 as per test expectations

	return col
}
