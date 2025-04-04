package blendmodes

import (
	"github.com/toxyl/gfx/core/blendmodes/constants"
	"github.com/toxyl/gfx/core/color"
)

// MultiplyMode implements the multiply blend mode
type Multiply struct {
	BlendMode
}

// NewMultiply creates a new multiply blend mode
func NewMultiply() *Multiply {
	return &Multiply{
		BlendMode: BlendMode{
			name: "multiply",
		},
	}
}

// Blend implements the multiply blend mode
func (m *Multiply) Blend(base, blend *color.RGBA64) *color.RGBA64 {
	if base == nil || blend == nil {
		return nil
	}

	// If blend color has zero alpha, return base color
	if blend.A < 0.0001 {
		return base.Process(false, false, func(c *color.RGBA64) {})
	}

	// Prepare colors in linear space since multiply needs linear RGB
	bottom, top := prepareColors(base, blend, 1.0, AlphaNonPremultiplied)
	bottom.Process(true, true, func(c *color.RGBA64) {}) // Convert to linear
	top.Process(true, true, func(c *color.RGBA64) {})    // Convert to linear

	// Calculate final alpha using standard alpha compositing formula
	a := alphaComposite(bottom, top, 1.0)

	// Perform multiply blend
	result := &color.RGBA64{
		R: bottom.R * top.R,
		G: bottom.G * top.G,
		B: bottom.B * top.B,
		A: a,
	}

	// Convert back to sRGB space and ensure proper alpha handling
	return finalizeColor(result)
}

// Register the blend mode
func init() {
	Register("multiply", "Multiplies the color values of the base and blend colors", constants.CategoryDarken, func(bottom, top *color.RGBA64, alpha float64) *color.RGBA64 {
		if bottom == nil || top == nil {
			return nil
		}

		// If blend color has zero alpha, return base color
		if top.A < 0.0001 {
			return bottom.Process(false, false, func(c *color.RGBA64) {})
		}

		// Prepare colors in linear space since multiply needs linear RGB
		bottom, top = prepareColors(bottom, top, alpha, AlphaNonPremultiplied)
		bottom.Process(true, true, func(c *color.RGBA64) {}) // Convert to linear
		top.Process(true, true, func(c *color.RGBA64) {})    // Convert to linear

		// Calculate final alpha using standard alpha compositing formula
		a := alphaComposite(bottom, top, alpha)

		// Perform multiply blend
		result := &color.RGBA64{
			R: bottom.R * top.R,
			G: bottom.G * top.G,
			B: bottom.B * top.B,
			A: a,
		}

		// Convert back to sRGB space and ensure proper alpha handling
		return finalizeColor(result)
	})
}
