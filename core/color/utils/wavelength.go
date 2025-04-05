package utils

import (
	"github.com/toxyl/gfx/core/color/constants"
	"github.com/toxyl/math"
)

// WavelengthToHue converts a wavelength in nanometers to a hue angle in degrees
func WavelengthToHue(w float64) float64 {
	// For black and white (saturation = 0), use the shortest wavelength
	if w < constants.WavelengthVioletMin {
		return constants.HueRed
	}

	// Map wavelength to hue using a piecewise linear function
	// The mapping is based on the visible spectrum:
	// Red: 620-750nm
	// Orange: 590-620nm
	// Yellow: 570-590nm
	// Green: 495-570nm
	// Blue: 450-495nm
	// Violet: 380-450nm
	var h float64
	if w >= constants.WavelengthRedMin { // Red
		h = constants.HueRed
	} else if w >= constants.WavelengthYellowMin { // Yellow
		h = constants.HueYellow + (w-constants.WavelengthYellowMin)/(constants.WavelengthRedMin-constants.WavelengthYellowMin)*constants.HueRange
	} else if w >= constants.WavelengthGreenMin { // Green
		h = constants.HueGreen
	} else if w >= constants.WavelengthCyanMin { // Cyan
		h = constants.HueCyan + (w-constants.WavelengthCyanMin)/(constants.WavelengthGreenMin-constants.WavelengthCyanMin)*constants.HueRange
	} else if w >= constants.WavelengthBlueMin { // Blue
		h = constants.HueBlue
	} else { // Magenta
		h = constants.HueMagenta + (w-constants.WavelengthVioletMin)/(constants.WavelengthBlueMin-constants.WavelengthVioletMin)*constants.HueRange
	}

	return math.Mod(h, constants.DegreesPerCircle)
}

// HueToWavelength converts a hue angle in degrees to a wavelength in nanometers
func HueToWavelength(h float64) float64 {
	// Normalize hue to [0,360]
	h = math.Mod(h, constants.DegreesPerCircle)
	if h < 0 {
		h += constants.DegreesPerCircle
	}

	// Map hue to wavelength using inverse of the above mapping
	var w float64
	if h < constants.HueYellow { // Red
		w = constants.WavelengthRedMin
	} else if h < constants.HueGreen { // Yellow
		w = constants.WavelengthYellowMin + (h-constants.HueYellow)/constants.HueRange*(constants.WavelengthRedMin-constants.WavelengthYellowMin)
	} else if h < constants.HueCyan { // Green
		w = constants.WavelengthGreenMin
	} else if h < constants.HueBlue { // Cyan
		w = constants.WavelengthCyanMin + (h-constants.HueCyan)/constants.HueRange*(constants.WavelengthGreenMin-constants.WavelengthCyanMin)
	} else if h < constants.HueMagenta { // Blue
		w = constants.WavelengthBlueMin
	} else { // Magenta
		w = constants.WavelengthVioletMin + (h-constants.HueMagenta)/(constants.DegreesPerCircle-constants.HueMagenta)*(constants.WavelengthRedMin-constants.WavelengthVioletMin)
	}

	return math.Clamp(w, constants.WavelengthVioletMin, constants.WavelengthRedMin)
}
