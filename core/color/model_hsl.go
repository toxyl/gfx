// core/color/base_hsl.go
package color

import (
	"fmt"

	"github.com/toxyl/gfx/core/color/constants"
	"github.com/toxyl/math"
)

//////////////////////////////////////////////////////
// Conversion utilities
//////////////////////////////////////////////////////

func rgbToHsl(r, g, b float64) (h, s, l float64) {
	max := math.Max(r, math.Max(g, b))
	min := math.Min(r, math.Min(g, b))

	// Calculate lightness
	l = (max + min) / 2

	// If max and min are equal, it's a shade of gray
	if max == min {
		return 0, 0, l
	}

	// Calculate saturation
	if l > 0.5 {
		s = (max - min) / (2 - max - min)
	} else {
		s = (max - min) / (max + min)
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
	if h < 0 {
		h += 360
	}
	h = math.Mod(h, 360) // Ensure hue is in [0,360] range

	return h, s, l
}

func hslToRgb(h, s, l float64) (r, g, b float64) {
	if s == 0 {
		return l, l, l
	}

	// Convert hue to normalized value
	h = h * constants.HSL_HueNormalization

	var q float64
	if l < 0.5 {
		q = l * (1 + s)
	} else {
		q = l + s - l*s
	}
	p := 2*l - q

	r = hueToRgb(p, q, h+constants.HSL_HueOffset1)
	g = hueToRgb(p, q, h)
	b = hueToRgb(p, q, h-constants.HSL_HueOffset1)

	return r, g, b
}

func hueToRgb(p, q, t float64) float64 {
	if t < 0 {
		t += 1
	}
	if t > 1 {
		t -= 1
	}
	if t < constants.HSL_HueOffset3 {
		return p + (q-p)*6*t
	}
	if t < constants.HSL_HueOffset4 {
		return q
	}
	if t < constants.HSL_HueOffset5 {
		return p + (q-p)*(constants.HSL_HueOffset5-t)*6
	}
	return p
}

//////////////////////////////////////////////////////
// Implementation check
//////////////////////////////////////////////////////

var _ iColor = (*HSL)(nil) // Ensure HSL implements the ColorModel interface.

//////////////////////////////////////////////////////
// Constructors
//////////////////////////////////////////////////////

// NewHSL creates a new HSL instance.
// It accepts channel values in the [0,1] range for H, S, L and A.
func NewHSL(h, s, l, alpha float64) (*HSL, error) {
	hsl := &HSL{
		H:     h,
		S:     s,
		L:     l,
		Alpha: alpha,
	}
	return hsl, nil
}

// HSLFromRGB converts an RGBA64 (RGB) to an HSL color.
func HSLFromRGB(c *RGBA64) *HSL {
	h, s, l := rgbToHsl(c.R, c.G, c.B)
	return &HSL{
		H:     h,
		S:     s,
		L:     l,
		Alpha: c.A,
	}
}

//////////////////////////////////////////////////////
// Type
//////////////////////////////////////////////////////

// HSL is a helper struct representing a color in the HSL color model with an alpha channel.
type HSL struct {
	H, S, L, Alpha float64
}

func (hsl *HSL) Meta() *ColorModelMeta {
	return NewModelMeta(
		"HSL",
		"Hue, Saturation, Lightness color model.",
		NewChannelMeta("H", 0, 360, "°", "Hue in degrees."),
		NewChannelMeta("S", 0, 100, "%", "Saturation percentage."),
		NewChannelMeta("L", 0, 100, "%", "Lightness percentage."),
		NewChannelMeta("Alpha", 0, 1, "", "Alpha channel."),
	)
}

//////////////////////////////////////////////////////
// Conversion
//////////////////////////////////////////////////////

func (hsl *HSL) ToRGB() *RGBA64 {
	r, g, b := hslToRgb(hsl.H, hsl.S, hsl.L)
	return &RGBA64{
		R: r,
		G: g,
		B: b,
		A: hsl.Alpha,
	}
}

// FromSlice initializes the color from a slice of float64 values.
func (hsl *HSL) FromSlice(values []float64) error {
	if len(values) != 4 {
		return fmt.Errorf("HSL requires exactly 4 values: H, S, L, Alpha")
	}

	hsl.H = values[0]
	hsl.S = values[1]
	hsl.L = values[2]
	hsl.Alpha = values[3]

	return nil
}

// FromRGBA64 converts an RGBA64 color to this color model.
func (hsl *HSL) FromRGBA64(rgba *RGBA64) iColor {
	return HSLFromRGB(rgba)
}

// ToRGBA64 converts the color to RGBA64.
func (hsl *HSL) ToRGBA64() *RGBA64 {
	return hsl.ToRGB()
}
