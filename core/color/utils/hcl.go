package utils

import (
	"github.com/toxyl/gfx/core/color/constants"
	"github.com/toxyl/math"
)

// LABToHCL converts LAB to HCL color space.
func LABToHCL(l, a, b float64) (h, c, l2 float64) {
	// Handle exact white
	if math.Abs(l-100) < constants.HCL_WhiteThreshold && math.Abs(a) < constants.HCL_WhiteThreshold && math.Abs(b) < constants.HCL_WhiteThreshold {
		return 0, 0, 100
	}

	// Handle near-white colors
	if math.Abs(l-100) < constants.HCL_NearWhiteThreshold && math.Abs(a) < constants.HCL_WhiteThreshold && math.Abs(b) < constants.HCL_WhiteThreshold {
		return 0, 0, l
	}

	// Handle gray colors (a and b close to zero)
	if math.Abs(a) < constants.HCL_WhiteThreshold && math.Abs(b) < constants.HCL_WhiteThreshold {
		return 0, 0, l
	}

	c = math.Sqrt(a*a + b*b)

	// Handle very small chroma values
	if c < constants.HCL_ChromaThreshold {
		return 0, 0, l
	}

	h = math.Atan2(b, a) * constants.HCL_DegreesPerRadian

	if h < 0 {
		h += constants.HCL_DegreesPerCircle
	}

	return h, c, l
}

// HCLToLAB converts HCL to LAB color space.
func HCLToLAB(h, c, l float64) (l2, a, b float64) {
	// Handle exact white
	if math.Abs(l-100) < constants.HCL_WhiteThreshold && math.Abs(c) < constants.HCL_WhiteThreshold {
		return 100, 0, 0
	}

	// Handle near-white colors
	if math.Abs(l-100) < constants.HCL_NearWhiteThreshold && math.Abs(c) < constants.HCL_WhiteThreshold {
		return l, 0, 0
	}

	// Handle gray colors (chroma close to zero)
	if math.Abs(c) < constants.HCL_ChromaThreshold {
		return l, 0, 0
	}

	h = h * constants.HCL_RadiansPerDegree
	a = c * math.Cos(h)
	b = c * math.Sin(h)

	return l, a, b
}

// HCLToRGB converts HCL to RGB color space.
func HCLToRGB(h, c, l float64) (r, g, b float64) {
	// First convert to LAB
	l2, a, b := HCLToLAB(l, c, h)

	// Then convert to XYZ
	x, y, z := LABToXYZ(l2, a, b)

	// Finally convert to RGB
	r, g, b = XYZToRGB(x, y, z)
	return r, g, b
}

// RGBToHCL converts RGB to HCL color space.
func RGBToHCL(r, g, b float64) (h, c, l float64) {
	// First convert to XYZ
	x, y, z := RGBToXYZ(r, g, b)

	// Then convert to LAB
	l2, a, b := XYZToLAB(x, y, z)

	// Finally convert to HCL
	h, c, l = LABToHCL(l2, a, b)
	return h, c, l
}
