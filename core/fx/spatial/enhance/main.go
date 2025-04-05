package fx

import (
	"image"
	stdcolor "image/color"
	"image/draw"
	"math"

	"github.com/toxyl/gfx/core/color"
	"github.com/toxyl/gfx/core/fx"
)

// EnhanceFunction represents a function that enhances image details and contrast
type EnhanceFunction struct {
	*fx.BaseFunction
	amount float64
}

// Function arguments
var enhanceArgs = []fx.FunctionArg{
	{
		Name:        "amount",
		Type:        fx.TypeFloat,
		Description: "enhance adjustment value",
		Required:    true,
		Min:         -1.0,
		Max:         1.0,
		Step:        0.1,
	},
}

func NewEnhance(amount float64) *EnhanceFunction {
	return &EnhanceFunction{
		BaseFunction: fx.NewBaseFunction("enhance", "Applies enhance transformation to an image", color.New(0, 0, 0, 1), enhanceArgs),
		amount:       amount,
	}
}

// Apply implements the Function interface
func (f *EnhanceFunction) Apply(img *image.Image) (*image.Image, error) {
	if img == nil || *img == nil {
		return nil, fx.ErrInvalidImage
	}
	bounds := (*img).Bounds()
	dst := image.NewRGBA(bounds)

	// Create temporary image for processing
	temp := image.NewRGBA(bounds)
	draw.Draw(temp, bounds, *img, bounds.Min, draw.Src)

	// Apply enhancement
	for y := bounds.Min.Y + 1; y < bounds.Max.Y-1; y++ {
		for x := bounds.Min.X + 1; x < bounds.Max.X-1; x++ {
			// Get 3x3 neighborhood
			var r, g, b float64
			var count int

			// Calculate average color in neighborhood
			for dy := -1; dy <= 1; dy++ {
				for dx := -1; dx <= 1; dx++ {
					col := color.New(0, 0, 0, 1)
					col.SetUint8(temp.At(x+dx, y+dy).(stdcolor.RGBA))

					r += col.R
					g += col.G
					b += col.B
					count++
				}
			}

			// Calculate average
			avgR := r / float64(count)
			avgG := g / float64(count)
			avgB := b / float64(count)

			// Get current pixel color
			col := color.New(0, 0, 0, 1)
			col.SetUint8(temp.At(x, y).(stdcolor.RGBA))

			// Enhance by moving away from average
			newR := col.R + (col.R-avgR)*f.amount
			newG := col.G + (col.G-avgG)*f.amount
			newB := col.B + (col.B-avgB)*f.amount

			// Clamp values
			newR = math.Max(0, math.Min(1, newR))
			newG = math.Max(0, math.Min(1, newG))
			newB = math.Max(0, math.Min(1, newB))

			// Set the pixel
			dst.Set(x, y, color.New(newR, newG, newB, col.A).ToUint8())
		}
	}

	result := image.Image(dst)
	return &result, nil
}
