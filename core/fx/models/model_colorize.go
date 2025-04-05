package models

import (
	"image"
	"image/color"

	"github.com/toxyl/gfx/core/fx"
	"github.com/toxyl/gfx/core/meta"
	"github.com/toxyl/math"
)

// Colorize represents a colorize effect.
type Colorize struct {
	Hue        float64 // Hue value (0.0 to 360.0)
	Saturation float64 // Saturation value (0.0 to 1.0)
	Strength   float64 // Effect strength (0.0 to 1.0)
	meta       *fx.EffectMeta
}

// NewColorizeEffect creates a new colorize effect.
func NewColorizeEffect(hue, saturation, strength float64) *Colorize {
	c := &Colorize{
		Hue:        hue,
		Saturation: saturation,
		Strength:   strength,
		meta: fx.NewEffectMeta(
			"Colorize",
			"Applies a color tint to an image while preserving luminance",
			meta.NewChannelMeta("Hue", 0.0, 360.0, "°", "Hue value (0.0 to 360.0)"),
			meta.NewChannelMeta("Saturation", 0.0, 1.0, "", "Saturation value (0.0 to 1.0)"),
			meta.NewChannelMeta("Strength", 0.0, 1.0, "", "Effect strength (0.0 to 1.0)"),
		),
	}
	c.Hue = fx.ClampParameter(hue, c.meta.Parameters[0])
	c.Saturation = fx.ClampParameter(saturation, c.meta.Parameters[1])
	c.Strength = fx.ClampParameter(strength, c.meta.Parameters[2])
	return c
}

// hueToRGB converts hue to RGB values.
func (c *Colorize) hueToRGB(h float64) (r, g, b float64) {
	h = math.Mod(h, 360.0)
	if h < 0 {
		h += 360.0
	}

	h /= 60.0
	i := math.Floor(h)
	f := h - i
	p := 0.0
	q := 1.0 - f
	t := f

	switch int(i) {
	case 0:
		r, g, b = 1, t, p
	case 1:
		r, g, b = q, 1, p
	case 2:
		r, g, b = p, 1, t
	case 3:
		r, g, b = p, q, 1
	case 4:
		r, g, b = t, p, 1
	default:
		r, g, b = 1, p, q
	}

	return r, g, b
}

// Apply applies the colorize effect to an image.
func (c *Colorize) Apply(img image.Image) image.Image {
	bounds := img.Bounds()
	dst := image.NewRGBA(bounds)

	// Convert hue to RGB
	targetR, targetG, targetB := c.hueToRGB(c.Hue)

	for y := bounds.Min.Y; y < bounds.Max.Y; y++ {
		for x := bounds.Min.X; x < bounds.Max.X; x++ {
			r, g, b, a := img.At(x, y).RGBA()

			// Convert to float64 for calculations
			rf := float64(r) / 0xFFFF
			gf := float64(g) / 0xFFFF
			bf := float64(b) / 0xFFFF

			// Calculate luminance
			luminance := 0.299*rf + 0.587*gf + 0.114*bf

			// Blend with target color
			rf = luminance + (targetR-luminance)*c.Saturation*c.Strength
			gf = luminance + (targetG-luminance)*c.Saturation*c.Strength
			bf = luminance + (targetB-luminance)*c.Saturation*c.Strength

			// Convert back to uint32
			r = uint32(math.Max(0, math.Min(0xFFFF, rf*0xFFFF)))
			g = uint32(math.Max(0, math.Min(0xFFFF, gf*0xFFFF)))
			b = uint32(math.Max(0, math.Min(0xFFFF, bf*0xFFFF)))

			dst.Set(x, y, color.RGBA64{
				R: uint16(r),
				G: uint16(g),
				B: uint16(b),
				A: uint16(a),
			})
		}
	}

	return dst
}

// Meta returns the effect metadata.
func (c *Colorize) Meta() *fx.EffectMeta {
	return c.meta
}
