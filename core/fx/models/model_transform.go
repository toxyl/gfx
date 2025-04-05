package models

import (
	"image"
	"image/color"

	"github.com/toxyl/gfx/core/fx"
	"github.com/toxyl/gfx/core/meta"
	"github.com/toxyl/math"
)

// Transform represents a combined transformation effect.
type Transform struct {
	ScaleX  float64 // Horizontal scale factor
	ScaleY  float64 // Vertical scale factor
	Angle   float64 // Rotation angle in degrees
	OffsetX float64 // Horizontal offset
	OffsetY float64 // Vertical offset
	meta    *fx.EffectMeta
}

// NewTransformEffect creates a new transformation effect.
func NewTransformEffect(scaleX, scaleY, angle, offsetX, offsetY float64) *Transform {
	t := &Transform{
		ScaleX:  scaleX,
		ScaleY:  scaleY,
		Angle:   angle,
		OffsetX: offsetX,
		OffsetY: offsetY,
		meta: fx.NewEffectMeta(
			"Transform",
			"Applies scale, rotation, and translation to an image",
			meta.NewChannelMeta("ScaleX", 0.1, 10.0, "", "Horizontal scale factor"),
			meta.NewChannelMeta("ScaleY", 0.1, 10.0, "", "Vertical scale factor"),
			meta.NewChannelMeta("Angle", -360.0, 360.0, "°", "Rotation angle in degrees"),
			meta.NewChannelMeta("OffsetX", -1000.0, 1000.0, "px", "Horizontal offset in pixels"),
			meta.NewChannelMeta("OffsetY", -1000.0, 1000.0, "px", "Vertical offset in pixels"),
		),
	}
	t.ScaleX = fx.ClampParameter(scaleX, t.meta.Parameters[0])
	t.ScaleY = fx.ClampParameter(scaleY, t.meta.Parameters[1])
	t.Angle = fx.ClampParameter(angle, t.meta.Parameters[2])
	t.OffsetX = fx.ClampParameter(offsetX, t.meta.Parameters[3])
	t.OffsetY = fx.ClampParameter(offsetY, t.meta.Parameters[4])
	return t
}

// Apply applies the transformation effect to an image.
func (t *Transform) Apply(img image.Image) image.Image {
	bounds := img.Bounds()
	width := bounds.Max.X - bounds.Min.X
	height := bounds.Max.Y - bounds.Min.Y

	// Convert angle to radians
	angle := t.Angle * math.Pi / 180.0
	cos := math.Cos(angle)
	sin := math.Sin(angle)

	// Calculate scaled dimensions
	scaledWidth := float64(width) * t.ScaleX
	scaledHeight := float64(height) * t.ScaleY

	// Calculate the corners of the transformed image
	corners := [][2]float64{
		{0, 0},
		{scaledWidth, 0},
		{scaledWidth, scaledHeight},
		{0, scaledHeight},
	}

	// Rotate corners and find min/max coordinates
	minX, minY := math.Inf(1), math.Inf(1)
	maxX, maxY := math.Inf(-1), math.Inf(-1)

	for _, corner := range corners {
		// Apply rotation
		rotX := corner[0]*cos - corner[1]*sin
		rotY := corner[0]*sin + corner[1]*cos

		// Apply translation
		rotX += t.OffsetX
		rotY += t.OffsetY

		minX = math.Min(minX, rotX)
		minY = math.Min(minY, rotY)
		maxX = math.Max(maxX, rotX)
		maxY = math.Max(maxY, rotY)
	}

	newWidth := int(maxX - minX)
	newHeight := int(maxY - minY)

	// Create new image with calculated dimensions
	dst := image.NewRGBA(image.Rect(0, 0, newWidth, newHeight))

	// Calculate center of new image
	centerX := float64(newWidth) / 2
	centerY := float64(newHeight) / 2

	for y := 0; y < newHeight; y++ {
		for x := 0; x < newWidth; x++ {
			// Convert to coordinates relative to center
			relX := float64(x) - centerX
			relY := float64(y) - centerY

			// Apply inverse translation
			relX -= t.OffsetX
			relY -= t.OffsetY

			// Apply inverse rotation
			srcX := relX*cos + relY*sin
			srcY := -relX*sin + relY*cos

			// Apply inverse scale
			srcX /= t.ScaleX
			srcY /= t.ScaleY

			// Ensure source coordinates are within bounds
			if srcX >= 0 && srcX < float64(width) && srcY >= 0 && srcY < float64(height) {
				dst.Set(x, y, img.At(int(srcX)+bounds.Min.X, int(srcY)+bounds.Min.Y))
			} else {
				dst.Set(x, y, color.RGBA64{})
			}
		}
	}

	return dst
}

// Meta returns the effect metadata.
func (t *Transform) Meta() *fx.EffectMeta {
	return t.meta
}
