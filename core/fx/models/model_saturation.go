package models

import (
	"image"
	"image/color"

	"github.com/toxyl/gfx/core/fx"
	"github.com/toxyl/gfx/core/meta"
	"github.com/toxyl/math"
)

// Saturation represents a saturation adjustment effect.
type Saturation struct {
	Amount float64 // Saturation adjustment value (-1.0 to 1.0)
	meta   *fx.EffectMeta
}

// NewSaturationEffect creates a new saturation adjustment effect.
func NewSaturationEffect(amount float64) *Saturation {
	s := &Saturation{
		Amount: amount,
		meta: fx.NewEffectMeta(
			"Saturation",
			"Adjusts the saturation of an image",
			meta.NewChannelMeta("Amount", -1.0, 1.0, "", "Saturation adjustment value (-1.0 to 1.0)"),
		),
	}
	s.Amount = fx.ClampParameter(amount, s.meta.Parameters[0])
	return s
}

// Apply applies the saturation adjustment effect to an image.
func (s *Saturation) Apply(img image.Image) image.Image {
	bounds := img.Bounds()
	dst := image.NewRGBA(bounds)

	for y := bounds.Min.Y; y < bounds.Max.Y; y++ {
		for x := bounds.Min.X; x < bounds.Max.X; x++ {
			r, green, blue, a := img.At(x, y).RGBA()

			// Convert to float64 for calculations
			rf := float64(r) / 0xFFFF
			gf := float64(green) / 0xFFFF
			bf := float64(blue) / 0xFFFF

			// Calculate grayscale value (luminance)
			gray := 0.299*rf + 0.587*gf + 0.114*bf

			// Apply saturation adjustment
			rf = gray + (rf-gray)*(1.0+s.Amount)
			gf = gray + (gf-gray)*(1.0+s.Amount)
			bf = gray + (bf-gray)*(1.0+s.Amount)

			// Convert back to uint32
			r = uint32(math.Max(0, math.Min(0xFFFF, rf*0xFFFF)))
			green = uint32(math.Max(0, math.Min(0xFFFF, gf*0xFFFF)))
			blue = uint32(math.Max(0, math.Min(0xFFFF, bf*0xFFFF)))

			dst.Set(x, y, color.RGBA64{
				R: uint16(r),
				G: uint16(green),
				B: uint16(blue),
				A: uint16(a),
			})
		}
	}

	return dst
}

// Meta returns the effect metadata.
func (s *Saturation) Meta() *fx.EffectMeta {
	return s.meta
}
