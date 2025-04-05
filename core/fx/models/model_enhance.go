package models

import (
	"image"
	"image/color"

	"github.com/toxyl/gfx/core/fx"
	"github.com/toxyl/gfx/core/meta"
	"github.com/toxyl/math"
)

// Enhance represents an image enhancement effect.
type Enhance struct {
	Amount float64 // Enhancement amount (0.0 to 1.0)
	meta   *fx.EffectMeta
}

// NewEnhanceEffect creates a new enhance effect.
func NewEnhanceEffect(amount float64) *Enhance {
	e := &Enhance{
		Amount: amount,
		meta: fx.NewEffectMeta(
			"Enhance",
			"Enhances image details and contrast",
			meta.NewChannelMeta("Amount", 0.0, 1.0, "", "Enhancement amount (0.0 to 1.0)"),
		),
	}
	e.Amount = fx.ClampParameter(amount, e.meta.Parameters[0])
	return e
}

// Apply applies the enhance effect to an image.
func (e *Enhance) Apply(img image.Image) (image.Image, error) {
	bounds := img.Bounds()
	dst := image.NewRGBA(bounds)

	for y := bounds.Min.Y; y < bounds.Max.Y; y++ {
		for x := bounds.Min.X; x < bounds.Max.X; x++ {
			r, g, b, a := img.At(x, y).RGBA()

			// Convert to float64 for calculations
			rf := float64(r) / 0xFFFF
			gf := float64(g) / 0xFFFF
			bf := float64(b) / 0xFFFF

			// Apply enhancement
			rf = math.Max(0, math.Min(1, rf*e.Amount))
			gf = math.Max(0, math.Min(1, gf*e.Amount))
			bf = math.Max(0, math.Min(1, bf*e.Amount))

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
func (e *Enhance) Meta() *fx.EffectMeta {
	return e.meta
}
