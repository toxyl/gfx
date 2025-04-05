package fx

import (
	"image"
	stdcolor "image/color"
	"image/draw"

	"github.com/toxyl/gfx/core/color"
	"github.com/toxyl/gfx/core/fx"
)

// SharpenFunction represents a function that sharpens an image
type SharpenFunction struct {
	*fx.BaseFunction
	amount float64 // Sharpening amount (0.0 to 1.0)
}

// Function arguments
var sharpenArgs = []fx.FunctionArg{
	{
		Name:        "amount",
		Type:        fx.TypeFloat,
		Description: "sharpen adjustment value",
		Required:    true,
		Min:         -1.0,
		Max:         1.0,
		Step:        0.1,
	},
}

func NewSharpen(amount float64) *SharpenFunction {
	return &SharpenFunction{
		BaseFunction: fx.NewBaseFunction("sharpen", "Applies sharpen transformation to an image", color.New(0, 0, 0, 1), sharpenArgs),
		amount:       amount,
	}
}

// Apply implements the Function interface
func (f *SharpenFunction) Apply(img image.Image) (image.Image, error) {
	if img == nil {
		return nil, fx.ErrInvalidImage
	}
	bounds := img.Bounds()
	dst := image.NewRGBA(bounds)

	// Create temporary image for processing
	temp := image.NewRGBA(bounds)
	draw.Draw(temp, bounds, img, bounds.Min, draw.Src)

	// Define sharpening kernel
	kernel := [3][3]float64{
		{0, -1, 0},
		{-1, 5, -1},
		{0, -1, 0},
	}

	// Process each pixel
	for y := bounds.Min.Y; y < bounds.Max.Y; y++ {
		for x := bounds.Min.X; x < bounds.Max.X; x++ {
			// Initialize accumulators for each channel
			var rSum, gSum, bSum float64

			// Apply kernel
			for ky := -1; ky <= 1; ky++ {
				for kx := -1; kx <= 1; kx++ {
					// Get source coordinates
					srcX := x + kx
					srcY := y + ky

					// Handle edge cases by clamping coordinates
					if srcX < bounds.Min.X {
						srcX = bounds.Min.X
					} else if srcX >= bounds.Max.X {
						srcX = bounds.Max.X - 1
					}
					if srcY < bounds.Min.Y {
						srcY = bounds.Min.Y
					} else if srcY >= bounds.Max.Y {
						srcY = bounds.Max.Y - 1
					}

					// Get source pixel color
					col := color.New(0, 0, 0, 1)
					col.SetUint8(temp.At(srcX, srcY).(stdcolor.RGBA))

					// Apply kernel weight
					weight := kernel[ky+1][kx+1]
					rSum += col.R * weight
					gSum += col.G * weight
					bSum += col.B * weight
				}
			}

			// Get original color
			origCol := color.New(0, 0, 0, 1)
			origCol.SetUint8(temp.At(x, y).(stdcolor.RGBA))

			// Blend sharpened result with original based on amount
			r := clamp(rSum*f.amount + origCol.R*(1.0-f.amount))
			g := clamp(gSum*f.amount + origCol.G*(1.0-f.amount))
			b := clamp(bSum*f.amount + origCol.B*(1.0-f.amount))

			// Set the pixel
			dst.Set(x, y, color.New(r, g, b, origCol.A).ToUint8())
		}
	}

	return dst, nil
}

// clamp ensures a value is between 0 and 1
func clamp(v float64) float64 {
	if v < 0 {
		return 0
	}
	if v > 1 {
		return 1
	}
	return v
}

func init() {
	fx.DefaultRegistry.Register(NewSharpen(0.5))
}
