package blendmodes

import (
	"github.com/toxyl/gfx/core/blendmodes/constants"
	"github.com/toxyl/gfx/core/color"
)

// Luminosity blend mode implementation.
// This blend mode takes the lightness (luminosity) from the blend color while preserving
// the hue and saturation from the base color. It works in HSL color space for accurate
// luminosity calculations.
type Luminosity struct {
	BlendMode
}

// NewLuminosity creates a new Luminosity blend mode.
func NewLuminosity() *Luminosity {
	return &Luminosity{
		BlendMode: BlendMode{
			name: "luminosity",
		},
	}
}

// blendLuminosity is a helper function that performs the actual luminosity blend operation.
// It takes two colors in linear space and returns the blended result.
func blendLuminosity(base, blend *color.RGBA64, alpha float64) *color.RGBA64 {
	// Convert to HSL
	baseHSL := color.HSLFromRGB(base)
	blendHSL := color.HSLFromRGB(blend)

	// Interpolate lightness based on alpha
	finalL := baseHSL.L*(1-alpha) + blendHSL.L*alpha

	// Take lightness from blend color, keep hue and saturation from base color
	result := &color.HSL{
		H:     baseHSL.H,
		S:     baseHSL.S,
		L:     finalL,
		Alpha: base.A,
	}

	// Convert back to RGBA64
	return result.ToRGBA64()
}

// Blend implements the luminosity blend mode.
func (l *Luminosity) Blend(base, blend *color.RGBA64) *color.RGBA64 {
	// Process both colors in linear space
	baseLinear := base.Process(false, true, func(c *color.RGBA64) {})
	blendLinear := blend.Process(false, true, func(c *color.RGBA64) {})

	// Perform the blend operation
	result := blendLuminosity(baseLinear, blendLinear, 1.0)

	// Convert back to sRGB space
	return result.Process(false, false, func(c *color.RGBA64) {})
}

// Register the blend mode
func init() {
	Register("luminosity", "Preserves the hue and saturation of the base color while using the luminosity of the blend color", constants.CategoryComponent, func(bottom, top *color.RGBA64, alpha float64) *color.RGBA64 {
		// Process both colors in linear space
		bottomLinear := bottom.Process(false, true, func(c *color.RGBA64) {})
		topLinear := top.Process(false, true, func(c *color.RGBA64) {})

		// Perform the blend operation
		result := blendLuminosity(bottomLinear, topLinear, alpha)

		// Convert back to sRGB space
		return result.Process(false, false, func(c *color.RGBA64) {})
	})
}
