package utils

import (
	"github.com/toxyl/gfx/core/color/constants"
	"github.com/toxyl/math"
)

// Constants for color space conversion
const (
	// D65 illuminant
	Xn = 0.95047
	Yn = 1.00000
	Zn = 1.08883

	// Constants for color conversion
	Kappa   = 903.3
	Epsilon = 0.008856

	// sRGB conversion constants
	SRGB_LinearThreshold = 0.0031308
	SRGB_LinearScale     = 12.92
	SRGB_Gamma           = 2.4
	SRGB_GammaOffset     = 0.055
	SRGB_GammaScale      = 1.055

	// sRGB to XYZ conversion matrix
	RGB_X1 = 0.4124564
	RGB_X2 = 0.3575761
	RGB_X3 = 0.1804375
	RGB_Y1 = 0.2126729
	RGB_Y2 = 0.7151522
	RGB_Y3 = 0.0721750
	RGB_Z1 = 0.0193339
	RGB_Z2 = 0.1191920
	RGB_Z3 = 0.9503041

	// XYZ to RGB conversion matrix
	XYZ_R1 = 3.2404542
	XYZ_R2 = -1.5371385
	XYZ_R3 = -0.4985314
	XYZ_G1 = -0.9692660
	XYZ_G2 = 1.8760108
	XYZ_G3 = 0.0415560
	XYZ_B1 = 0.0556434
	XYZ_B2 = -0.2040259
	XYZ_B3 = 1.0572252
)

// RGBToXYZ converts RGB to XYZ color space.
// RGB values should be in range [0,1].
// Returns XYZ values in range [0,1].
func RGBToXYZ(r, g, b float64) (x, y, z float64) {
	// Convert to linear RGB
	r = linearizeRGB(r)
	g = linearizeRGB(g)
	b = linearizeRGB(b)

	// Convert to XYZ using D65 illuminant
	x = r*RGB_X1 + g*RGB_X2 + b*RGB_X3
	y = r*RGB_Y1 + g*RGB_Y2 + b*RGB_Y3
	z = r*RGB_Z1 + g*RGB_Z2 + b*RGB_Z3

	return x, y, z
}

// XYZToRGB converts XYZ to RGB color space.
// XYZ values should be in range [0,1].
// Returns RGB values in range [0,1].
func XYZToRGB(x, y, z float64) (r, g, b float64) {
	// Convert to linear RGB using D65 illuminant
	r = x*XYZ_R1 + y*XYZ_R2 + z*XYZ_R3
	g = x*XYZ_G1 + y*XYZ_G2 + z*XYZ_G3
	b = x*XYZ_B1 + y*XYZ_B2 + z*XYZ_B3

	// Convert to sRGB
	r = delinearizeRGB(r)
	g = delinearizeRGB(g)
	b = delinearizeRGB(b)

	return r, g, b
}

// XYZToLAB converts XYZ to LAB color space.
func XYZToLAB(x, y, z float64) (l, a, b float64) {
	// Convert to LAB using D65 reference white
	fx := f(x / Xn)
	fy := f(y / Yn)
	fz := f(z / Zn)

	l = 116*fy - 16
	a = 500 * (fx - fy)
	b = 200 * (fy - fz)

	return l, a, b
}

// RGBToLAB converts RGB to LAB color space.
// RGB values should be in range [0,1].
// Returns L (lightness) in range [0,100], a and b in range [-128,128].
func RGBToLAB(r, g, b float64) (l, a, bVal float64) {
	// First convert RGB to XYZ
	x, y, z := RGBToXYZ(r, g, b)

	// Convert XYZ to LAB
	fx := f(x / Xn)
	fy := f(y / Yn)
	fz := f(z / Zn)

	l = 116*fy - 16
	a = 500 * (fx - fy)
	bVal = 200 * (fy - fz)

	return l, a, bVal
}

// LABToRGB converts LAB to RGB color space.
// L should be in range [0,100], a and b in range [-128,128].
// Returns RGB values in range [0,1].
func LABToRGB(l, a, b float64) (r, g, bVal float64) {
	// Convert LAB to XYZ
	y := (l + 16) / 116
	x := a/500 + y
	z := y - b/200

	// Convert XYZ to RGB
	return XYZToRGB(
		Xn*fInv(x),
		Yn*fInv(y),
		Zn*fInv(z),
	)
}

// LABToXYZ converts LAB to XYZ color space.
func LABToXYZ(l, a, b float64) (x, y, z float64) {
	// Convert to XYZ
	fy := (l + 16) / 116
	fx := a/500 + fy
	fz := fy - b/200

	// D65 reference white
	x = Xn * labInverseTransform(fx)
	y = Yn * labInverseTransform(fy)
	z = Zn * labInverseTransform(fz)

	return x, y, z
}

// LABToLCH converts LAB to LCH color space.
func LABToLCH(l, a, b float64) (l2, c, h float64) {
	// Handle exact white
	if math.Abs(l-100) < constants.WhiteThreshold && math.Abs(a) < constants.WhiteThreshold && math.Abs(b) < constants.WhiteThreshold {
		return 100, 0, 0
	}

	// Handle near-white colors
	if math.Abs(l-100) < constants.NearWhiteThreshold && math.Abs(a) < constants.WhiteThreshold && math.Abs(b) < constants.WhiteThreshold {
		return l, 0, 0
	}

	// Handle gray colors (a and b close to zero)
	if math.Abs(a) < constants.WhiteThreshold && math.Abs(b) < constants.WhiteThreshold {
		return l, 0, 0
	}

	c = math.Sqrt(a*a + b*b)

	// Handle very small chroma values
	if c < constants.ChromaThreshold {
		return l, 0, 0
	}

	h = math.Atan2(b, a) * constants.DegreesPerRadian

	if h < 0 {
		h += constants.DegreesPerCircle
	}

	return l, c, h
}

