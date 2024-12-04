package filter

import (
	"github.com/toxyl/gfx/color/hsla"
	"github.com/toxyl/gfx/math"
)

type ColorFilter struct {
	Col, Min, MinThres, Max, MaxThres *hsla.HSLA
}

func NewColorFilter(h, s, l FilterValue) *ColorFilter {
	// thres_min - min - max - thres_max
	var (
		hThresMin = h.Min - h.Tolerance
		hMin      = h.Min
		hMax      = h.Max
		hThresMax = h.Max + h.Tolerance
		cThresMin = hsla.New(hThresMin, math.Clamp(s.Min-s.Tolerance, 0.0, 1.0), math.Clamp(l.Min-l.Tolerance, 0.0, 1.0), 1.0)
		cMin      = hsla.New(hMin, math.Clamp(s.Min, 0.0, 1.0), math.Clamp(l.Min, 0.0, 1.0), 1.0)
		cMax      = hsla.New(hMax, math.Clamp(s.Max, 0.0, 1.0), math.Clamp(l.Max, 0.0, 1.0), 1.0)
		cThresMax = hsla.New(hThresMax, math.Clamp(s.Max+s.Tolerance, 0.0, 1.0), math.Clamp(l.Max+l.Tolerance, 0.0, 1.0), 1.0)
	)

	for hMax <= hMin {
		hThresMax += 360
		hMax += 360
		cThresMax.SetH(hThresMax)
		cMax.SetH(hMax)
	}

	cf := ColorFilter{
		Col:      hsla.New(hMin+(hMax-hMin)/2, 1.0, 0.5, 1.0),
		MinThres: cThresMin,
		Min:      cMin,
		Max:      cMax,
		MaxThres: cThresMax,
	}
	return &cf
}

func ToColorFilter(hue, hueTolerance, hueFeather,
	sat, satTolerance, satFeather,
	lum, lumTolerance, lumFeather float64) *ColorFilter {
	if hueTolerance >= 180 {
		hue = 0
		hueTolerance = 180
		hueFeather = 0
	}
	if hueTolerance+hueFeather > 180 {
		hueFeather = 180 - hueTolerance
	}

	return NewColorFilter(
		ToFilterValue(hue, hueTolerance, hueFeather),
		ToFilterValue(sat, satTolerance, satFeather),
		ToFilterValue(lum, lumTolerance, lumFeather),
	)
}
