package color

import (
	"fmt"
	"regexp"
	"strconv"
	"strings"
)

// ParseColor parses a color string in various formats
func ParseColor(format string) (iColor, error) {
	format = strings.TrimSpace(format)

	// Try different formats
	if color, err := parseHex(format); err == nil {
		return color, nil
	}
	if color, err := parseRGB(format); err == nil {
		return color, nil
	}
	if color, err := parseHSL(format); err == nil {
		return color, nil
	}
	if color, err := parseNamedColor(format); err == nil {
		return color, nil
	}

	return nil, fmt.Errorf("invalid color format: %s", format)
}

// parseHex parses a hex color string
func parseHex(hex string) (iColor, error) {
	hex = strings.TrimPrefix(hex, "#")

	// Handle shorthand hex (e.g., #F00)
	if len(hex) == 3 {
		hex = string([]byte{
			hex[0], hex[0],
			hex[1], hex[1],
			hex[2], hex[2],
		})
	}

	if len(hex) != 6 {
		return nil, fmt.Errorf("invalid hex color length")
	}

	r, err := strconv.ParseInt(hex[0:2], 16, 64)
	if err != nil {
		return nil, err
	}
	g, err := strconv.ParseInt(hex[2:4], 16, 64)
	if err != nil {
		return nil, err
	}
	b, err := strconv.ParseInt(hex[4:6], 16, 64)
	if err != nil {
		return nil, err
	}

	return &RGB8{
		R:     float64(r),
		G:     float64(g),
		B:     float64(b),
		Alpha: 1.0,
	}, nil
}

// parseRGB parses an RGB color string
func parseRGB(rgb string) (iColor, error) {
	re := regexp.MustCompile(`rgb\((\d+),\s*(\d+),\s*(\d+)\)`)
	matches := re.FindStringSubmatch(rgb)
	if matches == nil {
		return nil, fmt.Errorf("invalid RGB format")
	}

	r, _ := strconv.ParseFloat(matches[1], 64)
	g, _ := strconv.ParseFloat(matches[2], 64)
	b, _ := strconv.ParseFloat(matches[3], 64)

	return &RGB8{
		R:     r,
		G:     g,
		B:     b,
		Alpha: 1.0,
	}, nil
}

// parseHSL parses an HSL color string
func parseHSL(hsl string) (iColor, error) {
	re := regexp.MustCompile(`hsl\((\d+),\s*(\d+)%,\s*(\d+)%\)`)
	matches := re.FindStringSubmatch(hsl)
	if matches == nil {
		return nil, fmt.Errorf("invalid HSL format")
	}

	h, _ := strconv.ParseFloat(matches[1], 64)
	s, _ := strconv.ParseFloat(matches[2], 64)
	l, _ := strconv.ParseFloat(matches[3], 64)

	// Create HSL with proper field names from model_hsl.go
	return NewHSL(h, s/100, l/100, 1.0)
}

// parseNamedColor parses a named color
func parseNamedColor(name string) (iColor, error) {
	switch strings.ToLower(name) {
	case "red":
		return &RGB8{R: 255, G: 0, B: 0, Alpha: 1.0}, nil
	case "green":
		return &RGB8{R: 0, G: 255, B: 0, Alpha: 1.0}, nil
	case "blue":
		return &RGB8{R: 0, G: 0, B: 255, Alpha: 1.0}, nil
	// Add more named colors as needed
	default:
		return nil, fmt.Errorf("unknown color name: %s", name)
	}
}

// Format formats a color in the specified format
func Format(c iColor, format string) (string, error) {
	switch format {
	case "hex":
		return ToHex(c), nil
	case "rgb":
		return ToRGBString(c), nil
	case "hsl":
		return ToHSLString(c), nil
	default:
		return "", fmt.Errorf("unsupported format: %s", format)
	}
}

// ToHex converts a color to hex format
func ToHex(c iColor) string {
	rgb := c.ToRGBA64()
	r := uint8(rgb.R * 255)
	g := uint8(rgb.G * 255)
	b := uint8(rgb.B * 255)
	return fmt.Sprintf("#%02X%02X%02X", r, g, b)
}

// ToRGBString converts a color to RGB string format
func ToRGBString(c iColor) string {
	rgb := c.ToRGBA64()
	r := uint8(rgb.R * 255)
	g := uint8(rgb.G * 255)
	b := uint8(rgb.B * 255)
	return fmt.Sprintf("rgb(%d, %d, %d)", r, g, b)
}

// ToHSLString converts a color to HSL string format
func ToHSLString(c iColor) string {
	hsl := HSLFromRGB(c.ToRGBA64())
	h := hsl.H
	s := hsl.S * 100
	l := hsl.L * 100
	return fmt.Sprintf("hsl(%.0f, %.0f%%, %.0f%%)", h, s, l)
}
