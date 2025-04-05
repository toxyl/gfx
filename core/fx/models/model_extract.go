package models

import (
	"image"
	"image/color"

	"github.com/toxyl/gfx/core/fx"
	"github.com/toxyl/gfx/core/meta"
)

// Extract represents an image extraction effect.
type Extract struct {
	Channel string // Channel to extract ("red", "green", "blue", "alpha", "luminance")
	Invert  bool   // Whether to invert the extracted channel
	meta    *fx.EffectMeta
}

// NewExtractEffect creates a new extract effect.
func NewExtractEffect(channel string, invert bool) *Extract {
	e := &Extract{
		Channel: channel,
		Invert:  invert,
		meta: fx.NewEffectMeta(
			"Extract",
			"Extracts a specific channel from an image",
			meta.NewChannelMeta("Channel", 0.0, 1.0, "", "Channel to extract (red, green, blue, alpha, luminance)"),
		),
	}
	return e
}

// Apply applies the extract effect to an image.
func (e *Extract) Apply(img image.Image) (image.Image, error) {
	bounds := img.Bounds()
	dst := image.NewRGBA(bounds)

	for y := bounds.Min.Y; y < bounds.Max.Y; y++ {
		for x := bounds.Min.X; x < bounds.Max.X; x++ {
			r, g, b, a := img.At(x, y).RGBA()

			// Convert to float64 and normalize
			rf := float64(r) / 0xFFFF
			gf := float64(g) / 0xFFFF
			bf := float64(b) / 0xFFFF
			af := float64(a) / 0xFFFF

			var value float64

			// Extract specified channel
			switch e.Channel {
			case "red":
				value = rf
			case "green":
				value = gf
			case "blue":
				value = bf
			case "alpha":
				value = af
			case "luminance":
				// Calculate luminance using standard weights
				value = 0.299*rf + 0.587*gf + 0.114*bf
			default:
				// Default to luminance if channel is unknown
				value = 0.299*rf + 0.587*gf + 0.114*bf
			}

			// Invert if requested
			if e.Invert {
				value = 1.0 - value
			}

			// Set all channels to the extracted value
			vi := uint16(value * 0xFFFF)
			dst.Set(x, y, color.RGBA64{
				R: vi,
				G: vi,
				B: vi,
				A: uint16(a),
			})
		}
	}

	return dst, nil
}

// Meta returns the effect metadata.
func (e *Extract) Meta() *fx.EffectMeta {
	return e.meta
}
