package fx

import (
	"image"
	"image/draw"

	"github.com/toxyl/gfx/core/color"
	"github.com/toxyl/gfx/core/fx"
)

// FlipVFunction represents a function that flips an image vertically
type FlipVFunction struct {
	*fx.BaseFunction
}

// Function arguments
var flipvArgs = []fx.FunctionArg{
	{
		Name:        "amount",
		Type:        fx.TypeFloat,
		Description: "flipv adjustment value",
		Required:    true,
		Min:         -1.0,
		Max:         1.0,
		Step:        0.1,
	},
}

// NewFlipV creates a new vertical flip function
func NewFlipV() *FlipVFunction {
	return &FlipVFunction{
		BaseFunction: fx.NewBaseFunction("flipv", "Applies vertical flip transformation to an image", color.New(0, 0, 0, 1), flipvArgs),
	}
}

// Apply implements the Function interface
func (f *FlipVFunction) Apply(img image.Image) (image.Image, error) {
	if img == nil {
		return nil, fx.ErrInvalidImage
	}
	bounds := img.Bounds()
	dst := image.NewRGBA(bounds)

	// Create temporary image for processing
	temp := image.NewRGBA(bounds)
	draw.Draw(temp, bounds, img, bounds.Min, draw.Src)

	// Flip vertically
	for y := bounds.Min.Y; y < bounds.Max.Y; y++ {
		for x := bounds.Min.X; x < bounds.Max.X; x++ {
			dst.Set(x, bounds.Max.Y-1-y, temp.At(x, y))
		}
	}

	return dst, nil
}

func init() {
	fx.DefaultRegistry.Register(NewFlipV())
}
