# Color Package

A comprehensive Go package for working with colors in various color spaces. This package provides a rich set of color models and conversion utilities, supporting both device-dependent and device-independent color spaces.

## Features

- Multiple color model support:
  - RGB-based: RGB8, Hex, Grayscale
  - HSL/HSV-based: HSL, HSB
  - Wavelength-based: λSL (Wavelength, Saturation, Lightness), λSB (Wavelength, Saturation, Brightness)
  - Device-independent: LAB, XYZ, LCH, HCL
  - Video/Television: YUV, YIQ, YCbCr
  - Printing: CMY, CMYK
- Consistent alpha channel handling across all models
- Comprehensive metadata for each color channel
- Range validation and automatic clamping
- Efficient conversion utilities
- Type-safe implementations
- Color manipulation and transformation
- Color palettes and constants
- Color parsing and formatting
- Color comparison and operations

## Installation

```bash
go get github.com/toxyl/gfx/core/color
```

## Basic Usage

### Creating Colors

```go
import (
    stdcolor "image"
    "github.com/toxyl/gfx/core/color"
)
// Using the builder pattern
color := color.NewColor().
    RGB(0.5, 0.3, 0.8).
    Alpha(1.0).
    Build()

// Create a color from standard library color.RGBA
stdColor := stdcolor.RGBA{R: 255, G: 0, B: 0, A: 255}
rgbColor := color.New(0, 0, 0, 0)
rgbColor.SetUint8(stdColor)

// Parse from string
color, err := color.ParseColor("#FF0000") // hex
color, err = color.ParseColor("rgb(255, 0, 0)") // rgb
color, err = color.ParseColor("hsl(0, 100%, 50%)") // hsl
color, err = color.ParseColor("red") // named color
```

### Color Manipulation

```go
// Lighten a color
lighter := color.Lighten(0.2)

// Darken a color
darker := color.Darken(0.2)

// Change saturation
moreSaturated := color.Saturate(0.2)
lessSaturated := color.Desaturate(0.2)

// Rotate hue
rotated := color.RotateHue(90)

// Mix colors
mixed := color.Mix(otherColor, 0.5) // 50% mix
```

### Color Conversions

```go
// Convert between color spaces
rgb := color.New(0.5, 0.3, 0.8, 1.0)

// Convert to LAB
lab, err := color.ConvertTo[*color.LAB](rgb)

// Convert to HSL
hsl, err := color.ConvertTo[*color.HSL](rgb)

// Convert back to RGB
rgb2 := lab.ToRGB()
```

### Working with Palettes

```go
// Create a custom palette
palette := color.NewPalette("My Palette").
    Add(color.NewColor().RGB(1, 0, 0).Build()).
    Add(color.NewColor().RGB(0, 1, 0).Build()).
    Add(color.NewColor().RGB(0, 0, 1).Build())

// Use predefined palettes
color := color.WebSafe.Get(0) // Get first color from web-safe palette
colors := color.MaterialDesign.Colors() // Get all material design colors
```

### Color Formatting

```go
// Format colors
hex := color.ToHex(rgbColor)
rgbStr := color.ToRGBString(rgbColor)
hslStr := color.ToHSLString(rgbColor)

// Format with specific format
str, err := color.Format(rgbColor, "hex")
str, err = color.Format(rgbColor, "rgb")
str, err = color.Format(rgbColor, "hsl")
```

### Color Comparison

```go
// Check if colors are equal
isEqual := color1.Equals(color2)

// Calculate color distance
distance := color1.Distance(color2)

// Check if colors are similar
isSimilar := color1.IsSimilar(color2, 0.1) // within 10% tolerance
```

### Channel Manipulation

```go
// Get channel value
value, err := color.GetChannel("r")

// Set channel value
err := color.SetChannel("g", 0.5)

// Get all channels
channels := color.Channels()
```

## Color Model Ranges

- RGB: [0, 1] for all channels
- HSL/HSB: H [0, 360], S/L/B [0, 1]
- λSL/λSB: λ [380, 750] nm, S/L/B [0, 1]
- LAB: L [0, 100], a/b [-128, 127]
- XYZ: [0, 1] for all channels
- Alpha: [0, 1] across all models

## Error Handling

All constructors return `(ColorModel, error)` pairs. Input values are validated against channel ranges, and conversion errors are properly propagated with descriptive error messages.

```go
// Example of error handling
hsl, err := color.NewHSL(400, 0.5, 0.7, 1.0) // Invalid hue value
if err != nil {
    fmt.Printf("Error creating HSL color: %v\n", err)
}
```

## Supported Color Models

1. **RGB-based Models**
   - RGB8 (8-bit RGB)
   - RGBA64 (64-bit float RGB)
   - Hex (Hexadecimal RGB)
   - Grayscale

2. **HSL/HSV-based Models**
   - HSL (Hue, Saturation, Lightness)
   - HSB (Hue, Saturation, Brightness)

3. **Wavelength-based Models**
   - λSL (Wavelength, Saturation, Lightness)
   - λSB (Wavelength, Saturation, Brightness)
4. **Device-independent Models**
   - LAB (CIELAB)
   - XYZ (CIE XYZ)
   - LCH (Lightness, Chroma, Hue)
   - HCL (Hue, Chroma, Lightness)

5. **Video/Television Models**
   - YUV (Luma, Blue Projection, Red Projection)
   - YIQ (Luma, In-phase, Quadrature)
   - YCbCr (Luma, Blue-difference, Red-difference)

6. **Printing Models**
   - CMY (Cyan, Magenta, Yellow)
   - CMYK (Cyan, Magenta, Yellow, Key/Black)

## Notes

- All color conversions use standard conversion matrices and formulas as defined by the International Commission on Illumination (CIE) and other relevant standards organizations.
- The package uses `RGBA64` as the canonical representation for all color conversions.
- All color channels are automatically clamped to their valid ranges during operations.
- Alpha channel is consistently handled across all color models.
- Wavelength-based models (λSL/λSB) use the visible spectrum range of 380-750 nanometers.
