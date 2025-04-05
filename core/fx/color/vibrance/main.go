package fx

import (
	"image"
	stdcolor "image/color"

	"github.com/toxyl/gfx/core/color"
	"github.com/toxyl/gfx/core/color/utils"
	"github.com/toxyl/gfx/core/fx"
)

// VibranceFunction represents a function that adjusts the vibrance of an image
type VibranceFunction struct {
	*fx.BaseFunction
	amount float64 // Vibrance adjustment amount (-1.0 to 1.0)
}

// Function arguments
var vibranceArgs = []fx.FunctionArg{
	{
		Name:        "amount",
		Type:        fx.TypeFloat,
		Description: "Amount of vibrance adjustment (-1.0 to 1.0)",
		Required:    true,
		Min:         -1.0,
		Max:         1.0,
		Step:        0.1,
	},
}

// NewVibrance creates a new vibrance function
func NewVibrance(amount float64) *VibranceFunction {
	return &VibranceFunction{
		BaseFunction: fx.NewBaseFunction("vibrance", "Adjusts the vibrance of an image", color.New(0, 0, 0, 1), vibranceArgs),
		amount:       amount,
	}
}

// Apply implements the ImageFunction interface
func (f *VibranceFunction) Apply(img image.Image) (image.Image, error) {
	if img == nil {
		return nil, fx.ErrInvalidImage
	}

	bounds := img.Bounds()
	result := image.NewRGBA(bounds)

	for y := bounds.Min.Y; y < bounds.Max.Y; y++ {
		for x := bounds.Min.X; x < bounds.Max.X; x++ {
			col, err := f.ProcessPixel(x, y, img)
			if err != nil {
				return nil, err
			}

			result.Set(x, y, stdcolor.RGBA{
				R: uint8(col.R * 255),
				G: uint8(col.G * 255),
				B: uint8(col.B * 255),
				A: uint8(col.A * 255),
			})
		}
	}

	return result, nil
}

// ProcessPixel applies the vibrance effect to a single pixel
func (f *VibranceFunction) ProcessPixel(x, y int, img image.Image) (*color.Color64, error) {
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

	// Convert to HSL
	h, s, l := utils.RGBToHSL(col.R, col.G, col.B)

	// Adjust saturation based on vibrance
	// For vibrance, we adjust saturation more for less saturated colors
	// and less for already saturated colors
	saturationAdjustment := (1 - s) * f.amount
	s = clamp(s+saturationAdjustment, 0.0, 1.0)

	// Convert back to RGB
	r, g, b := utils.HSLToRGB(h, s, l)
	col.R = r
	col.G = g
	col.B = b

	return col, nil
}

// clamp restricts a value to the range [min, max]
func clamp(value, min, max float64) float64 {
	if value < min {
		return min
	}
	if value > max {
		return max
	}
	return value
}

func init() {
	fx.DefaultRegistry.Register(NewVibrance(0))
}
