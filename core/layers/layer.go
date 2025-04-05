package layers

import (
	stdImage "image"
	"image/draw"

	"github.com/toxyl/gfx/core/blendmodes"
	gfxColor "github.com/toxyl/gfx/core/color"
	"github.com/toxyl/gfx/core/fx"
)

// Layer represents a single layer in a composition
type Layer struct {
	width   int
	height  int
	blend   blendmodes.BlendMode
	alpha   float64
	effects []fx.Effect
	content stdImage.Image
}

// New creates a new layer with the specified properties
func New(width, height int, blend blendmodes.BlendMode, alpha float64, effects ...fx.Effect) *Layer {
	return &Layer{
		width:   width,
		height:  height,
		blend:   blend,
		alpha:   alpha,
		effects: effects,
		content: stdImage.NewRGBA(stdImage.Rect(0, 0, width, height)),
	}
}

// SetContent sets the layer's content
func (l *Layer) SetContent(img stdImage.Image) {
	l.content = img
}

// ApplyEffects applies all effects to the layer's content
func (l *Layer) ApplyEffects() stdImage.Image {
	result := l.content
	for _, effect := range l.effects {
		img, err := effect.Apply(result)
		if err != nil {
			continue // Skip effect on error
		}
		result = img
	}
	return result
}

// BlendWith blends this layer with another image
func (l *Layer) BlendWith(base stdImage.Image) stdImage.Image {
	content := l.ApplyEffects()

	// Convert images to RGBA
	baseRGBA := stdImage.NewRGBA(base.Bounds())
	draw.Draw(baseRGBA, baseRGBA.Bounds(), base, base.Bounds().Min, draw.Src)

	contentRGBA := stdImage.NewRGBA(content.Bounds())
	draw.Draw(contentRGBA, contentRGBA.Bounds(), content, content.Bounds().Min, draw.Src)

	// Create result image
	result := stdImage.NewRGBA(base.Bounds())

	// Blend each pixel
	for y := 0; y < result.Bounds().Dy(); y++ {
		for x := 0; x < result.Bounds().Dx(); x++ {
			// Get colors
			baseColor := baseRGBA.RGBA64At(x, y)
			contentColor := contentRGBA.RGBA64At(x, y)

			// Convert to our RGBA64
			baseRGBA64, err := gfxColor.NewRGBA64(
				float64(baseColor.R)/65535.0,
				float64(baseColor.G)/65535.0,
				float64(baseColor.B)/65535.0,
				float64(baseColor.A)/65535.0,
			)
			if err != nil {
				continue
			}

			contentRGBA64, err := gfxColor.NewRGBA64(
				float64(contentColor.R)/65535.0,
				float64(contentColor.G)/65535.0,
				float64(contentColor.B)/65535.0,
				float64(contentColor.A)/65535.0,
			)
			if err != nil {
				continue
			}

			// Blend colors
			blendMode, err := blendmodes.Get("normal") // Default to normal blend mode
			if err != nil {
				continue
			}
			blended, err := blendMode.Blend(baseRGBA64, contentRGBA64, l.alpha)
			if err != nil {
				continue
			}

			// Convert back to RGBA64
			rgba64 := blended.To16bit()
			result.SetRGBA64(x, y, rgba64)
		}
	}

	return result
}

// GetDimensions returns the layer's dimensions
func (l *Layer) GetDimensions() (int, int) {
	return l.width, l.height
}

// GetBlendMode returns the layer's blend mode
func (l *Layer) GetBlendMode() blendmodes.BlendMode {
	return l.blend
}

// GetAlpha returns the layer's alpha value
func (l *Layer) GetAlpha() float64 {
	return l.alpha
}

// GetEffects returns the layer's effects
func (l *Layer) GetEffects() []fx.Effect {
	return l.effects
}

// AddEffect adds an effect to the layer
func (l *Layer) AddEffect(effect fx.Effect) {
	l.effects = append(l.effects, effect)
}

// RemoveEffect removes an effect from the layer
func (l *Layer) RemoveEffect(index int) {
	if index >= 0 && index < len(l.effects) {
		l.effects = append(l.effects[:index], l.effects[index+1:]...)
	}
}

// SetBlendMode sets the layer's blend mode
func (l *Layer) SetBlendMode(blend blendmodes.BlendMode) {
	l.blend = blend
}

// SetAlpha sets the layer's alpha value
func (l *Layer) SetAlpha(alpha float64) {
	l.alpha = alpha
}
