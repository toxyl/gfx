package layers

import (
	stdImage "image"
	"image/color"
	"testing"

	"github.com/toxyl/gfx/core/blendmodes"
	"github.com/toxyl/gfx/core/fx"
	"github.com/toxyl/gfx/core/fx/models"
	"github.com/toxyl/math"
)

// TestLayerCreation tests the creation of new layers
func TestLayerCreation(t *testing.T) {
	tests := []struct {
		name    string
		width   int
		height  int
		blend   blendmodes.BlendMode
		alpha   float64
		effects []fx.Effect
		wantErr bool
	}{
		{
			name:    "valid_layer",
			width:   100,
			height:  100,
			blend:   blendmodes.NewNormal().BlendMode,
			alpha:   1.0,
			effects: []fx.Effect{},
			wantErr: false,
		},
		{
			name:    "invalid_dimensions",
			width:   -1,
			height:  -1,
			blend:   blendmodes.NewNormal().BlendMode,
			alpha:   1.0,
			effects: []fx.Effect{},
			wantErr: true,
		},
		{
			name:    "invalid_alpha",
			width:   100,
			height:  100,
			blend:   blendmodes.NewNormal().BlendMode,
			alpha:   -0.1,
			effects: []fx.Effect{},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			layer := New(tt.width, tt.height, tt.blend, tt.alpha, tt.effects...)
			if layer == nil && !tt.wantErr {
				t.Error("New() returned nil for valid input")
			}
			if layer != nil && tt.wantErr {
				t.Error("New() returned non-nil for invalid input")
			}
		})
	}
}

// TestLayerContent tests setting and getting layer content
func TestLayerContent(t *testing.T) {
	layer := New(100, 100, blendmodes.NewNormal().BlendMode, 1.0)
	img := stdImage.NewRGBA(stdImage.Rect(0, 0, 100, 100))

	// Test SetContent
	layer.SetContent(img)
	if layer.content != img {
		t.Error("SetContent() did not set the content correctly")
	}

	// Test ApplyEffects with no effects
	result := layer.ApplyEffects()
	if result != img {
		t.Error("ApplyEffects() with no effects should return the original image")
	}
}

// TestLayerEffects tests adding and removing effects
func TestLayerEffects(t *testing.T) {
	layer := New(100, 100, blendmodes.NewNormal().BlendMode, 1.0)
	effect := models.NewVibranceEffect(0.5)

	// Test AddEffect
	layer.AddEffect(effect)
	if len(layer.effects) != 1 {
		t.Error("AddEffect() did not add the effect")
	}

	// Test GetEffects
	effects := layer.GetEffects()
	if len(effects) != 1 || effects[0] != fx.Effect(effect) {
		t.Error("GetEffects() did not return the correct effects")
	}

	// Test RemoveEffect
	layer.RemoveEffect(0)
	if len(layer.effects) != 0 {
		t.Error("RemoveEffect() did not remove the effect")
	}
}

// TestLayerBlendMode tests setting and getting blend modes
func TestLayerBlendMode(t *testing.T) {
	layer := New(100, 100, blendmodes.NewNormal().BlendMode, 1.0)
	newBlend := blendmodes.NewMultiply().BlendMode

	// Test SetBlendMode
	layer.SetBlendMode(newBlend)
	if layer.blend != newBlend {
		t.Error("SetBlendMode() did not set the blend mode correctly")
	}

	// Test GetBlendMode
	if layer.GetBlendMode() != newBlend {
		t.Error("GetBlendMode() did not return the correct blend mode")
	}
}

// TestLayerAlpha tests setting and getting alpha values
func TestLayerAlpha(t *testing.T) {
	layer := New(100, 100, blendmodes.NewNormal().BlendMode, 1.0)
	newAlpha := 0.5

	// Test SetAlpha
	layer.SetAlpha(newAlpha)
	if !math.ApproxEqual(layer.alpha, newAlpha, 0.001) {
		t.Error("SetAlpha() did not set the alpha value correctly")
	}

	// Test GetAlpha
	if !math.ApproxEqual(layer.GetAlpha(), newAlpha, 0.001) {
		t.Error("GetAlpha() did not return the correct alpha value")
	}
}

