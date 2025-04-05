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
func (c *CropCircle) Apply(img image.Image) image.Image {
	bounds := img.Bounds()
	width := bounds.Max.X - bounds.Min.X
	height := bounds.Max.Y - bounds.Min.Y

	// Calculate center and radius in pixels
	centerX := float64(width) * c.CenterX
	centerY := float64(height) * c.CenterY
	radius := float64(width) * c.Radius

	// Calculate bounding box of circle
	minX := int(centerX - radius)
	minY := int(centerY - radius)
	maxX := int(centerX + radius)
	maxY := int(centerY + radius)

	// Ensure bounding box is within image bounds
	if minX < 0 {
		minX = 0
	}
	if minY < 0 {
		minY = 0
	}
	if maxX > width {
		maxX = width
	}
	if maxY > height {
		maxY = height
	}

	// Create new image with circle dimensions
	dst := image.NewRGBA(image.Rect(0, 0, maxX-minX, maxY-minY))

	// Copy pixels from source to destination, only within the circle
	for y := minY; y < maxY; y++ {
		for x := minX; x < maxX; x++ {
			// Calculate distance from center
			dx := float64(x) - centerX
			dy := float64(y) - centerY
			distance := math.Sqrt(dx*dx + dy*dy)

			// Only copy pixels within the circle
			if distance <= radius {
				dst.Set(x-minX, y-minY, img.At(x+bounds.Min.X, y+bounds.Min.Y))
			} else {
				dst.Set(x-minX, y-minY, color.RGBA64{})
			}
		}
	}

	return dst
}

// Meta returns the effect metadata.
func (c *CropCircle) Meta() *fx.EffectMeta {
	return c.meta
}
