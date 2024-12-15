package alphamap

import (
	"strings"

	"github.com/toxyl/gfx/color/hsla"
	"github.com/toxyl/gfx/filters/meta"
	"github.com/toxyl/gfx/image"
)

var Meta = meta.New("alpha-map", []*meta.FilterMetaDataArg{
	{Name: "source", Default: "l"},
	{Name: "lower", Default: 0.0},
	{Name: "upper", Default: 0.0},
})

func Apply(i *image.Image, source string, lowerThreshold, upperThreshold float64) *image.Image {
	alphaSrc := strings.ToLower(source)
	minVal := lowerThreshold
	maxVal := upperThreshold
	invert := minVal > maxVal
	if invert {
		minVal, maxVal = maxVal, minVal
	}
	return i.ProcessHSLA(0, 0, i.W(), i.H(), func(x, y int, col *hsla.HSLA) (x2 int, y2 int, col2 *hsla.HSLA) {
		val := 0.0
		switch alphaSrc {
		case "s":
			val = col.S()
		case "l":
			val = col.L()
		case "s*l":
			val = col.S() * col.L()
		default:
			panic("invalid alpha source, available options are: s, l, s*l")
		}

		if (invert && val <= minVal) || (!invert && val >= maxVal) {
			return x, y, col.SetA(1)
		}

		if (invert && val >= maxVal) || (!invert && val <= minVal) {
			return x, y, col.SetA(0)
		}

		rngLum := maxVal - minVal
		a := (val - minVal) / rngLum
		if invert {
			a = 1 - a
		}
		return x, y, col.SetA(a + minVal)
	})
}
