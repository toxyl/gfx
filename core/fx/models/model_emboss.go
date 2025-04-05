package models

import (
	"image"
	"image/color"

	"github.com/toxyl/gfx/core/fx"
	"github.com/toxyl/gfx/core/meta"
	"github.com/toxyl/math"
)

// Emboss represents an emboss effect.
type Emboss struct {
	Amount float64 // Emboss amount (0.0 to 1.0)
	meta   *fx.EffectMeta
}

// NewEmbossEffect creates a new emboss effect.
func NewEmbossEffect(amount float64) *Emboss {
	e := &Emboss{
		Amount: amount,
		meta: fx.NewEffectMeta(
			"Emboss",
			"Applies an emboss effect to an image",
			meta.NewChannelMeta("Amount", 0.0, 1.0, "", "Emboss amount (0.0 to 1.0)"),
		),
	}
	e.Amount = fx.ClampParameter(amount, e.meta.Parameters[0])
	return e
}

// Apply applies the emboss effect to an image.
func (e *Emboss) Apply(img image.Image) image.Image {
	bounds := img.Bounds()

	// Create new image with same dimensions
	dst := image.NewRGBA(bounds)

	// Emboss kernel
	kernel := [3][3]float64{
		{-1, -1, 0},
		{-1, 1, 1},
		{0, 1, 1},
	}

	for y := bounds.Min.Y + 1; y < bounds.Max.Y-1; y++ {
		for x := bounds.Min.X + 1; x < bounds.Max.X-1; x++ {
			var r, g, b float64

			// Apply kernel
			for ky := -1; ky <= 1; ky++ {
				for kx := -1; kx <= 1; kx++ {
					px := x + kx
					py := y + ky
					pr, pg, pb, _ := img.At(px, py).RGBA()

					// Convert to float64 and normalize
					prf := float64(pr) / 0xFFFF
					pgf := float64(pg) / 0xFFFF
					pbf := float64(pb) / 0xFFFF

					// Apply kernel weight
					weight := kernel[ky+1][kx+1]
					r += prf * weight
					g += pgf * weight
					b += pbf * weight
				}
			}

			// Normalize and apply amount
			r = (r + 1.0) / 2.0
			g = (g + 1.0) / 2.0
			b = (b + 1.0) / 2.0

			// Blend with original
			or, og, ob, oa := img.At(x, y).RGBA()
			orf := float64(or) / 0xFFFF
			ogf := float64(og) / 0xFFFF
			obf := float64(ob) / 0xFFFF

			r = orf + (r-orf)*e.Amount
			g = ogf + (g-ogf)*e.Amount
			b = obf + (b-obf)*e.Amount

			// Convert back to uint32
			ri := uint32(math.Max(0, math.Min(0xFFFF, r*0xFFFF)))
			gi := uint32(math.Max(0, math.Min(0xFFFF, g*0xFFFF)))
			bi := uint32(math.Max(0, math.Min(0xFFFF, b*0xFFFF)))

			dst.Set(x, y, color.RGBA64{
				R: uint16(ri),
				G: uint16(gi),
				B: uint16(bi),
				A: uint16(oa),
			})
		}
	}

	return dst
}

// Meta returns the effect metadata.
func (e *Emboss) Meta() *fx.EffectMeta {
	return e.meta
}
