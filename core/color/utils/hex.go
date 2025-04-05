package utils

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/toxyl/math"
)

// Constants for hex color conversion
const (
	HexPrefix      = "#"
	HexMax         = 255.0
	HexBase        = 16
	HexBits        = 8
	HexShortLength = 3
	HexLongLength  = 6
	HexShortOffset = 1
	HexLongOffset  = 2
)

// RGBToHex converts RGB values to a hex color string.
// RGB values should be in range [0,1].
// Returns a hex color string in format "#RRGGBB".
func RGBToHex(r, g, b float64) string {
	// Clamp and convert to 8-bit values
	r8 := math.Round(r * HexMax)
	g8 := math.Round(g * HexMax)
	b8 := math.Round(b * HexMax)

	// Convert to hex string
	return fmt.Sprintf("%s%02X%02X%02X", HexPrefix, uint8(r8), uint8(g8), uint8(b8))
}

// HexToRGB converts a hex color string to RGB values.
// Accepts formats: "#RGB", "#RRGGBB"
// Returns RGB values in range [0,1].
func HexToRGB(hex string) (r, g, b float64, err error) {
	// Remove prefix and convert to uppercase
	hex = strings.ToUpper(strings.TrimPrefix(hex, HexPrefix))

	// Validate length
	if len(hex) != HexShortLength && len(hex) != HexLongLength {
		return 0, 0, 0, fmt.Errorf("invalid hex length: %d", len(hex))
	}

	// Handle short format (#RGB)
	if len(hex) == HexShortLength {
		// Expand to long format
		hex = string([]byte{
			hex[0], hex[0],
			hex[1], hex[1],
			hex[2], hex[2],
		})
	}

	// Parse RGB components
	rVal, err := strconv.ParseUint(hex[0:2], HexBase, HexBits)
	if err != nil {
		return 0, 0, 0, fmt.Errorf("invalid red component: %v", err)
	}

	gVal, err := strconv.ParseUint(hex[2:4], HexBase, HexBits)
	if err != nil {
		return 0, 0, 0, fmt.Errorf("invalid green component: %v", err)
	}

	bVal, err := strconv.ParseUint(hex[4:6], HexBase, HexBits)
	if err != nil {
		return 0, 0, 0, fmt.Errorf("invalid blue component: %v", err)
	}

	// Convert to [0,1] range
	r = float64(rVal) / HexMax
	g = float64(gVal) / HexMax
	b = float64(bVal) / HexMax

	return r, g, b, nil
}
