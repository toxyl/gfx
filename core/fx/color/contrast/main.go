package fx

import (
	"image"
	stdcolor "image/color"
	"image/draw"

	"github.com/toxyl/gfx/core/color"
	"github.com/toxyl/gfx/core/fx"
)

// ContrastFunction represents a function that adjusts the contrast of an image
type ContrastFunction struct {
	*fx.BaseFunction
	amount float64
}

// NewContrast creates a new contrast adjustment function
// Function arguments
var contrastArgs = []fx.FunctionArg{
	{
		Name:        "amount",
		Type:        fx.TypeFloat,
		Description: "contrast adjustment value",
		Required:    true,
		Min:         -1.0,
		Max:         1.0,
		Step:        0.1,
	},
}

func NewContrast(amount float64) *ContrastFunction {
	return &ContrastFunction{
		BaseFunction: fx.NewBaseFunction("contrast", "Applies contrast transformation to an image", color.New(0, 0, 0, 1), contrastArgs),
		amount:       amount,
	}
}

// Apply implements the Function interface
func (f *ContrastFunction) Apply(img *image.Image) (*image.Image, error) {
	if img == nil || *img == nil {
		return nil, fx.ErrInvalidImage
	}
	bounds := (*img).Bounds()
	dst := image.NewRGBA(bounds)
	draw.Draw(dst, bounds, *img, bounds.Min, draw.Src)

	for y := bounds.Min.Y; y < bounds.Max.Y; y++ {
		for x := bounds.Min.X; x < bounds.Max.X; x++ {
			col, err := f.ProcessPixel(x, y, img)
			if err != nil {
				return nil, err
			}
			dst.Set(x, y, col.ToUint8())
		}
	}

	result := image.Image(dst)
	return &result, nil
}

// ProcessPixel implements the ImageFunction interface
func (f *ContrastFunction) ProcessPixel(x, y int, img *image.Image) (*color.Color64, error) {
	if img == nil || *img == nil {
		return nil, fx.ErrInvalidImage
	}

	bounds := (*img).Bounds()
	if x < bounds.Min.X || x >= bounds.Max.X || y < bounds.Min.Y || y >= bounds.Max.Y {
		return nil, fx.ErrOutOfBounds
	}

	// Get the original color from the input image
	col := color.New(0, 0, 0, 1)
	col.SetUint8((*img).At(x, y).(stdcolor.RGBA))

	// Convert to HSL for contrast adjustment
	rgb := &color.RGBA64{
		R: col.R,
		G: col.G,
		B: col.B,
		A: col.A,
	}
	hsl := color.HSLFromRGB(rgb)

	// Calculate contrast adjustment in HSL space
	// The formula: L' = L + (L - 0.5) * amount
	// This ensures that 0.5 (middle gray) stays at 0.5
	adjustedLightness := hsl.L + (hsl.L-0.5)*f.amount

	// Clamp the lightness value
	adjustedLightness = max(0.0, min(1.0, adjustedLightness))

	// Convert back to RGB
	hsl.L = adjustedLightness
	result := hsl.ToRGB()

	return color.New(result.R, result.G, result.B, col.A), nil
}

func max(a, b float64) float64 {
	if a > b {
		return a
	}
	return b
}

func min(a, b float64) float64 {
	if a < b {
		return a
	}
	return b
}
