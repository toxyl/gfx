package models

import (
	"image"
	"image/color"

	"github.com/toxyl/gfx/core/fx"
	"github.com/toxyl/gfx/core/meta"
	"github.com/toxyl/math"
)

// Invert represents an invert effect.
type Invert struct {
	Amount float64 // Invert amount (0.0 to 1.0)
	meta   *fx.EffectMeta
}

// NewInvertEffect creates a new invert effect.
func NewInvertEffect(amount float64) *Invert {
	i := &Invert{
		Amount: amount,
		meta: fx.NewEffectMeta(
			"Invert",
			"Inverts the colors of an image",
			meta.NewChannelMeta("Amount", 0.0, 1.0, "", "Invert amount (0.0 to 1.0)"),
		),
	}
	i.Amount = fx.ClampParameter(amount, i.meta.Parameters[0])
	return i
}

// Apply applies the invert effect to an image.
func (i *Invert) Apply(img image.Image) (image.Image, error) {
	bounds := img.Bounds()
	dst := image.NewRGBA(bounds)

	for y := bounds.Min.Y; y < bounds.Max.Y; y++ {
		for x := bounds.Min.X; x < bounds.Max.X; x++ {
			r, g, b, a := img.At(x, y).RGBA()

			// Convert to float64 for calculations
			rf := float64(r) / 0xFFFF
			gf := float64(g) / 0xFFFF
			bf := float64(b) / 0xFFFF

			// Invert colors
			rf = 1.0 - rf
			gf = 1.0 - gf
			bf = 1.0 - bf

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

	return dst, nil
}

// Meta returns the effect metadata.
func (i *Invert) Meta() *fx.EffectMeta {
	return i.meta
}
