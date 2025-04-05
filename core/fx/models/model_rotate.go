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
func (r *Rotate) Apply(img image.Image) (image.Image, error) {
	bounds := img.Bounds()
	width := bounds.Dx()
	height := bounds.Dy()

	// Convert angle to radians
	angle := r.Angle * math.Pi / 180.0

	// Calculate new dimensions
	cos := math.Cos(angle)
	sin := math.Sin(angle)

	// Calculate new image dimensions
	newWidth := int(math.Abs(float64(width)*cos) + math.Abs(float64(height)*sin))
	newHeight := int(math.Abs(float64(width)*sin) + math.Abs(float64(height)*cos))

	// Create new image
	dst := image.NewRGBA(image.Rect(0, 0, newWidth, newHeight))

	// Calculate center points
	centerX := float64(width) / 2
	centerY := float64(height) / 2
	newCenterX := float64(newWidth) / 2
	newCenterY := float64(newHeight) / 2

	// Rotate each pixel
	for y := 0; y < newHeight; y++ {
		for x := 0; x < newWidth; x++ {
			// Calculate original coordinates
			dx := float64(x) - newCenterX
			dy := float64(y) - newCenterY

			// Rotate coordinates
			srcX := int(centerX + dx*cos + dy*sin)
			srcY := int(centerY - dx*sin + dy*cos)

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
func (r *Rotate) Meta() *fx.EffectMeta {
	return r.meta
}
