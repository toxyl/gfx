package fx

import (
	"image"
	stdcolor "image/color"

	"github.com/toxyl/gfx/core/color"
	"github.com/toxyl/gfx/core/fx"
)

// PastelizeFunction represents a function that creates a pastel effect
type PastelizeFunction struct {
	*fx.BaseFunction
	amount float64 // Pastel effect amount (0.0 to 1.0)
}

// Function arguments
var pastelizeArgs = []fx.FunctionArg{
	{
		Name:        "amount",
		Type:        fx.TypeFloat,
		Description: "pastelization adjustment value",
		Required:    true,
		Min:         0.0,
		Max:         1.0,
		Step:        0.1,
	},
}

func NewPastelize(amount float64) *PastelizeFunction {
	return &PastelizeFunction{
		BaseFunction: fx.NewBaseFunction("pastelize", "Applies pastelization transformation to an image", color.New(0, 0, 0, 1), pastelizeArgs),
		amount:       amount,
	}
}

// Apply implements the Function interface
func (f *PastelizeFunction) Apply(img image.Image) (image.Image, error) {
	if img == nil {
		return nil, fx.ErrInvalidImage
	}

	bounds := img.Bounds()
	result := image.NewRGBA(bounds)

	for y := bounds.Min.Y; y < bounds.Max.Y; y++ {
		for x := bounds.Min.X; x < bounds.Max.X; x++ {
			// Get the original color from the input image
			col := color.New(0, 0, 0, 1)
			col.SetUint8(img.At(x, y).(stdcolor.RGBA))

			// Convert to HSL for pastelization
			rgb := &color.RGBA64{
				R: col.R,
				G: col.G,
				B: col.B,
				A: col.A,
			}
			hsl := color.HSLFromRGB(rgb)

			// Adjust saturation and lightness for pastel effect
			// The formula:
			// S' = S * (1 - amount)
			// L' = L + (1 - L) * amount
			adjustedSaturation := hsl.S * (1.0 - f.amount)
			adjustedLightness := hsl.L + (1.0-hsl.L)*f.amount

			// Clamp values to [0,1] range
			adjustedSaturation = max(0.0, min(1.0, adjustedSaturation))
			adjustedLightness = max(0.0, min(1.0, adjustedLightness))

			// Convert back to RGB
			hsl.S = adjustedSaturation
			hsl.L = adjustedLightness
			rgb = hsl.ToRGB()

			// Set the result pixel
			result.Set(x, y, stdcolor.RGBA{
				R: uint8(rgb.R * 255),
				G: uint8(rgb.G * 255),
				B: uint8(rgb.B * 255),
				A: uint8(rgb.A * 255),
			})
		}
	}

	return result, nil
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

func init() {
	fx.DefaultRegistry.Register(NewPastelize(0.5))
}

func (f *PastelizeFunction) ProcessPixel(x, y int, img image.Image) (*color.Color64, error) {
	if img == nil {
		return nil, fx.ErrInvalidImage
	}

	bounds := img.Bounds()
	if x < bounds.Min.X || x >= bounds.Max.X || y < bounds.Min.Y || y >= bounds.Max.Y {
		return nil, fx.ErrOutOfBounds
	}

	// Get the original color from the input image
	col := color.New(0, 0, 0, 1)
	col.SetUint8(img.At(x, y).(stdcolor.RGBA))

	// Convert to HSL for pastelization
	rgb := &color.RGBA64{
		R: col.R,
		G: col.G,
		B: col.B,
		A: col.A,
	}
	hsl := color.HSLFromRGB(rgb)

	// Adjust saturation and lightness for pastel effect
	// The formula:
	// S' = S * (1 - amount)
	// L' = L + (1 - L) * amount
	adjustedSaturation := hsl.S * (1.0 - f.amount)
	adjustedLightness := hsl.L + (1.0-hsl.L)*f.amount

	// Clamp values to [0,1] range
	adjustedSaturation = max(0.0, min(1.0, adjustedSaturation))
	adjustedLightness = max(0.0, min(1.0, adjustedLightness))

	// Convert back to RGB
	hsl.S = adjustedSaturation
	hsl.L = adjustedLightness
	rgb = hsl.ToRGB()

	// Return new color
	return color.New(rgb.R, rgb.G, rgb.B, col.A), nil
}
