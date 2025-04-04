// core/color/base_rgb8.go
package color

import (
	"fmt"

	"github.com/toxyl/gfx/core/color/constants"
)

//////////////////////////////////////////////////////
// Conversion utilities
//////////////////////////////////////////////////////

func rgbToRgb8(r, g, b float64) (r8, g8, b8 uint8) {
	// Convert normalized [0,1] RGB to 8-bit RGB
	r8 = uint8(clamp(r, 0, 1) * constants.RGB8_Max)
	g8 = uint8(clamp(g, 0, 1) * constants.RGB8_Max)
	b8 = uint8(clamp(b, 0, 1) * constants.RGB8_Max)
	return r8, g8, b8
}

func rgb8ToRgb(r8, g8, b8 uint8) (r, g, b float64) {
	// Convert 8-bit RGB to normalized [0,1] RGB
	r = float64(r8) / constants.RGB8_Max
	g = float64(g8) / constants.RGB8_Max
	b = float64(b8) / constants.RGB8_Max
	return r, g, b
}

func clamp(x, min, max float64) float64 {
	if x < min {
		return min
	}
	if x > max {
		return max
	}
	return x
}

//////////////////////////////////////////////////////
// Implementation check
//////////////////////////////////////////////////////

var _ iColor = (*RGB8)(nil) // Ensure RGB8 implements the ColorModel interface.

//////////////////////////////////////////////////////
// Constructors
//////////////////////////////////////////////////////

// NewRGB8 creates a new RGB8 instance.
// It accepts 8-bit channel values in the [0,255] range for R, G, B and normalized [0,1] range for Alpha.
func NewRGB8[N Number](r, g, b, alpha N) (*RGB8, error) {
	return newColor(func() *RGB8 { return &RGB8{} }, r, g, b, alpha)
}

// RGB8FromRGB converts an RGBA64 (RGB) to an RGB8 color.
func RGB8FromRGB(c *RGBA64) *RGB8 {
	r8, g8, b8 := rgbToRgb8(c.R, c.G, c.B)
	return &RGB8{
		R:     float64(r8),
		G:     float64(g8),
		B:     float64(b8),
		Alpha: c.A,
	}
}

//////////////////////////////////////////////////////
// Type
//////////////////////////////////////////////////////

// RGB8 is a helper struct representing a color in the RGB color model with 8-bit channels and an alpha channel.
type RGB8 struct {
	R, G, B, Alpha float64
}

func (rgb8 *RGB8) Meta() *ColorModelMeta {
	return NewModelMeta(
		"RGB8",
		"8-bit RGB color model.",
		NewChannelMeta("R", 0, 255, "", "Red channel."),
		NewChannelMeta("G", 0, 255, "", "Green channel."),
		NewChannelMeta("B", 0, 255, "", "Blue channel."),
		NewChannelMeta("Alpha", 0, 1, "", "Alpha channel."),
	)
}

//////////////////////////////////////////////////////
// Conversion
//////////////////////////////////////////////////////

func (rgb8 *RGB8) ToRGB() *RGBA64 {
	r, g, b := rgb8ToRgb(uint8(rgb8.R), uint8(rgb8.G), uint8(rgb8.B))
	return &RGBA64{R: r, G: g, B: b, A: rgb8.Alpha}
}

func (rgb8 *RGB8) ToRGBA64() *RGBA64 {
	return rgb8.ToRGB()
}

func (rgb8 *RGB8) FromRGBA64(rgba *RGBA64) iColor {
	r8, g8, b8 := rgbToRgb8(rgba.R, rgba.G, rgba.B)
	return &RGB8{
		R:     float64(r8),
		G:     float64(g8),
		B:     float64(b8),
		Alpha: rgba.A,
	}
}

// FromSlice initializes an RGB8 instance from a slice of float64 values.
// The slice must contain exactly 4 values in the order: red, green, blue, alpha.
// The RGB values must be in the range [0, 255] and alpha in [0, 1].
func (rgb8 *RGB8) FromSlice(values []float64) error {
	if len(values) != 4 {
		return fmt.Errorf("expected 4 values for RGB8, got %d", len(values))
	}

	for i, v := range values {
		if err := ValidateChannelValue(v, rgb8.Meta().Channels()[i]); err != nil {
			return err
		}
	}

	rgb8.R = values[0]
	rgb8.G = values[1]
	rgb8.B = values[2]
	rgb8.Alpha = values[3]

	return nil
}
