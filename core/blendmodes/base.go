/*
Package blendmodes provides a comprehensive set of blend modes for image processing and composition.

This package implements various blend modes commonly found in professional image editing software,
allowing for sophisticated layer compositing and color blending operations. Each blend mode defines
how pixels from a top layer interact with pixels from a bottom layer to produce a final result.

The package includes several categories of blend modes:

  - Basic: Normal, Erase
  - Darken: Darken, Multiply, Color Burn, Linear Burn, Darker Color
  - Lighten: Lighten, Screen, Color Dodge, Linear Dodge (Add), Lighter Color
  - Contrast: Overlay, Soft Light, Hard Light, Vivid Light, Linear Light, Pin Light, Hard Mix
  - Comparative: Difference, Exclusion, Subtract, Divide, Negation, Contrast Negate
  - Component: Hue, Saturation, Color, Luminosity
  - Special Effects: Reflect, Glow, Average

Key features:

  - Thread-safe blend mode registry
  - Consistent alpha channel handling
  - High precision color operations using RGBA64
  - Support for custom blend mode registration
  - Comprehensive metadata for each blend mode

Basic usage:

	// Get a blend mode
	multiply, err := blendmodes.Get("multiply")
	if err != nil {
	    // Handle error
	}

	// Blend two colors with 50% opacity
	result := multiply.Blend(bottomColor, topColor, 0.5)

The package uses a registry pattern to manage blend modes, allowing for runtime registration
of custom blend modes while maintaining thread safety through mutex locks.
*/
package blendmodes

import (
	"fmt"
	"sort"
	"sync"

	"github.com/toxyl/gfx/core/blendmodes/constants"
	"github.com/toxyl/gfx/core/color"
)

// Registry is a registry for blend modes.
type Registry struct {
	modes map[string]*IBlendMode
}

// AlphaMode represents how alpha should be handled in a blend operation
type AlphaMode int

const (
	// AlphaPremultiplied indicates the colors are in premultiplied alpha format
	AlphaPremultiplied AlphaMode = iota
	// AlphaNonPremultiplied indicates the colors are in non-premultiplied alpha format
	AlphaNonPremultiplied
)

// IBlendMode represents a blend mode in the registry.
type IBlendMode struct {
	meta  *IMeta
	blend func(bottom, top *color.RGBA64, alpha float64) *color.RGBA64
}

// Blend performs the blending operation between two colors.
func (b *IBlendMode) Blend(bottom, top *color.RGBA64, alpha float64) (*color.RGBA64, error) {
	if err := validateBlendParams(bottom, top, alpha); err != nil {
		return nil, err
	}
	return b.blend(bottom, top, alpha), nil
}

// Meta returns metadata about the blend mode.
func (b *IBlendMode) Meta() *IMeta {
	return b.meta
}

// registry is the internal registry instance.
var registry = struct {
	modes map[string]*IBlendMode
	mu    sync.RWMutex
}{
	modes: make(map[string]*IBlendMode),
}

// NewRegistry creates a new blend mode registry.
func NewRegistry() *Registry {
	return &Registry{
		modes: make(map[string]*IBlendMode),
	}
}

// Register registers a new blend mode with the given name, description, and category.
func Register(name, description, category string, blend func(bottom, top *color.RGBA64, alpha float64) *color.RGBA64) error {
	if err := validateRegistration(name, description, category, blend); err != nil {
		return fmt.Errorf("invalid blend mode registration: %w", err)
	}

	registry.mu.Lock()
	defer registry.mu.Unlock()

	if _, exists := registry.modes[name]; exists {
		return fmt.Errorf("blend mode %s already exists", name)
	}

	meta, err := NewIMeta(name, description, category)
	if err != nil {
		return fmt.Errorf("failed to create blend mode metadata: %w", err)
	}

	registry.modes[name] = &IBlendMode{
		meta:  meta,
		blend: blend,
	}
	return nil
}

// Get returns the blend mode with the given name.
func Get(name string) (*IBlendMode, error) {
	registry.mu.RLock()
	defer registry.mu.RUnlock()

	if mode, exists := registry.modes[name]; exists {
		return mode, nil
	}
	return nil, fmt.Errorf("blend mode %s not found", name)
}

