// core/color/base_hsb.go
package color

import (
	"fmt"

	"github.com/toxyl/math"
)

//////////////////////////////////////////////////////
// Conversion utilities
//////////////////////////////////////////////////////

func rgbToHsb(r, g, b float64) (h, s, v float64) {
	max := math.Max(r, math.Max(g, b))
	min := math.Min(r, math.Min(g, b))

	// Calculate value (brightness)
	v = max

	// If max is 0, it's black
	if max == 0 {
		return 0, 0, 0
	}

	// Calculate saturation
	s = (max - min) / max

	// If max and min are equal, it's a shade of gray
	if max == min {
		return 0, 0, v
	}

	// Calculate hue in degrees
	switch max {
	case r:
		h = (g - b) / (max - min)
		if g < b {
			h += 6
		}
	case g:
		h = (b-r)/(max-min) + 2
	case b:
		h = (r-g)/(max-min) + 4
	}
	h *= 60 // Convert to degrees

	return h, s, v
}

func hsbToRgb(h, s, v float64) (r, g, b float64) {
	if s == 0 {
		return v, v, v
	}

	// Convert hue to normalized value
	h = h / 360.0

	i := math.Floor(h * 6)
	f := h*6 - i
	p := v * (1 - s)
	q := v * (1 - s*f)
	t := v * (1 - s*(1-f))

	switch int(i) % 6 {
	case 0:
		r, g, b = v, t, p
	case 1:
		r, g, b = q, v, p
	case 2:
		r, g, b = p, v, t
	case 3:
		r, g, b = p, q, v
	case 4:
		r, g, b = t, p, v
	case 5:
		r, g, b = v, p, q
	}

	return r, g, b
}

//////////////////////////////////////////////////////
// Implementation check
//////////////////////////////////////////////////////

var _ iColor = (*HSB)(nil) // Ensure HSB implements the ColorModel interface.

//////////////////////////////////////////////////////
// Constructors
//////////////////////////////////////////////////////

// NewHSB creates a new HSB instance.
// It accepts channel values in the [0,1] range for H, S, B and A.
func NewHSB[N math.Number](h, s, b, alpha N) (*HSB, error) {
	return newColor(func() *HSB { return &HSB{} }, h, s, b, alpha)
}

// HSBFromRGB converts an RGBA64 (RGB) to an HSB color.
func HSBFromRGB(c *RGBA64) *HSB {
	h, s, b := rgbToHsb(c.R, c.G, c.B)
	return &HSB{
		H:     h,
		S:     s,
		B:     b,
		Alpha: c.A,
	}
}

//////////////////////////////////////////////////////
// Type
//////////////////////////////////////////////////////

// HSB is a helper struct representing a color in the HSB color model with an alpha channel.
type HSB struct {
	H, S, B, Alpha float64
}

func (hsb *HSB) Meta() *ColorModelMeta {
	return NewModelMeta(
		"HSB",
		"Hue, Saturation, Brightness color model.",
		NewChannelMeta("H", 0, 360, "°", "Hue in degrees."),
		NewChannelMeta("S", 0, 100, "%", "Saturation percentage."),
		NewChannelMeta("B", 0, 100, "%", "Brightness percentage."),
		NewChannelMeta("Alpha", 0, 1, "", "Alpha channel."),
	)
}

//////////////////////////////////////////////////////
// Conversion
//////////////////////////////////////////////////////

func (hsb *HSB) ToRGB() *RGBA64 {
	r, g, b := hsbToRgb(hsb.H, hsb.S, hsb.B)
	return &RGBA64{
		R: r,
		G: g,
		B: b,
		A: hsb.Alpha,
	}
}

// FromSlice initializes the color from a slice of float64 values.
func (hsb *HSB) FromSlice(values []float64) error {
	if len(values) != 4 {
		return fmt.Errorf("HSB requires exactly 4 values: H, S, B, Alpha")
	}

	hsb.H = values[0]
	hsb.S = values[1]
	hsb.B = values[2]
	hsb.Alpha = values[3]

	return nil
}

// FromRGBA64 converts an RGBA64 color to this color model.
func (hsb *HSB) FromRGBA64(rgba *RGBA64) iColor {
	return HSBFromRGB(rgba)
}

// ToRGBA64 converts the color to RGBA64.
func (hsb *HSB) ToRGBA64() *RGBA64 {
	return hsb.ToRGB()
}
