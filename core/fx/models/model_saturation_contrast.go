package models

import (
	"image"
	"image/color"

	"github.com/toxyl/gfx/core/fx"
	"github.com/toxyl/gfx/core/meta"
	"github.com/toxyl/math"
)

// SaturationContrast represents a saturation-based contrast adjustment effect.
type SaturationContrast struct {
	Amount float64 // Contrast amount (-1.0 to 1.0)
	meta   *fx.EffectMeta
}

// NewSaturationContrastEffect creates a new saturation-contrast effect.
func NewSaturationContrastEffect(amount float64) *SaturationContrast {
	sc := &SaturationContrast{
		Amount: amount,
		meta: fx.NewEffectMeta(
			"SaturationContrast",
			"Adjusts contrast based on saturation values",
			meta.NewChannelMeta("Amount", -1.0, 1.0, "", "Contrast adjustment amount (-1.0 to 1.0)"),
		),
	}
	sc.Amount = fx.ClampParameter(amount, sc.meta.Parameters[0])
	return sc
}

// Apply applies the saturation-contrast effect to an image.
func (sc *SaturationContrast) Apply(img image.Image) image.Image {
	bounds := img.Bounds()
	dst := image.NewRGBA(bounds)

	// Calculate average saturation
	var avgSat float64
	count := 0
	for y := bounds.Min.Y; y < bounds.Max.Y; y++ {
		for x := bounds.Min.X; x < bounds.Max.X; x++ {
			r, g, b, _ := img.At(x, y).RGBA()
			rf := float64(r) / 0xFFFF
			gf := float64(g) / 0xFFFF
			bf := float64(b) / 0xFFFF

			// Calculate saturation
			max := math.Max(math.Max(rf, gf), bf)
			min := math.Min(math.Min(rf, gf), bf)
			if max != min {
				avgSat += (max - min) / max
			}
			count++
		}
	}
	avgSat /= float64(count)

	// Apply contrast adjustment
	for y := bounds.Min.Y; y < bounds.Max.Y; y++ {
		for x := bounds.Min.X; x < bounds.Max.X; x++ {
			r, g, b, a := img.At(x, y).RGBA()

			// Convert to float64 and normalize
			rf := float64(r) / 0xFFFF
			gf := float64(g) / 0xFFFF
			bf := float64(b) / 0xFFFF

			// Calculate saturation
			max := math.Max(math.Max(rf, gf), bf)
			min := math.Min(math.Min(rf, gf), bf)
			sat := 0.0
			if max != min {
				sat = (max - min) / max
			}

			// Calculate contrast factor
			factor := 1.0 + sc.Amount
			if sc.Amount < 0 {
				factor = 1.0 / (1.0 - sc.Amount)
			}

			// Adjust saturation
			newSat := avgSat + (sat-avgSat)*factor

			// Convert back to RGB
			if sat == 0 {
				// Grayscale, no change needed
			} else {
				// Calculate new RGB values
				scale := newSat / sat
				rf = max - (max-rf)*scale
				gf = max - (max-gf)*scale
				bf = max - (max-bf)*scale
			}

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
func (sc *SaturationContrast) Meta() *fx.EffectMeta {
	return sc.meta
}
