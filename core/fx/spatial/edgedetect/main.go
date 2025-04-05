package fx

import (
	"image"
	stdcolor "image/color"
	"image/draw"
	"math"

	"github.com/toxyl/gfx/core/color"
	"github.com/toxyl/gfx/core/fx"
)

// EdgeDetectFunction represents a function that detects edges in an image
type EdgeDetectFunction struct {
	*fx.BaseFunction
	threshold float64
}

// Function arguments
var edgedetectArgs = []fx.FunctionArg{
	{
		Name:        "amount",
		Type:        fx.TypeFloat,
		Description: "edge detection threshold",
		Required:    true,
		Min:         -1.0,
		Max:         1.0,
		Step:        0.1,
	},
}

// NewEdgeDetect creates a new edge detection function
func NewEdgeDetect(threshold float64) *EdgeDetectFunction {
	return &EdgeDetectFunction{
		BaseFunction: fx.NewBaseFunction("edgedetect", "Applies edge detection to an image", color.New(0, 0, 0, 1), edgedetectArgs),
		threshold:    threshold,
	}
}

// Apply implements the Function interface
func (f *EdgeDetectFunction) Apply(img image.Image) (image.Image, error) {
	if img == nil {
		return nil, fx.ErrInvalidImage
	}
	bounds := img.Bounds()
	dst := image.NewRGBA(bounds)

	// Create temporary grayscale image
	gray := image.NewRGBA(bounds)
	draw.Draw(gray, bounds, img, bounds.Min, draw.Src)

	// Apply Sobel operator
	for y := bounds.Min.Y + 1; y < bounds.Max.Y-1; y++ {
		for x := bounds.Min.X + 1; x < bounds.Max.X-1; x++ {
			// Get 3x3 neighborhood
			var gx, gy float64

			// Calculate horizontal gradient (Gx)
			gx += float64(getGray(gray.At(x-1, y-1))) * -1
			gx += float64(getGray(gray.At(x-1, y))) * -2
			gx += float64(getGray(gray.At(x-1, y+1))) * -1
			gx += float64(getGray(gray.At(x+1, y-1))) * 1
			gx += float64(getGray(gray.At(x+1, y))) * 2
			gx += float64(getGray(gray.At(x+1, y+1))) * 1

			// Calculate vertical gradient (Gy)
			gy += float64(getGray(gray.At(x-1, y-1))) * -1
			gy += float64(getGray(gray.At(x, y-1))) * -2
			gy += float64(getGray(gray.At(x+1, y-1))) * -1
			gy += float64(getGray(gray.At(x-1, y+1))) * 1
			gy += float64(getGray(gray.At(x, y+1))) * 2
			gy += float64(getGray(gray.At(x+1, y+1))) * 1

			// Calculate gradient magnitude
			magnitude := math.Sqrt(gx*gx + gy*gy)

			// Apply threshold
			if magnitude > f.threshold {
				dst.Set(x, y, color.New(1, 1, 1, 1).ToUint8())
			} else {
				dst.Set(x, y, color.New(0, 0, 0, 1).ToUint8())
			}
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
	fx.DefaultRegistry.Register(NewEdgeDetect(128))
}
