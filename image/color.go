package image

import (
	"image"
	"image/draw"

	"github.com/toxyl/gfx/color/convert"
	"github.com/toxyl/gfx/color/hsla"
	"github.com/toxyl/gfx/color/rgba"
)

func (i *Image) GetRGBA(x, y int) *rgba.RGBA { return convert.RGBAPremulToRGBA(i.raw.RGBAAt(x, y)) }
func (i *Image) GetHSLA(x, y int) *hsla.HSLA { return convert.RGBAPremulToHSLA(i.raw.RGBAAt(x, y)) }

func (i *Image) SetRGBA(x, y int, c *rgba.RGBA) *Image {
	i.raw.Set(x, y, convert.RGBAToRGBAPremul(c))
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
