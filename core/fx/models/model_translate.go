package models

import (
	"image"
	"image/color"

	"github.com/toxyl/gfx/core/fx"
	"github.com/toxyl/gfx/core/meta"
	"github.com/toxyl/math"
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

// Apply applies the translation effect to an image.
func (t *Translate) Apply(img image.Image) image.Image {
	bounds := img.Bounds()
	width := bounds.Max.X - bounds.Min.X
	height := bounds.Max.Y - bounds.Min.Y

	// Calculate new dimensions
	newWidth := width + int(math.Abs(t.OffsetX))
	newHeight := height + int(math.Abs(t.OffsetY))

	// Create new image with calculated dimensions
	dst := image.NewRGBA(image.Rect(0, 0, newWidth, newHeight))

	// Calculate offset in pixels
	offsetX := int(t.OffsetX)
	offsetY := int(t.OffsetY)

	// Adjust offset to ensure positive coordinates
	if offsetX < 0 {
		offsetX = -offsetX
	}
	if offsetY < 0 {
		offsetY = -offsetY
	}

	for y := 0; y < newHeight; y++ {
		for x := 0; x < newWidth; x++ {
			// Calculate source coordinates
			srcX := x - offsetX
			srcY := y - offsetY

			// Ensure source coordinates are within bounds
			if srcX >= 0 && srcX < width && srcY >= 0 && srcY < height {
				dst.Set(x, y, img.At(srcX+bounds.Min.X, srcY+bounds.Min.Y))
			} else {
				dst.Set(x, y, color.RGBA64{})
			}
		}
	}

	return dst
}

// Meta returns the effect metadata.
func (t *Translate) Meta() *fx.EffectMeta {
	return t.meta
}
