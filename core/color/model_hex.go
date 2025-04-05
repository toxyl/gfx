// core/color/base_hex.go
package color

import (
	"fmt"

	"github.com/toxyl/gfx/core/color/utils"
	"github.com/toxyl/math"
)

//////////////////////////////////////////////////////
// Implementation check
//////////////////////////////////////////////////////

var _ iColor = (*Hex)(nil) // Ensure Hex implements the ColorModel interface.

//////////////////////////////////////////////////////
// Constructors
//////////////////////////////////////////////////////

// NewHex creates a new Hex instance.
// Alpha is in range [0,1].
func NewHex[N math.Number](hex string, alpha N) (*Hex, error) {
	// Validate hex string
	if _, _, _, err := utils.HexToRGB(hex); err != nil {
		return nil, fmt.Errorf("invalid hex color: %v", err)
	}

	// Create new instance
	h := &Hex{
		Hex:   hex,
		Alpha: float64(alpha),
	}

	return h, nil
}

// HexFromRGB converts an RGBA64 (RGB) to a Hex color.
func HexFromRGB(c *RGBA64) *Hex {
	// Convert RGB to hex string
	hex := utils.RGBToHex(c.R, c.G, c.B)

	return &Hex{
		Hex:   hex,
		Alpha: c.A,
	}
}

//////////////////////////////////////////////////////
// Type
//////////////////////////////////////////////////////

// Hex is a helper struct representing a color in the Hex color model.
// Hex represents colors using a hexadecimal string in format "#RRGGBB".
type Hex struct {
	Hex   string  // Hex color string in format "#RRGGBB"
	Alpha float64 // [0,1] Alpha
}

func (h *Hex) Meta() *ColorModelMeta {
	return NewModelMeta(
		"Hex",
		"Hexadecimal color model.",
		NewChannelMeta("Hex", 0, 1, "", "Hexadecimal color string."),
		NewChannelMeta("Alpha", 0, 1, "", "Alpha channel."),
	)
}

//////////////////////////////////////////////////////
// Conversion
//////////////////////////////////////////////////////

func (h *Hex) ToRGB() *RGBA64 {
	// Convert hex to RGB
	r, g, b, err := utils.HexToRGB(h.Hex)
	if err != nil {
		// Return black on error
		return &RGBA64{
			R: 0,
			G: 0,
			B: 0,
			A: h.Alpha,
		}
	}

	return &RGBA64{
		R: r,
		G: g,
		B: b,
		A: h.Alpha,
	}
}

// FromSlice initializes a Hex instance from a slice of float64 values.
// The slice must contain exactly 1 value: Alpha.
func (h *Hex) FromSlice(vals []float64) error {
	if len(vals) != 1 {
		return fmt.Errorf("Hex requires 1 value, got %d", len(vals))
	}

	h.Alpha = vals[0]

	return nil
}

// FromRGBA64 converts an RGBA64 color to this color model.
func (h *Hex) FromRGBA64(rgba *RGBA64) iColor {
	return HexFromRGB(rgba)
}

// ToRGBA64 converts the color to RGBA64.
func (h *Hex) ToRGBA64() *RGBA64 {
	return h.ToRGB()
}
