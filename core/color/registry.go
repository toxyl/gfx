package color

import (
	"github.com/toxyl/gfx/core/meta"
)

// ColorModelRegistry is a registry for color models.
type ColorModelRegistry struct {
	models map[string]*ColorModel
}

// ColorModel represents a color model in the registry.
type ColorModel struct {
	Meta *meta.ColorModelMeta
	To   func(src iColor) *RGBA64
	From func(rgba *RGBA64) iColor
}

// NewColorModelRegistry creates a new color model registry.
func NewColorModelRegistry() *ColorModelRegistry {
	return &ColorModelRegistry{
		models: make(map[string]*ColorModel),
	}
}

// Register adds a color model to the registry.
func (r *ColorModelRegistry) Register(name string, model *ColorModel) {
	r.models[name] = model
}

// Get retrieves a color model from the registry by name.
func (r *ColorModelRegistry) Get(name string) (*ColorModel, bool) {
	model, ok := r.models[name]
	return model, ok
}

// List returns a list of all registered color model names.
func (r *ColorModelRegistry) List() []string {
	names := make([]string, 0, len(r.models))
	for name := range r.models {
		names = append(names, name)
	}
	return names
}

// DefaultColorModelRegistry is the default registry for color models.
var DefaultColorModelRegistry = NewColorModelRegistry()
