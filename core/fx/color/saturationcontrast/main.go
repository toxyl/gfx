package fx

import (
	"image"
	stdcolor "image/color"

	"github.com/toxyl/gfx/core/color"
	"github.com/toxyl/gfx/core/color/utils"
	"github.com/toxyl/gfx/core/fx"
)

// SaturationContrastFunction represents a function that adjusts the saturation contrast of an image
type SaturationContrastFunction struct {
	*fx.BaseFunction
	amount float64 // Saturation contrast adjustment amount (-1.0 to 1.0)
}

// Function arguments
var satContrastArgs = []fx.FunctionArg{
	{
		Name:        "amount",
		Type:        fx.TypeFloat,
		Description: "saturation contrast adjustment value",
		Required:    true,
		Min:         -1.0,
		Max:         1.0,
		Step:        0.1,
	},
}

func NewSaturationContrast(amount float64) *SaturationContrastFunction {
	return &SaturationContrastFunction{
		BaseFunction: fx.NewBaseFunction("saturation_contrast", "Applies saturation contrast transformation to an image", color.New(0, 0, 0, 1), satContrastArgs),
		amount:       amount,
	}
}

// Apply implements the ImageFunction interface
func (f *SaturationContrastFunction) Apply(img image.Image) (image.Image, error) {
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

			// Convert RGB to HSL
			h, s, l := utils.RGBToHSL(col.R, col.G, col.B)

			// Adjust saturation contrast
			// The formula: S' = S + (0.5 - S) * amount
			// This ensures that 0.5 stays at 0.5
			adjustedSaturation := s + (0.5-s)*f.amount

			// Clamp the saturation value
			adjustedSaturation = clamp(adjustedSaturation)

			// Convert back to RGB
			r, g, b := utils.HSLToRGB(h, adjustedSaturation, l)

			// Set the result pixel
			result.Set(x, y, stdcolor.RGBA{
				R: uint8(r * 255),
				G: uint8(g * 255),
				B: uint8(b * 255),
				A: uint8(col.A * 255),
			})
		}
	}

	return result, nil
}

// clamp ensures a value is between 0 and 1
func clamp(v float64) float64 {
	if v < 0 {
		return 0
	}
	if v > 1 {
		return 1
	}
	return v
}

func init() {
	fx.DefaultRegistry.Register(NewSaturationContrast(0))
}
