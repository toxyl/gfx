package lum

import (
	"github.com/toxyl/gfx/color/rgba"
	"github.com/toxyl/gfx/image"
	"github.com/toxyl/gfx/math"
)

func Apply(img *image.Image, luminance float64) *image.Image {
	return img.ProcessRGBA(0, 0, img.W(), img.H(), func(x, y int, col *rgba.RGBA) (x2 int, y2 int, col2 *rgba.RGBA) {
		r := float64(col.R())
		g := float64(col.G())
		b := float64(col.B())

		// Calculate luminance (Y) using the formula Y = 0.299*R + 0.587*G + 0.114*B
		Y := 0.299*r + 0.587*g + 0.114*b

		// Adjust luminance
		Y = math.Clamp(Y+luminance*255.0, 0.0, 255.0)

		// Calculate the new RGB values based on the adjusted luminance
		factor := Y / (0.299*r + 0.587*g + 0.114*b)
		r = math.Clamp(r*factor, 0.0, 255.0)
		g = math.Clamp(g*factor, 0.0, 255.0)
		b = math.Clamp(b*factor, 0.0, 255.0)

		return x, y, rgba.New(uint8(r), uint8(g), uint8(b), col.A())
	})
}
