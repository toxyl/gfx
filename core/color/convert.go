package color

import "fmt"

// ConvertTo converts a color to the specified model
func ConvertTo[M iColor](from iColor) (M, error) {
	// Generic conversion helper
	switch v := from.(type) {
	case *RGB8:
		return convertRGB8To[M](v)
	case *LAB:
		return convertLABTo[M](v)
	case *HSL:
		return convertHSLTo[M](v)
	// ... other cases
	default:
		var zero M
		return zero, fmt.Errorf("unsupported conversion from %T", from)
	}
}

// convertRGB8To converts an RGB8 color to the specified model
func convertRGB8To[M iColor](rgb *RGB8) (M, error) {
	var zero M
	switch any(zero).(type) {
	case *LAB:
		return any(LABFromRGB(rgb.ToRGBA64())).(M), nil
	case *HSL:
		return any(HSLFromRGB(rgb.ToRGBA64())).(M), nil
	default:
		return zero, fmt.Errorf("unsupported conversion from RGB8 to %T", zero)
	}
}

// convertLABTo converts a LAB color to the specified model
func convertLABTo[M iColor](lab *LAB) (M, error) {
	var zero M
	switch any(zero).(type) {
	case *RGB8:
		return any(RGB8FromRGB(lab.ToRGBA64())).(M), nil
	case *HSL:
		rgb := lab.ToRGBA64()
		return any(HSLFromRGB(rgb)).(M), nil
	default:
		return zero, fmt.Errorf("unsupported conversion from LAB to %T", zero)
	}
}

// convertHSLTo converts an HSL color to the specified model
func convertHSLTo[M iColor](hsl *HSL) (M, error) {
	var zero M
	switch any(zero).(type) {
	case *RGB8:
		return any(RGB8FromRGB(hsl.ToRGBA64())).(M), nil
	case *LAB:
		rgb := hsl.ToRGBA64()
		return any(LABFromRGB(rgb)).(M), nil
	default:
		return zero, fmt.Errorf("unsupported conversion from HSL to %T", zero)
	}
}
