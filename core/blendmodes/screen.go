package blendmodes

import (
	"github.com/toxyl/gfx/core/blendmodes/constants"
	"github.com/toxyl/gfx/core/color"
)

// Screen blend mode implementation
type Screen struct {
	BlendMode
}

// NewScreen creates a new Screen blend mode
func NewScreen() *Screen {
	return &Screen{
		BlendMode: BlendMode{
			name: "screen",
		},
	}
}

// Blend implements the screen blend mode
func (s *Screen) Blend(base, blend *color.RGBA64) *color.RGBA64 {
	if base == nil || blend == nil {
		return nil
	}

	// Prepare colors in linear space since screen needs linear RGB
	bottom, top := prepareColors(base, blend, 1.0, AlphaNonPremultiplied)
	bottom.Process(true, true, func(c *color.RGBA64) {}) // Convert to linear
	top.Process(true, true, func(c *color.RGBA64) {})    // Convert to linear

	// Calculate final alpha using standard alpha compositing formula
	a := alphaComposite(bottom, top, 1.0)

	// Perform screen blend: 1 - (1-a)(1-b)
	result := &color.RGBA64{
		R: 1 - (1-bottom.R)*(1-top.R),
		G: 1 - (1-bottom.G)*(1-top.G),
		B: 1 - (1-bottom.B)*(1-top.B),
		A: a,
	}

	// Convert back to sRGB space and ensure proper alpha handling
	return finalizeColor(result)
}

// Register the blend mode
func init() {
	Register("screen", "Multiplies the complements of the base and blend colors, then takes the complement of the result", constants.CategoryLighten, func(bottom, top *color.RGBA64, alpha float64) *color.RGBA64 {
		if bottom == nil || top == nil {
			return nil
		}

		// Prepare colors in linear space since screen needs linear RGB
		bottom, top = prepareColors(bottom, top, alpha, AlphaNonPremultiplied)
		bottom.Process(true, true, func(c *color.RGBA64) {}) // Convert to linear
		top.Process(true, true, func(c *color.RGBA64) {})    // Convert to linear

		// Calculate final alpha using standard alpha compositing formula
		a := alphaComposite(bottom, top, alpha)

		// Perform screen blend: 1 - (1-a)(1-b)
		result := &color.RGBA64{
			R: 1 - (1-bottom.R)*(1-top.R),
			G: 1 - (1-bottom.G)*(1-top.G),
			B: 1 - (1-bottom.B)*(1-top.B),
			A: a,
		}

		// Convert back to sRGB space and ensure proper alpha handling
		return finalizeColor(result)
	})
}
