package models

import (
	"image"
	"image/color"

	"github.com/toxyl/gfx/core/fx"
	"github.com/toxyl/gfx/core/meta"
	"github.com/toxyl/math"
)

// ColorBalance represents a color balance adjustment effect.
type ColorBalance struct {
	Red   float64 // Red channel adjustment (-1.0 to 1.0)
	Green float64 // Green channel adjustment (-1.0 to 1.0)
	Blue  float64 // Blue channel adjustment (-1.0 to 1.0)
	meta  *fx.EffectMeta
}

// NewColorBalanceEffect creates a new color balance adjustment effect.
func NewColorBalanceEffect(red, green, blue float64) *ColorBalance {
	cb := &ColorBalance{
		Red:   red,
		Green: green,
		Blue:  blue,
		meta: fx.NewEffectMeta(
			"Color Balance",
			"Adjusts the color balance of an image",
			meta.NewChannelMeta("Red", -1.0, 1.0, "", "Red channel adjustment (-1.0 to 1.0)"),
			meta.NewChannelMeta("Green", -1.0, 1.0, "", "Green channel adjustment (-1.0 to 1.0)"),
			meta.NewChannelMeta("Blue", -1.0, 1.0, "", "Blue channel adjustment (-1.0 to 1.0)"),
		),
	}
	cb.Red = fx.ClampParameter(red, cb.meta.Parameters[0])
	cb.Green = fx.ClampParameter(green, cb.meta.Parameters[1])
	cb.Blue = fx.ClampParameter(blue, cb.meta.Parameters[2])
	return cb
}

// Apply applies the color balance adjustment effect to an image.
func (cb *ColorBalance) Apply(img image.Image) image.Image {
	bounds := img.Bounds()
	dst := image.NewRGBA(bounds)

	for y := bounds.Min.Y; y < bounds.Max.Y; y++ {
		for x := bounds.Min.X; x < bounds.Max.X; x++ {
			r, g, b, a := img.At(x, y).RGBA()

			// Convert to float64 for calculations
			rf := float64(r) / 0xFFFF
			gf := float64(g) / 0xFFFF
			bf := float64(b) / 0xFFFF

			// Apply color balance adjustments
			rf = math.Max(0, math.Min(1, rf+cb.Red))
			gf = math.Max(0, math.Min(1, gf+cb.Green))
			bf = math.Max(0, math.Min(1, bf+cb.Blue))

			// Convert back to uint32
			r = uint32(math.Max(0, math.Min(0xFFFF, rf*0xFFFF)))
			g = uint32(math.Max(0, math.Min(0xFFFF, gf*0xFFFF)))
			b = uint32(math.Max(0, math.Min(0xFFFF, bf*0xFFFF)))

			dst.Set(x, y, color.RGBA64{
				R: uint16(r),
				G: uint16(g),
				B: uint16(b),
				A: uint16(a),
			})
		}
	}

	return dst
}

// Meta returns the effect metadata.
func (cb *ColorBalance) Meta() *fx.EffectMeta {
	return cb.meta
}
