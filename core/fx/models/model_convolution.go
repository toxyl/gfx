package models

import (
	"image"
	"image/color"

	"github.com/toxyl/gfx/core/fx"
	"github.com/toxyl/gfx/core/meta"
	"github.com/toxyl/math"
)

// Convolution represents a convolution matrix effect.
type Convolution struct {
	Matrix  [][]float64 // Convolution matrix
	Divisor float64     // Divisor for normalization
	Offset  float64     // Offset to add after division
	meta    *fx.EffectMeta
}

// NewConvolutionEffect creates a new convolution matrix effect.
func NewConvolutionEffect(matrix [][]float64, divisor, offset float64) *Convolution {
	c := &Convolution{
		Matrix:  matrix,
		Divisor: divisor,
		Offset:  offset,
		meta: fx.NewEffectMeta(
			"Convolution",
			"Applies a convolution matrix to an image",
			meta.NewChannelMeta("Divisor", 0.0, 100.0, "", "Divisor for normalization"),
			meta.NewChannelMeta("Offset", -1.0, 1.0, "", "Offset to add after division"),
		),
	}
	c.Divisor = fx.ClampParameter(divisor, c.meta.Parameters[0])
	c.Offset = fx.ClampParameter(offset, c.meta.Parameters[1])
	return c
}

// Apply applies the convolution matrix effect to an image.
func (c *Convolution) Apply(img image.Image) image.Image {
	bounds := img.Bounds()
	dst := image.NewRGBA(bounds)

	// Get matrix dimensions
	rows := len(c.Matrix)
	cols := len(c.Matrix[0])
	halfRows := rows / 2
	halfCols := cols / 2

	for y := bounds.Min.Y + halfRows; y < bounds.Max.Y-halfRows; y++ {
		for x := bounds.Min.X + halfCols; x < bounds.Max.X-halfCols; x++ {
			var r, g, b float64

			// Apply convolution matrix
			for ky := -halfRows; ky <= halfRows; ky++ {
				for kx := -halfCols; kx <= halfCols; kx++ {
					px := x + kx
					py := y + ky
					pr, pg, pb, _ := img.At(px, py).RGBA()

					// Convert to float64 and normalize
					prf := float64(pr) / 0xFFFF
					pgf := float64(pg) / 0xFFFF
					pbf := float64(pb) / 0xFFFF

					// Apply matrix weight
					weight := c.Matrix[ky+halfRows][kx+halfCols]
					r += prf * weight
					g += pgf * weight
					b += pbf * weight
				}
			}

			// Apply divisor and offset
			r = r/c.Divisor + c.Offset
			g = g/c.Divisor + c.Offset
			b = b/c.Divisor + c.Offset

			// Clamp values
			r = math.Max(0, math.Min(1, r))
			g = math.Max(0, math.Min(1, g))
			b = math.Max(0, math.Min(1, b))

			// Get original alpha
			_, _, _, oa := img.At(x, y).RGBA()

			dst.Set(x, y, color.RGBA64{
				R: uint16(r * 0xFFFF),
				G: uint16(g * 0xFFFF),
				B: uint16(b * 0xFFFF),
				A: uint16(oa),
			})
		}
	}

	return dst
}

// Meta returns the effect metadata.
func (c *Convolution) Meta() *fx.EffectMeta {
	return c.meta
}
