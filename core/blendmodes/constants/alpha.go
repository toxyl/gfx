package constants

// Alpha value constants
const (
	AlphaTransparent float64 = 0.0
	AlphaOpaque      float64 = 1.0
	AlphaHalf        float64 = 0.5
)

// Epsilon constants for floating point comparisons
const (
	// Epsilon for general floating point comparisons
	Epsilon float64 = 1e-10
	// AlphaEpsilon is specifically for alpha comparisons
	// Using a slightly larger epsilon for alpha since it's often
	// the result of multiple operations
	AlphaEpsilon float64 = 1e-8
)

// IsAlmostEqual checks if two float64 values are equal within epsilon
func IsAlmostEqual(a, b float64) bool {
	return Abs(a-b) < Epsilon
}

// IsAlmostZero checks if a float64 value is almost zero
func IsAlmostZero(a float64) bool {
	return Abs(a) < Epsilon
}

// IsAlmostOne checks if a float64 value is almost one
func IsAlmostOne(a float64) bool {
	return IsAlmostEqual(a, 1.0)
}

// IsAlphaAlmostEqual checks if two alpha values are equal within AlphaEpsilon
func IsAlphaAlmostEqual(a, b float64) bool {
	return Abs(a-b) < AlphaEpsilon
}

// IsAlphaAlmostZero checks if an alpha value is almost zero
func IsAlphaAlmostZero(a float64) bool {
	return Abs(a) < AlphaEpsilon
}

// IsAlphaAlmostOne checks if an alpha value is almost one
func IsAlphaAlmostOne(a float64) bool {
	return IsAlphaAlmostEqual(a, 1.0)
}

// Abs returns the absolute value of a float64
func Abs(x float64) float64 {
	if x < 0 {
		return -x
	}
	return x
}
