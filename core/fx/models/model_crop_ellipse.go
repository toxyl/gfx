package models

import (
	"image"
	"image/color"
)

// CropEllipse represents the elliptical crop effect.
type CropEllipse struct {
	CenterX float64
	CenterY float64
	RadiusX float64
	RadiusY float64
}

// Apply applies the elliptical crop effect to an image.
func (c *CropEllipse) Apply(img image.Image) (image.Image, error) {
	bounds := img.Bounds()
	width := bounds.Max.X - bounds.Min.X
	height := bounds.Max.Y - bounds.Min.Y

	// Calculate center and radii in pixels
	centerX := int(float64(width) * c.CenterX)
	centerY := int(float64(height) * c.CenterY)
	radiusX := int(float64(width) * c.RadiusX)
	radiusY := int(float64(height) * c.RadiusY)

	// Create new image with same dimensions
	dst := image.NewRGBA(bounds)

	for y := bounds.Min.Y; y < bounds.Max.Y; y++ {
		for x := bounds.Min.X; x < bounds.Max.X; x++ {
			// Calculate normalized distance from center
			dx := float64(x-centerX) / float64(radiusX)
			dy := float64(y-centerY) / float64(radiusY)
			distance := dx*dx + dy*dy

			// If pixel is within ellipse, copy from source
			if distance <= 1.0 {
				dst.Set(x, y, img.At(x, y))
			} else {
				// Otherwise set to transparent
				dst.Set(x, y, color.RGBA64{})
			}
		}
	}

	return dst, nil
}
