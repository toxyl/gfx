package fx

import (
	"image"
	stdcolor "image/color"

	"github.com/toxyl/gfx/core/color"
	"github.com/toxyl/gfx/core/color/utils"
	"github.com/toxyl/gfx/core/fx"
)

// HueFunction represents a function that adjusts the hue of an image
type HueFunction struct {
	*fx.BaseFunction
	amount float64 // Hue adjustment amount (-1.0 to 1.0)
}

// Function arguments
var hueArgs = []fx.FunctionArg{
	{
		Name:        "amount",
		Type:        fx.TypeFloat,
		Description: "hue adjustment value",
		Required:    true,
		Min:         -1.0,
		Max:         1.0,
		Step:        0.1,
	},
}

func NewHue(amount float64) *HueFunction {
	return &HueFunction{
		BaseFunction: fx.NewBaseFunction("hue", "Applies hue transformation to an image", color.New(0, 0, 0, 1), hueArgs),
		amount:       amount,
	}
}

// Apply implements the ImageFunction interface
func (f *HueFunction) Apply(img image.Image) (image.Image, error) {
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

			// Adjust hue
			// The formula: H' = H + amount * 360
			// This ensures that hue wraps around correctly
			h = h + f.amount*360
			if h < 0 {
				h += 360
			}
			if h >= 360 {
				h -= 360
			}

			// Convert back to RGB
			r, g, b := utils.HSLToRGB(h, s, l)

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

func init() {
	fx.DefaultRegistry.Register(NewHue(0))
}
