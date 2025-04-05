package fx

import (
	"image"
	stdcolor "image/color"
	"image/draw"

	"github.com/toxyl/gfx/core/color"
	"github.com/toxyl/gfx/core/fx"
)

// GrayFunction represents a function that converts an image to grayscale
type GrayFunction struct {
	*fx.BaseFunction
	amount float64
}

// Function arguments
var grayArgs = []fx.FunctionArg{
	{
		Name:        "amount",
		Type:        fx.TypeFloat,
		Description: "grayscale adjustment value",
		Required:    true,
		Min:         0.0,
		Max:         1.0,
		Step:        0.1,
	},
}

func NewGray(amount float64) *GrayFunction {
	return &GrayFunction{
		BaseFunction: fx.NewBaseFunction("gray", "Applies grayscale transformation to an image", color.New(0, 0, 0, 1), grayArgs),
		amount:       amount,
	}
}

// Apply implements the Function interface
func (f *GrayFunction) Apply(img image.Image) (image.Image, error) {
	if img == nil {
		return nil, fx.ErrInvalidImage
	}
	bounds := img.Bounds()
	dst := image.NewRGBA(bounds)

	// Create temporary image for processing
	temp := image.NewRGBA(bounds)
	draw.Draw(temp, bounds, img, bounds.Min, draw.Src)

	// Convert to grayscale
	for y := bounds.Min.Y; y < bounds.Max.Y; y++ {
		for x := bounds.Min.X; x < bounds.Max.X; x++ {
			// Get current pixel color
			col := color.New(0, 0, 0, 1)
			col.SetUint8(temp.At(x, y).(stdcolor.RGBA))

			// Calculate grayscale value using luminance formula
			gray := 0.299*col.R + 0.587*col.G + 0.114*col.B

			// Set the pixel
			dst.Set(x, y, color.New(gray, gray, gray, col.A).ToUint8())
		}
	}

	return dst, nil
}

func (f *GrayFunction) ProcessPixel(x, y int, img image.Image) (*color.Color64, error) {
	if img == nil {
		return nil, fx.ErrInvalidImage
	}

	bounds := img.Bounds()
	if x < bounds.Min.X || x >= bounds.Max.X || y < bounds.Min.Y || y >= bounds.Max.Y {
		return nil, fx.ErrOutOfBounds
	}

	// Get current pixel color
	col := color.New(0, 0, 0, 1)
	col.SetUint8(img.At(x, y).(stdcolor.RGBA))

	// Calculate grayscale value using luminance formula
	gray := 0.299*col.R + 0.587*col.G + 0.114*col.B

	return color.New(gray, gray, gray, col.A), nil
}

func init() {
	fx.DefaultRegistry.Register(NewGray(0.5))
}
