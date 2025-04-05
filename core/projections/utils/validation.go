package utils

import (
	"fmt"

	"github.com/toxyl/math"
)

// ValidateLatitude checks if a latitude value is valid.
// Returns an error if the latitude is outside [-90, 90].
func ValidateLatitude(lat float64) error {
	return ValidateRange(lat, -90, 90, "latitude")
}

// ValidateLongitude checks if a longitude value is valid.
// Returns an error if the longitude is outside [-180, 180].
func ValidateLongitude(lon float64) error {
	return ValidateRange(lon, -180, 180, "longitude")
}

// ValidateScale checks if a scale value is valid.
// Returns an error if the scale is not positive.
func ValidateScale(scale float64) error {
	return ValidatePositive(scale, "scale")
}

// ValidateRange checks if a value is within the specified range.
// Returns an error if the value is outside the range.
func ValidateRange(value, min, max float64, name string) error {
	if value < min || value > max {
		return fmt.Errorf("%s value %f is outside valid range [%f, %f]", name, value, min, max)
	}
	return nil
}

// ClampRange restricts value to be within the range [min, max].
func ClampRange(value, min, max float64) float64 {
	return math.Clamp(value, min, max)
}

// ValidatePositive checks if a value is positive.
// Returns an error if the value is not positive.
func ValidatePositive(value float64, name string) error {
	if value <= 0 {
		return fmt.Errorf("%s value %f must be positive", name, value)
	}
	return nil
}
