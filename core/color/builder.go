package color

import (
	"fmt"
)

// ColorBuilder provides a fluent interface for creating colors
type ColorBuilder struct {
	model  string
	values []float64
	alpha  float64
}

// NewColor creates a new ColorBuilder instance
func NewColor() *ColorBuilder {
	return &ColorBuilder{
		alpha: 1.0, // Default to fully opaque
	}
}

// RGB sets the color to RGB values
func (cb *ColorBuilder) RGB(r, g, b float64) *ColorBuilder {
	cb.model = "RGB"
	cb.values = []float64{r, g, b}
	return cb
}

// LAB sets the color to LAB values
func (cb *ColorBuilder) LAB(l, a, b float64) *ColorBuilder {
	cb.model = "LAB"
	cb.values = []float64{l, a, b}
	return cb
}

// HSL sets the color to HSL values
func (cb *ColorBuilder) HSL(h, s, l float64) *ColorBuilder {
	cb.model = "HSL"
	cb.values = []float64{h, s, l}
	return cb
}

// Alpha sets the alpha value
func (cb *ColorBuilder) Alpha(alpha float64) *ColorBuilder {
	cb.alpha = alpha
	return cb
}

// Build creates the color based on the builder configuration
func (cb *ColorBuilder) Build() (iColor, error) {
	if len(cb.values) == 0 {
		return nil, fmt.Errorf("no color values specified")
	}

	switch cb.model {
	case "RGB":
		return NewRGB8(cb.values[0], cb.values[1], cb.values[2], cb.alpha)
	case "LAB":
		return NewLAB(cb.values[0], cb.values[1], cb.values[2], cb.alpha)
	case "HSL":
		return NewHSL(cb.values[0], cb.values[1], cb.values[2], cb.alpha)
	default:
		return nil, fmt.Errorf("unsupported color model: %s", cb.model)
	}
}

// ColorManipulator provides methods for manipulating colors
type ColorManipulator interface {
	Lighten(amount float64) iColor
	Darken(amount float64) iColor
	Saturate(amount float64) iColor
	Desaturate(amount float64) iColor
	RotateHue(degrees float64) iColor
	Mix(other iColor, ratio float64) iColor
	Equals(other iColor) bool
	Distance(other iColor) float64
	IsSimilar(other iColor, tolerance float64) bool
	GetChannel(name string) (float64, error)
	SetChannel(name string, value float64) error
	Channels() map[string]float64
}

// ColorFormatter provides methods for formatting colors
type ColorFormatter interface {
	Format(format string) (string, error)
	ToHex() string
	ToRGBString() string
	ToHSLString() string
}

// ExtendedColorModel extends the base iColor interface with additional functionality
type ExtendedColorModel interface {
	iColor
	ColorManipulator
	ColorFormatter
}
