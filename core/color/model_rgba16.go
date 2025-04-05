package color

import (
	"fmt"

	"github.com/toxyl/math"
)

// RGBA16 represents a 16-bit RGBA color.
type RGBA16 struct {
	R     uint16 // [0,65535] Red
	G     uint16 // [0,65535] Green
	B     uint16 // [0,65535] Blue
	Alpha uint16 // [0,65535] Alpha
}

// NewRGBA16 creates a new RGBA16 instance.
func NewRGBA16[N math.Number](r, g, b, alpha N) (*RGBA16, error) {
	return newColor(func() *RGBA16 { return &RGBA16{} }, r, g, b, alpha)
}

// RGBA16FromRGB converts an RGBA64 (RGB) to an RGBA16 color.
func RGBA16FromRGB(c *RGBA64) *RGBA16 {
	return &RGBA16{
		R:     uint16(c.R * 65535),
		G:     uint16(c.G * 65535),
		B:     uint16(c.B * 65535),
		Alpha: uint16(c.A * 65535),
	}
}

// Meta returns the color model metadata.
func (c *RGBA16) Meta() *ColorModelMeta {
	return NewModelMeta(
		"RGBA16",
		"16-bit RGBA color model.",
		NewChannelMeta("R", 0, 65535, "", "Red channel."),
		NewChannelMeta("G", 0, 65535, "", "Green channel."),
		NewChannelMeta("B", 0, 65535, "", "Blue channel."),
		NewChannelMeta("Alpha", 0, 65535, "", "Alpha channel."),
	)
}

// ToRGB converts the color to RGBA64.
func (c *RGBA16) ToRGB() *RGBA64 {
	return &RGBA64{
		R: float64(c.R) / 65535,
		G: float64(c.G) / 65535,
		B: float64(c.B) / 65535,
		A: float64(c.Alpha) / 65535,
	}
}

// FromSlice initializes an RGBA16 instance from a slice of float64 values.
// The slice must contain exactly 4 values in the order: R, G, B, Alpha.
func (c *RGBA16) FromSlice(vals []float64) error {
	if len(vals) != 4 {
		return fmt.Errorf("RGBA16 requires 4 values, got %d", len(vals))
	}

	c.R = uint16(vals[0] * 65535)
	c.G = uint16(vals[1] * 65535)
	c.B = uint16(vals[2] * 65535)
	c.Alpha = uint16(vals[3] * 65535)

	return nil
}

// FromRGBA64 converts an RGBA64 color to this color model.
func (c *RGBA16) FromRGBA64(rgba *RGBA64) iColor {
	return RGBA16FromRGB(rgba)
}

// ToRGBA64 converts the color to RGBA64.
func (c *RGBA16) ToRGBA64() *RGBA64 {
	return c.ToRGB()
}
