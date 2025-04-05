// core/color/base_lambdasl.go
package color

import (
	"fmt"

	"github.com/toxyl/gfx/core/color/constants"
	"github.com/toxyl/gfx/core/color/utils"
	"github.com/toxyl/math"
)

//////////////////////////////////////////////////////
// Implementation check
//////////////////////////////////////////////////////

var _ iColor = (*LSL)(nil) // Ensure LSL implements the ColorModel interface.

//////////////////////////////////////////////////////
// Constructors
//////////////////////////////////////////////////////

// NewLSL creates a new LSL instance.
// Wavelength is in nanometers [380-750]
// Saturation and Lightness are in range [0,1]
// Alpha is in range [0,1]
func NewLSL[N math.Number](wavelength, saturation, lightness, alpha N) (*LSL, error) {
	return newColor(func() *LSL { return &LSL{} }, wavelength, saturation, lightness, alpha)
}

// LSLFromRGB converts an RGBA64 (RGB) to an LSL color.
func LSLFromRGB(c *RGBA64) *LSL {
	// Convert RGB to HSL
	h, s, l := utils.RGBToHSL(c.R, c.G, c.B)

	// For black and white (saturation = 0), use the shortest wavelength
	if s < constants.WhiteThreshold {
		return &LSL{
			Wavelength: constants.WavelengthVioletMin, // Use shortest wavelength for black/white
			Saturation: s,
			Lightness:  l,
			Alpha:      c.A,
		}
	}

	// Convert hue to wavelength
	w := utils.HueToWavelength(h)

	return &LSL{
		Wavelength: w,
		Saturation: s,
		Lightness:  l,
		Alpha:      c.A,
	}
}

//////////////////////////////////////////////////////
// Type
//////////////////////////////////////////////////////

// LSL is a helper struct representing a color in the wavelength-based color model.
type LSL struct {
	Wavelength float64 // in nanometers [380-750]
	Saturation float64 // [0-1]
	Lightness  float64 // [0-1]
	Alpha      float64 // [0-1]
}

func (l *LSL) Meta() *ColorModelMeta {
	return NewModelMeta(
		"λSL",
		"Wavelength-based color model (physical spectrum).",
		NewChannelMeta("λ", constants.WavelengthVioletMin, constants.WavelengthRedMin, "nm", "Wavelength in nanometers."),
		NewChannelMeta("S", 0, 1, "", "Saturation."),
		NewChannelMeta("L", 0, 1, "", "Lightness."),
		NewChannelMeta("Alpha", 0, 1, "", "Alpha channel."),
	)
}

//////////////////////////////////////////////////////
// Conversion
//////////////////////////////////////////////////////

func (l *LSL) ToRGB() *RGBA64 {
	// For black and white (saturation = 0), wavelength doesn't matter
	if l.Saturation < constants.WhiteThreshold {
		return &RGBA64{
			R: l.Lightness,
			G: l.Lightness,
			B: l.Lightness,
			A: l.Alpha,
		}
	}

	// Convert wavelength to hue
	h := utils.WavelengthToHue(l.Wavelength)

	// Convert back to RGB
	r, g, b := utils.HSLToRGB(h, l.Saturation, l.Lightness)

	return &RGBA64{
		R: r,
		G: g,
		B: b,
		A: l.Alpha,
	}
}

// FromSlice initializes a LSL instance from a slice of float64 values.
// The slice must contain exactly 4 values in the order: Wavelength, Saturation, Lightness, Alpha.
func (l *LSL) FromSlice(vals []float64) error {
	if len(vals) != 4 {
		return fmt.Errorf("LSL requires 4 values, got %d", len(vals))
	}

	l.Wavelength = vals[0]
	l.Saturation = vals[1]
	l.Lightness = vals[2]
	l.Alpha = vals[3]

	return nil
}

// FromRGBA64 converts an RGBA64 color to this color model.
func (l *LSL) FromRGBA64(rgba *RGBA64) iColor {
	return LSLFromRGB(rgba)
}

// ToRGBA64 converts the color to RGBA64.
func (l *LSL) ToRGBA64() *RGBA64 {
	return l.ToRGB()
}
