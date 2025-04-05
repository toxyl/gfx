package models

import (
	"image"
	"image/color"

	"github.com/toxyl/gfx/core/fx"
	"github.com/toxyl/gfx/core/meta"
	"github.com/toxyl/math"
)

// Hue represents a hue adjustment effect.
type Hue struct {
	Amount float64 // Hue adjustment value (-180.0 to 180.0)
	meta   *fx.EffectMeta
}

// NewHue creates a new hue adjustment effect.
func NewHue(amount float64) *Hue {
	h := &Hue{
		Amount: amount,
		meta: fx.NewEffectMeta(
			"Hue",
			"Adjusts the hue of an image",
			meta.NewChannelMeta("Amount", -180.0, 180.0, "°", "Hue adjustment value in degrees (-180.0 to 180.0)"),
		),
	}
	h.Amount = fx.ClampParameter(amount, h.meta.Parameters[0])
	return h
}

// Apply applies the hue adjustment effect to an image.
func (h *Hue) Apply(img image.Image) image.Image {
	bounds := img.Bounds()
	dst := image.NewRGBA(bounds)

	// Convert degrees to radians
	radians := h.Amount * math.Pi / 180.0
	cos := math.Cos(radians)
	sin := math.Sin(radians)

	for y := bounds.Min.Y; y < bounds.Max.Y; y++ {
		for x := bounds.Min.X; x < bounds.Max.X; x++ {
			r, green, blue, a := img.At(x, y).RGBA()

			// Convert to float64 for calculations
			rf := float64(r) / 0xFFFF
			gf := float64(green) / 0xFFFF
			bf := float64(blue) / 0xFFFF

			// Convert RGB to YIQ
			luma := 0.299*rf + 0.587*gf + 0.114*bf
			i := 0.596*rf - 0.274*gf - 0.322*bf
			q := 0.211*rf - 0.523*gf + 0.312*bf

			// Rotate I and Q components
			iNew := i*cos - q*sin
			qNew := i*sin + q*cos

			// Convert back to RGB
			rf = luma + 0.956*iNew + 0.621*qNew
			gf = luma - 0.272*iNew - 0.647*qNew
			bf = luma - 1.106*iNew + 1.703*qNew

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
func (h *Hue) Meta() *fx.EffectMeta {
	return h.meta
}
