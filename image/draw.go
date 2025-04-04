package image

import (
	"github.com/toxyl/gfx/core/blendmodes"
	"github.com/toxyl/gfx/core/color"
)

func (i *Image) Draw(src *Image, srcX, srcY, srcW, srcH, dstX, dstY, dstW, dstH int, mode *blendmodes.IBlendMode, alpha float64) *Image {
	s := src.Crop(srcX, srcY, srcW, srcH, true)
	if srcW != dstW || srcH != dstH {
		s = s.Resize(dstW, dstH)
	}
	return i.mergeHSLA(s, 0, 0, dstW, dstH, func(x, y int, srcCol *color.HSL) (x2 int, y2 int, col2 *color.HSL) {
		rgb := srcCol.ToRGB()
		srcRGBA, _ := color.NewRGBA64(rgb.R, rgb.G, rgb.B, rgb.A)
		dstRGB := i.GetRGBA64(x+dstX, y+dstY)
		blended, _ := mode.Blend(dstRGB, srcRGBA, alpha)
		return x + dstX, y + dstY, color.HSLFromRGB(blended)
	})
}

func (i *Image) DrawLineV(x, yStart, yEnd, thickness int, col *color.HSL, mode *blendmodes.IBlendMode) *Image {
	xStart := x - thickness/2
	xEnd := x + thickness/2
	if thickness%2 != 0 {
		xStart = x - thickness/2
		xEnd = x + thickness/2 + 1
	}
	i.ProcessHSLA(xStart, yStart, xEnd, yEnd, func(x, y int, colSrc *color.HSL) (x2 int, y2 int, col2 *color.HSL) {
		if col.Alpha < 1 {
			rgb := colSrc.ToRGB()
			srcRGBA, _ := color.NewRGBA64(rgb.R, rgb.G, rgb.B, rgb.A)
			rgb = col.ToRGB()
			colRGBA, _ := color.NewRGBA64(rgb.R, rgb.G, rgb.B, rgb.A)
			blended, _ := mode.Blend(srcRGBA, colRGBA, 1.0)
			return x, y, color.HSLFromRGB(blended)
		}
		return x, y, col
	})
	return i
}

func (i *Image) DrawLineH(y, xStart, xEnd, thickness int, col *color.HSL, mode *blendmodes.IBlendMode) *Image {
	yStart := y - thickness/2
	yEnd := y + thickness/2
	if thickness%2 != 0 {
		yStart = y - thickness/2
		yEnd = y + thickness/2 + 1
	}
	i.ProcessHSLA(xStart, yStart, xEnd, yEnd, func(x, y int, colSrc *color.HSL) (x2 int, y2 int, col2 *color.HSL) {
		if col.Alpha < 1 {
			rgb := colSrc.ToRGB()
			srcRGBA, _ := color.NewRGBA64(rgb.R, rgb.G, rgb.B, rgb.A)
			rgb = col.ToRGB()
			colRGBA, _ := color.NewRGBA64(rgb.R, rgb.G, rgb.B, rgb.A)
			blended, _ := mode.Blend(srcRGBA, colRGBA, 1.0)
			return x, y, color.HSLFromRGB(blended)
		}
		return x, y, col
	})
	return i
}

func (i *Image) DrawRect(x, y, w, h, thickness int, colBorder, colFill *color.HSL, mode *blendmodes.IBlendMode) *Image {
	i.ProcessHSLA(x, y, x+w, y+h, func(x1, y1 int, colSrc *color.HSL) (x2 int, y2 int, col2 *color.HSL) {
		col := colFill
		if x1-x < thickness || x1-x >= w-thickness || y1-y < thickness || y1-y >= h-thickness {
			col = colBorder
		}
		if col.Alpha > 0 {
			rgb := colSrc.ToRGB()
			srcRGBA, _ := color.NewRGBA64(rgb.R, rgb.G, rgb.B, rgb.A)
			rgb = col.ToRGB()
			colRGBA, _ := color.NewRGBA64(rgb.R, rgb.G, rgb.B, rgb.A)
			blended, _ := mode.Blend(srcRGBA, colRGBA, 1.0)
			return x1, y1, color.HSLFromRGB(blended)
		}
		return x1, y1, colSrc
	})
	return i
}
