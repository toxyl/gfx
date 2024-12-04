package image

import (
	_ "embed"
	"regexp"
	"strconv"
	"strings"

	"github.com/toxyl/gfx/color/blend"
	"github.com/toxyl/gfx/color/hsla"
	"github.com/toxyl/gfx/vars"
)

//go:embed font.png
var fontBytes []byte

type spriteFont struct {
	image      *Image
	charWidth  int
	charHeight int
	columns    int
}

var (
	font = &spriteFont{
		image:      NewFromBytes("png", fontBytes),
		charWidth:  vars.CHAR_W,
		charHeight: vars.CHAR_H,
		columns:    vars.SPRITESHEET_COLS,
	}
)

func (i *Image) drawFontWithOutline(src *Image, dstX, dstY int, col hsla.HSLA, glow bool, mode blend.BlendMode) {
	outlineCol := col
	if glow {
		outlineCol.SetL(col.L() * 0.40)
		outlineCol.SetS(col.S() * 0.75)
		col.SetS(col.S() * 1.25)
	} else {
		outlineCol.SetL(0.1)
	}
	fnIsSet := func(x, y int) bool {
		if x < 0 || y < 0 {
			return false // out of bounds
		}
		if x >= src.W() || y >= src.H() {
			return false // out of bounds
		}
		return src.GetRGBA(x, y).A() != 0
	}
	outlinedSrc := src.Clone()
	outlinedSrc.ProcessHSLA(0, 0, src.W(), src.H(), func(x, y int, colSrc *hsla.HSLA) (x2 int, y2 int, col2 *hsla.HSLA) {
		if colSrc.A() == 0 {
			colSrc.SetA(0)
			// we have to check if there is a non-transparent pixel around us,
			// if so the current pixel needs to be colored
			if fnIsSet(x-1, y-1) || fnIsSet(x, y-1) || fnIsSet(x+1, y-1) ||
				fnIsSet(x-1, y) || fnIsSet(x+1, y) ||
				fnIsSet(x-1, y+1) || fnIsSet(x, y+1) || fnIsSet(x+1, y+1) {
				colSrc = &outlineCol
			}
		} else {
			colSrc.SetH(col.H())
			colSrc.SetS(col.S())
			colSrc.SetL(col.L())
		}
		return x, y, colSrc
	})

	i.mergeHSLA(outlinedSrc, 0, 0, outlinedSrc.W(), outlinedSrc.H(), func(x, y int, colSrc *hsla.HSLA) (x2 int, y2 int, col2 *hsla.HSLA) {
		if colSrc.A() > 0 {
			return dstX + x, dstY + y, blend.HSLA(&col, colSrc, mode, 1.0)
		}
		return dstX + x, dstY + y, i.GetHSLA(dstX+x, dstY+y)
	})

}

func (sf *spriteFont) render(text string, col hsla.HSLA, glow bool) *Image {
	reColor := regexp.MustCompile(`^\[:(gon|goff|white|black|gray|color|r|g|b|c|m|y|):\]`)
	reColorHSL := regexp.MustCompile(`^\[:(\d{1,3}):(\d\.\d+):(\d\.\d+):\]`)
	cleanText := reColorHSL.ReplaceAllString(text, "")
	cleanText = reColor.ReplaceAllString(cleanText, "")
	initialColor := col

	w := int(float64(len(cleanText)) * float64(sf.charWidth))
	h := int(float64(sf.charHeight)) + strings.Count(text, "\n")*sf.charHeight

	res := New(w, h)
	yOffset := 0
	xOffset := 0
	nextI := 0
	for i, char := range text {
		if i < nextI {
			continue
		}
		if char == '\n' {
			yOffset += sf.charHeight
			xOffset = 0
			continue
		}
		if char == '[' {
			if reColor.MatchString(text[i:]) {
				matches := reColor.FindAllStringSubmatch(text[i:], 1)
				if len(matches) > 0 {
					m := matches[0]
					switch m[1] {
					case "":
						col = initialColor
					case "white":
						col.SetS(0)
						col.SetL(1)
					case "black":
						col.SetS(0)
						col.SetL(0)
					case "gray":
						col.SetS(0)
					case "color":
						col.SetS(0.5)
					case "r":
						col.SetH(0)
					case "g":
						col.SetH(120)
					case "b":
						col.SetH(240)
					case "c":
						col.SetH(180)
					case "m":
						col.SetH(300)
					case "y":
						col.SetH(60)
					case "gon":
						glow = true
					case "goff":
						glow = false
					}
					nextI = i + len(m[0])
					continue
				}
			}
			if reColorHSL.MatchString(text[i:]) {
				matches := reColorHSL.FindAllStringSubmatch(text[i:], 1)
				if len(matches) > 0 {
					m := matches[0]
					if h, err := strconv.Atoi(m[1]); err == nil {
						col.SetH(float64(h))
					}
					if s, err := strconv.ParseFloat(m[2], 64); err == nil {
						col.SetS(float64(s))
					}
					if l, err := strconv.ParseFloat(m[3], 64); err == nil {
						col.SetL(float64(l))
					}
					nextI = i + len(m[0])
					continue
				}
			}
		}
		charIndex := int(char) - 33
		if charIndex <= 0 || charIndex >= sf.columns*6 {
			charIndex = 95 // Default to space character for unsupported characters
		}

		sx := (charIndex % sf.columns) * sf.charWidth
		sy := (charIndex / sf.columns) * sf.charHeight
		charImg := New(sf.charWidth, sf.charHeight)
		charImg.Draw(sf.image, sx, sy, sf.charWidth, sf.charHeight, 0, 0, sf.charWidth, sf.charHeight, blend.NORMAL, 1)
		res.drawFontWithOutline(charImg, xOffset, yOffset, col, glow, blend.NORMAL)

		xOffset += sf.charWidth
	}

	return res
}

func (i *Image) DrawText(text string, x, y int, col hsla.HSLA, glow bool, mode blend.BlendMode) *Image {
	src := font.render(text, col, glow)
	return i.Draw(src, 0, 0, src.W(), src.H(), x, y, src.W(), src.H(), mode, col.A())
}
