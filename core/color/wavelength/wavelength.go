// core/color/wavelength/wavelength.go
package wavelength

import (
	"github.com/toxyl/gfx/core/color/constants"
	"github.com/toxyl/math"
)

// ToRGB converts a wavelength in nanometers to RGB values.
// Based on the algorithm from http://www.physics.sfasu.edu/astro/color/spectra.html
func ToRGB(wavelength float64) (r, g, b float64) {
	// Normalize wavelength to [0,1] range
	lambda := math.Clamp(wavelength, constants.WavelengthMin, constants.WavelengthMax)
	lambda = (lambda - constants.WavelengthMin) / (constants.WavelengthMax - constants.WavelengthMin)

	// Calculate RGB values based on wavelength
	if lambda < constants.WavelengthRange1 {
		r = 0
		g = 0
		b = constants.WavelengthBlueCoeff1 + constants.WavelengthBlueCoeff2*lambda/constants.WavelengthNorm1
	} else if lambda < constants.WavelengthRange2 {
		r = 0
		g = constants.WavelengthGreenCoeff1 + constants.WavelengthGreenCoeff2*(lambda-constants.WavelengthRange1)/constants.WavelengthNorm2
		b = 1.0
	} else if lambda < constants.WavelengthRange3 {
		r = constants.WavelengthRedCoeff1 + constants.WavelengthRedCoeff2*(lambda-constants.WavelengthRange2)/constants.WavelengthNorm3
		g = 1.0
		b = 1.0 - (lambda-constants.WavelengthRange2)/constants.WavelengthNorm3
	} else if lambda < constants.WavelengthRange4 {
		r = 1.0
		g = 1.0 - (lambda-constants.WavelengthRange3)/constants.WavelengthNorm4
		b = 0
	} else {
		r = 1.0 - constants.WavelengthRedCoeff3*(lambda-constants.WavelengthRange4)/constants.WavelengthNorm5
		g = 0
		b = 0
	}

	// Apply gamma correction
	r = math.Pow(r, constants.WavelengthGamma)
	g = math.Pow(g, constants.WavelengthGamma)
	b = math.Pow(b, constants.WavelengthGamma)

	return r, g, b
}

// RGBToHSV converts RGB values to HSV.
func RGBToHSV(r, g, b float64) (h, s, v float64) {
	max := math.Max(r, math.Max(g, b))
	min := math.Min(r, math.Min(g, b))

	v = max

	if max == min {
		h = 0
		s = 0
		return
	}

	s = (max - min) / max

	switch max {
	case r:
		h = (g - b) / (max - min)
		if g < b {
			h += 6
		}
	case g:
		h = (b-r)/(max-min) + 2
	case b:
		h = (r-g)/(max-min) + 4
	}
	h /= 6

	return h, s, v
}
