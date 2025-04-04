// core/color/base_lch.go
package color

import (
	"fmt"

	"github.com/toxyl/gfx/core/color/constants"
	"github.com/toxyl/math"
)

//////////////////////////////////////////////////////
// Implementation check
//////////////////////////////////////////////////////

var _ iColor = (*LCH)(nil) // Ensure LCH implements the ColorModel interface.

//////////////////////////////////////////////////////
// Conversion utilities
//////////////////////////////////////////////////////

// labToLch converts LAB to LCH color space.
func labToLch(l, a, b float64) (l2, c, h float64) {
	// Handle exact white
	if math.Abs(l-100) < constants.LCH_WhiteThreshold && math.Abs(a) < constants.LCH_WhiteThreshold && math.Abs(b) < constants.LCH_WhiteThreshold {
		return 100, 0, 0
	}

	// Handle near-white colors
	if math.Abs(l-100) < constants.LCH_NearWhiteThreshold && math.Abs(a) < constants.LCH_WhiteThreshold && math.Abs(b) < constants.LCH_WhiteThreshold {
		return l, 0, 0
	}

	// Handle gray colors (a and b close to zero)
	if math.Abs(a) < constants.LCH_WhiteThreshold && math.Abs(b) < constants.LCH_WhiteThreshold {
		return l, 0, 0
	}

	c = math.Sqrt(a*a + b*b)

	// Handle very small chroma values
	if c < constants.LCH_ChromaThreshold {
		return l, 0, 0
	}

	h = math.Atan2(b, a) * constants.LCH_DegreesPerRadian

	if h < 0 {
		h += constants.LCH_DegreesPerCircle
	}

	return l, c, h
}

// lchToLab converts LCH to LAB color space.
func lchToLab(l, c, h float64) (l2, a, b float64) {
	// Handle exact white
	if math.Abs(l-100) < constants.LCH_WhiteThreshold && math.Abs(c) < constants.LCH_WhiteThreshold {
		return 100, 0, 0
	}

	// Handle near-white colors
	if math.Abs(l-100) < constants.LCH_NearWhiteThreshold && math.Abs(c) < constants.LCH_WhiteThreshold {
		return l, 0, 0
	}

	// Handle gray colors (chroma close to zero)
	if math.Abs(c) < constants.LCH_ChromaThreshold {
		return l, 0, 0
	}

	h = h * constants.LCH_RadiansPerDegree
	a = c * math.Cos(h)
	b = c * math.Sin(h)

	return l, a, b
}

//////////////////////////////////////////////////////
// Constructors
//////////////////////////////////////////////////////

// NewLCH creates a new LCH instance.
// L (lightness) is in range [0,1]
// C (chroma) is in range [0,1]
// H (hue) is in range [0,360]
// Alpha is in range [0,1]
func NewLCH[N math.Number](l, c, h, alpha N) (*LCH, error) {
	return newColor(func() *LCH { return &LCH{} }, l, c, h, alpha)
}

// LCHFromRGB converts an RGBA64 (RGB) to an LCH color.
func LCHFromRGB(c *RGBA64) *LCH {
	// First convert to LAB
	x, y, z := rgbToXyz(c.R, c.G, c.B)
	l, a, b := xyzToLab(x, y, z)

	// Then convert to LCH
	l2, chroma, h := labToLch(l, a, b)

	return &LCH{
		L:     l2,
		C:     chroma,
		H:     h,
		Alpha: c.A,
	}
}

//////////////////////////////////////////////////////
// Type
//////////////////////////////////////////////////////

// LCH is a helper struct representing a color in the LCH color model.
// LCH is similar to HCL but with a different ordering of components.
type LCH struct {
	L, C, H, Alpha float64
}

func (l *LCH) Meta() *ColorModelMeta {
	return NewModelMeta(
		"LCH",
		"Lightness, Chroma, Hue color model.",
		NewChannelMeta("L", 0, 100, "", "Lightness."),
		NewChannelMeta("C", 0, 1, "", "Chroma."),
		NewChannelMeta("H", 0, 360, "°", "Hue in degrees."),
		NewChannelMeta("Alpha", 0, 1, "", "Alpha channel."),
	)
}

//////////////////////////////////////////////////////
// Conversion
//////////////////////////////////////////////////////

func (l *LCH) ToRGB() *RGBA64 {
	// First convert to LAB
	l2, a, b := lchToLab(l.L, l.C, l.H)

	// Then convert to XYZ
	x, y, z := labToXyz(l2, a, b)

	// Finally convert to RGB
	r, g, b := xyzToRgb(x, y, z)

	return &RGBA64{R: r, G: g, B: b, A: l.Alpha}
}

// FromSlice initializes an LCH instance from a slice of float64 values.
// The slice must contain exactly 4 values in the order: L, C, H, Alpha.
func (l *LCH) FromSlice(vals []float64) error {
	if len(vals) != 4 {
		return fmt.Errorf("LCH requires 4 values, got %d", len(vals))
	}

	l.L = vals[0]
	l.C = vals[1]
	l.H = vals[2]
	l.Alpha = vals[3]

	return nil
}

// FromRGBA64 converts an RGBA64 color to this color model.
func (l *LCH) FromRGBA64(rgba *RGBA64) iColor {
	return LCHFromRGB(rgba)
}

// ToRGBA64 converts the color to RGBA64.
func (l *LCH) ToRGBA64() *RGBA64 {
	return l.ToRGB()
}
