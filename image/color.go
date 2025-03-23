package image

import (
	"image"
	"image/color"
	"image/draw"

	"github.com/toxyl/gfx/color/convert"
	"github.com/toxyl/gfx/color/hsla"
	"github.com/toxyl/gfx/color/rgba"
	"github.com/toxyl/gfx/math"
	"github.com/toxyl/gfx/vars"
)

func (i *Image) GetRGBA(x, y int) *rgba.RGBA {
	col := i.raw.RGBAAt(x, y)
	if col.A == 0 {
		return vars.COLOR_TRANSPARENT_RGBA // If alpha is zero, return fully transparent black (0, 0, 0, 0).
	}
	scale := 255.0 / float32(col.A)
	return rgba.New(
		float32(col.R)*scale,
		float32(col.G)*scale,
		float32(col.B)*scale,
		float32(col.A),
	)
}

func (i *Image) GetHSLA(x, y int) *hsla.HSLA {
	col := i.raw.RGBAAt(x, y)
	if col.A == 0 {
		return vars.COLOR_TRANSPARENT
	}
	scale := 255.0 / float32(col.A)
	cc := rgba.New(
		float32(col.R)*scale,
		float32(col.G)*scale,
		float32(col.B)*scale,
		float32(col.A),
	)

	rf := float64(cc.R()) / 255.0
	gf := float64(cc.G()) / 255.0
	bf := float64(cc.B()) / 255.0
	af := float64(cc.A()) / 255.0

	max := math.Max(rf, math.Max(gf, bf))
	min := math.Min(rf, math.Min(gf, bf))
	delta := max - min

	h := 0.0
	s := 0.0
	l := (max + min) / 2

	if delta != 0 {
		s = delta / (1 - math.Abs(2*l-1))
	}

	switch {
	case delta == 0:
		h = 0
	case max == rf:
		h = math.Mod((gf-bf)/delta, 6)
	case max == gf:
		h = (bf-rf)/delta + 2
	case max == bf:
		h = (rf-gf)/delta + 4
	}
	h *= 60
	return hsla.New(h, s, l, af)
}

func (i *Image) SetRGBA(x, y int, c *rgba.RGBA) *Image {
	res := color.RGBA{0, 0, 0, 0}
	a := c.A()
	if a > 0 { // Avoid division by zero if original alpha is zero
		scale := float32(a) / 255.0
		res = color.RGBA{
			uint8(float32(c.R()) * scale),
			uint8(float32(c.G()) * scale),
			uint8(float32(c.B()) * scale),
			a,
		}
	}
	i.raw.Set(x, y, res)
	return i
}
func (i *Image) SetHSLA(x, y int, c *hsla.HSLA) *Image { return i.SetRGBA(x, y, convert.HSLAToRGBA(c)) }

func (i *Image) FillRGBA(x, y, w, h int, col *rgba.RGBA) *Image {
	draw.Draw(i.raw, image.Rect(x, y, w, h), &image.Uniform{col.RGBA()}, image.Point{}, draw.Src)
	return i
}

func (i *Image) FillHSLA(x, y, w, h int, col *hsla.HSLA) *Image {
	return i.FillRGBA(x, y, w, h, convert.HSLAToRGBA(col))
}
