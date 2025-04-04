// core/color/base_lambdasl.go
package color

import (
	"fmt"

	"github.com/toxyl/gfx/core/color/constants"
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
// It accepts wavelength in nanometers [380-750], saturation [0-1], lightness [0-1], and alpha [0-1].
func NewLSL[N math.Number](wavelength, saturation, lightness, alpha N) (*LSL, error) {
	return newColor(func() *LSL { return &LSL{} }, wavelength, saturation, lightness, alpha)
}

// LSLFromRGB converts an RGBA64 (RGB) to a LSL color.
// This is an approximation since RGB can represent colors not in the visible spectrum.
func LSLFromRGB(c *RGBA64) *LSL {
	// Convert RGB to wavelength (approximate)
	// This is a simplified version - in reality, RGB can represent colors
	// that don't correspond to a single wavelength
	h, s, l := rgbToHsl(c.R, c.G, c.B)

	// For black and white (saturation = 0), use the shortest wavelength
	if s < constants.LSL_Threshold {
		return &LSL{
			Wavelength: constants.LSL_VioletMin, // Use shortest wavelength for black/white
			Saturation: s,
			Lightness:  l,
			Alpha:      c.A,
		}
	}

	// Map hue to wavelength using a piecewise linear function
	// The mapping is based on the visible spectrum:
	// Red: 620-750nm
	// Orange: 590-620nm
	// Yellow: 570-590nm
	// Green: 495-570nm
	// Blue: 450-495nm
	// Violet: 380-450nm
	var w float64
	if h < constants.LSL_HueYellow { // Red
		w = constants.LSL_RedMin + (h/constants.LSL_HueYellow)*(constants.LSL_YellowMin-constants.LSL_RedMin)
	} else if h < constants.LSL_HueGreen { // Yellow
		w = constants.LSL_YellowMin + ((h-constants.LSL_HueYellow)/constants.LSL_HueRange)*(constants.LSL_GreenMin-constants.LSL_YellowMin)
	} else if h < constants.LSL_HueCyan { // Green
		w = constants.LSL_GreenMin + ((h-constants.LSL_HueGreen)/constants.LSL_HueRange)*(constants.LSL_CyanMin-constants.LSL_GreenMin)
	} else if h < constants.LSL_HueBlue { // Cyan
		w = constants.LSL_CyanMin + ((h-constants.LSL_HueCyan)/constants.LSL_HueRange)*(constants.LSL_BlueMin-constants.LSL_CyanMin)
	} else if h < constants.LSL_HueMagenta { // Blue
		w = constants.LSL_BlueMin + ((h-constants.LSL_HueBlue)/constants.LSL_HueRange)*(constants.LSL_VioletMin-constants.LSL_BlueMin)
	} else { // Magenta
		w = constants.LSL_VioletMin + ((h-constants.LSL_HueMagenta)/(360-constants.LSL_HueMagenta))*(constants.LSL_RedMin-constants.LSL_VioletMin)
	}

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
		NewChannelMeta("λ", constants.WavelengthMin, constants.WavelengthMax, "nm", "Wavelength in nanometers."),
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
	if l.Saturation < constants.LSL_Threshold {
		return &RGBA64{
			R: l.Lightness,
			G: l.Lightness,
			B: l.Lightness,
			A: l.Alpha,
		}
	}

	// Convert wavelength to hue using inverse mapping
	var h float64
	w := l.Wavelength

	// Map wavelength to hue using inverse of the above mapping
	if w >= constants.LSL_RedMin { // Red
		h = constants.LSL_HueRed // Pure red
	} else if w >= constants.LSL_YellowMin { // Yellow
		h = constants.LSL_HueYellow + (w-constants.LSL_YellowMin)/(constants.LSL_RedMin-constants.LSL_YellowMin)*constants.LSL_HueRange
	} else if w >= constants.LSL_GreenMin { // Green
		h = constants.LSL_HueGreen // Pure green
	} else if w >= constants.LSL_CyanMin { // Cyan
		h = constants.LSL_HueCyan + (w-constants.LSL_CyanMin)/(constants.LSL_GreenMin-constants.LSL_CyanMin)*constants.LSL_HueRange
	} else if w >= constants.LSL_BlueMin { // Blue
		h = constants.LSL_HueBlue // Pure blue
	} else { // Magenta
		h = constants.LSL_HueMagenta + (w-constants.LSL_VioletMin)/(constants.LSL_BlueMin-constants.LSL_VioletMin)*constants.LSL_HueRange
	}

	// Convert back to RGB
	r, g, b := hslToRgb(h, l.Saturation, l.Lightness)

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
		return fmt.Errorf("LSL requires exactly 4 values: wavelength, saturation, lightness, alpha")
	}

	l.Wavelength = vals[0]
	l.Saturation = vals[1]
	l.Lightness = vals[2]
	l.Alpha = vals[3]

	return nil
}

func (l *LSL) FromRGBA64(rgba *RGBA64) iColor {
	return LSLFromRGB(rgba)
}

func (l *LSL) ToRGBA64() *RGBA64 {
	return l.ToRGB()
}
