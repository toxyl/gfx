package fx

import (
	"image"
	stdcolor "image/color"

	"github.com/toxyl/gfx/core/color"
	"github.com/toxyl/gfx/core/fx"
)

// InvertFunction represents a function that inverts the colors of an image
type InvertFunction struct {
	*fx.BaseFunction
	amount float64
}

// Function arguments
var invertArgs = []fx.FunctionArg{
	{
		Name:        "amount",
		Type:        fx.TypeFloat,
		Description: "inversion adjustment value",
		Required:    true,
		Min:         0.0,
		Max:         1.0,
		Step:        0.1,
	},
}

// NewInvert creates a new color inversion function
func NewInvert(amount float64) *InvertFunction {
	return &InvertFunction{
		BaseFunction: fx.NewBaseFunction("invert", "Applies color inversion transformation to an image", color.New(0, 0, 0, 1), invertArgs),
		amount:       amount,
	}
}

// Apply implements the Function interface
func (f *InvertFunction) Apply(img image.Image) (image.Image, error) {
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

			// Invert RGB values
			// The formula: C' = 1 - C
			// This ensures that 0 becomes 1 and 1 becomes 0
			r := 1.0 - col.R
			g := 1.0 - col.G
			b := 1.0 - col.B

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

// ProcessPixel implements the Function interface
func (f *InvertFunction) ProcessPixel(x, y int, img image.Image) (*color.Color64, error) {
	if img == nil {
		return nil, fx.ErrInvalidImage
	}

	bounds := img.Bounds()
	if x < bounds.Min.X || x >= bounds.Max.X || y < bounds.Min.Y || y >= bounds.Max.Y {
		return nil, fx.ErrOutOfBounds
	}

	col := color.New(0, 0, 0, 1)
	rgba := img.At(x, y).(stdcolor.RGBA)
	col.SetUint8(rgba)

	// Invert only RGB channels, keep alpha at 1.0
	col.R = 1.0 - col.R
	col.G = 1.0 - col.G
	col.B = 1.0 - col.B
	col.A = 1.0

	return col, nil
}

func init() {
	fx.DefaultRegistry.Register(NewInvert(1.0))
}
