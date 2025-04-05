package fx

import (
	"image"
	stdcolor "image/color"
	"math"

	"github.com/toxyl/gfx/core/color"
	"github.com/toxyl/gfx/core/fx"
)

// GammaFunction represents a function that adjusts the gamma value of an image
type GammaFunction struct {
	*fx.BaseFunction
	amount float64 // Gamma adjustment amount (0.0 to 5.0)
}

// Function arguments
var gammaArgs = []fx.FunctionArg{
	{
		Name:        "amount",
		Type:        fx.TypeFloat,
		Description: "gamma adjustment value",
		Required:    true,
		Min:         0.0,
		Max:         5.0,
		Step:        0.1,
	},
}

// NewGamma creates a new gamma correction function
func NewGamma(amount float64) *GammaFunction {
	return &GammaFunction{
		BaseFunction: fx.NewBaseFunction("gamma", "Applies gamma transformation to an image", color.New(0, 0, 0, 1), gammaArgs),
		amount:       amount,
	}
}

// Apply implements the ImageFunction interface
func (f *GammaFunction) Apply(img image.Image) (image.Image, error) {
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

// ProcessPixel applies the gamma correction to a single pixel
func (f *GammaFunction) ProcessPixel(x, y int, img image.Image) (*color.Color64, error) {
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

	// Calculate gamma value
	// amount of 1.0 means gamma of 2.0
	// amount of 0.0 means gamma of 1.0
	gamma := 1.0 + f.amount

	// Apply gamma correction to each channel
	// The formula: C' = C^(1/gamma)
	// This ensures that 0 stays at 0 and 1 stays at 1
	col.R = clamp(math.Pow(col.R, 1.0/gamma), 0.0, 1.0)
	col.G = clamp(math.Pow(col.G, 1.0/gamma), 0.0, 1.0)
	col.B = clamp(math.Pow(col.B, 1.0/gamma), 0.0, 1.0)

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
	fx.DefaultRegistry.Register(NewGamma(1.0))
}
