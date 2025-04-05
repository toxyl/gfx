package models

import (
	"image"
	"image/color"
	"math"
)

// Skew represents the skew effect.
type Skew struct {
	AngleX float64
	AngleY float64
}

// Apply applies the skew effect to an image.
func (s *Skew) Apply(img image.Image) (image.Image, error) {
	bounds := img.Bounds()
	width := bounds.Dx()
	height := bounds.Dy()

	// Convert angles to radians
	angleX := s.AngleX * math.Pi / 180.0
	angleY := s.AngleY * math.Pi / 180.0

	// Calculate skew factors
	skewX := math.Tan(angleX)
	skewY := math.Tan(angleY)

	// Calculate new dimensions
	newWidth := int(float64(width) + math.Abs(float64(height)*skewX))
	newHeight := int(float64(height) + math.Abs(float64(width)*skewY))

	// Create new image
	dst := image.NewRGBA(image.Rect(0, 0, newWidth, newHeight))

	// Skew each pixel
	for y := 0; y < newHeight; y++ {
		for x := 0; x < newWidth; x++ {
			// Calculate source coordinates
			srcX := int(float64(x) - float64(y)*skewX)
			srcY := int(float64(y) - float64(x)*skewY)

			// Check if source coordinates are within bounds
			if srcX >= 0 && srcX < width && srcY >= 0 && srcY < height {
				dst.Set(x, y, img.At(srcX, srcY))
			} else {
				dst.Set(x, y, color.Transparent)
			}
		}
	}

	return dst, nil
}
