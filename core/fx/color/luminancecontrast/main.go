package luminancecontrast

import (
	"image"
	stdcolor "image/color"

	"github.com/toxyl/gfx/core/color"
	"github.com/toxyl/gfx/core/fx"
)

// Function arguments
var luminancecontrastArgs = []fx.FunctionArg{
	{
		Name:        "amount",
		Type:        fx.TypeFloat,
		Description: "luminance contrast adjustment value",
		Required:    true,
		Min:         -1.0,
		Max:         1.0,
		Step:        0.1,
	},
}

// LuminanceContrastFunction represents a function to adjust the luminance contrast of an image
type LuminanceContrastFunction struct {
	*fx.BaseFunction
	amount float64
}

// NewLuminanceContrast creates a new luminance contrast function
func NewLuminanceContrast(amount float64) *LuminanceContrastFunction {
	return &LuminanceContrastFunction{
		BaseFunction: fx.NewBaseFunction("luminancecontrast", "Applies luminance contrast transformation to an image", color.New(0, 0, 0, 1), luminancecontrastArgs),
		amount:       amount,
	}
}

// Apply implements the Function interface
func (f *LuminanceContrastFunction) Apply(img *image.Image) (*image.Image, error) {
	if img == nil || *img == nil {
		return nil, fx.ErrInvalidImage
	}

	bounds := (*img).Bounds()
	result := image.NewRGBA(bounds)

	for y := bounds.Min.Y; y < bounds.Max.Y; y++ {
		for x := bounds.Min.X; x < bounds.Max.X; x++ {
			// Get the original color from the input image
			col := color.New(0, 0, 0, 1)
			col.SetUint8((*img).At(x, y).(stdcolor.RGBA))

			// Calculate luminance using standard weights
			// The formula: L = 0.299*R + 0.587*G + 0.114*B
			luminance := 0.299*col.R + 0.587*col.G + 0.114*col.B

			// Calculate contrast adjustment
			// The formula: L' = (L - 0.5) * (1 + amount) + 0.5
			// This ensures that 0.5 (middle gray) stays at 0.5
			adjustedLuminance := (luminance-0.5)*(1.0+f.amount) + 0.5

			// Clamp the luminance value
			adjustedLuminance = max(0.0, min(1.0, adjustedLuminance))

			// Calculate the ratio between adjusted and original luminance
			ratio := adjustedLuminance / luminance
			if luminance == 0 {
				ratio = 1.0
			}

			// Apply the ratio to each color channel
			r := col.R * ratio
			g := col.G * ratio
			b := col.B * ratio

			// Clamp values to [0,1] range
			r = max(0.0, min(1.0, r))
			g = max(0.0, min(1.0, g))
			b = max(0.0, min(1.0, b))

			// Set the result pixel
			result.Set(x, y, stdcolor.RGBA{
				R: uint8(r * 255),
				G: uint8(g * 255),
				B: uint8(b * 255),
				A: uint8(col.A * 255),
			})
		}
	}

	var resultImage image.Image = result
	return &resultImage, nil
}

func max(a, b float64) float64 {
	if a > b {
		return a
	}
	return b
}

func min(a, b float64) float64 {
	if a < b {
		return a
	}
	return b
}

func (f *LuminanceContrastFunction) ProcessPixel(x, y int, img *image.Image) (*color.Color64, error) {
	if img == nil || *img == nil {
		return nil, fx.ErrInvalidImage
	}

	bounds := (*img).Bounds()
	if x < bounds.Min.X || x >= bounds.Max.X || y < bounds.Min.Y || y >= bounds.Max.Y {
		return nil, fx.ErrOutOfBounds
	}

	// Get the original color from the input image
	col := color.New(0, 0, 0, 1)
	col.SetUint8((*img).At(x, y).(stdcolor.RGBA))

	// Calculate luminance using standard weights
	// The formula: L = 0.299*R + 0.587*G + 0.114*B
	luminance := 0.299*col.R + 0.587*col.G + 0.114*col.B

	// Calculate contrast adjustment
	// The formula: L' = (L - 0.5) * (1 + amount) + 0.5
	// This ensures that 0.5 (middle gray) stays at 0.5
	adjustedLuminance := (luminance-0.5)*(1.0+f.amount) + 0.5

	// Clamp the luminance value
	adjustedLuminance = max(0.0, min(1.0, adjustedLuminance))

	// Calculate the ratio between adjusted and original luminance
	ratio := adjustedLuminance / luminance
	if luminance == 0 {
		ratio = 1.0
	}

	// Apply the ratio to each color channel
	r := col.R * ratio
	g := col.G * ratio
	b := col.B * ratio

	// Clamp values to [0,1] range
	r = max(0.0, min(1.0, r))
	g = max(0.0, min(1.0, g))
	b = max(0.0, min(1.0, b))

	// Set the result pixel
	return color.New(r, g, b, col.A), nil
}
