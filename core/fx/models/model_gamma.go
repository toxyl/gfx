package models

import (
	"image"
	"image/color"

	"github.com/toxyl/gfx/core/fx"
	"github.com/toxyl/gfx/core/meta"
	"github.com/toxyl/math"
)

// Gamma represents a gamma correction effect.
type Gamma struct {
	Amount float64 // Gamma correction value (0.1 to 5.0)
	meta   *fx.EffectMeta
}

// NewGamma creates a new gamma correction effect.
func NewGamma(amount float64) *Gamma {
	g := &Gamma{
		Amount: amount,
		meta: fx.NewEffectMeta(
			"Gamma",
			"Applies gamma correction to an image",
			meta.NewChannelMeta("Amount", 0.1, 5.0, "", "Gamma correction value (0.1 to 5.0)"),
		),
	}
	g.Amount = fx.ClampParameter(amount, g.meta.Parameters[0])
	return g
}

// Apply applies the gamma correction effect to an image.
func (g *Gamma) Apply(img image.Image) image.Image {
	bounds := img.Bounds()
	dst := image.NewRGBA(bounds)

	for y := bounds.Min.Y; y < bounds.Max.Y; y++ {
		for x := bounds.Min.X; x < bounds.Max.X; x++ {
			r, green, blue, a := img.At(x, y).RGBA()

			// Convert to float64 for calculations
			rf := float64(r) / 0xFFFF
			gf := float64(green) / 0xFFFF
			bf := float64(blue) / 0xFFFF

			// Apply gamma correction
			rf = math.Pow(rf, 1.0/g.Amount)
			gf = math.Pow(gf, 1.0/g.Amount)
			bf = math.Pow(bf, 1.0/g.Amount)

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
func (g *Gamma) Meta() *fx.EffectMeta {
	return g.meta
}
