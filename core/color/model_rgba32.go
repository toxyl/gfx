// core/color/base_rgba32.go
package color

import (
	"fmt"

	"github.com/toxyl/gfx/core/color/utils"
	"github.com/toxyl/math"
)

//////////////////////////////////////////////////////
// Implementation check
//////////////////////////////////////////////////////

var _ iColor = (*RGBA32)(nil) // Ensure RGBA32 implements the ColorModel interface.

//////////////////////////////////////////////////////
// Constructors
//////////////////////////////////////////////////////

// NewRGBA32 creates a new RGBA32 instance.
// R, G, B and Alpha are in range [0,4294967295].
func NewRGBA32[N math.Number](r, g, b, alpha N) (*RGBA32, error) {
	return newColor(func() *RGBA32 { return &RGBA32{} }, r, g, b, alpha)
}

// RGBA32FromRGB converts an RGBA64 (RGB) to an RGBA32 color.
func RGBA32FromRGB(c *RGBA64) *RGBA32 {
	// Convert RGB to RGBA32
	r32, g32, b32 := utils.RGBToRGBA32(c.R, c.G, c.B)

	return &RGBA32{
		R:     r32,
		G:     g32,
		B:     b32,
		Alpha: c.A * utils.RGBA32Max,
	}
}

//////////////////////////////////////////////////////
// Type
//////////////////////////////////////////////////////

// RGBA32 is a helper struct representing a color in the RGBA32 color model.
// RGBA32 represents colors using 32-bit values for each channel.
type RGBA32 struct {
	R     float64 // [0,4294967295] Red component
	G     float64 // [0,4294967295] Green component
	B     float64 // [0,4294967295] Blue component
	Alpha float64 // [0,4294967295] Alpha
}

func (r *RGBA32) Meta() *ColorModelMeta {
	return NewModelMeta(
		"RGBA32",
		"32-bit RGBA color model.",
		NewChannelMeta("R", 0, utils.RGBA32Max, "", "Red component."),
		NewChannelMeta("G", 0, utils.RGBA32Max, "", "Green component."),
		NewChannelMeta("B", 0, utils.RGBA32Max, "", "Blue component."),
		NewChannelMeta("Alpha", 0, utils.RGBA32Max, "", "Alpha channel."),
	)
}

//////////////////////////////////////////////////////
// Conversion
//////////////////////////////////////////////////////

func (r *RGBA32) ToRGB() *RGBA64 {
	// Convert RGBA32 to RGB
	rVal, gVal, bVal := utils.RGBA32ToRGB(r.R, r.G, r.B)

	return &RGBA64{
		R: rVal,
		G: gVal,
		B: bVal,
		A: r.Alpha / utils.RGBA32Max,
	}
}

// FromSlice initializes an RGBA32 instance from a slice of float64 values.
// The slice must contain exactly 4 values in the order: R, G, B, Alpha.
func (r *RGBA32) FromSlice(vals []float64) error {
	if len(vals) != 4 {
		return fmt.Errorf("RGBA32 requires 4 values, got %d", len(vals))
	}

	r.R = vals[0]
	r.G = vals[1]
	r.B = vals[2]
	r.Alpha = vals[3]

	return nil
}

// FromRGBA64 converts an RGBA64 color to this color model.
func (r *RGBA32) FromRGBA64(rgba *RGBA64) iColor {
	return RGBA32FromRGB(rgba)
}

// ToRGBA64 converts the color to RGBA64.
func (r *RGBA32) ToRGBA64() *RGBA64 {
	return r.ToRGB()
}
