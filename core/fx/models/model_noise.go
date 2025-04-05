package models

import (
	"image"
	stdcolor "image/color"
	"image/draw"
	"math/rand"
	"time"

	"github.com/toxyl/gfx/core/color"
	"github.com/toxyl/gfx/core/fx"
)

// NoiseFunction represents a function that adds noise to an image
type NoiseFunction struct {
	*fx.BaseFunction
	amount float64
	rng    *rand.Rand
}

// NewNoise creates a new noise function
func NewNoise(amount float64) *NoiseFunction {
	return &NoiseFunction{
		BaseFunction: fx.NewBaseFunction("noise", "Adds random noise to image", color.New(0, 0, 0, 1), nil),
		amount:       amount,
		rng:          rand.New(rand.NewSource(time.Now().UnixNano())),
	}
}

// Apply implements the Function interface
func (f *NoiseFunction) Apply(img *image.Image) (*image.Image, error) {
	bounds := (*img).Bounds()
	dst := image.NewRGBA(bounds)
	draw.Draw(dst, bounds, *img, bounds.Min, draw.Src)

	for y := bounds.Min.Y; y < bounds.Max.Y; y++ {
		for x := bounds.Min.X; x < bounds.Max.X; x++ {
			col := f.ProcessPixel(x, y, img)
			dst.Set(x, y, col.ToUint8())
		}
	}

	*img = dst
	return img, nil
}

// ProcessPixel implements the ImageFunction interface
func (f *NoiseFunction) ProcessPixel(x, y int, img *image.Image) *color.Color64 {
	col := color.New(0, 0, 0, 1)
	col.SetUint8((*img).At(x, y).(stdcolor.RGBA))

	// Generate random noise
	noise := (f.rng.Float64()*2 - 1) * f.amount

	// Apply noise to each channel
	col.R = noiseClamp(col.R + noise)
	col.G = noiseClamp(col.G + noise)
	col.B = noiseClamp(col.B + noise)

	return col
}

// noiseClamp ensures a value is between 0 and 1
func noiseClamp(v float64) float64 {
	if v < 0 {
		return 0
	}
	if v > 1 {
		return 1
	}
	return v
}
