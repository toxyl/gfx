package models

import (
	"image"
	"image/color"

	"github.com/toxyl/gfx/core/fx"
	"github.com/toxyl/gfx/core/meta"
	"github.com/toxyl/math"
)

// Brightness represents a brightness adjustment effect.
type Brightness struct {
	Amount float64 // Amount of brightness adjustment (-1.0 to 1.0)
	meta   *fx.EffectMeta
}

// NewBrightness creates a new brightness effect.
func NewBrightness(amount float64) *Brightness {
	b := &Brightness{
		Amount: amount,
		meta: fx.NewEffectMeta(
			"Brightness",
			"Adjusts the brightness of an image",
			meta.NewChannelMeta("Amount", -1.0, 1.0, "", "Amount of brightness adjustment (-1.0 to 1.0)"),
		),
	}
	b.Amount = fx.ClampParameter(amount, b.meta.Parameters[0])
	return b
}

// Apply applies the brightness effect to an image.
func (b *Brightness) Apply(img image.Image) image.Image {
	bounds := img.Bounds()
	dst := image.NewRGBA(bounds)

	for y := bounds.Min.Y; y < bounds.Max.Y; y++ {
		for x := bounds.Min.X; x < bounds.Max.X; x++ {
			r, g, blue, a := img.At(x, y).RGBA()

			// Convert to float64 for calculations
			rf := float64(r) / 0xFFFF
			gf := float64(g) / 0xFFFF
			bf := float64(blue) / 0xFFFF

			// Apply brightness adjustment
			if b.Amount > 0 {
				rf = rf + (1.0-rf)*b.Amount
				gf = gf + (1.0-gf)*b.Amount
				bf = bf + (1.0-bf)*b.Amount
			} else {
				rf = rf * (1.0 + b.Amount)
				gf = gf * (1.0 + b.Amount)
				bf = bf * (1.0 + b.Amount)
			}

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
func (b *Brightness) Meta() *fx.EffectMeta {
	return b.meta
}
