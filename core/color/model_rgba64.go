// Package color provides a comprehensive set of color models and conversion utilities.
// This file implements the RGBA64 color model with stateful operations.
package color

import (
	"fmt"
	"image"
	"image/color"

	"github.com/toxyl/math"
)

// RGBA64 represents a color with float64 channels in the range [0,1].
// The default state is premultiplied alpha in the sRGB color space.
// No input validation (clamping) is performed during processing, so operations
// may temporarily exceed the [0,1] range. Conversions to image/color clamp values.
type RGBA64 struct {
	R, G, B, A      float64 // The normalized color channels in the range [0,1].
	isLinear        bool    // isLinear indicates whether the color is in linear space. If false, the color is in sRGB space.
	isPremultiplied bool    // isPremultiplied indicates whether the RGB channels have been premultiplied by alpha.
}

// NewRGBA64 creates a new RGBA64 color from the given values.
// The values are normalized to the [0,1] range.
func NewRGBA64[N math.Number](r, g, b, a N) (*RGBA64, error) {
	return newColor(func() *RGBA64 { return &RGBA64{} }, r, g, b, a)
}

// FromSlice initializes the color from a slice of float64 values.
func (c *RGBA64) FromSlice(values []float64) error {
	if len(values) != 4 {
		return fmt.Errorf("RGBA64 requires exactly 4 values, got %d", len(values))
	}
	c.R = values[0]
	c.G = values[1]
	c.B = values[2]
	c.A = values[3]
	return nil
}

// Meta returns metadata about the RGBA64 color model.
func (c *RGBA64) Meta() *ColorModelMeta {
	return NewModelMeta(
		"RGBA64",
		"RGBA color model with 64-bit float precision.",
		NewChannelMeta("R", 0, 1, "", "Red channel."),
		NewChannelMeta("G", 0, 1, "", "Green channel."),
		NewChannelMeta("B", 0, 1, "", "Blue channel."),
		NewChannelMeta("A", 0, 1, "", "Alpha channel."),
	)
}

// Copy creates and returns a deep copy of the current RGBA64 color.
func (c *RGBA64) Copy() *RGBA64 {
	return &RGBA64{
		R:               c.R,
		G:               c.G,
		B:               c.B,
		A:               c.A,
		isLinear:        c.isLinear,
		isPremultiplied: c.isPremultiplied,
	}
}

// String returns a string representation of the color.
func (c *RGBA64) String() string {
	return fmt.Sprintf("r: %f, g: %f, b: %f, a: %f, srgb: %v, pre: %v", c.R, c.G, c.B, c.A, !c.isLinear, c.isPremultiplied)
}

// To8bit converts the normalized color to an 8-bit per channel color.RGBA.
// The color is converted to sRGB space before clamping and quantization.
func (c *RGBA64) To8bit() color.RGBA {
	copy := c.Copy()
	copy.toSRGB().Clamp() // Convert to sRGB and clamp
	return color.RGBA{
		R: uint8(copy.R * math.MaxUint8),
		G: uint8(copy.G * math.MaxUint8),
		B: uint8(copy.B * math.MaxUint8),
		A: uint8(copy.A * math.MaxUint8),
	}
}

// To16bit converts the normalized color to a 16-bit per channel color.RGBA64.
// The color is converted to sRGB space before clamping and quantization.
func (c *RGBA64) To16bit() color.RGBA64 {
	copy := c.Copy()
	copy.toSRGB().Clamp() // Convert to sRGB and clamp
	return color.RGBA64{
		R: uint16(copy.R * math.MaxUint16),
		G: uint16(copy.G * math.MaxUint16),
		B: uint16(copy.B * math.MaxUint16),
		A: uint16(copy.A * math.MaxUint16),
	}
}

// Clamp clamps all channels to valid ranges.
func (c *RGBA64) Clamp() *RGBA64 {
	c.A = math.Clamp(c.A, 0.0, 1.0)
	a := c.A
	if c.isPremultiplied {
		// Premultiplied RGB must be ≤ A and ≥ 0
		c.R = math.Clamp(c.R, 0.0, a)
		c.G = math.Clamp(c.G, 0.0, a)
		c.B = math.Clamp(c.B, 0.0, a)
	} else {
		c.R = math.Clamp(c.R, 0.0, 1.0)
		c.G = math.Clamp(c.G, 0.0, 1.0)
		c.B = math.Clamp(c.B, 0.0, 1.0)
	}
	return c
}

