package fx

import (
	"errors"
	"image"
	stdcolor "image/color"
	"image/draw"

	"github.com/toxyl/gfx/core/color"
	"github.com/toxyl/gfx/core/fx"
)

// ErrInvalidChannel is returned when an invalid color channel is specified
var ErrInvalidChannel = errors.New("invalid color channel")

// Channel represents the color channel to extract
type Channel int

const (
	ChannelRed Channel = iota
	ChannelGreen
	ChannelBlue
	ChannelAlpha
)

// ExtractFunction represents a function that extracts a specific color channel from an image
type ExtractFunction struct {
	*fx.BaseFunction
	channel Channel
	min     float64
	max     float64
}

// Function arguments
var extractArgs = []fx.FunctionArg{
	{
		Name:        "amount",
		Type:        fx.TypeFloat,
		Description: "extraction adjustment value",
		Required:    true,
		Min:         -1.0,
		Max:         1.0,
		Step:        0.1,
	},
}

func NewExtract(channel Channel, min, max float64) *ExtractFunction {
	return &ExtractFunction{
		BaseFunction: fx.NewBaseFunction("extract", "Extracts a specific color channel from an image", color.New(0, 0, 0, 1), extractArgs),
		channel:      channel,
		min:          min,
		max:          max,
	}
}

// Apply implements the Function interface
func (f *ExtractFunction) Apply(img image.Image) (image.Image, error) {
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
func (f *ExtractFunction) ProcessPixel(x, y int, img image.Image) (*color.Color64, error) {
	if img == nil {
		return nil, fx.ErrInvalidImage
	}

	bounds := img.Bounds()
	if x < bounds.Min.X || x >= bounds.Max.X || y < bounds.Min.Y || y >= bounds.Max.Y {
		return nil, fx.ErrOutOfBounds
	}

	col := color.New(0, 0, 0, 1)
	col.SetUint8(img.At(x, y).(stdcolor.RGBA))

	switch f.channel {
	case ChannelRed:
		return color.New(col.R, 0, 0, 1), nil
	case ChannelGreen:
		return color.New(0, col.G, 0, 1), nil
	case ChannelBlue:
		return color.New(0, 0, col.B, 1), nil
	case ChannelAlpha:
		return color.New(0, 0, 0, col.A), nil
	default:
		return nil, ErrInvalidChannel
	}
}

func init() {
	fx.DefaultRegistry.Register(NewExtract(ChannelRed, 0, 1))
}
