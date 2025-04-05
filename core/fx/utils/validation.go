package utils

import (
	"fmt"

	"github.com/toxyl/math"
)

// ValidateRange checks if a value is within the specified range.
// Returns an error if the value is outside the range.
func ValidateRange(value, min, max float64, name string) error {
	if value < min || value > max {
		return fmt.Errorf("%s value %f is outside valid range [%f, %f]", name, value, min, max)
	}
	return nil
}

// ClampRange clamps a value to the specified range.
// Returns the clamped value.
func ClampRange(value, min, max float64) float64 {
	return math.Max(min, math.Min(max, value))
}

// ValidatePositive checks if a value is positive.
// Returns an error if the value is not positive.
func ValidatePositive(value float64, name string) error {
	if value <= 0 {
		return fmt.Errorf("%s value %f must be positive", name, value)
	}
	return nil
}

// ValidateNonNegative checks if a value is non-negative.
// Returns an error if the value is negative.
func ValidateNonNegative(value float64, name string) error {
	if value < 0 {
		return fmt.Errorf("%s value %f must be non-negative", name, value)
	}
	return nil
}
