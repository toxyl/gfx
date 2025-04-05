package models

import (
	"image"
	"image/color"

	"github.com/toxyl/gfx/core/fx"
	"github.com/toxyl/gfx/core/meta"
	"github.com/toxyl/math"
)

// Sepia represents a sepia effect.
type Sepia struct {
	Amount float64 // Sepia amount (0.0 to 1.0)
	meta   *fx.EffectMeta
}

// NewSepiaEffect creates a new sepia effect.
func NewSepiaEffect(amount float64) *Sepia {
	s := &Sepia{
		Amount: amount,
		meta: fx.NewEffectMeta(
			"Sepia",
			"Applies a sepia tone to an image",
			meta.NewChannelMeta("Amount", 0.0, 1.0, "", "Sepia amount (0.0 to 1.0)"),
		),
	}
	s.Amount = fx.ClampParameter(amount, s.meta.Parameters[0])
	return s
}

// Apply applies the sepia effect to an image.
func (s *Sepia) Apply(img image.Image) (image.Image, error) {
	bounds := img.Bounds()
	dst := image.NewRGBA(bounds)

	for y := bounds.Min.Y; y < bounds.Max.Y; y++ {
		for x := bounds.Min.X; x < bounds.Max.X; x++ {
			r, g, b, a := img.At(x, y).RGBA()

			// Convert to float64 for calculations
			rf := float64(r) / 0xFFFF
			gf := float64(g) / 0xFFFF
			bf := float64(b) / 0xFFFF

			// Calculate luminance
			luminance := 0.299*rf + 0.587*gf + 0.114*bf

			// Apply sepia tone
			rf = luminance + (0.393*rf-luminance)*s.Amount
			gf = luminance + (0.769*gf-luminance)*s.Amount
			bf = luminance + (0.189*bf-luminance)*s.Amount

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
func (s *Sepia) Meta() *fx.EffectMeta {
	return s.meta
}
