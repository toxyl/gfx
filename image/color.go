package image

import (
	"image"
	"image/draw"

	"github.com/toxyl/gfx/core/color"
)

func (i *Image) GetRGBA64(x, y int) *color.RGBA64 { return color.Get(i.raw, x, y) }
func (i *Image) SetRGBA(x, y int, c *color.RGBA64) *Image {
	img := i.Get()
	color.Set(img, x, y, c)
	i.Set(img)
	return i
}
func (i *Image) SetHSLA(x, y int, c *color.HSL) *Image {
	rgb := c.ToRGB()
	rgba64, _ := color.NewRGBA64(rgb.R, rgb.G, rgb.B, rgb.A)
	return i.SetRGBA(x, y, rgba64)
}

func (i *Image) FillRGBA(x, y, w, h int, col *color.RGBA64) *Image {
	draw.Draw(i.raw, image.Rect(x, y, w, h), &image.Uniform{col.To8bit()}, image.Point{}, draw.Src)
	return i
}

func (i *Image) FillHSLA(x, y, w, h int, col *color.HSL) *Image {
	rgb := col.ToRGB()
	rgba64, _ := color.NewRGBA64(rgb.R, rgb.G, rgb.B, rgb.A)
	return i.FillRGBA(x, y, w, h, rgba64)
}
