package fx

import (
	"image"
	stdcolor "image/color"

	"github.com/toxyl/gfx/core/color"
	"github.com/toxyl/gfx/core/color/utils"
	"github.com/toxyl/gfx/core/fx"
)

// LuminanceFunction represents a function to adjust the luminance of an image
type LuminanceFunction struct {
	*fx.BaseFunction
	amount float64 // Luminance adjustment amount (-1.0 to 1.0)
}

// Function arguments
var luminanceArgs = []fx.FunctionArg{
	{
		Name:        "amount",
		Type:        fx.TypeFloat,
		Description: "luminance adjustment value",
		Required:    true,
		Min:         -1.0,
		Max:         1.0,
		Step:        0.1,
	},
}

// NewLuminance creates a new luminance function
func NewLuminance(amount float64) *LuminanceFunction {
	return &LuminanceFunction{
		BaseFunction: fx.NewBaseFunction("luminance", "Applies luminance transformation to an image", color.New(0, 0, 0, 1), luminanceArgs),
		amount:       amount,
	}
}

// Apply implements the ImageFunction interface
func (f *LuminanceFunction) Apply(img image.Image) (image.Image, error) {
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

// ProcessPixel applies the luminance effect to a single pixel
func (f *LuminanceFunction) ProcessPixel(x, y int, img image.Image) (*color.Color64, error) {
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

	// Convert to HSL to adjust luminance
	h, s, l := utils.RGBToHSL(col.R, col.G, col.B)

	// Adjust luminance
	// For positive amount, we increase luminance while preserving headroom
	// For negative amount, we decrease luminance while preserving shadows
	if f.amount >= 0 {
		l = l + (1-l)*f.amount
	} else {
		l = l * (1 + f.amount)
	}

	// Convert back to RGB
	r, g, b := utils.HSLToRGB(h, s, l)
	col.R = r
	col.G = g
	col.B = b

	return col, nil
}

func init() {
	fx.DefaultRegistry.Register(NewLuminance(0))
}
