package blendmodes

import (
	"github.com/toxyl/gfx/core/blendmodes/constants"
	"github.com/toxyl/gfx/core/color"
)

// Overlay blend mode implementation
type Overlay struct {
	BlendMode
}

// NewOverlay creates a new Overlay blend mode
func NewOverlay() *Overlay {
	return &Overlay{
		BlendMode: BlendMode{
			name: "overlay",
		},
	}
}

// Blend implements the overlay blend mode
func (o *Overlay) Blend(base, blend *color.RGBA64) *color.RGBA64 {
	if base == nil || blend == nil {
		return nil
	}

	// If blend color has zero alpha, return base color
	if blend.A < 0.0001 {
		return base.Process(false, false, func(c *color.RGBA64) {})
	}

	// Prepare colors in linear space since overlay needs linear RGB
	bottom, top := prepareColors(base, blend, 1.0, AlphaNonPremultiplied)
	bottom.Process(true, true, func(c *color.RGBA64) {}) // Convert to linear
	top.Process(true, true, func(c *color.RGBA64) {})    // Convert to linear

	// Perform overlay blend
	result := overlayBlend(bottom, top, 1.0)

	// Convert back to sRGB space and ensure proper alpha handling
	return finalizeColor(result)
}

// overlayBlend performs the overlay blend operation for a single channel
func overlayBlend(bottom, top *color.RGBA64, alpha float64) *color.RGBA64 {
	// Calculate final alpha using standard alpha compositing formula
	a := alphaComposite(bottom, top, alpha)

	// Helper function for overlay calculation
	overlayComponent := func(b, t float64) float64 {
		if b <= 0.5 {
			return 2 * b * t
		}
		return 1 - 2*(1-b)*(1-t)
	}

	// Scale blend color contribution by alpha
	t := top.A * alpha

	// Perform overlay blend
	result := &color.RGBA64{
		R: bottom.R*(1-t) + overlayComponent(bottom.R, top.R)*t,
		G: bottom.G*(1-t) + overlayComponent(bottom.G, top.G)*t,
		B: bottom.B*(1-t) + overlayComponent(bottom.B, top.B)*t,
		A: a,
	}

	return result
}

// Register the blend mode
func init() {
	Register("overlay", "Multiplies or screens the colors, depending on the base color", constants.CategoryContrast, func(bottom, top *color.RGBA64, alpha float64) *color.RGBA64 {
		if bottom == nil || top == nil {
			return nil
		}

		// If blend color has zero alpha, return base color
		if top.A < 0.0001 {
			return bottom.Process(false, false, func(c *color.RGBA64) {})
		}

		// Prepare colors in linear space since overlay needs linear RGB
		bottom, top = prepareColors(bottom, top, alpha, AlphaNonPremultiplied)
		bottom.Process(true, true, func(c *color.RGBA64) {}) // Convert to linear
		top.Process(true, true, func(c *color.RGBA64) {})    // Convert to linear

		// Perform overlay blend
		result := overlayBlend(bottom, top, alpha)

		// Convert back to sRGB space and ensure proper alpha handling
		return finalizeColor(result)
	})
}
