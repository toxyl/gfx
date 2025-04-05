package models

import (
	"image"
	"image/color"

	"github.com/toxyl/gfx/core/fx"
	"github.com/toxyl/gfx/core/meta"
)

// Translate represents a translation effect.
type Translate struct {
	OffsetX float64 // Horizontal offset
	OffsetY float64 // Vertical offset
	meta    *fx.EffectMeta
}

// NewTranslateEffect creates a new translation effect.
func NewTranslateEffect(offsetX, offsetY float64) *Translate {
	t := &Translate{
		OffsetX: offsetX,
		OffsetY: offsetY,
		meta: fx.NewEffectMeta(
			"Translate",
			"Translates an image by the specified offsets",
			meta.NewChannelMeta("OffsetX", -1000.0, 1000.0, "px", "Horizontal offset in pixels"),
			meta.NewChannelMeta("OffsetY", -1000.0, 1000.0, "px", "Vertical offset in pixels"),
		),
	}
	t.OffsetX = fx.ClampParameter(offsetX, t.meta.Parameters[0])
	t.OffsetY = fx.ClampParameter(offsetY, t.meta.Parameters[1])
	return t
}

// Apply applies the translate effect to an image.
func (t *Translate) Apply(img image.Image) (image.Image, error) {
	bounds := img.Bounds()
	width := bounds.Max.X - bounds.Min.X
	height := bounds.Max.Y - bounds.Min.Y

	// Calculate translation in pixels
	dx := int(float64(width) * t.OffsetX)
	dy := int(float64(height) * t.OffsetY)

	// Create new image with same dimensions
	dst := image.NewRGBA(bounds)

	for y := bounds.Min.Y; y < bounds.Max.Y; y++ {
		for x := bounds.Min.X; x < bounds.Max.X; x++ {
			// Calculate source coordinates
			srcX := x - dx
			srcY := y - dy

			// Check if source coordinates are within bounds
			if srcX >= bounds.Min.X && srcX < bounds.Max.X &&
				srcY >= bounds.Min.Y && srcY < bounds.Max.Y {
				dst.Set(x, y, img.At(srcX, srcY))
			} else {
				// Set to transparent if outside bounds
				dst.Set(x, y, color.RGBA64{})
			}
		}
	}

	return dst, nil
}

// Meta returns the effect metadata.
func (t *Translate) Meta() *fx.EffectMeta {
	return t.meta
}
