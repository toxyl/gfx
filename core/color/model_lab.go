// core/color/base_lab.go
package color

import (
	"fmt"

	"github.com/toxyl/gfx/core/color/constants"
	"github.com/toxyl/math"
)

//////////////////////////////////////////////////////
// Conversion utilities
//////////////////////////////////////////////////////

func rgbToXyz(r, g, b float64) (x, y, z float64) {
	// Handle exact white
	if math.Abs(r-1) < 1e-10 && math.Abs(g-1) < 1e-10 && math.Abs(b-1) < 1e-10 {
		return constants.LAB_RefX, constants.LAB_RefY, constants.LAB_RefZ
	}

	// Handle near-white colors
	if math.Abs(r-1) < 0.0001 && math.Abs(g-1) < 0.0001 && math.Abs(b-1) < 0.0001 {
		// Convert to linear RGB
		r = srgbToLinear(r)
		// All channels are equal, so just use one
		return r * constants.LAB_RefX, r * constants.LAB_RefY, r * constants.LAB_RefZ
	}

	// Handle gray colors (equal RGB values)
	if math.Abs(r-g) < 1e-10 && math.Abs(g-b) < 1e-10 {
		// Convert to linear RGB
		r = srgbToLinear(r)
		// All channels are equal, so just use one
		return r * constants.LAB_RefX, r * constants.LAB_RefY, r * constants.LAB_RefZ
	}

	// Convert RGB to linear RGB
	r = srgbToLinear(r)
	g = srgbToLinear(g)
	b = srgbToLinear(b)

	// Convert linear RGB to XYZ
	x = r*constants.LAB_RGB_X1 + g*constants.LAB_RGB_X2 + b*constants.LAB_RGB_X3
	y = r*constants.LAB_RGB_Y1 + g*constants.LAB_RGB_Y2 + b*constants.LAB_RGB_Y3
	z = r*constants.LAB_RGB_Z1 + g*constants.LAB_RGB_Z2 + b*constants.LAB_RGB_Z3

	return x, y, z
}

func xyzToLab(x, y, z float64) (l, a, b float64) {
	// Handle exact white
	if math.Abs(x/constants.LAB_RefX-1) < 1e-10 && math.Abs(y/constants.LAB_RefY-1) < 1e-10 && math.Abs(z/constants.LAB_RefZ-1) < 1e-10 {
		return 100, 0, 0
	}

	// Handle near-white colors
	if math.Abs(x/constants.LAB_RefX-1) < 0.0001 && math.Abs(y/constants.LAB_RefY-1) < 0.0001 && math.Abs(z/constants.LAB_RefZ-1) < 0.0001 {
		l = 116*xyzToLabF(y/constants.LAB_RefY) - 16
		return l, 0, 0
	}

	// Handle gray colors (equal RGB values)
	if math.Abs(x/constants.LAB_RefX-y/constants.LAB_RefY) < 1e-10 && math.Abs(y/constants.LAB_RefY-z/constants.LAB_RefZ) < 1e-10 {
		l = 116*xyzToLabF(y/constants.LAB_RefY) - 16
		return l, 0, 0
	}

	// Convert XYZ to CIELAB
	fx := xyzToLabF(x / constants.LAB_RefX)
	fy := xyzToLabF(y / constants.LAB_RefY)
	fz := xyzToLabF(z / constants.LAB_RefZ)

	l = 116*fy - 16
	a = 500 * (fx - fy)
	b = 200 * (fy - fz)

	return l, a, b
}

func labToXyz(l, a, b float64) (x, y, z float64) {
	// Handle exact white
	if math.Abs(l-100) < 1e-10 && math.Abs(a) < 1e-10 && math.Abs(b) < 1e-10 {
		return constants.LAB_RefX, constants.LAB_RefY, constants.LAB_RefZ
	}

	// Handle near-white colors
	if math.Abs(l-100) < 0.01 && math.Abs(a) < 1e-10 && math.Abs(b) < 1e-10 {
		y = labToXyzF((l+16)/116) * constants.LAB_RefY
		return y, y, y
	}

	// Handle gray colors (a and b close to zero)
	if math.Abs(a) < 1e-10 && math.Abs(b) < 1e-10 {
		y = labToXyzF((l+16)/116) * constants.LAB_RefY
		return y, y, y
	}

	// Convert CIELAB to XYZ
	y = (l + 16) / 116
	x = a/500 + y
	z = y - b/200

	x = labToXyzF(x) * constants.LAB_RefX
	y = labToXyzF(y) * constants.LAB_RefY
	z = labToXyzF(z) * constants.LAB_RefZ

	return x, y, z
}

