package models

import (
	"image"

	"github.com/toxyl/gfx/core/fx"
)

// FlipV represents a vertical flip effect.
type FlipV struct {
	meta *fx.EffectMeta
}

// NewFlipVEffect creates a new vertical flip effect.
func NewFlipVEffect() *FlipV {
	return &FlipV{
		meta: fx.NewEffectMeta(
			"FlipV",
			"Flips an image vertically",
		),
	}
}

// Apply applies the vertical flip effect to an image.
func (f *FlipV) Apply(img image.Image) image.Image {
	bounds := img.Bounds()

	// Create new image with same dimensions
	dst := image.NewRGBA(bounds)

	for y := bounds.Min.Y; y < bounds.Max.Y; y++ {
		for x := bounds.Min.X; x < bounds.Max.X; x++ {
			// Calculate mirrored y coordinate
			mirrorY := bounds.Max.Y - 1 - (y - bounds.Min.Y)
			dst.Set(x, y, img.At(x, mirrorY))
		}
	}

	return dst
}

// Meta returns the effect metadata.
func (f *FlipV) Meta() *fx.EffectMeta {
	return f.meta
}
