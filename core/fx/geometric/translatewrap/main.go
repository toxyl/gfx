package fx

import (
	"image"
	"image/draw"

	"github.com/toxyl/gfx/core/color"
	"github.com/toxyl/gfx/core/fx"
)

// filters/translatewrap/main.go

// TranslateWrapFunction represents a function that translates an image with wrapping
type TranslateWrapFunction struct {
	*fx.BaseFunction
	offsetX, offsetY float64 // Translation offsets
}

// Function arguments
var translateWrapArgs = []fx.FunctionArg{
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

// NewTranslateWrap creates a new translation function with wrapping
func NewTranslateWrap(offsetX, offsetY float64) *TranslateWrapFunction {
	return &TranslateWrapFunction{
		BaseFunction: fx.NewBaseFunction("translatewrap", "Applies translation transformation to an image with wrapping", color.New(0, 0, 0, 1), translateWrapArgs),
		offsetX:      offsetX,
		offsetY:      offsetY,
	}
}

// Apply implements the Function interface
func (f *TranslateWrapFunction) Apply(img image.Image) (image.Image, error) {
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

			// Wrap around edges
			srcX = (srcX-bounds.Min.X)%bounds.Dx() + bounds.Min.X
			if srcX < bounds.Min.X {
				srcX += bounds.Dx()
			}
			srcY = (srcY-bounds.Min.Y)%bounds.Dy() + bounds.Min.Y
			if srcY < bounds.Min.Y {
				srcY += bounds.Dy()
			}

			// Copy pixel from source
			dst.Set(x, y, temp.At(srcX, srcY))
		}
	}

	return dst, nil
}

func init() {
	fx.DefaultRegistry.Register(NewTranslateWrap(0, 0))
}
