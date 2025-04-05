package models

import (
	"image"
	"image/color"

	"github.com/toxyl/gfx/core/fx"
	"github.com/toxyl/gfx/core/meta"
	"github.com/toxyl/math"
)

// LuminanceContrast represents a luminance-based contrast adjustment effect.
type LuminanceContrast struct {
	Amount float64 // Contrast amount (-1.0 to 1.0)
	meta   *fx.EffectMeta
}

// NewLuminanceContrastEffect creates a new luminance-contrast effect.
func NewLuminanceContrastEffect(amount float64) *LuminanceContrast {
	lc := &LuminanceContrast{
		Amount: amount,
		meta: fx.NewEffectMeta(
			"LuminanceContrast",
			"Adjusts contrast based on luminance values",
			meta.NewChannelMeta("Amount", -1.0, 1.0, "", "Contrast adjustment amount (-1.0 to 1.0)"),
		),
	}
	lc.Amount = fx.ClampParameter(amount, lc.meta.Parameters[0])
	return lc
}

// Apply applies the luminance-contrast effect to an image.
func (lc *LuminanceContrast) Apply(img image.Image) image.Image {
	bounds := img.Bounds()
	dst := image.NewRGBA(bounds)

	// Calculate average luminance
	var avgLum float64
	count := 0
	for y := bounds.Min.Y; y < bounds.Max.Y; y++ {
		for x := bounds.Min.X; x < bounds.Max.X; x++ {
			r, g, b, _ := img.At(x, y).RGBA()
			rf := float64(r) / 0xFFFF
			gf := float64(g) / 0xFFFF
			bf := float64(b) / 0xFFFF
			avgLum += 0.299*rf + 0.587*gf + 0.114*bf
			count++
		}
	}
	avgLum /= float64(count)

	// Apply contrast adjustment
	for y := bounds.Min.Y; y < bounds.Max.Y; y++ {
		for x := bounds.Min.X; x < bounds.Max.X; x++ {
			r, g, b, a := img.At(x, y).RGBA()

			// Convert to float64 and normalize
			rf := float64(r) / 0xFFFF
			gf := float64(g) / 0xFFFF
			bf := float64(b) / 0xFFFF

			// Calculate luminance
			lum := 0.299*rf + 0.587*gf + 0.114*bf

			// Calculate contrast factor
			factor := 1.0 + lc.Amount
			if lc.Amount < 0 {
				factor = 1.0 / (1.0 - lc.Amount)
			}

			// Adjust luminance
			newLum := avgLum + (lum-avgLum)*factor

			// Scale RGB values to maintain relative ratios
			scale := newLum / lum
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
func (lc *LuminanceContrast) Meta() *fx.EffectMeta {
	return lc.meta
}
