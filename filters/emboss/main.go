package emboss

import (
	"github.com/toxyl/gfx/filters/convolution"
	"github.com/toxyl/gfx/image"
)

func Apply(img *image.Image, amount float64) *image.Image {
	return convolution.NewConvolutionMatrix([][]float64{
		{-1 * amount, -1 * amount, 0},
		{-1 * amount, 1 * amount, 1 * amount},
		{0, 1 * amount, 1 * amount},
	}, 1, 0).Apply(img)
}