// TestCompositionCreation tests the creation of new compositions
func TestCompositionCreation(t *testing.T) {
	tests := []struct {
		name       string
		width      int
		height     int
		background color.Color
		effects    []fx.Effect
		layers     []*Layer
		wantErr    bool
	}{
		{
			name:       "valid_composition",
			width:      100,
			height:     100,
			background: color.White,
			effects:    []fx.Effect{},
			layers:     []*Layer{},
			wantErr:    false,
		},
		{
			name:       "invalid_dimensions",
			width:      -1,
			height:     -1,
			background: color.White,
			effects:    []fx.Effect{},
			layers:     []*Layer{},
			wantErr:    true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			comp := NewComposition(tt.width, tt.height, tt.background, tt.effects, tt.layers...)
			if comp == nil && !tt.wantErr {
				t.Error("NewComposition() returned nil for valid input")
			}
			if comp != nil && tt.wantErr {
				t.Error("NewComposition() returned non-nil for invalid input")
			}
		})
	}
}

// TestCompositionLayers tests adding and removing layers
func TestCompositionLayers(t *testing.T) {
	comp := NewComposition(100, 100, color.White, []fx.Effect{})
	layer := New(100, 100, blendmodes.NewNormal().BlendMode, 1.0)

	// Test AddLayer
	comp.AddLayer(layer)
	if len(comp.layers) != 1 {
		t.Error("AddLayer() did not add the layer")
	}

	// Test GetLayers
	layers := comp.GetLayers()
	if len(layers) != 1 || layers[0] != layer {
		t.Error("GetLayers() did not return the correct layers")
	}

	// Test RemoveLayer
	comp.RemoveLayer(0)
	if len(comp.layers) != 0 {
		t.Error("RemoveLayer() did not remove the layer")
	}
}

// TestCompositionEffects tests adding and removing effects
func TestCompositionEffects(t *testing.T) {
	comp := NewComposition(100, 100, color.White, []fx.Effect{})
	effect := models.NewVibranceEffect(0.5)

	// Test AddEffect
	comp.AddEffect(effect)
	if len(comp.effects) != 1 {
		t.Error("AddEffect() did not add the effect")
	}

	// Test GetEffects
	effects := comp.GetEffects()
	if len(effects) != 1 || effects[0] != fx.Effect(effect) {
		t.Error("GetEffects() did not return the correct effects")
	}

	// Test RemoveEffect
	comp.RemoveEffect(0)
	if len(comp.effects) != 0 {
		t.Error("RemoveEffect() did not remove the effect")
	}
}

// TestCompositionRender tests rendering a composition
func TestCompositionRender(t *testing.T) {
	comp := NewComposition(100, 100, color.White, []fx.Effect{})
	layer := New(100, 100, blendmodes.NewNormal().BlendMode, 1.0)
	img := stdImage.NewRGBA(stdImage.Rect(0, 0, 100, 100))
	layer.SetContent(img)
	comp.AddLayer(layer)

	// Test Render
	result := comp.Render()
	if result == nil {
		t.Error("Render() returned nil")
	}
	if result.Bounds().Dx() != 100 || result.Bounds().Dy() != 100 {
		t.Error("Render() returned image with incorrect dimensions")
	}
}

// TestCompositionRenderToFile tests saving a composition to a file
func TestCompositionRenderToFile(t *testing.T) {
	comp := NewComposition(100, 100, color.White, []fx.Effect{})
	layer := New(100, 100, blendmodes.NewNormal().BlendMode, 1.0)
	img := stdImage.NewRGBA(stdImage.Rect(0, 0, 100, 100))
	layer.SetContent(img)
	comp.AddLayer(layer)

	// Test RenderToFile
	err := comp.RenderToFile("test_output.png")
	if err != nil {
		t.Errorf("RenderToFile() returned error: %v", err)
	}
}
