package filter

import (
	"fmt"

	"github.com/toxyl/gfx/math"
)

type FilterValue struct {
	Min, Max, Tolerance float64
}

func (fv *FilterValue) String() string {
	return fmt.Sprintf("min: %f, max: %f, tolerance: %f", fv.Min, fv.Max, fv.Tolerance)
}

func NewFilterValue(min, max, tolerance float64) FilterValue {
	fv := FilterValue{
		Min:       min,
		Max:       max,
		Tolerance: tolerance,
	}
	return fv
}

func ToFilterValue(val, tolerance, feather float64) FilterValue {
	if feather > tolerance {
		fmt.Printf("Feather exceeds tolerance, decreasing feather from %f to %f", feather, tolerance)
		feather = tolerance
	}
	// Calculate min and max range based on value and tolerance
	mn := val - tolerance
	mx := val + tolerance

	// Ensure symmetry and handle inverted ranges for circular values
	if mn > mx {
		mn, mx = mx, mn
	}

	// Avoid edge case biases for exact matches
	if mn == val {
		mn -= 0.01
	}
	if mx == val {
		mx += 0.01
	}
	t := math.Min((mx-mn)/2, feather)

	return NewFilterValue(mn, mx, t)
}
