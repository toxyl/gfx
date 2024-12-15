package image

import (
	"github.com/toxyl/gfx/color/blend"
	"github.com/toxyl/gfx/color/hsla"
)

func (i *Image) Draw(src *Image, srcX, srcY, srcW, srcH, dstX, dstY, dstW, dstH int, mode blend.BlendMode, alpha float64) *Image {
	s := src.Crop(srcX, srcY, srcW, srcH, true)
	if srcW != dstW || srcH != dstH {
		s = s.Resize(dstW, dstH)
	}
	return i.mergeHSLA(s, 0, 0, dstW, dstH, func(x, y int, srcCol *hsla.HSLA) (x2 int, y2 int, col2 *hsla.HSLA) {
		return x + dstX, y + dstY, blend.HSLA(i.GetHSLA(x+dstX, y+dstY), srcCol, mode, alpha)
	})
}

func (i *Image) DrawLineV(x, yStart, yEnd, thickness int, col *hsla.HSLA, mode blend.BlendMode) {
	xStart := x - thickness/2
	xEnd := x + thickness/2
	if thickness%2 != 0 {
		xStart = x - thickness/2
		xEnd = x + thickness/2 + 1
	}
	i.ProcessHSLA(xStart, yStart, xEnd, yEnd, func(x, y int, colSrc *hsla.HSLA) (x2 int, y2 int, col2 *hsla.HSLA) {
		if col.A() < 1 {
			return x, y, blend.HSLA(colSrc, col, mode, 1.0)
		}
		return x, y, col
	})
}

func (i *Image) DrawLineH(y, xStart, xEnd, thickness int, col *hsla.HSLA, mode blend.BlendMode) {
	yStart := y - thickness/2
	yEnd := y + thickness/2
	if thickness%2 != 0 {
		yStart = y - thickness/2
		yEnd = y + thickness/2 + 1
	}
	i.ProcessHSLA(xStart, yStart, xEnd, yEnd, func(x, y int, colSrc *hsla.HSLA) (x2 int, y2 int, col2 *hsla.HSLA) {
		if col.A() < 1 {
			return x, y, blend.HSLA(colSrc, col, mode, 1.0)
		}
		return x, y, col
	})
}

func (i *Image) DrawRect(x, y, w, h, thickness int, colBorder, colFill *hsla.HSLA, mode blend.BlendMode) {
	i.ProcessHSLA(x, y, x+w, y+h, func(x1, y1 int, colSrc *hsla.HSLA) (x2 int, y2 int, col2 *hsla.HSLA) {
		col := colFill
		if x1-x < thickness || x1-x >= w-thickness || y1-y < thickness || y1-y >= h-thickness {
			col = colBorder
		}
		if col.A() > 0 {
			return x1, y1, blend.HSLA(colSrc, col, mode, 1.0)
		}
		return x1, y1, colSrc
	})
}
