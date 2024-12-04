package extract

import (
	"fmt"

	"github.com/toxyl/gfx/color/filter"
	"github.com/toxyl/gfx/color/hsla"
	"github.com/toxyl/gfx/image"
	"github.com/toxyl/gfx/math"
)

// prepHueRange normalizes and adjusts the boundaries of a FuzzyRange to operate
// within a hue range of 0 to 360 degrees. It handles edge cases where the range
// crosses the zero boundary (e.g., from 350 to 10 degrees) or spans the entire circle.
//
// This function ensures that the hue values are normalized to the [0,1] range
// and modifies the FuzzyRange's properties to correctly account for wrapping
// over the zero boundary, as well as detecting full-circle ranges.
//
// Parameters:
// - fr: The FuzzyRange representing the hue range, including thresholds and falloff function.
//
// Returns:
// - A FuzzyRange with normalized and adjusted properties for hue calculations.
func prepHueRange(fr FuzzyRange) FuzzyRange {
	// Normalize hue values to the range [0, 1]
	normalize := func(x float64) float64 {
		return x / 360
	}

	// Normalize input values
	fr.Min = normalize(fr.Min)
	fr.Max = normalize(fr.Max)
	fr.MinThres = normalize(fr.MinThres)
	fr.MaxThres = normalize(fr.MaxThres)

	// Adjust for range crossing zero
	fr = handleZeroCrossing(fr)

	// Determine if the range spans a full circle
	fr.fullCircle = isFullCircle(fr)

	return fr
}

// handleZeroCrossing adjusts ranges for cases where the range wraps over 0 (e.g., 350-10 degrees).
func handleZeroCrossing(fr FuzzyRange) FuzzyRange {
	switch {
	case fr.MinThres > fr.Min:
		// Threshold starts higher than minimum: move range forward
		fr.Min++
		fr.Max++
		fr.MaxThres++
		fr.crossesZero = true

	case fr.Min > fr.Max:
		// Minimum is greater than maximum: range crosses zero
		fr.Max++
		fr.MaxThres++
		fr.crossesZero = true

	case fr.MaxThres < fr.Max:
		// Maximum threshold is less than max: move range forward
		fr.MaxThres++
		fr.crossesZero = true
	}

	return fr
}

// isFullCircle checks if the range spans the entire circle (or wraps to cover all values).
func isFullCircle(fr FuzzyRange) bool {
	d := math.Wrap(fr.MaxThres-fr.MinThres, 0.0, 1.0)
	return d == 1.0 || d == 0.0
}

type FuzzyRange struct {
	MinThres, Min, Max, MaxThres float64
	crossesZero                  bool
	fullCircle                   bool
}

func (fr *FuzzyRange) String() string {
	return fmt.Sprintf("%f < %f - %f > %f (crosses zero: %v, full circle: %v)", fr.MinThres, fr.Min, fr.Max, fr.MaxThres, fr.crossesZero, fr.fullCircle)
}

func (fr *FuzzyRange) calc(v float64) float64 {
	v = math.Clamp(v, 0.0, 1.0)
	switch {
	case v < fr.MinThres || v > fr.MaxThres:
		// Outside the range
		return 0.0
	case fr.Min == fr.MinThres && v >= fr.MinThres && v <= fr.Max:
		// Inside range without fade
		return 1.0
	case fr.Max == fr.MaxThres && v >= fr.Min && v <= fr.Max:
		// Inside range without fade
		return 1.0

	case v >= fr.Min && v <= fr.Max:
		// Fully within the range
		return 1.0

	case v <= fr.Min:
		// Fade below minimum
		return fr.calculateFade(v, fr.Min, fr.MinThres, true)

	case v >= fr.Max:
		// Fade above maximum
		return fr.calculateFade(v, fr.MaxThres, fr.Max, false)

	default:
		// Fallback (shouldn't reach here)
		return 0.0
	}
}

// calculateFade handles fade calculations for values outside the core range.
// `isBelow` determines if the fade is below (`true`) or above (`false`) the range.
func (fr *FuzzyRange) calculateFade(v, start, end float64, isBelow bool) float64 {
	if isBelow {
		// Fade below the range (v approaches start from end)
		return 1 - (start-v)/math.Abs(start-end)
	}
	// Fade above the range (v moves away from start toward end)
	return (v - start) / math.Abs(end-start)
}

