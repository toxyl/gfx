package parser

import (
	"strconv"
	"strings"

	"github.com/toxyl/gfx/color/hsla"
)

func parseColor(value string) *hsla.HSLA {
	value = strings.TrimSuffix(strings.TrimPrefix(strings.TrimSpace(strings.ToLower(value)), "hsla"+STR_LPAREN), STR_RPAREN)
	parts := strings.Split(value, STR_SPACE)
	if len(parts) != 4 {
		panic("composition color must be given as `hsla" + STR_LPAREN + "hue" + STR_SPACE + "sat" + STR_SPACE + "lum" + STR_SPACE + "alpha" + STR_RPAREN + "`")
	}
	h, _ := strconv.ParseFloat(parts[0], 64)
	s, _ := strconv.ParseFloat(parts[1], 64)
	l, _ := strconv.ParseFloat(parts[2], 64)
	a, _ := strconv.ParseFloat(parts[3], 64)
	return hsla.New(h, s, l, a)
}
