// core/color/main.go
// Package color provides a comprehensive set of color models and conversion utilities
// for working with colors in various color spaces. It supports both device-dependent
// (RGB, HSL, HSB, etc.) and device-independent (LAB, XYZ, etc.) color models.
//
// The package is built around the Color64 type, which serves as the canonical
// representation of colors in RGBA format with 64-bit float precision. All color
// models can be converted to and from Color64, enabling seamless conversion
// between different color spaces.
//
// Features:
//   - Support for multiple color models:
//   - RGB-based: RGB8, Hex, Grayscale
//   - HSL/HSV-based: HSL, HSB, LSL, LSB
//   - Device-independent: LAB, XYZ, LCH, HCL
//   - Video/Television: YUV, YIQ, YCbCr
//   - Printing: CMY, CMYK
//   - Consistent alpha channel handling across all models
//   - Comprehensive metadata for each color channel
//   - Range validation and automatic clamping
//   - Efficient conversion utilities
//   - Type-safe implementations
//
// Usage:
//
//	// Create a new color in RGB space
//	rgb := color.New(0.5, 0.3, 0.8, 1.0)
//
//	// Convert to HSL
//	hsl := color.HSLFromRGB(rgb)
//
//	// Convert to LAB
//	lab := color.LABFromRGB(rgb)
//
//	// Convert back to RGB
//	rgb2 := lab.ToRGB()
//
// Each color model provides:
//   - Constructors for creating new colors
//   - Conversion methods to/from RGB
//   - Channel metadata (ranges, units, descriptions)
//   - Validation and error handling
//
// Error Handling:
//   - All constructors return (ColorModel, error) pairs
//   - Input values are validated against channel ranges
//   - Conversion errors are properly propagated
//   - Invalid operations return descriptive error messages
//
// Channel Ranges:
//   - RGB: [0, 1] for all channels
//   - HSL/HSB: H [0, 360], S/L/B [0, 1]
//   - LAB: L [0, 100], a/b [-128, 127]
//   - XYZ: [0, 1] for all channels
//   - Alpha: [0, 1] across all models
//
// Note: All color conversions use standard conversion matrices and formulas
// as defined by the International Commission on Illumination (CIE) and other
// relevant standards organizations.
package color

import (
	"fmt"
	"image/color"

	"github.com/toxyl/math"
)

// Color64 represents a color in RGBA format with 64-bit float precision.
// All channels are normalized to the range [0,1], where:
//   - R, G, B: Red, Green, Blue components
//   - A: Alpha (transparency) component
//
// This is the canonical representation used for all color conversions.
// All other color models in this package can be converted to and from Color64.
type Color64 struct {
	R, G, B, A float64
}

// New creates a new Color64 from RGBA components.
// All components are normalized to the range [0,1].
//
// Parameters:
//   - r, g, b: Red, Green, Blue components [0,1]
//   - a: Alpha (transparency) component [0,1]
//
// Returns:
//   - A pointer to a new Color64 instance
func New[N math.Number](r, g, b, a N) *Color64 {
	return &Color64{
		R: math.Clamp(float64(r), 0, 1),
		G: math.Clamp(float64(g), 0, 1),
		B: math.Clamp(float64(b), 0, 1),
		A: math.Clamp(float64(a), 0, 1),
	}
}

// SetUint8 sets the color using 8-bit components.
// The input values are automatically normalized to the range [0,1].
//
// Parameters:
//   - col: A standard library color.RGBA instance
func (c *Color64) SetUint8(col color.RGBA) {
	c.R = math.Clamp(float64(col.R)/math.MaxUint8, 0, 1)
	c.G = math.Clamp(float64(col.G)/math.MaxUint8, 0, 1)
	c.B = math.Clamp(float64(col.B)/math.MaxUint8, 0, 1)
	c.A = math.Clamp(float64(col.A)/math.MaxUint8, 0, 1)
}

// SetUint16 sets the color using 16-bit components.
// The input values are automatically normalized to the range [0,1].
//
// Parameters:
//   - col: A standard library color.RGBA64 instance
func (c *Color64) SetUint16(col color.RGBA64) {
	c.R = math.Clamp(float64(col.R)/math.MaxUint16, 0, 1)
	c.G = math.Clamp(float64(col.G)/math.MaxUint16, 0, 1)
	c.B = math.Clamp(float64(col.B)/math.MaxUint16, 0, 1)
	c.A = math.Clamp(float64(col.A)/math.MaxUint16, 0, 1)
}

// Set sets the color using 64-bit float components.
// The input values are automatically clamped to the range [0,1].
//
// Parameters:
//   - r, g, b: Red, Green, Blue components [0,1]
//   - a: Alpha (transparency) component [0,1]
func (c *Color64) Set(r, g, b, a float64) {
	c.R = math.Clamp(r, 0, 1)
	c.G = math.Clamp(g, 0, 1)
	c.B = math.Clamp(b, 0, 1)
	c.A = math.Clamp(a, 0, 1)
}

// ToUint8 converts the color to 8-bit RGBA format.
// The color channels are automatically clamped to [0,1] before conversion.
//
// Returns:
//   - A standard library color.RGBA instance
func (c *Color64) ToUint8() color.RGBA {
	return color.RGBA{
		R: uint8(math.Clamp(c.R, 0, 1) * math.MaxUint8),
		G: uint8(math.Clamp(c.G, 0, 1) * math.MaxUint8),
		B: uint8(math.Clamp(c.B, 0, 1) * math.MaxUint8),
		A: uint8(math.Clamp(c.A, 0, 1) * math.MaxUint8),
	}
}

// ToUint16 converts the color to 16-bit RGBA format.
// The color channels are automatically clamped to [0,1] before conversion.
//
// Returns:
//   - A standard library color.RGBA64 instance
func (c *Color64) ToUint16() color.RGBA64 {
	return color.RGBA64{
		R: uint16(math.Clamp(c.R, 0, 1) * math.MaxUint16),
		G: uint16(math.Clamp(c.G, 0, 1) * math.MaxUint16),
		B: uint16(math.Clamp(c.B, 0, 1) * math.MaxUint16),
		A: uint16(math.Clamp(c.A, 0, 1) * math.MaxUint16),
	}
}

// Get returns a copy of the current color.
// The returned color is a new instance with the same channel values.
//
// Returns:
//   - A copy of the current Color64 instance
func (c *Color64) Get() Color64 {
	return *c
}

// String returns a string representation of the color in RGBA format.
// The values are formatted with 6 decimal places of precision.
//
// Returns:
//   - A string in the format "RGBA(r: X.XXXXXX, g: X.XXXXXX, b: X.XXXXXX, a: X.XXXXXX)"
func (c *Color64) String() string {
	return fmt.Sprintf("RGBA(r: %.6f, g: %.6f, b: %.6f, a: %.6f)", c.R, c.G, c.B, c.A)
}
