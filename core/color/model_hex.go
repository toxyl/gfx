// core/color/base_hex.go
package color

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/toxyl/gfx/core/color/constants"
	"github.com/toxyl/math"
)

//////////////////////////////////////////////////////
// Conversion utilities
//////////////////////////////////////////////////////

func rgbToHex(r, g, b float64) string {
	// Convert normalized [0,1] RGB to hex string
	r8 := uint8(math.Clamp(r, 0, 1) * constants.HEX_Max)
	g8 := uint8(math.Clamp(g, 0, 1) * constants.HEX_Max)
	b8 := uint8(math.Clamp(b, 0, 1) * constants.HEX_Max)
	return fmt.Sprintf("%s%02x%02x%02x", constants.HEX_Prefix, r8, g8, b8)
}

func hexToRgb(hex string) (r, g, b float64, err error) {
	// Remove prefix and convert to lowercase
	hex = strings.TrimPrefix(strings.ToLower(hex), constants.HEX_Prefix)

	// Parse hex string
	var r8, g8, b8 uint64
	switch len(hex) {
	case constants.HEX_ShortLength: // #RGB
		r8, err = strconv.ParseUint(hex[0:constants.HEX_ShortOffset]+hex[0:constants.HEX_ShortOffset], constants.HEX_Base, constants.HEX_Bits)
		if err != nil {
			return
		}
		g8, err = strconv.ParseUint(hex[1:constants.HEX_ShortOffset+1]+hex[1:constants.HEX_ShortOffset+1], constants.HEX_Base, constants.HEX_Bits)
		if err != nil {
			return
		}
		b8, err = strconv.ParseUint(hex[2:constants.HEX_ShortOffset+2]+hex[2:constants.HEX_ShortOffset+2], constants.HEX_Base, constants.HEX_Bits)
		if err != nil {
			return
		}
	case constants.HEX_LongLength: // #RRGGBB
		r8, err = strconv.ParseUint(hex[0:constants.HEX_LongOffset], constants.HEX_Base, constants.HEX_Bits)
		if err != nil {
			return
		}
		g8, err = strconv.ParseUint(hex[2:constants.HEX_LongOffset+2], constants.HEX_Base, constants.HEX_Bits)
		if err != nil {
			return
		}
		b8, err = strconv.ParseUint(hex[4:constants.HEX_LongOffset+4], constants.HEX_Base, constants.HEX_Bits)
		if err != nil {
			return
		}
	default:
		err = fmt.Errorf("invalid hex color format: %s", hex)
		return
	}

	// Convert to normalized [0,1] range
	r = float64(r8) / constants.HEX_Max
	g = float64(g8) / constants.HEX_Max
	b = float64(b8) / constants.HEX_Max
	return
}

//////////////////////////////////////////////////////
// Implementation check
//////////////////////////////////////////////////////

var _ iColor = (*Hex)(nil) // Ensure Hex implements the ColorModel interface.

//////////////////////////////////////////////////////
// Constructors
//////////////////////////////////////////////////////

// NewHex creates a new Hex instance from a hex color string.
// The string should be in the format "#RRGGBB" or "#RGB".
func NewHex(hex string, alpha float64) (*Hex, error) {
	h := &Hex{
		Hex:   hex,
		Alpha: alpha,
	}
	return h, nil
}

// HexFromRGB converts an RGBA64 (RGB) to a Hex color.
func HexFromRGB(c *RGBA64) *Hex {
	hex := rgbToHex(c.R, c.G, c.B)
	return &Hex{
		Hex:   hex,
		Alpha: c.A,
	}
}

//////////////////////////////////////////////////////
// Type
//////////////////////////////////////////////////////

// Hex is a helper struct representing a color in the Hex color model with an alpha channel.
type Hex struct {
	Hex   string
	Alpha float64
}

func (h *Hex) Meta() *ColorModelMeta {
	return NewModelMeta(
		"Hex",
		"Hexadecimal color model.",
		NewChannelMeta("Hex", 0, 1, "", "Hexadecimal color value."),
		NewChannelMeta("Alpha", 0, 1, "", "Alpha channel."),
	)
}

//////////////////////////////////////////////////////
// Conversion
//////////////////////////////////////////////////////

func (h *Hex) ToRGB() *RGBA64 {
	r, g, b, _ := hexToRgb(h.Hex)
	return &RGBA64{
		R: r,
		G: g,
		B: b,
		A: h.Alpha,
	}
}

// FromSlice initializes the color from a slice of float64 values.
func (h *Hex) FromSlice(values []float64) error {
	if len(values) != 1 {
		return fmt.Errorf("Hex requires exactly 1 value: Alpha")
	}

	h.Alpha = values[0]

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
