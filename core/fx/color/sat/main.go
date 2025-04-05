package fx

import (
	"image"
	stdcolor "image/color"

	"github.com/toxyl/gfx/core/color"
	"github.com/toxyl/gfx/core/color/utils"
	"github.com/toxyl/gfx/core/fx"
)

// SaturationFunction represents a function that adjusts the saturation of an image
type SaturationFunction struct {
	*fx.BaseFunction
	amount float64 // Saturation adjustment amount (-1.0 to 1.0)
}

// Function arguments
var satArgs = []fx.FunctionArg{
	{
		Name:        "amount",
		Type:        fx.TypeFloat,
		Description: "sat adjustment value",
		Required:    true,
		Min:         -1.0,
		Max:         1.0,
		Step:        0.1,
	},
}

func NewSaturation(amount float64) *SaturationFunction {
	return &SaturationFunction{
		BaseFunction: fx.NewBaseFunction("saturation", "Applies saturation transformation to an image", color.New(0, 0, 0, 1), satArgs),
		amount:       amount,
	}
}

// Apply implements the ImageFunction interface
func (f *SaturationFunction) Apply(img image.Image) (image.Image, error) {
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

			// Adjust saturation
			// The formula: S' = S * (1 + amount)
			// This ensures that 0 stays at 0
			adjustedSaturation := s * (1.0 + f.amount)

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
	fx.DefaultRegistry.Register(NewSaturation(0))
}
