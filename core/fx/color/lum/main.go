package fx

import (
	"image"
	stdcolor "image/color"
	"image/draw"

	"github.com/toxyl/gfx/core/color"
	"github.com/toxyl/gfx/core/fx"
	"github.com/toxyl/math"
)

// LuminanceFunction represents a function that adjusts the luminance of an image
type LumFunction struct {
	*fx.BaseFunction
	amount float64 // Luminance adjustment amount (-1.0 to 1.0)
}

// NewLum creates a new lum function
func NewLum(amount float64) *LumFunction {
	return &LumFunction{
		BaseFunction: fx.NewBaseFunction("lum", "Applies luminance transformation to an image", color.New(0, 0, 0, 1), lumArgs),
		amount:       amount,
	}
}

// Function arguments
var lumArgs = []fx.FunctionArg{
	{
		Name:        "amount",
		Type:        fx.TypeFloat,
		Description: "luminance adjustment value",
		Required:    true,
		Min:         -1.0,
		Max:         1.0,
		Step:        0.1,
	},
}

// Apply implements the Function interface
func (f *LumFunction) Apply(img *image.Image) (*image.Image, error) {
	if img == nil || *img == nil {
		return nil, fx.ErrInvalidImage
	}

	bounds := (*img).Bounds()
	dst := image.NewRGBA(bounds)

	// Create temporary image for processing
	temp := image.NewRGBA(bounds)
	draw.Draw(temp, bounds, *img, bounds.Min, draw.Src)

	// Adjust luminance for each pixel
	for y := bounds.Min.Y; y < bounds.Max.Y; y++ {
		for x := bounds.Min.X; x < bounds.Max.X; x++ {
			// Get current pixel color
			col := color.New(0, 0, 0, 1)
			col.SetUint8(temp.At(x, y).(stdcolor.RGBA))

			// Convert to HSL
			h, s, l := rgbToHsl(col.R, col.G, col.B)

			// Adjust luminance - scale the adjustment based on current luminance
			if f.amount > 0 {
				l = l + (1-l)*f.amount
			} else {
				l = l + l*f.amount
			}

			// Convert back to RGB
			r, g, b := hslToRgb(h, s, l)

			// Set the pixel
			dst.Set(x, y, color.New(r, g, b, col.A).ToUint8())
		}
	}

	*img = dst
	return img, nil
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

func (f *LumFunction) ProcessPixel(x, y int, img *image.Image) (*color.Color64, error) {
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

	// Adjust luminance - scale the adjustment based on current luminance
	if f.amount > 0 {
		l = l + (1-l)*f.amount
	} else {
		l = l + l*f.amount
	}

	// Convert back to RGB
	r, g, b := hslToRgb(h, s, l)

	// Return new color
	return color.New(r, g, b, col.A), nil
}
