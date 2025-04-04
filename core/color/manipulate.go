package color

import (
	"fmt"
	"math"
	"strings"
)

// Lighten lightens a color by the specified amount
func (c *RGB8) Lighten(amount float64) iColor {
	hsl := HSLFromRGB(c.ToRGBA64())
	hsl.L = math.Min(1.0, hsl.L+amount)
	return hsl.ToRGB()
}

// Darken darkens a color by the specified amount
func (c *RGB8) Darken(amount float64) iColor {
	hsl := HSLFromRGB(c.ToRGBA64())
	hsl.L = math.Max(0.0, hsl.L-amount)
	return hsl.ToRGB()
}

// Saturate increases the saturation of a color by the specified amount
func (c *RGB8) Saturate(amount float64) iColor {
	hsl := HSLFromRGB(c.ToRGBA64())
	hsl.S = math.Min(1.0, hsl.S+amount)
	return hsl.ToRGB()
}

// Desaturate decreases the saturation of a color by the specified amount
func (c *RGB8) Desaturate(amount float64) iColor {
	hsl := HSLFromRGB(c.ToRGBA64())
	hsl.S = math.Max(0.0, hsl.S-amount)
	return hsl.ToRGB()
}

// RotateHue rotates the hue of a color by the specified number of degrees
func (c *RGB8) RotateHue(degrees float64) iColor {
	hsl := HSLFromRGB(c.ToRGBA64())
	hsl.H = math.Mod(hsl.H+degrees, 360)
	if hsl.H < 0 {
		hsl.H += 360
	}
	return hsl.ToRGB()
}

// Mix mixes two colors by the specified ratio
func (c *RGB8) Mix(other iColor, ratio float64) iColor {
	otherRGB := other.ToRGBA64()

	r := c.R*(1-ratio) + otherRGB.R*ratio
	g := c.G*(1-ratio) + otherRGB.G*ratio
	b := c.B*(1-ratio) + otherRGB.B*ratio
	alpha := c.Alpha*(1-ratio) + otherRGB.A*ratio

	return &RGB8{
		R:     r,
		G:     g,
		B:     b,
		Alpha: alpha,
	}
}

// Equals checks if two colors are equal
func (c *RGB8) Equals(other iColor) bool {
	otherRGB := other.ToRGBA64()
	return c.R == otherRGB.R &&
		c.G == otherRGB.G &&
		c.B == otherRGB.B &&
		c.Alpha == otherRGB.A
}

// Distance calculates the Euclidean distance between two colors
func (c *RGB8) Distance(other iColor) float64 {
	otherRGB := other.ToRGBA64()
	dr := c.R - otherRGB.R
	dg := c.G - otherRGB.G
	db := c.B - otherRGB.B
	return math.Sqrt(dr*dr + dg*dg + db*db)
}

// IsSimilar checks if two colors are similar within a tolerance
func (c *RGB8) IsSimilar(other iColor, tolerance float64) bool {
	return c.Distance(other) <= tolerance
}

// GetChannel gets the value of a specific channel
func (c *RGB8) GetChannel(name string) (float64, error) {
	switch strings.ToLower(name) {
	case "r":
		return c.R, nil
	case "g":
		return c.G, nil
	case "b":
		return c.B, nil
	case "alpha":
		return c.Alpha, nil
	default:
		return 0, fmt.Errorf("unknown channel: %s", name)
	}
}

// SetChannel sets the value of a specific channel
func (c *RGB8) SetChannel(name string, value float64) error {
	switch strings.ToLower(name) {
	case "r":
		c.R = value
	case "g":
		c.G = value
	case "b":
		c.B = value
	case "alpha":
		c.Alpha = value
	default:
		return fmt.Errorf("unknown channel: %s", name)
	}
	return nil
}

// Channels returns a map of all channel values
func (c *RGB8) Channels() map[string]float64 {
	return map[string]float64{
		"r":     c.R,
		"g":     c.G,
		"b":     c.B,
		"alpha": c.Alpha,
	}
}
