package blendmodes

import (
	"github.com/toxyl/gfx/core/blendmodes/constants"
	"github.com/toxyl/gfx/core/color"
)

// Hue represents the hue blend mode
type Hue struct {
	BlendMode
}

// NewHue creates a new hue blend mode
func NewHue() *Hue {
	return &Hue{
		BlendMode: BlendMode{name: constants.ModeHue},
	}
}

// Blend performs the hue blend operation
func (h *Hue) Blend(base, blend *color.RGBA64) *color.RGBA64 {
	if base == nil || blend == nil {
		return nil
	}

	// If blend color has zero alpha, return base color
	if blend.A < constants.AlphaEpsilon {
		return base.Process(false, false, func(c *color.RGBA64) {})
	}

	// Process both colors in linear space
	baseLinear := base.Process(false, true, func(c *color.RGBA64) {})
	blendLinear := blend.Process(false, true, func(c *color.RGBA64) {})

	// Convert to HSL
	baseHSL := color.HSLFromRGB(baseLinear)
	blendHSL := color.HSLFromRGB(blendLinear)

	// Check if base color is grayscale (very low saturation)
	if baseHSL.S < constants.Epsilon {
		// For grayscale colors, use the blend color's RGB values
		return blendLinear.Process(false, false, func(c *color.RGBA64) {})
	}

	// Create result with blend hue, base saturation and lightness
	result := &color.HSL{
		H:     blendHSL.H,
		S:     blendHSL.S,
		L:     blendHSL.L,
		Alpha: baseLinear.A,
	}

	// Convert back to RGBA64 and then to sRGB space
	return result.ToRGBA64().Process(false, false, func(c *color.RGBA64) {})
}

// Register the blend mode
func init() {
	Register("hue", "Takes the hue from the blend color while preserving saturation and lightness from the base color", constants.CategoryComponent, func(bottom, top *color.RGBA64, alpha float64) *color.RGBA64 {
		if bottom == nil || top == nil {
			return nil
		}

		// If blend color has zero alpha, return base color
		if top.A < constants.AlphaEpsilon {
			return bottom.Process(false, false, func(c *color.RGBA64) {})
		}

		// Process both colors in linear space
		bottomLinear := bottom.Process(false, true, func(c *color.RGBA64) {})
		topLinear := top.Process(false, true, func(c *color.RGBA64) {})

		// Convert to HSL
		bottomHSL := color.HSLFromRGB(bottomLinear)
		topHSL := color.HSLFromRGB(topLinear)

		// Check if base color is grayscale (very low saturation)
		if bottomHSL.S < constants.Epsilon {
			if alpha < constants.AlphaOpaque {
				// Interpolate between base and blend colors based on alpha
				return bottomLinear.Process(false, false, func(c *color.RGBA64) {
					c.R = bottomLinear.R*(1-alpha) + topLinear.R*alpha
					c.G = bottomLinear.G*(1-alpha) + topLinear.G*alpha
					c.B = bottomLinear.B*(1-alpha) + topLinear.B*alpha
				})
			}
			// For grayscale colors, use the blend color's RGB values
			return topLinear.Process(false, false, func(c *color.RGBA64) {})
		}

		// Create result with blend hue, base saturation and lightness
		result := &color.HSL{
			H:     topHSL.H,
			S:     topHSL.S,
			L:     topHSL.L,
			Alpha: bottomLinear.A,
		}

		if alpha < constants.AlphaOpaque {
			// Interpolate hue based on alpha
			result.H = topHSL.H*alpha + bottomHSL.H*(1-alpha)
		}

		// Convert back to RGBA64 and then to sRGB space
		rgba := result.ToRGBA64()
		return rgba.Process(false, false, func(c *color.RGBA64) {})
	})
}
