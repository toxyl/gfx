package blendmodes

import (
	"fmt"
	"strings"

	"github.com/toxyl/gfx/core/blendmodes/constants"
	"github.com/toxyl/gfx/core/color"
)

// validateBlendParams validates the parameters for blend operations
func validateBlendParams(bottom, top *color.RGBA64, alpha float64) error {
	if bottom == nil || top == nil {
		return fmt.Errorf("blend colors cannot be nil")
	}
	// Use epsilon comparison for alpha range check
	if alpha < constants.AlphaTransparent-constants.AlphaEpsilon || alpha > constants.AlphaOpaque+constants.AlphaEpsilon {
		return fmt.Errorf("alpha must be between %.1f and %.1f", constants.AlphaTransparent, constants.AlphaOpaque)
	}
	return nil
}

// validateCategory checks if the given category is valid
func validateCategory(category string) error {
	switch category {
	case constants.CategoryBasic, constants.CategoryDarken, constants.CategoryLighten, constants.CategoryContrast,
		constants.CategoryComparative, constants.CategoryComponent, constants.CategorySpecial:
		return nil
	default:
		return fmt.Errorf("invalid blend mode category: %s", category)
	}
}

// validateBlendModeName checks if the given blend mode name is valid
func validateBlendModeName(name string) error {
	if strings.TrimSpace(name) == "" {
		return fmt.Errorf("blend mode name cannot be empty")
	}
	return nil
}

// validateBlendFunction checks if the given blend function is valid
func validateBlendFunction(blend func(bottom, top *color.RGBA64, alpha float64) *color.RGBA64) error {
	if blend == nil {
		return fmt.Errorf("blend function cannot be nil")
	}
	return nil
}

// validateRegistration validates all parameters for registering a new blend mode
func validateRegistration(name, description, category string, blend func(bottom, top *color.RGBA64, alpha float64) *color.RGBA64) error {
	if err := validateBlendModeName(name); err != nil {
		return err
	}
	if err := validateCategory(category); err != nil {
		return err
	}
	if err := validateBlendFunction(blend); err != nil {
		return err
	}
	return nil
}