// LCHToLAB converts LCH to LAB color space.
func LCHToLAB(l, c, h float64) (l2, a, b float64) {
	// Handle exact white
	if math.Abs(l-100) < constants.WhiteThreshold && math.Abs(c) < constants.WhiteThreshold {
		return 100, 0, 0
	}

	// Handle near-white colors
	if math.Abs(l-100) < constants.NearWhiteThreshold && math.Abs(c) < constants.WhiteThreshold {
		return l, 0, 0
	}

	// Handle gray colors (chroma close to zero)
	if math.Abs(c) < constants.ChromaThreshold {
		return l, 0, 0
	}

	h = h * constants.RadiansPerDegree
	a = c * math.Cos(h)
	b = c * math.Sin(h)

	return l, a, b
}

// RGBToLCH converts RGB to LCH color space.
// RGB values should be in range [0,1].
// Returns L in range [0,100], C in range [0,100], H in range [0,360].
func RGBToLCH(r, g, b float64) (l, c, h float64) {
	// First convert to LAB
	labL, labA, labB := RGBToLAB(r, g, b)

	// Then convert to LCH
	return LABToLCH(labL, labA, labB)
}

// LCHToRGB converts LCH to RGB color space.
// L should be in range [0,100], C in range [0,100], H in range [0,360].
// Returns RGB values in range [0,1].
func LCHToRGB(l, c, h float64) (r, g, b float64) {
	// First convert to LAB
	labL, labA, labB := LCHToLAB(l, c, h)

	// Then convert to RGB
	return LABToRGB(labL, labA, labB)
}

// RGBToLUV converts RGB to LUV color space.
// RGB values should be in range [0,1].
// Returns L in range [0,100], u and v in range [-100,100].
func RGBToLUV(r, g, b float64) (l, u, v float64) {
	// First convert to XYZ
	x, y, z := RGBToXYZ(r, g, b)

	// Calculate u' and v' for the color
	uPrime := 4 * x / (x + 15*y + 3*z)
	vPrime := 9 * y / (x + 15*y + 3*z)

	// Calculate u' and v' for the white point
	uPrimeN := 4 * Xn / (Xn + 15*Yn + 3*Zn)
	vPrimeN := 9 * Yn / (Xn + 15*Yn + 3*Zn)

	// Calculate L*
	if y/Yn > Epsilon {
		l = 116*math.Pow(y/Yn, 1.0/3.0) - 16
	} else {
		l = Kappa * y / Yn
	}

	// Calculate u* and v*
	u = 13 * l * (uPrime - uPrimeN)
	v = 13 * l * (vPrime - vPrimeN)

	return l, u, v
}

// LUVToRGB converts LUV to RGB color space.
// L should be in range [0,100], u and v in range [-100,100].
// Returns RGB values in range [0,1].
func LUVToRGB(l, u, v float64) (r, g, b float64) {
	// Calculate u' and v' for the white point
	uPrimeN := 4 * Xn / (Xn + 15*Yn + 3*Zn)
	vPrimeN := 9 * Yn / (Xn + 15*Yn + 3*Zn)

	// Calculate u' and v' for the color
	uPrime := u/(13*l) + uPrimeN
	vPrime := v/(13*l) + vPrimeN

	// Calculate Y
	var y float64
	if l > Kappa*Epsilon {
		y = math.Pow((l+16)/116, 3) * Yn
	} else {
		y = l * Yn / Kappa
	}

	// Calculate X and Z
	x := y * 9 * uPrime / (4 * vPrime)
	z := y * (12 - 3*uPrime - 20*vPrime) / (4 * vPrime)

	// Convert to RGB
	return XYZToRGB(x, y, z)
}

// Helper functions for RGB conversion
func linearizeRGB(c float64) float64 {
	if c <= 0.04045 {
		return c / SRGB_LinearScale
	}
	return math.Pow((c+SRGB_GammaOffset)/SRGB_GammaScale, SRGB_Gamma)
}

func delinearizeRGB(c float64) float64 {
	if c <= SRGB_LinearThreshold {
		return SRGB_LinearScale * c
	}
	return SRGB_GammaScale*math.Pow(c, 1/SRGB_Gamma) - SRGB_GammaOffset
}

// Helper functions for LAB conversion
func labTransform(t float64) float64 {
	if t > 0.008856 {
		return math.Pow(t, 1.0/3.0)
	}
	return 7.787*t + 16.0/116.0
}

func labInverseTransform(t float64) float64 {
	if t > 0.206893 {
		return math.Pow(t, 3.0)
	}
	return (t - 16.0/116.0) / 7.787
}

// Helper function for LAB conversion
func f(t float64) float64 {
	if t > Epsilon {
		return math.Pow(t, 1.0/3.0)
	}
	return (Kappa*t + 16) / 116
}

// Inverse helper function for LAB conversion
func fInv(t float64) float64 {
	t3 := t * t * t
	if t3 > Epsilon {
		return t3
	}
	return (116*t - 16) / Kappa
}
