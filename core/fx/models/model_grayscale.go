package models

import (
	"image"
	"image/color"

	"github.com/toxyl/gfx/core/fx"
	"github.com/toxyl/gfx/core/meta"
	"github.com/toxyl/math"
)

// Grayscale represents a grayscale effect.
type Grayscale struct {
	Amount float64 // Grayscale amount (0.0 to 1.0)
	meta   *fx.EffectMeta
}

// NewGrayscaleEffect creates a new grayscale effect.
func NewGrayscaleEffect(amount float64) *Grayscale {
	gs := &Grayscale{
		Amount: amount,
		meta: fx.NewEffectMeta(
			"Grayscale",
			"Converts an image to grayscale",
			meta.NewChannelMeta("Amount", 0.0, 1.0, "", "Grayscale amount (0.0 to 1.0)"),
		),
	}
	gs.Amount = fx.ClampParameter(amount, gs.meta.Parameters[0])
	return gs
}

// Apply applies the grayscale effect to an image.
func (gs *Grayscale) Apply(img image.Image) (image.Image, error) {
	bounds := img.Bounds()
	dst := image.NewRGBA(bounds)

	for y := bounds.Min.Y; y < bounds.Max.Y; y++ {
		for x := bounds.Min.X; x < bounds.Max.X; x++ {
			r, g, b, a := img.At(x, y).RGBA()

			// Convert to float64 for calculations
			rf := float64(r) / 0xFFFF
			gf := float64(g) / 0xFFFF
			bf := float64(b) / 0xFFFF

			// Calculate luminance using standard weights
			luminance := 0.299*rf + 0.587*gf + 0.114*bf

			// Blend between original and grayscale based on amount
			rf = rf + (luminance-rf)*gs.Amount
			gf = gf + (luminance-gf)*gs.Amount
			bf = bf + (luminance-bf)*gs.Amount

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
func (gs *Grayscale) Meta() *fx.EffectMeta {
	return gs.meta
}