func xyzToRgb(x, y, z float64) (r, g, b float64) {
	// Handle exact white
	if math.Abs(x-constants.LAB_RefX) < 1e-10 && math.Abs(y-constants.LAB_RefY) < 1e-10 && math.Abs(z-constants.LAB_RefZ) < 1e-10 {
		return 1, 1, 1
	}

	// Handle near-white colors
	if math.Abs(x-constants.LAB_RefX) < 0.0001 && math.Abs(y-constants.LAB_RefY) < 0.0001 && math.Abs(z-constants.LAB_RefZ) < 0.0001 {
		// Convert Y to linear RGB (all channels will be equal)
		r = y / constants.LAB_RefY
		return linearToSrgb(r), linearToSrgb(r), linearToSrgb(r)
	}

	// Handle gray colors (equal XYZ values relative to reference white)
	if math.Abs(x/constants.LAB_RefX-y/constants.LAB_RefY) < 1e-10 && math.Abs(y/constants.LAB_RefY-z/constants.LAB_RefZ) < 1e-10 {
		// Convert Y to linear RGB (all channels will be equal)
		r = y / constants.LAB_RefY
		return linearToSrgb(r), linearToSrgb(r), linearToSrgb(r)
	}

	// Convert XYZ to linear RGB
	r = x*constants.LAB_XYZ_R1 + y*constants.LAB_XYZ_R2 + z*constants.LAB_XYZ_R3
	g = x*constants.LAB_XYZ_G1 + y*constants.LAB_XYZ_G2 + z*constants.LAB_XYZ_G3
	b = x*constants.LAB_XYZ_B1 + y*constants.LAB_XYZ_B2 + z*constants.LAB_XYZ_B3

	// Convert linear RGB to sRGB
	r = linearToSrgb(r)
	g = linearToSrgb(g)
	b = linearToSrgb(b)

	return r, g, b
}

func srgbToLinear(c float64) float64 {
	if c <= 0.04045 {
		return c / constants.LAB_SRGB_LinearScale
	}
	return math.Pow((c+constants.LAB_SRGB_GammaOffset)/constants.LAB_SRGB_GammaScale, constants.LAB_SRGB_Gamma)
}

func linearToSrgb(c float64) float64 {
	if c <= constants.LAB_SRGB_LinearThreshold {
		return constants.LAB_SRGB_LinearScale * c
	}
	return constants.LAB_SRGB_GammaScale*math.Pow(c, 1/constants.LAB_SRGB_Gamma) - constants.LAB_SRGB_GammaOffset
}

func xyzToLabF(t float64) float64 {
	if t > constants.LAB_Epsilon {
		return math.Pow(t, 1.0/3.0)
	}
	return (constants.LAB_Kappa*t + 16) / 116
}

func labToXyzF(t float64) float64 {
	t3 := t * t * t
	if t3 > constants.LAB_Epsilon {
		return t3
	}
	return (116*t - 16) / constants.LAB_Kappa
}

//////////////////////////////////////////////////////
// Implementation check
//////////////////////////////////////////////////////

var _ iColor = (*LAB)(nil) // Ensure LAB implements the ColorModel interface.

//////////////////////////////////////////////////////
// Constructors
//////////////////////////////////////////////////////

// NewLAB creates a new LAB instance.
// It accepts channel values in the [0,100] range for L and [-128,128] range for a and b.
func NewLAB[N math.Number](l, a, b, alpha N) (*LAB, error) {
	return newColor(func() *LAB { return &LAB{} }, l, a, b, alpha)
}

// LABFromRGB converts an RGBA64 (RGB) to a LAB color.
func LABFromRGB(c *RGBA64) *LAB {
	x, y, z := rgbToXyz(c.R, c.G, c.B)
	l, a, b := xyzToLab(x, y, z)
	return &LAB{
		L:     l,
		A:     a,
		B:     b,
		Alpha: c.A,
	}
}

//////////////////////////////////////////////////////
// Type
//////////////////////////////////////////////////////

// LAB is a helper struct representing a color in the CIELAB color model with an alpha channel.
type LAB struct {
	L, A, B, Alpha float64
}

func (lab *LAB) Meta() *ColorModelMeta {
	return NewModelMeta(
		"LAB",
		"CIE L*a*b* color model.",
		NewChannelMeta("L", 0, 100, "", "Lightness."),
		NewChannelMeta("A", -128, 127, "", "Green to Red."),
		NewChannelMeta("B", -128, 127, "", "Blue to Yellow."),
		NewChannelMeta("Alpha", 0, 1, "", "Alpha channel."),
	)
}

//////////////////////////////////////////////////////
// Conversion
//////////////////////////////////////////////////////

func (lab *LAB) ToRGB() *RGBA64 {
	x, y, z := labToXyz(lab.L, lab.A, lab.B)
	r, g, b := xyzToRgb(x, y, z)
	return &RGBA64{R: r, G: g, B: b, A: lab.Alpha}
}

// ToRGBA64 converts the color to RGBA64.
func (lab *LAB) ToRGBA64() *RGBA64 {
	return lab.ToRGB()
}

// FromSlice initializes the color from a slice of float64 values.
func (lab *LAB) FromSlice(vals []float64) error {
	if len(vals) != 4 {
		return fmt.Errorf("LAB requires exactly 4 values: L, A, B, Alpha")
	}

	lab.L = vals[0]
	lab.A = vals[1]
	lab.B = vals[2]
	lab.Alpha = vals[3]

	return nil
}

// FromRGBA64 converts an RGBA64 color to this color model.
func (lab *LAB) FromRGBA64(rgba *RGBA64) iColor {
	return LABFromRGB(rgba)
}
