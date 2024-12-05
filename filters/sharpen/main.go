package sharpen

import (
	"github.com/toxyl/gfx/filters/convolution"
	"github.com/toxyl/gfx/image"
)

func Apply(img *image.Image, amount float64) *image.Image {
	return convolution.NewConvolutionMatrix([][]float64{
		{-amount, -amount, -amount},
		{-amount, 1 + 8*amount, -amount},
		{-amount, -amount, -amount},
	}, 1, 0).Apply(img)
}
