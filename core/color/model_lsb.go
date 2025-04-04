// core/color/base_lambdasb.go
package color

import (
	"fmt"

	"github.com/toxyl/gfx/core/color/constants"
	"github.com/toxyl/math"
)

//////////////////////////////////////////////////////
// Implementation check
//////////////////////////////////////////////////////

var _ iColor = (*LSB)(nil) // Ensure LSB implements the ColorModel interface.

//////////////////////////////////////////////////////
// Constructors
//////////////////////////////////////////////////////

// NewLSB creates a new LSB instance.
// It accepts wavelength in nanometers [380-750], saturation [0-1], brightness [0-1], and alpha [0-1].
func NewLSB[N math.Number](wavelength, saturation, brightness, alpha N) (*LSB, error) {
	return newColor(func() *LSB { return &LSB{} }, wavelength, saturation, brightness, alpha)
}

// LSBFromRGB converts an RGBA64 (RGB) to a LSB color.
// This is an approximation since RGB can represent colors not in the visible spectrum.
func LSBFromRGB(c *RGBA64) *LSB {
	// Convert RGB to wavelength (approximate)
	// This is a simplified version - in reality, RGB can represent colors
	// that don't correspond to a single wavelength
	h, s, b := rgbToHsb(c.R, c.G, c.B)

	// For black and white (saturation = 0), use the shortest wavelength
	if s < 1e-10 {
		return &LSB{
			Wavelength: 380, // Use shortest wavelength for black/white
			Saturation: s,
			Brightness: b,
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
	if h < 60 { // Red
		w = 700 // Pure red wavelength
	} else if h < 120 { // Yellow
		w = 570 + (590-570)*((h-60)/60)
	} else if h < 180 { // Green
		w = 520 // Pure green wavelength
	} else if h < 240 { // Cyan
		w = 495 + (520-495)*((h-180)/60)
	} else if h < 300 { // Blue
		w = 450 // Pure blue wavelength
	} else { // Magenta
		w = 380 + (450-380)*((h-300)/60)
	}

	return &LSB{
		Wavelength: w,
		Saturation: s,
		Brightness: b,
		Alpha:      c.A,
	}
}

//////////////////////////////////////////////////////
// Type
//////////////////////////////////////////////////////

// LSB is a helper struct representing a color in the wavelength-based color model.
type LSB struct {
	Wavelength float64 // in nanometers [380-750]
	Saturation float64 // [0-1]
	Brightness float64 // [0-1]
	Alpha      float64 // [0-1]
}

func (l *LSB) Meta() *ColorModelMeta {
	return NewModelMeta(
		"λSB",
		"Wavelength-based color model (physical spectrum) with brightness.",
		NewChannelMeta("λ", constants.WavelengthMin, constants.WavelengthMax, "nm", "Wavelength in nanometers."),
		NewChannelMeta("S", 0, 1, "", "Saturation."),
		NewChannelMeta("B", 0, 1, "", "Brightness."),
		NewChannelMeta("Alpha", 0, 1, "", "Alpha channel."),
	)
}

//////////////////////////////////////////////////////
// Conversion
//////////////////////////////////////////////////////

func (l *LSB) ToRGB() *RGBA64 {
	// For black and white (saturation = 0), wavelength doesn't matter
	if l.Saturation < 1e-10 {
		return &RGBA64{
			R: l.Brightness,
			G: l.Brightness,
			B: l.Brightness,
			A: l.Alpha,
		}
	}

	// Convert wavelength to hue using inverse mapping
	var h float64
	w := l.Wavelength

	// Map wavelength to hue using inverse of the above mapping
	if w >= 620 { // Red
		h = 0 // Pure red
	} else if w >= 590 { // Yellow
		h = 60 + (w-590)/(620-590)*60
	} else if w >= 520 { // Green
		h = 120 // Pure green
	} else if w >= 495 { // Cyan
		h = 180 + (w-495)/(520-495)*60
	} else if w >= 450 { // Blue
		h = 240 // Pure blue
	} else { // Magenta
		h = 300 + (w-380)/(450-380)*60
	}

	// Convert back to RGB
	r, g, b := hsbToRgb(h, l.Saturation, l.Brightness)

	return &RGBA64{
		R: r,
		G: g,
		B: b,
		A: l.Alpha,
	}
}

// FromSlice initializes a LSB instance from a slice of float64 values.
// The slice must contain exactly 4 values in the order: Wavelength, Saturation, Brightness, Alpha.
func (l *LSB) FromSlice(vals []float64) error {
	if len(vals) != 4 {
		return fmt.Errorf("expected 4 values for LSB, got %d", len(vals))
	}

	l.Wavelength = vals[0]
	l.Saturation = vals[1]
	l.Brightness = vals[2]
	l.Alpha = vals[3]

	return nil
}

// FromRGBA64 converts an RGBA64 color to this color model.
func (l *LSB) FromRGBA64(rgba *RGBA64) iColor {
	return LSBFromRGB(rgba)
}

// ToRGBA64 converts the color to RGBA64.
func (l *LSB) ToRGBA64() *RGBA64 {
	return l.ToRGB()
}
