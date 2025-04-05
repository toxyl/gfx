package fx

import (
	"image"
	"image/draw"
	"math/rand"

	"github.com/toxyl/gfx/core/color"
	"github.com/toxyl/gfx/core/fx"
)

// NoiseFunction represents a function that adds noise to an image
type NoiseFunction struct {
	*fx.BaseFunction
	amount float64 // Noise amount (0.0 to 1.0)
}

// Function arguments
var noiseArgs = []fx.FunctionArg{
	{
		Name:        "amount",
		Type:        fx.TypeFloat,
		Description: "noise adjustment value",
		Required:    true,
		Min:         -1.0,
		Max:         1.0,
		Step:        0.1,
	},
}

func NewNoise(amount float64) *NoiseFunction {
	return &NoiseFunction{
		BaseFunction: fx.NewBaseFunction("noise", "Applies noise transformation to an image", color.New(0, 0, 0, 1), noiseArgs),
		amount:       amount,
	}
}

// Apply implements the Function interface
func (f *NoiseFunction) Apply(img image.Image) (image.Image, error) {
	if img == nil {
		return nil, fx.ErrInvalidImage
	}
	bounds := img.Bounds()
	dst := image.NewRGBA(bounds)

	// Create temporary image for processing
	temp := image.NewRGBA(bounds)
	draw.Draw(temp, bounds, img, bounds.Min, draw.Src)

	// Process each pixel
	for y := bounds.Min.Y; y < bounds.Max.Y; y++ {
		for x := bounds.Min.X; x < bounds.Max.X; x++ {
			col, err := f.ProcessPixel(x, y, temp)
			if err != nil {
				return nil, err
			}
			dst.Set(x, y, col.ToUint8())
		}
	}

	return dst, nil
}

// ProcessPixel implements the ImageFunction interface
func (f *NoiseFunction) ProcessPixel(x, y int, img image.Image) (*color.Color64, error) {
	if img == nil {
		return nil, fx.ErrInvalidImage
	}

	bounds := img.Bounds()
	if x < bounds.Min.X || x >= bounds.Max.X || y < bounds.Min.Y || y >= bounds.Max.Y {
		return nil, fx.ErrOutOfBounds
	}

	// Get the original color
	original := img.At(x, y)
	r, g, b, a := original.RGBA()

	// Add noise to each channel
	noise := (rand.Float64()*2 - 1) * f.amount
	r = uint32(float64(r) * (1 + noise))
	g = uint32(float64(g) * (1 + noise))
	b = uint32(float64(b) * (1 + noise))

	// Process the color
	return color.New(
		float64(r)/65535.0,
		float64(g)/65535.0,
		float64(b)/65535.0,
		float64(a)/65535.0,
	), nil
}

func init() {
	fx.DefaultRegistry.Register(NewNoise(0.1))
}
