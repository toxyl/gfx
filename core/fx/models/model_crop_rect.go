package models

import (
	"image"

	"github.com/toxyl/gfx/core/fx"
	"github.com/toxyl/gfx/core/meta"
)

// CropRect represents a rectangular crop effect.
type CropRect struct {
	X      float64 // X coordinate of top-left corner (0.0 to 1.0)
	Y      float64 // Y coordinate of top-left corner (0.0 to 1.0)
	Width  float64 // Width of crop area (0.0 to 1.0)
	Height float64 // Height of crop area (0.0 to 1.0)
	meta   *fx.EffectMeta
}

// NewCropRectEffect creates a new rectangular crop effect.
func NewCropRectEffect(x, y, width, height float64) *CropRect {
	c := &CropRect{
		X:      x,
		Y:      y,
		Width:  width,
		Height: height,
		meta: fx.NewEffectMeta(
			"CropRect",
			"Crops an image to a rectangular area",
			meta.NewChannelMeta("X", 0.0, 1.0, "", "X coordinate of top-left corner (0.0 to 1.0)"),
			meta.NewChannelMeta("Y", 0.0, 1.0, "", "Y coordinate of top-left corner (0.0 to 1.0)"),
			meta.NewChannelMeta("Width", 0.0, 1.0, "", "Width of crop area (0.0 to 1.0)"),
			meta.NewChannelMeta("Height", 0.0, 1.0, "", "Height of crop area (0.0 to 1.0)"),
		),
	}
	c.X = fx.ClampParameter(x, c.meta.Parameters[0])
	c.Y = fx.ClampParameter(y, c.meta.Parameters[1])
	c.Width = fx.ClampParameter(width, c.meta.Parameters[2])
	c.Height = fx.ClampParameter(height, c.meta.Parameters[3])
	return c
}

// Apply applies the rectangular crop effect to an image.
func (c *CropRect) Apply(img image.Image) image.Image {
	bounds := img.Bounds()
	width := bounds.Max.X - bounds.Min.X
	height := bounds.Max.Y - bounds.Min.Y

	// Calculate crop rectangle in pixels
	cropX := int(float64(width) * c.X)
	cropY := int(float64(height) * c.Y)
	cropWidth := int(float64(width) * c.Width)
	cropHeight := int(float64(height) * c.Height)

	// Ensure crop rectangle is within image bounds
	if cropX < 0 {
		cropX = 0
	}
	if cropY < 0 {
		cropY = 0
	}
	if cropX+cropWidth > width {
		cropWidth = width - cropX
	}
	if cropY+cropHeight > height {
		cropHeight = height - cropY
	}

	// Create new image with crop dimensions
	dst := image.NewRGBA(image.Rect(0, 0, cropWidth, cropHeight))

	// Copy pixels from source to destination
	for y := 0; y < cropHeight; y++ {
		for x := 0; x < cropWidth; x++ {
			srcX := cropX + x + bounds.Min.X
			srcY := cropY + y + bounds.Min.Y
			dst.Set(x, y, img.At(srcX, srcY))
		}
	}

	return dst
}

// Meta returns the effect metadata.
func (c *CropRect) Meta() *fx.EffectMeta {
	return c.meta
}
