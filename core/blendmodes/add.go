package blendmodes

import (
	"github.com/toxyl/gfx/core/blendmodes/constants"
	"github.com/toxyl/gfx/core/color"
)

// AddMode implements the add (linear dodge) blend mode
type AddMode struct {
	BlendMode
}

// NewAdd creates a new add blend mode
func NewAdd() *AddMode {
	return &AddMode{
		BlendMode: BlendMode{
			name: "add",
		},
	}
}

// Blend implements the add blend mode
func (m *AddMode) Blend(base, blend *color.RGBA64) *color.RGBA64 {
	if base == nil || blend == nil {
		return nil
	}

	// If blend color has zero alpha, return base color
	if blend.A < 0.0001 {
		return base.Process(false, false, func(c *color.RGBA64) {})
	}

	// Prepare colors in linear space since add needs linear RGB
	bottom, top := prepareColors(base, blend, 1.0, AlphaNonPremultiplied)
	bottom.Process(true, true, func(c *color.RGBA64) {}) // Convert to linear
	top.Process(true, true, func(c *color.RGBA64) {})    // Convert to linear

	// Calculate final alpha using standard alpha compositing formula
	finalAlpha := alphaComposite(bottom, top, 1.0)

	// Scale blend color contribution by its alpha
	topR := top.R * top.A
	topG := top.G * top.A
	topB := top.B * top.A

	// Perform add blend
	result := &color.RGBA64{
		R: bottom.R + topR,
		G: bottom.G + topG,
		B: bottom.B + topB,
		A: finalAlpha,
	}

	// Convert back to sRGB space and ensure proper alpha handling
	return finalizeColor(result)
}

// Register the blend mode
func init() {
	Register("add", "Adds the color values of the base and blend colors (also known as Linear Dodge)", constants.CategoryLighten, func(bottom, top *color.RGBA64, alpha float64) *color.RGBA64 {
		if bottom == nil || top == nil {
			return nil
		}

		// If blend color has zero alpha, return base color
		if top.A < 0.0001 {
			return bottom.Process(false, false, func(c *color.RGBA64) {})
		}

		// Prepare colors in linear space since add needs linear RGB
		bottom, top = prepareColors(bottom, top, alpha, AlphaNonPremultiplied)
		bottom.Process(true, true, func(c *color.RGBA64) {}) // Convert to linear
		top.Process(true, true, func(c *color.RGBA64) {})    // Convert to linear

		// Calculate final alpha using standard alpha compositing formula
		finalAlpha := alphaComposite(bottom, top, alpha)

		// Scale blend color contribution by its alpha
		topR := top.R * top.A * alpha
		topG := top.G * top.A * alpha
		topB := top.B * top.A * alpha

		// Perform add blend
		result := &color.RGBA64{
			R: bottom.R + topR,
			G: bottom.G + topG,
			B: bottom.B + topB,
			A: finalAlpha,
		}

		// Convert back to sRGB space and ensure proper alpha handling
		return finalizeColor(result)
	})
}
