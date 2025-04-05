package fx

import (
	"image"
	stdcolor "image/color"
	"image/draw"

	"github.com/toxyl/gfx/core/color"
	"github.com/toxyl/gfx/core/fx"
)

// ColorShiftFunction represents a function that shifts colors in an image
type ColorShiftFunction struct {
	*fx.BaseFunction
	redShift   float64
	greenShift float64
	blueShift  float64
}

// Function arguments
var colorshiftArgs = []fx.FunctionArg{
	{
		Name:        "amount",
		Type:        fx.TypeFloat,
		Description: "colorshift adjustment value",
		Required:    true,
		Min:         -1.0,
		Max:         1.0,
		Step:        0.1,
	},
}

// NewColorShift creates a new color shift function
func NewColorShift(redShift, greenShift, blueShift float64) *ColorShiftFunction {
	return &ColorShiftFunction{
		BaseFunction: fx.NewBaseFunction("colorshift", "Applies colorshift transformation to an image", color.New(0, 0, 0, 1), colorshiftArgs),
		redShift:     redShift,
		greenShift:   greenShift,
		blueShift:    blueShift,
	}
}

// Apply implements the Function interface
func (f *ColorShiftFunction) Apply(img image.Image) (image.Image, error) {
	if img == nil {
		return nil, fx.ErrInvalidImage
	}
	bounds := img.Bounds()
	dst := image.NewRGBA(bounds)
	draw.Draw(dst, bounds, img, bounds.Min, draw.Src)

	for y := bounds.Min.Y; y < bounds.Max.Y; y++ {
		for x := bounds.Min.X; x < bounds.Max.X; x++ {
			col, err := f.ProcessPixel(x, y, img)
			if err != nil {
				return nil, err
			}
			dst.Set(x, y, col.ToUint8())
		}
	}

	return dst, nil
}

// ProcessPixel implements the ImageFunction interface
func (f *ColorShiftFunction) ProcessPixel(x, y int, img image.Image) (*color.Color64, error) {
	if img == nil {
		return nil, fx.ErrInvalidImage
	}

	bounds := img.Bounds()
	if x < bounds.Min.X || x >= bounds.Max.X || y < bounds.Min.Y || y >= bounds.Max.Y {
		return nil, fx.ErrOutOfBounds
	}

	col := color.New(0, 0, 0, 1)
	col.SetUint8(img.At(x, y).(stdcolor.RGBA))

	// Adjust each color channel independently
	col.R = clamp(col.R + f.redShift)
	col.G = clamp(col.G + f.greenShift)
	col.B = clamp(col.B + f.blueShift)

	return col, nil
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
	fx.DefaultRegistry.Register(NewColorShift(0, 0, 0))
}