// unpremultiply converts a premultiplied color to an unpremultiplied color.
func (c *RGBA64) unpremultiply() *RGBA64 {
	if !c.isPremultiplied {
		return c
	}
	c.isPremultiplied = false
	a := c.A
	if a == 0 {
		c.R, c.G, c.B = 0, 0, 0
		return c
	}
	c.R /= a
	c.G /= a
	c.B /= a
	return c
}

// repremultiply converts an unpremultiplied color to a premultiplied color.
func (c *RGBA64) repremultiply() *RGBA64 {
	if c.isPremultiplied {
		return c
	}
	c.isPremultiplied = true
	a := c.A
	if a == 0 {
		c.R, c.G, c.B = 0, 0, 0
		return c
	}
	c.R *= a
	c.G *= a
	c.B *= a
	return c
}

// toLinear converts the color to linear space.
func (c *RGBA64) toLinear() *RGBA64 {
	if c.isLinear {
		return c
	}
	origPremultiplied := c.isPremultiplied
	c.unpremultiply()

	c.R = srgbToLinear(c.R)
	c.G = srgbToLinear(c.G)
	c.B = srgbToLinear(c.B)

	c.isLinear = true

	if origPremultiplied {
		c.repremultiply()
	}
	return c
}

// toSRGB converts the color to sRGB space.
func (c *RGBA64) toSRGB() *RGBA64 {
	if !c.isLinear {
		return c
	}
	origPremultiplied := c.isPremultiplied
	c.unpremultiply()

	c.R = linearToSrgb(c.R)
	c.G = linearToSrgb(c.G)
	c.B = linearToSrgb(c.B)

	c.isLinear = false

	if origPremultiplied {
		c.repremultiply()
	}
	return c
}

// Process converts the color to a specified working state, applies a processing
// function, and then restores the original color space and premultiplication state.
func (c *RGBA64) Process(unpremultiplied, linearSpace bool, fn func(c *RGBA64)) *RGBA64 {
	if fn == nil {
		panic("Process function: provided processing function is nil")
	}

	origPremultiplied := c.isPremultiplied
	origLinear := c.isLinear

	defer func() {
		if origLinear {
			c.toLinear()
		} else {
			c.toSRGB()
		}
		if origPremultiplied {
			c.repremultiply()
		} else {
			c.unpremultiply()
		}
	}()

	if unpremultiplied {
		c.unpremultiply()
	} else {
		c.repremultiply()
	}
	if linearSpace {
		c.toLinear()
	} else {
		c.toSRGB()
	}

	fn(c)

	return c
}

// Set sets the pixel at (x, y) in an image.RGBA to the specified RGBA64 color.
func Set(i *image.RGBA, x, y int, rgba *RGBA64) {
	i.SetRGBA64(x, y, rgba.To16bit())
}

// Get retrieves the RGBA64 color from the pixel at (x, y) in an image.RGBA.
func Get(i *image.RGBA, x, y int) *RGBA64 {
	r, g, b, a := i.RGBA64At(x, y).RGBA()
	return &RGBA64{
		R:               float64(r) / float64(math.MaxUint16),
		G:               float64(g) / float64(math.MaxUint16),
		B:               float64(b) / float64(math.MaxUint16),
		A:               float64(a) / float64(math.MaxUint16),
		isLinear:        false,
		isPremultiplied: true,
	}
}

// FromRGBA64 converts an RGBA64 color to this color model.
func (c *RGBA64) FromRGBA64(rgba *RGBA64) iColor {
	return rgba.Copy()
}

// ToRGBA64 converts the color to RGBA64.
func (c *RGBA64) ToRGBA64() *RGBA64 {
	return c.Copy()
}

// SetUint8 sets the color using 8-bit components.
// The input values are automatically normalized to the range [0,1].
//
// Parameters:
//   - col: A standard library color.RGBA instance
func (c *RGBA64) SetUint8(col color.RGBA) {
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
func (c *RGBA64) SetUint16(col color.RGBA64) {
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
func (c *RGBA64) Set(r, g, b, a float64) {
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
func (c *RGBA64) ToUint8() color.RGBA {
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
func (c *RGBA64) ToUint16() color.RGBA64 {
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
//   - A copy of the current RGBA64 instance
func (c *RGBA64) Get() *RGBA64 {
	return c.Copy()
}
