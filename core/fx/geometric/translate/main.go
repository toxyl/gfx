package fx

import (
	"image"
	"image/draw"

	"github.com/toxyl/gfx/core/color"
	"github.com/toxyl/gfx/core/fx"
)

// filters/translate/main.go

// TranslateFunction represents a function that translates an image
type TranslateFunction struct {
	*fx.BaseFunction
	offsetX, offsetY float64 // Translation offsets
}

// Function arguments
var translateArgs = []fx.FunctionArg{
	{
		Name:        "amount",
		Type:        fx.TypeFloat,
		Description: "translate adjustment value",
		Required:    true,
		Min:         -1.0,
		Max:         1.0,
		Step:        0.1,
	},
}

func NewTranslate(offsetX, offsetY float64) *TranslateFunction {
	return &TranslateFunction{
		BaseFunction: fx.NewBaseFunction("translate", "Applies translation transformation to an image", color.New(0, 0, 0, 1), translateArgs),
		offsetX:      offsetX,
		offsetY:      offsetY,
	}
}

// Apply implements the Function interface
func (f *TranslateFunction) Apply(img image.Image) (image.Image, error) {
	if img == nil {
		return nil, fx.ErrInvalidImage
	}
	bounds := img.Bounds()

	// Calculate translation in pixels
	dx := int(float64(bounds.Dx()) * f.offsetX)
	dy := int(float64(bounds.Dy()) * f.offsetY)

	// Create new image with same dimensions
	dst := image.NewRGBA(bounds)

	// Create temporary image for processing
	temp := image.NewRGBA(bounds)
	draw.Draw(temp, bounds, img, bounds.Min, draw.Src)

	// Process each pixel
	for y := bounds.Min.Y; y < bounds.Max.Y; y++ {
		for x := bounds.Min.X; x < bounds.Max.X; x++ {
			// Calculate source coordinates
			srcX := x - dx
			srcY := y - dy

			// Skip if source coordinates are out of bounds
			if srcX < bounds.Min.X || srcX >= bounds.Max.X ||
				srcY < bounds.Min.Y || srcY >= bounds.Max.Y {
				continue
			}

			// Copy pixel from source
			dst.Set(x, y, temp.At(srcX, srcY))
		}
	}

	return dst, nil
}

func init() {
	fx.DefaultRegistry.Register(NewTranslate(0, 0))
}
