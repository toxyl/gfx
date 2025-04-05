package models

import (
	"image"

	"github.com/toxyl/gfx/core/fx"
)

// FlipH represents a horizontal flip effect.
type FlipH struct {
	meta *fx.EffectMeta
}

// NewFlipHEffect creates a new horizontal flip effect.
func NewFlipHEffect() *FlipH {
	return &FlipH{
		meta: fx.NewEffectMeta(
			"FlipH",
			"Flips an image horizontally",
		),
	}
}

// Apply applies the horizontal flip effect to an image.
func (f *FlipH) Apply(img image.Image) image.Image {
	bounds := img.Bounds()

	// Create new image with same dimensions
	dst := image.NewRGBA(bounds)

	for y := bounds.Min.Y; y < bounds.Max.Y; y++ {
		for x := bounds.Min.X; x < bounds.Max.X; x++ {
			// Calculate mirrored x coordinate
			mirrorX := bounds.Max.X - 1 - (x - bounds.Min.X)
			dst.Set(x, y, img.At(mirrorX, y))
		}
	}

	return dst
}

// Meta returns the effect metadata.
func (f *FlipH) Meta() *fx.EffectMeta {
	return f.meta
}
