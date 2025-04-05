package models

import (
	"image"
	"image/color"

	"github.com/toxyl/gfx/core/fx"
	"github.com/toxyl/gfx/core/meta"
)

// Threshold represents a threshold effect.
type Threshold struct {
	Level float64 // Threshold level (0.0 to 1.0)
	meta  *fx.EffectMeta
}

// NewThresholdEffect creates a new threshold effect.
func NewThresholdEffect(level float64) *Threshold {
	t := &Threshold{
		Level: level,
		meta: fx.NewEffectMeta(
			"Threshold",
			"Converts an image to black and white based on a threshold level",
			meta.NewChannelMeta("Level", 0.0, 1.0, "", "Threshold level (0.0 to 1.0)"),
		),
	}
	t.Level = fx.ClampParameter(level, t.meta.Parameters[0])
	return t
}

// Apply applies the threshold effect to an image.
func (t *Threshold) Apply(img image.Image) (image.Image, error) {
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

			// Apply threshold
			var value float64
			if lum >= t.Level {
				value = 1.0
			}

			// Set all channels to the threshold value
			vi := uint16(value * 0xFFFF)
			dst.Set(x, y, color.RGBA64{
				R: vi,
				G: vi,
				B: vi,
				A: uint16(a),
			})
		}
	}

	return dst, nil
}

// Meta returns the effect metadata.
func (t *Threshold) Meta() *fx.EffectMeta {
	return t.meta
}
