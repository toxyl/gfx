# Blend Modes Package

A comprehensive Go package for working with image blend modes, providing a rich set of compositing operations commonly found in professional image editing software. This package implements various blend modes for combining layers or images with precise control over opacity and composition.

## Features

- Multiple blend mode categories:
  - Basic: Normal, Erase
  - Darken: Darken, Multiply, Color Burn, Linear Burn, Darker Color
  - Lighten: Lighten, Screen, Color Dodge, Linear Dodge (Add), Lighter Color
  - Contrast: Overlay, Soft Light, Hard Light, Vivid Light, Linear Light, Pin Light, Hard Mix
  - Comparative: Difference, Exclusion, Subtract, Divide, Negation, Contrast Negate
  - Component: Hue, Saturation, Color, Luminosity
  - Special Effects: Reflect, Glow, Average
- Thread-safe blend mode registry
- Consistent alpha channel handling
- Comprehensive metadata for each blend mode
- Type-safe implementations
- Easy-to-use interface
- Support for custom blend modes

## Installation

```bash
go get github.com/toxyl/gfx/core/blendmodes
```

## Basic Usage

### Using Blend Modes

```go
import (
    "github.com/toxyl/gfx/core/blendmodes"
    "github.com/toxyl/gfx/core/color"
)

// Get a blend mode
multiply, err := blendmodes.Get("multiply")
if err != nil {
    // Handle error
}

// Create two colors to blend
bottom := color.NewColor().RGB(1.0, 0.0, 0.0).Build() // Red
top := color.NewColor().RGB(0.0, 1.0, 0.0).Build()    // Green

// Blend the colors with 50% opacity
result := multiply.Blend(bottom, top, 0.5)
```

### Registering Custom Blend Modes

```go
// Define a custom blend function
customBlend := func(bottom, top *color.RGBA64, alpha float64) *color.RGBA64 {
    // Your blend implementation here
    return result
}

// Register the custom blend mode
err := blendmodes.Register("custom", "Description of custom blend mode", customBlend)
if err != nil {
    // Handle error
}
```

### Listing Available Blend Modes

```go
// Get all registered blend mode names
modes := blendmodes.List()
for _, name := range modes {
    mode, _ := blendmodes.Get(name)
    fmt.Printf("Mode: %s - %s\n", mode.Meta().Name(), mode.Meta().Desc())
}
```

## Blend Mode Categories

### Basic Blend Modes
- **Normal**: Standard alpha compositing
- **Erase**: Removes the bottom layer based on top layer's alpha

### Darken Blend Modes
- **Darken**: Selects the darker of the blend and base colors
- **Multiply**: Multiplies the base and blend colors
- **Color Burn**: Darkens the base color by increasing contrast with the blend color
- **Linear Burn**: Darkens the base color by decreasing brightness
- **Darker Color**: Compares total of all channel values and displays darker color

### Lighten Blend Modes
- **Lighten**: Selects the lighter of the blend and base colors
- **Screen**: Multiplies the inverse of the blend and base colors
- **Color Dodge**: Brightens the base color by decreasing contrast with the blend color
- **Linear Dodge (Add)**: Brightens the base color by increasing brightness
- **Lighter Color**: Compares total of all channel values and displays lighter color

### Contrast Blend Modes
- **Overlay**: Multiplies or screens colors depending on the base color
- **Soft Light**: Darkens or lightens colors depending on the blend color
- **Hard Light**: Multiplies or screens colors depending on the blend color
- **Vivid Light**: Combination of Color Dodge and Color Burn
- **Linear Light**: Combination of Linear Dodge and Linear Burn
- **Pin Light**: Replaces colors depending on blend color
- **Hard Mix**: Creates posterization effect

### Comparative Blend Modes
- **Difference**: Subtracts the darker from the lighter color
- **Exclusion**: Similar to Difference but with lower contrast
- **Subtract**: Subtracts blend color from base color
- **Divide**: Divides base color by blend color
- **Negation**: Similar to Difference but with softer contrast
- **Contrast Negate**: Inverts colors based on blend color contrast

### Component Blend Modes
- **Hue**: Uses hue of blend color with saturation and luminosity of base
- **Saturation**: Uses saturation of blend color with hue and luminosity of base
- **Color**: Uses hue and saturation of blend color with luminosity of base
- **Luminosity**: Uses luminosity of blend color with hue and saturation of base

### Special Effects
- **Reflect**: Creates a glowing effect based on blend color
- **Glow**: Creates a glowing effect based on base color
- **Average**: Averages the base and blend colors

## Implementation Notes

- All blend modes operate on RGBA64 color values for maximum precision
- Alpha values are handled consistently across all blend modes
- Blend operations are performed in linear color space
- Results are automatically clamped to valid color ranges
- Thread-safe operations through synchronized registry access
- Efficient implementation with minimal memory allocations

## Error Handling

All functions that can fail return appropriate errors:

```go
// Example of error handling
mode, err := blendmodes.Get("nonexistent")
if err != nil {
    fmt.Printf("Error getting blend mode: %v\n", err)
}

err = blendmodes.Register("multiply", "...", myBlendFunc)
if err != nil {
    fmt.Printf("Error registering blend mode: %v\n", err)
}
```

## Performance Considerations

- Blend operations are optimized for performance
- Registry operations use read-write mutex for thread safety
- Blend functions are designed to minimize allocations
- Results are cached where appropriate
- Parallel processing friendly

## Contributing

When adding new blend modes:

1. Implement the blend function following the standard signature
2. Register the blend mode with a clear, descriptive name
3. Provide comprehensive documentation
4. Include test cases
5. Ensure thread safety
6. Follow existing naming conventions

## References

- Adobe Photoshop blend mode specifications
- W3C Compositing and Blending Level 1 specification
- Porter-Duff compositing operations 