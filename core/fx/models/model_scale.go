package models

import (
	"image"
	"image/color"

	"github.com/toxyl/gfx/core/fx"
	"github.com/toxyl/gfx/core/meta"
)

// Scale represents a scale effect.
type Scale struct {
	ScaleX float64 // Horizontal scale factor
	ScaleY float64 // Vertical scale factor
	meta   *fx.EffectMeta
}

// NewScaleEffect creates a new scale effect.
func NewScaleEffect(scaleX, scaleY float64) *Scale {
	s := &Scale{
		ScaleX: scaleX,
		ScaleY: scaleY,
		meta: fx.NewEffectMeta(
			"Scale",
			"Scales an image by the specified factors",
			meta.NewChannelMeta("ScaleX", 0.1, 10.0, "", "Horizontal scale factor"),
			meta.NewChannelMeta("ScaleY", 0.1, 10.0, "", "Vertical scale factor"),
		),
	}
	s.ScaleX = fx.ClampParameter(scaleX, s.meta.Parameters[0])
	s.ScaleY = fx.ClampParameter(scaleY, s.meta.Parameters[1])
	return s
}

// Apply applies the scaling effect to an image.
func (s *Scale) Apply(img image.Image) (image.Image, error) {
	bounds := img.Bounds()
	width := bounds.Dx()
	height := bounds.Dy()

	// Calculate new dimensions
	newWidth := int(float64(width) * s.ScaleX)
	newHeight := int(float64(height) * s.ScaleY)

	// Create new image
	dst := image.NewRGBA(image.Rect(0, 0, newWidth, newHeight))

	// Scale each pixel
	for y := 0; y < newHeight; y++ {
		for x := 0; x < newWidth; x++ {
			// Calculate source coordinates
			srcX := int(float64(x) / s.ScaleX)
			srcY := int(float64(y) / s.ScaleY)

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

// Meta returns the effect metadata.
func (s *Scale) Meta() *fx.EffectMeta {
	return s.meta
}
