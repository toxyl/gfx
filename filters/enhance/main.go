package enhance

import (
	"github.com/toxyl/gfx/filters/convolution"
	"github.com/toxyl/gfx/image"
	"github.com/toxyl/gfx/math"
)

func Apply(img *image.Image, amount float64) *image.Image {
	return convolution.NewCustomFilter(
		amount, 1, 0,
		func(a float64) (matrix [][]float64) {
			return [][]float64{
				{-0.5 / a / 4.0, 1 / a / 6.0, -0.5 / a / 4.0},
				{1 / a / 8.0, math.Clamp(a*a, 0.0, 1.5), 1 / a / 8.0},
				{-0.5 / a / 4.0, 1 / a / 6.0, -0.5 / a / 4.0},
			}
		},
	).Apply(img)
}