// GetByCategory returns all blend modes in the specified category.
func GetByCategory(category string) ([]*IBlendMode, error) {
	registry.mu.RLock()
	defer registry.mu.RUnlock()

	var modes []*IBlendMode
	for _, mode := range registry.modes {
		if mode.meta.Category() == category {
			modes = append(modes, mode)
		}
	}

	if len(modes) == 0 {
		return nil, fmt.Errorf("no blend modes found for category: %s", category)
	}

	// Sort modes by name for consistent ordering
	sort.Slice(modes, func(i, j int) bool {
		return modes[i].meta.Name() < modes[j].meta.Name()
	})

	return modes, nil
}

// List returns a list of all registered blend mode names.
func List() []string {
	registry.mu.RLock()
	defer registry.mu.RUnlock()

	names := make([]string, 0, len(registry.modes))
	for name := range registry.modes {
		names = append(names, name)
	}
	sort.Strings(names)
	return names
}

// ListCategories returns a list of all available blend mode categories.
func ListCategories() []string {
	return []string{
		constants.CategoryBasic,
		constants.CategoryDarken,
		constants.CategoryLighten,
		constants.CategoryContrast,
		constants.CategoryComparative,
		constants.CategoryComponent,
		constants.CategorySpecial,
	}
}

// BlendMode provides common functionality for all blend modes
type BlendMode struct {
	name string
}

// Name returns the name of the blend mode
func (b *BlendMode) Name() string {
	return b.name
}

// alphaComposite performs standard alpha compositing
func alphaComposite(bottom, top *color.RGBA64, alpha float64) float64 {
	return bottom.A + alpha*top.A*(1-bottom.A)
}

// toPremultiplied converts a color to premultiplied alpha format if it isn't already
func toPremultiplied(c *color.RGBA64) *color.RGBA64 {
	if c == nil {
		return nil
	}

	// If alpha is almost 0 or 1, handle as special case to avoid precision issues
	if constants.IsAlphaAlmostZero(c.A) {
		return &color.RGBA64{R: 0, G: 0, B: 0, A: 0}
	}
	if constants.IsAlphaAlmostOne(c.A) {
		return &color.RGBA64{R: c.R, G: c.G, B: c.B, A: 1.0}
	}

	// Create new color and process it
	result := c.Copy()
	result.Process(false, false, func(c *color.RGBA64) {
		// Process will handle premultiplication state
	})
	return result
}

// toNonPremultiplied converts a color to non-premultiplied alpha format if it isn't already
func toNonPremultiplied(c *color.RGBA64) *color.RGBA64 {
	if c == nil {
		return nil
	}

	// If alpha is almost 0 or 1, handle as special case to avoid precision issues
	if constants.IsAlphaAlmostZero(c.A) {
		return &color.RGBA64{R: 0, G: 0, B: 0, A: 0}
	}
	if constants.IsAlphaAlmostOne(c.A) {
		return &color.RGBA64{R: c.R, G: c.G, B: c.B, A: 1.0}
	}

	// Create new color and process it
	result := c.Copy()
	result.Process(true, false, func(c *color.RGBA64) {
		// Process will handle premultiplication state
	})
	return result
}

// prepareColors prepares colors for blending by converting them to the desired alpha mode
// and applying the blend opacity
func prepareColors(bottom, top *color.RGBA64, alpha float64, mode AlphaMode) (*color.RGBA64, *color.RGBA64) {
	if bottom == nil || top == nil {
		return nil, nil
	}

	// Create new colors to avoid modifying originals
	bottomNew := bottom.Copy()
	topNew := top.Copy()

	// Apply blend opacity to top color's alpha
	topNew.A *= alpha

	// Convert to desired alpha mode
	if mode == AlphaPremultiplied {
		bottomNew = toPremultiplied(bottomNew)
		topNew = toPremultiplied(topNew)
	} else {
		bottomNew = toNonPremultiplied(bottomNew)
		topNew = toNonPremultiplied(topNew)
	}

	return bottomNew, topNew
}

// finalizeColor ensures the final color is in premultiplied format (Go's native format)
// and has valid values
func finalizeColor(c *color.RGBA64) *color.RGBA64 {
	if c == nil {
		return nil
	}

	// Create new color and process it
	result := c.Copy()
	result.Process(false, false, func(c *color.RGBA64) {
		// Process will handle premultiplication state
		c.Clamp() // Ensure values are in valid range
	})
	return result
}
