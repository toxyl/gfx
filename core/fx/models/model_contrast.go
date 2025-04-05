package models

import (
	"image"
	"image/color"

	"github.com/toxyl/gfx/core/fx"
	"github.com/toxyl/gfx/core/meta"
	"github.com/toxyl/math"
)

// Contrast represents a contrast adjustment effect.
type Contrast struct {
	Amount float64 // Amount of contrast adjustment (-1.0 to 1.0)
	meta   *fx.EffectMeta
}

// NewContrast creates a new contrast effect.
func NewContrast(amount float64) *Contrast {
	c := &Contrast{
		Amount: amount,
		meta: fx.NewEffectMeta(
			"Contrast",
			"Adjusts the contrast of an image",
			meta.NewChannelMeta("Amount", -1.0, 1.0, "", "Amount of contrast adjustment (-1.0 to 1.0)"),
		),
	}
	c.Amount = fx.ClampParameter(amount, c.meta.Parameters[0])
	return c
}

// Apply applies the contrast effect to an image.
func (c *Contrast) Apply(img image.Image) image.Image {
	bounds := img.Bounds()
	dst := image.NewRGBA(bounds)

	// Calculate contrast factor
	factor := (259.0 * (c.Amount + 255.0)) / (255.0 * (259.0 - c.Amount))

	for y := bounds.Min.Y; y < bounds.Max.Y; y++ {
		for x := bounds.Min.X; x < bounds.Max.X; x++ {
			r, g, blue, a := img.At(x, y).RGBA()

			// Convert to float64 for calculations
			rf := float64(r) / 0xFFFF
			gf := float64(g) / 0xFFFF
			bf := float64(blue) / 0xFFFF

			// Apply contrast adjustment
			rf = factor*(rf-0.5) + 0.5
			gf = factor*(gf-0.5) + 0.5
			bf = factor*(bf-0.5) + 0.5

			// Convert back to uint32
			r = uint32(math.Max(0, math.Min(0xFFFF, rf*0xFFFF)))
			g = uint32(math.Max(0, math.Min(0xFFFF, gf*0xFFFF)))
			blue = uint32(math.Max(0, math.Min(0xFFFF, bf*0xFFFF)))

			dst.Set(x, y, color.RGBA64{
				R: uint16(r),
				G: uint16(g),
				B: uint16(blue),
				A: uint16(a),
			})
		}
	}

	return dst
}

// Meta returns the effect metadata.
func (c *Contrast) Meta() *fx.EffectMeta {
	return c.meta
}
