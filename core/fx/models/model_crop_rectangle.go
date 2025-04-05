package models

import (
	"image"
)

// CropRectangle represents a rectangular crop effect.
type CropRectangle struct {
	X      float64
	Y      float64
	Width  float64
	Height float64
}

// Apply applies the rectangular crop effect to an image.
func (c *CropRectangle) Apply(img image.Image) (image.Image, error) {
	bounds := img.Bounds()
	width := bounds.Max.X - bounds.Min.X
	height := bounds.Max.Y - bounds.Min.Y

	// Calculate crop rectangle in pixels
	x := int(float64(width) * c.X)
	y := int(float64(height) * c.Y)
	w := int(float64(width) * c.Width)
	h := int(float64(height) * c.Height)

	// Ensure crop rectangle is within image bounds
	if x < 0 {
		x = 0
	}
	if y < 0 {
		y = 0
	}
	if x+w > width {
		w = width - x
	}
	if y+h > height {
		h = height - y
	}

	// Create new image with crop dimensions
	dst := image.NewRGBA(image.Rect(0, 0, w, h))

	// Copy pixels from source to destination
	for dy := 0; dy < h; dy++ {
		for dx := 0; dx < w; dx++ {
			dst.Set(dx, dy, img.At(x+dx, y+dy))
		}
	}

	return dst, nil
}
