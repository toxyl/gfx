package fx

import (
	"image"
	stdcolor "image/color"
	"image/draw"

	"github.com/toxyl/gfx/core/color"
	"github.com/toxyl/gfx/core/fx"
	"github.com/toxyl/math"
)

// LuminanceContrastFunction represents a function that adjusts the contrast of luminance values
type LumContrastFunction struct {
	*fx.BaseFunction
	amount float64 // Contrast adjustment amount (-1.0 to 1.0)
}

// Function arguments
var lumcontrastArgs = []fx.FunctionArg{
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

// NewLumContrast creates a new luminance contrast function
func NewLumContrast(amount float64) *LumContrastFunction {
	return &LumContrastFunction{
		BaseFunction: fx.NewBaseFunction("lumcontrast", "Applies luminance contrast transformation to an image", color.New(0, 0, 0, 1), lumcontrastArgs),
		amount:       amount,
	}
}

// Apply implements the Function interface
func (f *LumContrastFunction) Apply(img *image.Image) (*image.Image, error) {
	if img == nil || *img == nil {
		return nil, fx.ErrInvalidImage
	}

	bounds := (*img).Bounds()
	dst := image.NewRGBA(bounds)

	// Create temporary image for processing
	temp := image.NewRGBA(bounds)
	draw.Draw(temp, bounds, *img, bounds.Min, draw.Src)

	// Adjust luminance contrast for each pixel
	for y := bounds.Min.Y; y < bounds.Max.Y; y++ {
		for x := bounds.Min.X; x < bounds.Max.X; x++ {
			// Get current pixel color
			col := color.New(0, 0, 0, 1)
			col.SetUint8(temp.At(x, y).(stdcolor.RGBA))

			// Convert to HSL
			h, s, l := rgbToHsl(col.R, col.G, col.B)

			// Adjust luminance contrast using the contrast curve
			l = contrastCurve(l, f.amount)

			// Convert back to RGB
			r, g, b := hslToRgb(h, s, l)

			// Set the pixel
			dst.Set(x, y, color.New(r, g, b, col.A).ToUint8())
		}
	}

	*img = dst
	return img, nil
}

// contrastCurve applies a contrast curve to a value
func contrastCurve(value, amount float64) float64 {
	// Center the value around 0.5
	value = value - 0.5

	// Apply contrast curve
	// When amount is positive, values are pushed away from 0.5
	// When amount is negative, values are pulled towards 0.5
	value = value * (1.0 + amount)

	// Move back to [0,1] range
	value = value + 0.5

	// Clamp to valid range
	return math.Clamp(value, 0.0, 1.0)
}

// rgbToHsl converts RGB values to HSL
func rgbToHsl(r, g, b float64) (h, s, l float64) {
	max := math.Max(r, math.Max(g, b))
	min := math.Min(r, math.Min(g, b))
	l = (max + min) / 2

	if max == min {
		return 0, 0, l
	}

	if l > 0.5 {
		s = (max - min) / (2 - max - min)
	} else {
		s = (max - min) / (max + min)
	}

	switch max {
	case r:
		h = (g - b) / (max - min)
		if g < b {
			h += 6
		}
	case g:
		h = (b-r)/(max-min) + 2
	case b:
		h = (r-g)/(max-min) + 4
	}
	h /= 6

	return h, s, l
}

// hslToRgb converts HSL values to RGB
func hslToRgb(h, s, l float64) (r, g, b float64) {
	if s == 0 {
		return l, l, l
	}

	var q float64
	if l < 0.5 {
		q = l * (1 + s)
	} else {
		q = l + s - l*s
	}
	p := 2*l - q

	r = hueToRgb(p, q, h+1.0/3.0)
	g = hueToRgb(p, q, h)
	b = hueToRgb(p, q, h-1.0/3.0)

	return r, g, b
}

// hueToRgb is a helper function for hslToRgb
func hueToRgb(p, q, t float64) float64 {
	if t < 0 {
		t += 1
	}
	if t > 1 {
		t -= 1
	}
	if t < 1.0/6.0 {
		return p + (q-p)*6*t
	}
	if t < 1.0/2.0 {
		return q
	}
	if t < 2.0/3.0 {
		return p + (q-p)*(2.0/3.0-t)*6
	}
	return p
}

func (f *LumContrastFunction) ProcessPixel(x, y int, img *image.Image) (*color.Color64, error) {
	if img == nil || *img == nil {
		return nil, fx.ErrInvalidImage
	}

	bounds := (*img).Bounds()
	if x < bounds.Min.X || x >= bounds.Max.X || y < bounds.Min.Y || y >= bounds.Max.Y {
		return nil, fx.ErrOutOfBounds
	}

	// Get current pixel color
	col := color.New(0, 0, 0, 1)
	col.SetUint8((*img).At(x, y).(stdcolor.RGBA))

	// Convert to HSL
	h, s, l := rgbToHsl(col.R, col.G, col.B)

	// Adjust luminance contrast using the contrast curve
	l = contrastCurve(l, f.amount)

	// Convert back to RGB
	r, g, b := hslToRgb(h, s, l)

	// Return new color
	return color.New(r, g, b, col.A), nil
}
