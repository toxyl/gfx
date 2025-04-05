package fx

import (
	"image"
	stdcolor "image/color"

	"github.com/toxyl/gfx/core/color"
	"github.com/toxyl/gfx/core/color/utils"
	"github.com/toxyl/gfx/core/fx"
)

// ThresholdFunction represents a function that applies thresholding to an image
type ThresholdFunction struct {
	*fx.BaseFunction
	amount float64 // Threshold value (0.0 to 1.0)
}

// Function arguments
var thresholdArgs = []fx.FunctionArg{
	{
		Name:        "amount",
		Type:        fx.TypeFloat,
		Description: "threshold adjustment value",
		Required:    true,
		Min:         0.0,
		Max:         1.0,
		Step:        0.1,
	},
}

// NewThreshold creates a new threshold function
func NewThreshold(amount float64) *ThresholdFunction {
	return &ThresholdFunction{
		BaseFunction: fx.NewBaseFunction("threshold", "Applies threshold transformation to an image", color.New(0, 0, 0, 1), thresholdArgs),
		amount:       amount,
	}
}

// Apply implements the ImageFunction interface
func (f *ThresholdFunction) Apply(img image.Image) (image.Image, error) {
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

// ProcessPixel applies the threshold effect to a single pixel
func (f *ThresholdFunction) ProcessPixel(x, y int, img image.Image) (*color.Color64, error) {
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

	// Convert to HSL to get luminance
	_, _, l := utils.RGBToHSL(col.R, col.G, col.B)

	// Apply threshold
	var value float64
	if l > f.amount {
		value = 1.0
	} else {
		value = 0.0
	}

	// Set all channels to the threshold value
	col.R = value
	col.G = value
	col.B = value

	return col, nil
}

func init() {
	fx.DefaultRegistry.Register(NewThreshold(0.5))
}
