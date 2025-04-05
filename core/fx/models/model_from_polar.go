package models

import (
	"image"
	"image/color"

	"github.com/toxyl/gfx/core/fx"
	"github.com/toxyl/gfx/core/meta"
	"github.com/toxyl/math"
)

// FromPolar represents a transformation from polar to Cartesian coordinates.
type FromPolar struct {
	Fisheye float64 // Fisheye correction factor (0.0 to 1.0)
	meta    *fx.EffectMeta
}

// NewFromPolarEffect creates a new from-polar transformation effect.
func NewFromPolarEffect(fisheye float64) *FromPolar {
	f := &FromPolar{
		Fisheye: fisheye,
		meta: fx.NewEffectMeta(
			"FromPolar",
			"Transforms an image from polar to Cartesian coordinates",
			meta.NewChannelMeta("Fisheye", 0.0, 1.0, "", "Fisheye correction factor (0.0 to 1.0)"),
		),
	}
	f.Fisheye = fx.ClampParameter(fisheye, f.meta.Parameters[0])
	return f
}

// Apply applies the from-polar transformation effect to an image.
func (f *FromPolar) Apply(img image.Image) image.Image {
	bounds := img.Bounds()
	width := bounds.Max.X - bounds.Min.X
	height := bounds.Max.Y - bounds.Min.Y

	// Calculate center point
	centerX := float64(width) / 2
	centerY := float64(height) / 2

	// Calculate maximum radius
	maxRadius := math.Min(centerX, centerY)

	// Create new image with same dimensions
	dst := image.NewRGBA(bounds)

	// Calculate fisheye correction factor
	fisheye := 1.0 - f.Fisheye*0.5 // Range from 1.0 to 0.5

	for y := bounds.Min.Y; y < bounds.Max.Y; y++ {
		for x := bounds.Min.X; x < bounds.Max.X; x++ {
			// Convert to polar coordinates
			dx := float64(x) - centerX
			dy := float64(y) - centerY
			radius := math.Sqrt(dx*dx + dy*dy)
			angle := math.Atan2(dy, dx)

			// Apply inverse fisheye correction
			correctedRadius := radius * math.Pow(radius/maxRadius, -fisheye)

			// Convert back to Cartesian coordinates
			srcX := centerX + correctedRadius*math.Cos(angle)
			srcY := centerY + correctedRadius*math.Sin(angle)

			// Ensure source coordinates are within bounds
			if srcX >= 0 && srcX < float64(width) && srcY >= 0 && srcY < float64(height) {
				dst.Set(x, y, img.At(int(srcX)+bounds.Min.X, int(srcY)+bounds.Min.Y))
			} else {
				dst.Set(x, y, color.RGBA64{})
			}
		}
	}

	return dst
}

// Meta returns the effect metadata.
func (f *FromPolar) Meta() *fx.EffectMeta {
	return f.meta
}
