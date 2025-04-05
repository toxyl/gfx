package fx

import (
	"image"
	"image/draw"

	"github.com/toxyl/gfx/core/color"
	"github.com/toxyl/gfx/core/fx"
)

// FlipHFunction represents a function that flips an image horizontally
type FlipHFunction struct {
	*fx.BaseFunction
}

// Function arguments
var fliphArgs = []fx.FunctionArg{
	{
		Name:        "amount",
		Type:        fx.TypeFloat,
		Description: "fliph adjustment value",
		Required:    true,
		Min:         -1.0,
		Max:         1.0,
		Step:        0.1,
	},
}

// NewFlipH creates a new horizontal flip function
func NewFlipH() *FlipHFunction {
	return &FlipHFunction{
		BaseFunction: fx.NewBaseFunction("fliph", "Applies horizontal flip transformation to an image", color.New(0, 0, 0, 1), fliphArgs),
	}
}

// Apply implements the Function interface
func (f *FlipHFunction) Apply(img image.Image) (image.Image, error) {
	if img == nil {
		return nil, fx.ErrInvalidImage
	}
	bounds := img.Bounds()
	dst := image.NewRGBA(bounds)

	// Create temporary image for processing
	temp := image.NewRGBA(bounds)
	draw.Draw(temp, bounds, img, bounds.Min, draw.Src)

	// Flip horizontally
	for y := bounds.Min.Y; y < bounds.Max.Y; y++ {
		for x := bounds.Min.X; x < bounds.Max.X; x++ {
			dst.Set(bounds.Max.X-1-x, y, temp.At(x, y))
		}
	}

	return dst, nil
}

func init() {
	fx.DefaultRegistry.Register(NewFlipH())
}
