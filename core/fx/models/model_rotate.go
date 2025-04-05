package models

import (
	"image"
	"image/color"

	"github.com/toxyl/gfx/core/fx"
	"github.com/toxyl/gfx/core/meta"
	"github.com/toxyl/math"
)

// Rotate represents a rotation effect.
type Rotate struct {
	Angle float64 // Rotation angle in degrees
	meta  *fx.EffectMeta
}

// NewRotateEffect creates a new rotation effect.
func NewRotateEffect(angle float64) *Rotate {
	r := &Rotate{
		Angle: angle,
		meta: fx.NewEffectMeta(
			"Rotate",
			"Rotates an image by the specified angle",
			meta.NewChannelMeta("Angle", -360.0, 360.0, "°", "Rotation angle in degrees"),
		),
	}
	r.Angle = fx.ClampParameter(angle, r.meta.Parameters[0])
	return r
}

// Apply applies the rotation effect to an image.
func (r *Rotate) Apply(img image.Image) image.Image {
	bounds := img.Bounds()
	width := bounds.Max.X - bounds.Min.X
	height := bounds.Max.Y - bounds.Min.Y

	// Convert angle to radians
	angle := r.Angle * math.Pi / 180.0
	cos := math.Cos(angle)
	sin := math.Sin(angle)

	// Calculate new dimensions
	// We need to find the bounding box of the rotated image
	halfWidth := float64(width) / 2
	halfHeight := float64(height) / 2

	// Calculate the corners of the rotated image
	corners := [][2]float64{
		{-halfWidth, -halfHeight},
		{halfWidth, -halfHeight},
		{halfWidth, halfHeight},
		{-halfWidth, halfHeight},
	}

	// Rotate corners and find min/max coordinates
	minX, minY := math.Inf(1), math.Inf(1)
	maxX, maxY := math.Inf(-1), math.Inf(-1)

	for _, corner := range corners {
		rotX := corner[0]*cos - corner[1]*sin
		rotY := corner[0]*sin + corner[1]*cos
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

			// Rotate back to original coordinates
			srcX := relX*cos + relY*sin + halfWidth
			srcY := -relX*sin + relY*cos + halfHeight

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
func (r *Rotate) Meta() *fx.EffectMeta {
	return r.meta
}
