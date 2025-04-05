// core/color/base_hsb.go
package color

import (
	"fmt"

	"github.com/toxyl/gfx/core/color/utils"
	"github.com/toxyl/math"
)

//////////////////////////////////////////////////////
// Implementation check
//////////////////////////////////////////////////////

var _ iColor = (*HSB)(nil) // Ensure HSB implements the ColorModel interface.

//////////////////////////////////////////////////////
// Constructors
//////////////////////////////////////////////////////

// NewHSB creates a new HSB instance.
// Hue is in degrees [0,360]
// Saturation is in range [0,1]
// Brightness is in range [0,1]
// Alpha is in range [0,1]
func NewHSB[N math.Number](hue, saturation, brightness, alpha N) (*HSB, error) {
	return newColor(func() *HSB { return &HSB{} }, hue, saturation, brightness, alpha)
}

// HSBFromRGB converts an RGBA64 (RGB) to an HSB color.
func HSBFromRGB(c *RGBA64) *HSB {
	// Convert RGB to HSB
	h, s, b := utils.RGBToHSB(c.R, c.G, c.B)

	return &HSB{
		Hue:        h,
		Saturation: s,
		Brightness: b,
		Alpha:      c.A,
	}
}

//////////////////////////////////////////////////////
// Type
//////////////////////////////////////////////////////

// HSB is a helper struct representing a color in the HSB color model.
type HSB struct {
	Hue        float64 // in degrees [0,360]
	Saturation float64 // [0,1]
	Brightness float64 // [0,1]
	Alpha      float64 // [0,1]
}

func (h *HSB) Meta() *ColorModelMeta {
	return NewModelMeta(
		"HSB",
		"Hue, Saturation, Brightness color model.",
		NewChannelMeta("H", 0, 360, "°", "Hue in degrees."),
		NewChannelMeta("S", 0, 1, "", "Saturation."),
		NewChannelMeta("B", 0, 1, "", "Brightness."),
		NewChannelMeta("Alpha", 0, 1, "", "Alpha channel."),
	)
}

//////////////////////////////////////////////////////
// Conversion
//////////////////////////////////////////////////////

func (h *HSB) ToRGB() *RGBA64 {
	// Convert HSB to RGB
	r, g, b := utils.HSBToRGB(h.Hue, h.Saturation, h.Brightness)

	return &RGBA64{
		R: r,
		G: g,
		B: b,
		A: h.Alpha,
	}
}

// FromSlice initializes a HSB instance from a slice of float64 values.
// The slice must contain exactly 4 values in the order: Hue, Saturation, Brightness, Alpha.
func (h *HSB) FromSlice(vals []float64) error {
	if len(vals) != 4 {
		return fmt.Errorf("HSB requires 4 values, got %d", len(vals))
	}

	h.Hue = vals[0]
	h.Saturation = vals[1]
	h.Brightness = vals[2]
	h.Alpha = vals[3]

	return nil
}

// FromRGBA64 converts an RGBA64 color to this color model.
func (h *HSB) FromRGBA64(rgba *RGBA64) iColor {
	return HSBFromRGB(rgba)
}

// ToRGBA64 converts the color to RGBA64.
func (h *HSB) ToRGBA64() *RGBA64 {
	return h.ToRGB()
}
