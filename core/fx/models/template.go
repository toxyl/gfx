package models

import (
	"image"

	"github.com/toxyl/gfx/core/color"
	"github.com/toxyl/gfx/core/fx"
)

// TemplateFunction represents a template for implementing image processing functions
type TemplateFunction struct {
	*fx.BaseFunction
	// Add any additional fields needed for the function
}

// NewTemplate creates a new template function
func NewTemplate() *TemplateFunction {
	args := []fx.FunctionArg{
		{
			Name:        "example_arg",
			Type:        fx.TypeFloat,
			Description: "Example argument description",
			Required:    true,
			Min:         0.0,
			Max:         1.0,
			Step:        0.1,
		},
	}

	return &TemplateFunction{
		BaseFunction: fx.NewBaseFunction(
			"template",
			"Template function description",
			color.New(1, 1, 1, 1),
			args,
		),
	}
}

// Apply implements the Function interface
func (f *TemplateFunction) Apply(img *image.Image) (*image.Image, error) {
	if err := f.ValidateArgs(); err != nil {
		return nil, err
	}

	if img == nil || *img == nil {
		return nil, fx.ErrInvalidImage
	}

	bounds := (*img).Bounds()
	dst := image.NewRGBA(bounds)

	for y := bounds.Min.Y; y < bounds.Max.Y; y++ {
		for x := bounds.Min.X; x < bounds.Max.X; x++ {
			c, err := f.ProcessPixel(x, y, img)
			if err != nil {
				return nil, err
			}
			dst.Set(x, y, c.ToUint8())
		}
	}

	result := image.Image(dst)
	return &result, nil
}

// ProcessPixel implements the ImageFunction interface
func (f *TemplateFunction) ProcessPixel(x, y int, img *image.Image) (*color.Color64, error) {
	if img == nil || *img == nil {
		return nil, fx.ErrInvalidImage
	}

	bounds := (*img).Bounds()
	if x < bounds.Min.X || x >= bounds.Max.X || y < bounds.Min.Y || y >= bounds.Max.Y {
		return nil, fx.ErrOutOfBounds
	}

	// Get the original color
	original := (*img).At(x, y)
	r, g, b, a := original.RGBA()

	// Process the color (example: convert to grayscale)
	gray := uint32((0.299*float64(r) + 0.587*float64(g) + 0.114*float64(b)) / 65535.0)

	// Create and return the new color
	return color.New(
		float64(gray),
		float64(gray),
		float64(gray),
		float64(a)/65535.0,
	), nil
}
