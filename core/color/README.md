# Color Package

The color package provides color space conversions and utilities for working with different color representations.

## Structure

- Color space conversions (RGB, HSL, HSV, CMYK, etc.)
- Color manipulation utilities
- Color model implementations

## Usage

```go
// Example: Converting between color spaces
rgb := color.NewRGB(1.0, 0.0, 0.0) // Red
hsl := rgb.ToHSL()
hsv := hsl.ToHSV()
```

## Color Spaces

- RGB (Red, Green, Blue)
- HSL (Hue, Saturation, Lightness)
- HSV (Hue, Saturation, Value)
- CMYK (Cyan, Magenta, Yellow, Key/Black)
- YCbCr (Luminance, Blue Chroma, Red Chroma)
- XYZ (CIE 1931)
- LAB (CIE L*a*b*)

## Implementation Details

Each color space:
1. Provides conversion methods to other color spaces
2. Implements proper color space transformations
3. Handles color gamut and range limitations
4. Preserves color accuracy during conversions

## Testing

Run tests with:
```bash
go test ./...
```
