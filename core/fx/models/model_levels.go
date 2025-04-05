package models

import (
	"image"
	"image/color"

	"github.com/toxyl/gfx/core/fx"
	"github.com/toxyl/gfx/core/meta"
	"github.com/toxyl/math"
)

// Levels represents a levels adjustment effect.
type Levels struct {
	BlackPoint float64 // Black point adjustment (0.0 to 1.0)
	WhitePoint float64 // White point adjustment (0.0 to 1.0)
	Gamma      float64 // Gamma adjustment (0.1 to 5.0)
	meta       *fx.EffectMeta
}

// NewLevelsEffect creates a new levels adjustment effect.
func NewLevelsEffect(blackPoint, whitePoint, gamma float64) *Levels {
	l := &Levels{
		BlackPoint: blackPoint,
		WhitePoint: whitePoint,
		Gamma:      gamma,
		meta: fx.NewEffectMeta(
			"Levels",
			"Adjusts the black point, white point, and gamma of an image",
			meta.NewChannelMeta("BlackPoint", 0.0, 1.0, "", "Black point adjustment (0.0 to 1.0)"),
			meta.NewChannelMeta("WhitePoint", 0.0, 1.0, "", "White point adjustment (0.0 to 1.0)"),
			meta.NewChannelMeta("Gamma", 0.1, 5.0, "", "Gamma adjustment (0.1 to 5.0)"),
		),
	}
	l.BlackPoint = fx.ClampParameter(blackPoint, l.meta.Parameters[0])
	l.WhitePoint = fx.ClampParameter(whitePoint, l.meta.Parameters[1])
	l.Gamma = fx.ClampParameter(gamma, l.meta.Parameters[2])
	return l
}

// Apply applies the levels adjustment effect to an image.
func (l *Levels) Apply(img image.Image) image.Image {
	bounds := img.Bounds()
	dst := image.NewRGBA(bounds)

	// Calculate the range for normalization
	rangeMin := l.BlackPoint
	rangeMax := l.WhitePoint
	rangeDiff := rangeMax - rangeMin

	for y := bounds.Min.Y; y < bounds.Max.Y; y++ {
		for x := bounds.Min.X; x < bounds.Max.X; x++ {
			r, g, b, a := img.At(x, y).RGBA()

			// Convert to float64 for calculations
			rf := float64(r) / 0xFFFF
			gf := float64(g) / 0xFFFF
			bf := float64(b) / 0xFFFF

			// Apply levels adjustment
			rf = math.Pow((rf-rangeMin)/rangeDiff, 1.0/l.Gamma)
			gf = math.Pow((gf-rangeMin)/rangeDiff, 1.0/l.Gamma)
			bf = math.Pow((bf-rangeMin)/rangeDiff, 1.0/l.Gamma)

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

	return dst
}

// Meta returns the effect metadata.
func (l *Levels) Meta() *fx.EffectMeta {
	return l.meta
}
