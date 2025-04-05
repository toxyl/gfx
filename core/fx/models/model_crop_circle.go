package models

import (
	"image"
	"image/color"

	"github.com/toxyl/gfx/core/fx"
	"github.com/toxyl/gfx/core/meta"
	"github.com/toxyl/math"
)

// CropCircle represents a circular crop effect.
type CropCircle struct {
	CenterX float64 // X coordinate of center (0.0 to 1.0)
	CenterY float64 // Y coordinate of center (0.0 to 1.0)
	Radius  float64 // Radius of circle (0.0 to 1.0)
	meta    *fx.EffectMeta
}

// NewCropCircleEffect creates a new circular crop effect.
func NewCropCircleEffect(centerX, centerY, radius float64) *CropCircle {
	c := &CropCircle{
		CenterX: centerX,
		CenterY: centerY,
		Radius:  radius,
		meta: fx.NewEffectMeta(
			"CropCircle",
			"Crops an image to a circular area",
			meta.NewChannelMeta("CenterX", 0.0, 1.0, "", "X coordinate of center (0.0 to 1.0)"),
			meta.NewChannelMeta("CenterY", 0.0, 1.0, "", "Y coordinate of center (0.0 to 1.0)"),
			meta.NewChannelMeta("Radius", 0.0, 1.0, "", "Radius of circle (0.0 to 1.0)"),
		),
	}
	c.CenterX = fx.ClampParameter(centerX, c.meta.Parameters[0])
	c.CenterY = fx.ClampParameter(centerY, c.meta.Parameters[1])
	c.Radius = fx.ClampParameter(radius, c.meta.Parameters[2])
	return c
}

// Apply applies the circular crop effect to an image.
func (c *CropCircle) Apply(img image.Image) (image.Image, error) {
	bounds := img.Bounds()
	width := bounds.Max.X - bounds.Min.X
	height := bounds.Max.Y - bounds.Min.Y

	// Calculate center and radius in pixels
	centerX := int(float64(width) * c.CenterX)
	centerY := int(float64(height) * c.CenterY)
	radius := int(float64(width) * c.Radius)

	// Create new image with same dimensions
	dst := image.NewRGBA(bounds)

	for y := bounds.Min.Y; y < bounds.Max.Y; y++ {
		for x := bounds.Min.X; x < bounds.Max.X; x++ {
			// Calculate distance from center
			dx := x - centerX
			dy := y - centerY
			distance := math.Sqrt(float64(dx*dx + dy*dy))

			// If pixel is within circle, copy from source
			if distance <= float64(radius) {
				dst.Set(x, y, img.At(x, y))
			} else {
				// Otherwise set to transparent
				dst.Set(x, y, color.RGBA64{})
			}
		}
	}

	return dst, nil
}

// Meta returns the effect metadata.
func (c *CropCircle) Meta() *fx.EffectMeta {
	return c.meta
}
