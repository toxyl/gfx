package models

import (
	"image"
	"image/color"

	"github.com/toxyl/gfx/core/fx"
	"github.com/toxyl/gfx/core/meta"
	"github.com/toxyl/math"
)

// Pastelize represents a pastelize effect.
type Pastelize struct {
	Amount float64 // Pastelize amount (0.0 to 1.0)
	meta   *fx.EffectMeta
}

// NewPastelizeEffect creates a new pastelize effect.
func NewPastelizeEffect(amount float64) *Pastelize {
	p := &Pastelize{
		Amount: amount,
		meta: fx.NewEffectMeta(
			"Pastelize",
			"Applies a pastel effect to an image",
			meta.NewChannelMeta("Amount", 0.0, 1.0, "", "Pastelize amount (0.0 to 1.0)"),
		),
	}
	p.Amount = fx.ClampParameter(amount, p.meta.Parameters[0])
	return p
}

// Apply applies the pastelize effect to an image.
func (p *Pastelize) Apply(img image.Image) image.Image {
	bounds := img.Bounds()
	dst := image.NewRGBA(bounds)

	for y := bounds.Min.Y; y < bounds.Max.Y; y++ {
		for x := bounds.Min.X; x < bounds.Max.X; x++ {
			r, g, b, a := img.At(x, y).RGBA()

			// Convert to float64 and normalize
			rf := float64(r) / 0xFFFF
			gf := float64(g) / 0xFFFF
			bf := float64(b) / 0xFFFF

			// Calculate luminance
			lum := 0.299*rf + 0.587*gf + 0.114*bf

			// Calculate saturation
			max := math.Max(math.Max(rf, gf), bf)
			min := math.Min(math.Min(rf, gf), bf)
			sat := 0.0
			if max != min {
				sat = (max - min) / max
			}

			// Reduce saturation and increase luminance
			newSat := sat * (1.0 - p.Amount)
			newLum := lum + (1.0-lum)*p.Amount*0.5

			// Convert back to RGB
			if sat == 0 {
				// Grayscale
				rf = newLum
				gf = newLum
				bf = newLum
			} else {
				// Calculate new RGB values
				scale := newSat / sat
				rf = max - (max-rf)*scale
				gf = max - (max-gf)*scale
				bf = max - (max-bf)*scale

				// Adjust luminance
				currLum := 0.299*rf + 0.587*gf + 0.114*bf
				scale = newLum / currLum
				rf *= scale
				gf *= scale
				bf *= scale
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
func (p *Pastelize) Meta() *fx.EffectMeta {
	return p.meta
}
