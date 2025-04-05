# Blend Modes Package

The blend modes package provides a comprehensive set of image blending operations, implementing various compositing algorithms for combining images.

## Structure

- Individual blend mode implementations
- Common utilities for blend mode operations

## Usage

```go
// Example: Applying a blend mode
base := image.NewRGBA(...)
overlay := image.NewRGBA(...)
result := blendmodes.Normal(base, overlay)
```

## Blend Modes

- Normal
- Multiply
- Screen
- Overlay
- Darken
- Lighten
- ColorDodge
- ColorBurn
- HardLight
- SoftLight
- Difference
- Exclusion
- Hue
- Saturation
- Color
- Luminosity

## Implementation Details

Each blend mode:
1. Takes two images as input (base and overlay)
2. Returns a new image with the blended result
3. Preserves alpha channel when appropriate
4. Handles image bounds properly

## Testing

Run tests with:
```bash
go test ./...
``` 