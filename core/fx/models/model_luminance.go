package models

import (
	"image"
	"image/color"

	"github.com/toxyl/gfx/core/fx"
	"github.com/toxyl/gfx/core/meta"
	"github.com/toxyl/math"
)

// Luminance represents a luminance adjustment effect.
type Luminance struct {
	Amount float64 // Luminance amount (-1.0 to 1.0)
	meta   *fx.EffectMeta
}

// NewLuminanceEffect creates a new luminance effect.
func NewLuminanceEffect(amount float64) *Luminance {
	l := &Luminance{
		Amount: amount,
		meta: fx.NewEffectMeta(
			"Luminance",
			"Adjusts the luminance of an image",
			meta.NewChannelMeta("Amount", -1.0, 1.0, "", "Luminance adjustment amount (-1.0 to 1.0)"),
		),
	}
	l.Amount = fx.ClampParameter(amount, l.meta.Parameters[0])
	return l
}

// Apply applies the luminance effect to an image.
func (l *Luminance) Apply(img image.Image) image.Image {
	bounds := img.Bounds()
	dst := image.NewRGBA(bounds)

	for y := bounds.Min.Y; y < bounds.Max.Y; y++ {
		for x := bounds.Min.X; x < bounds.Max.X; x++ {
			r, g, b, a := img.At(x, y).RGBA()

			// Convert to float64 and normalize
			rf := float64(r) / 0xFFFF
			gf := float64(g) / 0xFFFF
			bf := float64(b) / 0xFFFF

			// Calculate luminance using standard weights
			lum := 0.299*rf + 0.587*gf + 0.114*bf

			// Adjust luminance
			if l.Amount > 0 {
				lum = lum + (1.0-lum)*l.Amount
			} else {
				lum = lum * (1.0 + l.Amount)
			}

			// Scale RGB values to maintain relative ratios
			scale := lum / (0.299*rf + 0.587*gf + 0.114*bf)
			rf *= scale
			gf *= scale
			bf *= scale

			// Clamp values
			rf = math.Max(0, math.Min(1, rf))
			gf = math.Max(0, math.Min(1, gf))
			bf = math.Max(0, math.Min(1, bf))

			dst.Set(x, y, color.RGBA64{
				R: uint16(rf * 0xFFFF),
				G: uint16(gf * 0xFFFF),
				B: uint16(bf * 0xFFFF),
				A: uint16(a),
			})
		}
	}

	return dst
}

// Meta returns the effect metadata.
func (l *Luminance) Meta() *fx.EffectMeta {
	return l.meta
}
