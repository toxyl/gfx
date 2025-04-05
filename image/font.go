package image

import (
	_ "embed"
	"regexp"
	"strconv"
	"strings"

	"github.com/toxyl/gfx/config"
	"github.com/toxyl/gfx/core/blendmodes"
	"github.com/toxyl/gfx/core/color"
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
		charWidth:  config.CHAR_W,
		charHeight: config.CHAR_H,
		columns:    config.SPRITESHEET_COLS,
	}
)

func (i *Image) drawFontWithOutline(src *Image, dstX, dstY int, col color.HSL, glow bool, mode *blendmodes.IBlendMode) {
	outlineCol := col
	if glow {
		outlineCol.L = col.L * 0.40
		outlineCol.S = col.S * 0.75
		col.S = col.S * 1.25
	} else {
		outlineCol.L = 0.1
	}
	fnIsSet := func(x, y int) bool {
		if x < 0 || y < 0 {
			return false // out of bounds
		}
		if x >= src.W() || y >= src.H() {
			return false // out of bounds
		}
		return src.GetRGBA64(x, y).A != 0
	}
	outlinedSrc := src.Clone()
	outlinedSrc.ProcessHSLA(0, 0, src.W(), src.H(), func(x, y int, colSrc *color.HSL) (x2 int, y2 int, col2 *color.HSL) {
		if colSrc.Alpha == 0 {
			colSrc.Alpha = 0
			// we have to check if there is a non-transparent pixel around us,
			// if so the current pixel needs to be colored
			if fnIsSet(x-1, y-1) || fnIsSet(x, y-1) || fnIsSet(x+1, y-1) ||
				fnIsSet(x-1, y) || fnIsSet(x+1, y) ||
				fnIsSet(x-1, y+1) || fnIsSet(x, y+1) || fnIsSet(x+1, y+1) {
				colSrc = &outlineCol
			}
		} else {
			colSrc.H = col.H
			colSrc.S = col.S
			colSrc.L = col.L
		}
		return x, y, colSrc
	})

	i.mergeHSLA(outlinedSrc, 0, 0, outlinedSrc.W(), outlinedSrc.H(), func(x, y int, colSrc *color.HSL) (x2 int, y2 int, col2 *color.HSL) {
		if colSrc.Alpha > 0 {
			bottom := i.GetRGBA64(dstX+x, dstY+y)
			top := colSrc.ToRGB()
			topRGBA64, err := color.NewRGBA64(top.R, top.G, top.B, top.A)
			if err != nil {
				return dstX + x, dstY + y, colSrc
			}
			blended, err := mode.Blend(bottom, topRGBA64, 1.0)
			if err != nil {
				return dstX + x, dstY + y, colSrc
			}
			return dstX + x, dstY + y, color.HSLFromRGB(blended)
		}
		return dstX + x, dstY + y, color.HSLFromRGB(i.GetRGBA64(dstX+x, dstY+y))
	})
}

func (sf *spriteFont) render(text string, col color.HSL, glow bool) *Image {
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
						col.S = 0
						col.L = 1
					case "black":
						col.S = 0
						col.L = 0
					case "gray":
						col.S = 0
					case "color":
						col.S = 0.5
					case "r":
						col.H = 0
					case "g":
						col.H = 120
					case "b":
						col.H = 240
					case "c":
						col.H = 180
					case "m":
						col.H = 300
					case "y":
						col.H = 60
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
						col.H = float64(h)
					}
					if s, err := strconv.ParseFloat(m[2], 64); err == nil {
						col.S = float64(s)
					}
					if l, err := strconv.ParseFloat(m[3], 64); err == nil {
						col.L = float64(l)
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
		blendNormal, err := blendmodes.Get("normal")
		if err != nil {
			panic("There must be a 'normal' blendmode!")
		}
		sx := (charIndex % sf.columns) * sf.charWidth
		sy := (charIndex / sf.columns) * sf.charHeight
		charImg := New(sf.charWidth, sf.charHeight)
		charImg.Draw(sf.image, sx, sy, sf.charWidth, sf.charHeight, 0, 0, sf.charWidth, sf.charHeight, blendNormal, 1)
		res.drawFontWithOutline(charImg, xOffset, yOffset, col, glow, blendNormal)

		xOffset += sf.charWidth
	}

	return res
}

func (i *Image) DrawText(text string, x, y int, col color.HSL, glow bool, mode *blendmodes.IBlendMode) *Image {
	src := font.render(text, col, glow)
	return i.Draw(src, 0, 0, src.W(), src.H(), x, y, src.W(), src.H(), mode, col.Alpha)
}