func (fr *FuzzyRange) calcWrapped(v float64) float64 {
	minThres := fr.MinThres
	min := fr.Min
	max := fr.Max
	maxThres := fr.MaxThres
	v /= 360.0

	for max <= min {
		maxThres += 1.0
		max += 1.0
	}

	if fr.fullCircle {
		if v >= max && v <= maxThres {
			return 1.0 - (v-max)/(maxThres-max)
		}
		if (v+1.0) >= max && (v+1.0) <= maxThres {
			return 1.0 - ((v+1.0)-max)/(maxThres-max)
		}
		if minThres < 1.0 && min > 1.0 && (v > minThres || (v+1.0) < min) {
			if v+1.0 < min {
				return ((v + 1.0) - minThres) / (min - minThres)
			} else if v > minThres {
				return (v - minThres) / (min - minThres)
			}
			return 1.0
		}
		if v >= minThres && v <= min {
			return (v - minThres) / (min - minThres)
		}
		return 1.0
	}

	// this is a partial circle then

	// out of range:

	if (!fr.crossesZero && (v < minThres || v > maxThres)) ||
		v >= (maxThres-1.0) && v <= minThres {
		return 0.0
	}

	// in range:

	if (v >= min && v <= max) || (v >= (min-1.0) && v <= (max-1.0)) {
		return 1.0
	}

	// in left fade:

	if ((v+1 >= maxThres || v >= maxThres) && v <= max) ||
		(fr.crossesZero && (minThres < min && min > max && v >= 0.0 && v <= min) ||
			(fr.crossesZero && minThres < 0.0 && min < 1.0 && v >= 0.0 && v <= min) ||
			(fr.crossesZero && v >= 0.0 && v <= minThres && v <= (min-1.0))) {

		if v <= min && (v >= minThres || v >= 0) {
			// left fade part

			if minThres == 0 && min <= 1.0 && v >= 0.0 && v <= min {
				return 1.0 - (1.0 - (v / min))
			}
			if v >= minThres && v <= min && min <= 1.0 {
				return 1.0 - (1.0 - (v-minThres)/(min-minThres))
			}
			if v >= (minThres-1.0) && v <= 1 && min >= 1.0 && min < 2.0 { // crosses 0
				v -= (minThres - 1.0)
				if v > 1.0 {
					v--
				}
				if minThres > min {
					return 1.0 - (1.0 - (v / ((min + 1.0) - minThres)))
				}
				return 1.0 - (1.0 - (v / math.Abs(minThres-min)))
			}
			if v >= minThres {
				return (v - 1.0) / (min - minThres)
			}
			return 1.0
		}

		return 0.0
	}

	// in right fade:

	if v >= (min-1.0) && v <= minThres && max >= 1.0 {
		return 1.0 - (v-(max-1.0))/(maxThres-max)
	}

	if v >= 0.0 && v <= 1.0 && (v >= (maxThres-1.0) || (v <= minThres && maxThres > 1.0)) {
		// right fade over zero crossing
		if v <= minThres {
			v++
		}
		return 1.0 - (v-max)/(maxThres-max)
	}

	return 0.0 // if we get here the range spans more than one circle and this is the overlap
}

type FuzzyRangeHSLA struct {
	hue FuzzyRange
	sat FuzzyRange
	lum FuzzyRange
}

func (f *FuzzyRangeHSLA) Calc(c *hsla.HSLA) *hsla.HSLA {
	h := f.hue.calcWrapped(c.H())
	s := f.sat.calc(c.S())
	l := f.lum.calc(c.L())
	alpha := h * l * s
	if alpha <= 0 {
		return hsla.New(0, 0, 0, 0) // Fully transparent if any component is out of range
	}
	c.SetA(alpha)
	return c
}

func Extract(i *image.Image, cf *filter.ColorFilter) *image.Image {
	fr := FuzzyRangeHSLA{
		hue: prepHueRange(FuzzyRange{cf.MinThres.H(), cf.Min.H(), cf.Max.H(), cf.MaxThres.H(), false, false}),
		sat: FuzzyRange{math.Clamp(cf.MinThres.S(), 0.0, 1.0), math.Clamp(cf.Min.S(), 0.0, 1.0), math.Clamp(cf.Max.S(), 0.0, 1.0), math.Clamp(cf.MaxThres.S(), 0.0, 1.0), false, false},
		lum: FuzzyRange{math.Clamp(cf.MinThres.L(), 0.0, 1.0), math.Clamp(cf.Min.L(), 0.0, 1.0), math.Clamp(cf.Max.L(), 0.0, 1.0), math.Clamp(cf.MaxThres.L(), 0.0, 1.0), false, false},
	}
	return i.ProcessHSLA(0, 0, i.W(), i.H(), func(x, y int, col *hsla.HSLA) (x2 int, y2 int, col2 *hsla.HSLA) { return x, y, fr.Calc(col) })
}
