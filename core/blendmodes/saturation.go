package blendmodes

import (
	"github.com/toxyl/gfx/core/blendmodes/constants"
	"github.com/toxyl/gfx/core/color"
)

// Saturation blend mode implementation
type Saturation struct {
	BlendMode
}

// NewSaturation creates a new Saturation blend mode
func NewSaturation() *Saturation {
	return &Saturation{
		BlendMode: BlendMode{
			name: "saturation",
		},
	}
}

// Blend implements the saturation blend mode
func (s *Saturation) Blend(base, blend *color.RGBA64) *color.RGBA64 {
	if base == nil || blend == nil {
		return nil
	}

	// Process both colors in linear space
	baseLinear := base.Process(false, true, func(c *color.RGBA64) {})
	blendLinear := blend.Process(false, true, func(c *color.RGBA64) {})

	// Convert to HSL
	baseHSL := color.HSLFromRGB(baseLinear)
	blendHSL := color.HSLFromRGB(blendLinear)

	// Check if base color is grayscale (very low saturation)
	if baseHSL.S < 0.0001 {
		// For grayscale colors, preserve the base color
		return baseLinear.Process(false, false, func(c *color.RGBA64) {})
	}

	// Create result with base hue and lightness, blend saturation
	result := &color.HSL{
		H:     baseHSL.H,
		S:     blendHSL.S,
		L:     baseHSL.L,
		Alpha: baseLinear.A,
	}

	// Convert back to RGBA64 and then to sRGB space
	return result.ToRGBA64().Process(false, false, func(c *color.RGBA64) {})
}

// Register the blend mode
func init() {
	Register("saturation", "Preserves the hue and lightness of the base color while using the saturation of the blend color", constants.CategoryComponent, func(bottom, top *color.RGBA64, alpha float64) *color.RGBA64 {
		if bottom == nil || top == nil {
			return nil
		}

		// Process both colors in linear space
		bottomLinear := bottom.Process(false, true, func(c *color.RGBA64) {})
		topLinear := top.Process(false, true, func(c *color.RGBA64) {})

		// Convert to HSL
		bottomHSL := color.HSLFromRGB(bottomLinear)
		topHSL := color.HSLFromRGB(topLinear)

		// Check if base color is grayscale (very low saturation)
		if bottomHSL.S < 0.0001 {
			// For grayscale colors, preserve the base color
			return bottomLinear.Process(false, false, func(c *color.RGBA64) {})
		}

		// Create result with base hue and lightness, blend saturation
		result := &color.HSL{
			H:     bottomHSL.H,
			S:     topHSL.S,
			L:     bottomHSL.L,
			Alpha: bottomLinear.A,
		}

		if alpha < 1.0 {
			// Interpolate saturation based on alpha
			result.S = topHSL.S*alpha + bottomHSL.S*(1-alpha)
		}

		// Convert back to RGBA64 and then to sRGB space
		rgba := result.ToRGBA64()
		return rgba.Process(false, false, func(c *color.RGBA64) {})
	})
}
