package layers

import (
	stdImage "image"
	"image/color"
	"image/draw"

	"github.com/toxyl/gfx/core/fx"
	gfxImage "github.com/toxyl/gfx/core/image"
)

// Composition represents a collection of layers that can be rendered together
type Composition struct {
	width      int
	height     int
	background color.Color
	effects    []fx.Effect
	layers     []*Layer
}

// NewComposition creates a new composition with the specified properties
func NewComposition(width, height int, background color.Color, effects []fx.Effect, layers ...*Layer) *Composition {
	return &Composition{
		width:      width,
		height:     height,
		background: background,
		effects:    effects,
		layers:     layers,
	}
}

// AddLayer adds a layer to the composition
func (c *Composition) AddLayer(layer *Layer) {
	c.layers = append(c.layers, layer)
}

// RemoveLayer removes a layer from the composition
func (c *Composition) RemoveLayer(index int) {
	if index >= 0 && index < len(c.layers) {
		c.layers = append(c.layers[:index], c.layers[index+1:]...)
	}
}

// GetLayers returns all layers in the composition
func (c *Composition) GetLayers() []*Layer {
	return c.layers
}

// GetDimensions returns the composition's dimensions
func (c *Composition) GetDimensions() (int, int) {
	return c.width, c.height
}

// GetBackground returns the composition's background color
func (c *Composition) GetBackground() color.Color {
	return c.background
}

// GetEffects returns the composition's effects
func (c *Composition) GetEffects() []fx.Effect {
	return c.effects
}

// AddEffect adds an effect to the composition
func (c *Composition) AddEffect(effect fx.Effect) {
	c.effects = append(c.effects, effect)
}

// RemoveEffect removes an effect from the composition
func (c *Composition) RemoveEffect(index int) {
	if index >= 0 && index < len(c.effects) {
		c.effects = append(c.effects[:index], c.effects[index+1:]...)
	}
}

// Render renders the composition to an image
func (c *Composition) Render() stdImage.Image {
	// Create base image with background
	base := stdImage.NewRGBA(stdImage.Rect(0, 0, c.width, c.height))
	if c.background != nil {
		for y := 0; y < c.height; y++ {
			for x := 0; x < c.width; x++ {
				base.Set(x, y, c.background)
			}
		}
	}

	// Apply layers
	for _, layer := range c.layers {
		blended := layer.BlendWith(base)
		if rgba, ok := blended.(*stdImage.RGBA); ok {
			base = rgba
		}
	}

	// Apply composition effects
	result := base
	for _, effect := range c.effects {
		img, err := effect.Apply(result)
		if err != nil {
			continue // Skip effect on error
		}
		if rgba, ok := img.(*stdImage.RGBA); ok {
			result = rgba
		} else {
			// Convert to RGBA if not already
			rgba := stdImage.NewRGBA(img.Bounds())
			draw.Draw(rgba, rgba.Bounds(), img, img.Bounds().Min, draw.Src)
			result = rgba
		}
	}

	return result
}

// RenderToFile renders the composition and saves it to a file
func (c *Composition) RenderToFile(filename string) error {
	img := c.Render()

	// Create a new Image instance from the rendered image
	gfxImg, err := gfxImage.FromImage(img)
	if err != nil {
		return err
	}

	// Save with default options
	return gfxImg.Save(filename, nil)
}
