package blendmodes

import (
	"github.com/toxyl/gfx/core/blendmodes/constants"
	"github.com/toxyl/gfx/core/color"
)

// NormalMode implements the normal (standard alpha compositing) blend mode
type NormalMode struct {
	BlendMode
}

// NewNormal creates a new normal blend mode
func NewNormal() *NormalMode {
	return &NormalMode{
		BlendMode: BlendMode{
			name: "normal",
		},
	}
}

// Blend implements the normal blend mode
func (n *NormalMode) Blend(base, blend *color.RGBA64) *color.RGBA64 {
	if base == nil || blend == nil {
		return nil
	}

	// Prepare colors in linear space for accurate blending
	bottom, top := prepareColors(base, blend, 1.0, AlphaNonPremultiplied)
	bottom.Process(true, true, func(c *color.RGBA64) {})
	top.Process(true, true, func(c *color.RGBA64) {})

	// Calculate final alpha using standard alpha compositing formula
	a := alphaComposite(bottom, top, 1.0)

	// For normal blend mode, we interpolate between bottom and top based on top's alpha
	result := &color.RGBA64{
		R: bottom.R*(1-top.A) + top.R*top.A,
		G: bottom.G*(1-top.A) + top.G*top.A,
		B: bottom.B*(1-top.A) + top.B*top.A,
		A: a,
	}

	// Convert back to sRGB space and ensure proper alpha handling
	return finalizeColor(result)
}

// Register the blend mode
func init() {
	Register("normal", "Standard alpha compositing blend mode", constants.CategoryBasic, func(bottom, top *color.RGBA64, alpha float64) *color.RGBA64 {
		if bottom == nil || top == nil {
			return nil
		}

		// Prepare colors in linear space for accurate blending
		bottom, top = prepareColors(bottom, top, alpha, AlphaNonPremultiplied)
		bottom.Process(true, true, func(c *color.RGBA64) {})
		top.Process(true, true, func(c *color.RGBA64) {})

		// Calculate final alpha using standard alpha compositing formula
		a := alphaComposite(bottom, top, alpha)

		// For normal blend mode, we interpolate between bottom and top based on top's alpha
		result := &color.RGBA64{
			R: bottom.R*(1-top.A) + top.R*top.A,
			G: bottom.G*(1-top.A) + top.G*top.A,
			B: bottom.B*(1-top.A) + top.B*top.A,
			A: a,
		}

		// Convert back to sRGB space and ensure proper alpha handling
		return finalizeColor(result)
	})
}
