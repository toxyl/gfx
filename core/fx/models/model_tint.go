package models

import (
	"image"
	"image/color"

	"github.com/toxyl/gfx/core/fx"
	"github.com/toxyl/gfx/core/meta"
	"github.com/toxyl/math"
)

// Tint represents a tint effect.
type Tint struct {
	Color  color.RGBA64 // Tint color
	Amount float64      // Tint amount (0.0 to 1.0)
	meta   *fx.EffectMeta
}

// NewTintEffect creates a new tint effect.
func NewTintEffect(tintColor color.RGBA64, amount float64) *Tint {
	t := &Tint{
		Color:  tintColor,
		Amount: amount,
		meta: fx.NewEffectMeta(
			"Tint",
			"Applies a color tint to an image",
			meta.NewChannelMeta("Amount", 0.0, 1.0, "", "Tint amount (0.0 to 1.0)"),
		),
	}
	t.Amount = fx.ClampParameter(amount, t.meta.Parameters[0])
	return t
}

// Apply applies the tint effect to an image.
func (t *Tint) Apply(img image.Image) image.Image {
	bounds := img.Bounds()
	dst := image.NewRGBA(bounds)

	// Convert tint color to float64
	tr := float64(t.Color.R) / 0xFFFF
	tg := float64(t.Color.G) / 0xFFFF
	tb := float64(t.Color.B) / 0xFFFF

	for y := bounds.Min.Y; y < bounds.Max.Y; y++ {
		for x := bounds.Min.X; x < bounds.Max.X; x++ {
			r, g, b, a := img.At(x, y).RGBA()

			// Convert to float64 for calculations
			rf := float64(r) / 0xFFFF
			gf := float64(g) / 0xFFFF
			bf := float64(b) / 0xFFFF

			// Calculate luminance
			luminance := 0.299*rf + 0.587*gf + 0.114*bf

			// Apply tint
			rf = luminance + (tr-luminance)*t.Amount
			gf = luminance + (tg-luminance)*t.Amount
			bf = luminance + (tb-luminance)*t.Amount

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
func (t *Tint) Meta() *fx.EffectMeta {
	return t.meta
}
