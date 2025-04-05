package fx

import (
	"image"
	stdcolor "image/color"

	"github.com/toxyl/gfx/core/color"
	"github.com/toxyl/gfx/core/fx"
)

// SepiaFunction represents a function that applies a sepia tone effect
type SepiaFunction struct {
	*fx.BaseFunction
	amount float64 // Sepia effect amount (0.0 to 1.0)
}

// Function arguments
var sepiaArgs = []fx.FunctionArg{
	{
		Name:        "amount",
		Type:        fx.TypeFloat,
		Description: "sepia adjustment value",
		Required:    true,
		Min:         0.0,
		Max:         1.0,
		Step:        0.1,
	},
}

// NewSepia creates a new sepia function
func NewSepia(amount float64) *SepiaFunction {
	return &SepiaFunction{
		BaseFunction: fx.NewBaseFunction("sepia", "Applies sepia transformation to an image", color.New(0, 0, 0, 1), sepiaArgs),
		amount:       amount,
	}
}

// Apply implements the ImageFunction interface
func (f *SepiaFunction) Apply(img image.Image) (image.Image, error) {
	if img == nil {
		return nil, fx.ErrInvalidImage
	}

	bounds := img.Bounds()
	result := image.NewRGBA(bounds)

	for y := bounds.Min.Y; y < bounds.Max.Y; y++ {
		for x := bounds.Min.X; x < bounds.Max.X; x++ {
			col, err := f.ProcessPixel(x, y, img)
			if err != nil {
				return nil, err
			}

			result.Set(x, y, stdcolor.RGBA{
				R: uint8(col.R * 255),
				G: uint8(col.G * 255),
				B: uint8(col.B * 255),
				A: uint8(col.A * 255),
			})
		}
	}

	return result, nil
}

// ProcessPixel applies the sepia effect to a single pixel
func (f *SepiaFunction) ProcessPixel(x, y int, img image.Image) (*color.Color64, error) {
	if img == nil {
		return nil, fx.ErrInvalidImage
	}

	bounds := img.Bounds()
	if x < bounds.Min.X || x >= bounds.Max.X || y < bounds.Min.Y || y >= bounds.Max.Y {
		return nil, fx.ErrOutOfBounds
	}

	// Get the original color from the input image
	col := color.New(0, 0, 0, 1)
	col.SetUint8(img.At(x, y).(stdcolor.RGBA))

	// Calculate sepia tone using matrix multiplication
	// The formula:
	// R' = (R * 0.393) + (G * 0.769) + (B * 0.189)
	// G' = (R * 0.349) + (G * 0.686) + (B * 0.168)
	// B' = (R * 0.272) + (G * 0.534) + (B * 0.131)
	r := col.R*0.393 + col.G*0.769 + col.B*0.189
	g := col.R*0.349 + col.G*0.686 + col.B*0.168
	b := col.R*0.272 + col.G*0.534 + col.B*0.131

	// Interpolate between original and sepia colors based on amount
	col.R = lerp(col.R, clamp(r, 0.0, 1.0), f.amount)
	col.G = lerp(col.G, clamp(g, 0.0, 1.0), f.amount)
	col.B = lerp(col.B, clamp(b, 0.0, 1.0), f.amount)

	return col, nil
}

// lerp performs linear interpolation between a and b based on t
func lerp(a, b, t float64) float64 {
	return a + t*(b-a)
}

// clamp restricts a value to the range [min, max]
func clamp(value, min, max float64) float64 {
	if value < min {
		return min
	}
	if value > max {
		return max
	}
	return value
}

func init() {
	fx.DefaultRegistry.Register(NewSepia(0.5))
}
