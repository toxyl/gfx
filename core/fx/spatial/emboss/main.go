package fx

import (
	"image"
	stdcolor "image/color"
	"image/draw"
	"math"

	"github.com/toxyl/gfx/core/color"
	"github.com/toxyl/gfx/core/fx"
)

// EmbossFunction represents a function that creates an embossed effect
type EmbossFunction struct {
	*fx.BaseFunction
	amount float64
}

// Function arguments
var embossArgs = []fx.FunctionArg{
	{
		Name:        "amount",
		Type:        fx.TypeFloat,
		Description: "emboss adjustment value",
		Required:    true,
		Min:         -1.0,
		Max:         1.0,
		Step:        0.1,
	},
}

func NewEmboss(amount float64) *EmbossFunction {
	return &EmbossFunction{
		BaseFunction: fx.NewBaseFunction("emboss", "Applies emboss transformation to an image", color.New(0.5, 0.5, 0.5, 1), embossArgs),
		amount:       amount,
	}
}

// Apply implements the Function interface
func (f *EmbossFunction) Apply(img image.Image) (image.Image, error) {
	if img == nil {
		return nil, fx.ErrInvalidImage
	}
	bounds := img.Bounds()
	dst := image.NewRGBA(bounds)

	// Create temporary grayscale image
	gray := image.NewRGBA(bounds)
	draw.Draw(gray, bounds, img, bounds.Min, draw.Src)

	// Apply emboss effect
	for y := bounds.Min.Y + 1; y < bounds.Max.Y-1; y++ {
		for x := bounds.Min.X + 1; x < bounds.Max.X-1; x++ {
			// Get 3x3 neighborhood
			var gx, gy float64

			// Calculate horizontal gradient (Gx)
			gx += float64(getGray(gray.At(x-1, y-1))) * -1
			gx += float64(getGray(gray.At(x-1, y))) * -1
			gx += float64(getGray(gray.At(x-1, y+1))) * -1
			gx += float64(getGray(gray.At(x+1, y-1))) * 1
			gx += float64(getGray(gray.At(x+1, y))) * 1
			gx += float64(getGray(gray.At(x+1, y+1))) * 1

			// Calculate vertical gradient (Gy)
			gy += float64(getGray(gray.At(x-1, y-1))) * -1
			gy += float64(getGray(gray.At(x, y-1))) * -1
			gy += float64(getGray(gray.At(x+1, y-1))) * -1
			gy += float64(getGray(gray.At(x-1, y+1))) * 1
			gy += float64(getGray(gray.At(x, y+1))) * 1
			gy += float64(getGray(gray.At(x+1, y+1))) * 1

			// Calculate gradient magnitude and apply amount
			magnitude := math.Sqrt(gx*gx+gy*gy) * f.amount

			// Add 0.5 to center the values around gray
			value := magnitude + 0.5

			// Clamp the value
			value = math.Max(0, math.Min(1, value))

			// Set the pixel
			dst.Set(x, y, color.New(value, value, value, 1).ToUint8())
		}
	}

	return dst, nil
}

// getGray returns the grayscale value of a color
func getGray(c stdcolor.Color) uint8 {
	r, g, b, _ := c.RGBA()
	return uint8((0.299*float64(r) + 0.587*float64(g) + 0.114*float64(b)) / 256)
}

func init() {
	fx.DefaultRegistry.Register(NewEmboss(0.5))
}
