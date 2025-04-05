// core/color/base_rgb8.go
package color

import (
	"fmt"

	"github.com/toxyl/gfx/core/color/utils"
	"github.com/toxyl/math"
)

//////////////////////////////////////////////////////
// Implementation check
//////////////////////////////////////////////////////

var _ iColor = (*RGB8)(nil) // Ensure RGB8 implements the ColorModel interface.

//////////////////////////////////////////////////////
// Constructors
//////////////////////////////////////////////////////

// NewRGB8 creates a new RGB8 instance.
// R, G, B and Alpha are in range [0,255].
func NewRGB8[N math.Number](r, g, b, alpha N) (*RGB8, error) {
	return newColor(func() *RGB8 { return &RGB8{} }, r, g, b, alpha)
}

// RGB8FromRGB converts an RGBA64 (RGB) to an RGB8 color.
func RGB8FromRGB(c *RGBA64) *RGB8 {
	// Convert RGB to RGB8
	r8, g8, b8 := utils.RGBToRGB8(c.R, c.G, c.B)

	return &RGB8{
		R:     r8,
		G:     g8,
		B:     b8,
		Alpha: c.A * utils.RGB8Max,
	}
}

//////////////////////////////////////////////////////
// Type
//////////////////////////////////////////////////////

// RGB8 is a helper struct representing a color in the RGB8 color model.
// RGB8 represents colors using 8-bit values for each channel.
type RGB8 struct {
	R     float64 // [0,255] Red component
	G     float64 // [0,255] Green component
	B     float64 // [0,255] Blue component
	Alpha float64 // [0,255] Alpha
}

func (r *RGB8) Meta() *ColorModelMeta {
	return NewModelMeta(
		"RGB8",
		"8-bit RGB color model.",
		NewChannelMeta("R", 0, utils.RGB8Max, "", "Red component."),
		NewChannelMeta("G", 0, utils.RGB8Max, "", "Green component."),
		NewChannelMeta("B", 0, utils.RGB8Max, "", "Blue component."),
		NewChannelMeta("Alpha", 0, utils.RGB8Max, "", "Alpha channel."),
	)
}

//////////////////////////////////////////////////////
// Conversion
//////////////////////////////////////////////////////

func (r *RGB8) ToRGB() *RGBA64 {
	// Convert RGB8 to RGB
	rVal, gVal, bVal := utils.RGB8ToRGB(r.R, r.G, r.B)

	return &RGBA64{
		R: rVal,
		G: gVal,
		B: bVal,
		A: r.Alpha / utils.RGB8Max,
	}
}

// FromSlice initializes an RGB8 instance from a slice of float64 values.
// The slice must contain exactly 4 values in the order: R, G, B, Alpha.
func (r *RGB8) FromSlice(vals []float64) error {
	if len(vals) != 4 {
		return fmt.Errorf("RGB8 requires 4 values, got %d", len(vals))
	}

	r.R = vals[0]
	r.G = vals[1]
	r.B = vals[2]
	r.Alpha = vals[3]

	return nil
}

// FromRGBA64 converts an RGBA64 color to this color model.
func (r *RGB8) FromRGBA64(rgba *RGBA64) iColor {
	return RGB8FromRGB(rgba)
}

// ToRGBA64 converts the color to RGBA64.
func (r *RGB8) ToRGBA64() *RGBA64 {
	return r.ToRGB()
}
