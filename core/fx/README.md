# FX Package

The fx package provides a comprehensive set of image effects and filters.

## Structure

- Effect models
- Filter implementations
- Effect utilities

## Usage

```go
// Example: Applying an effect
img := image.NewRGBA(...)
effect := fx.NewVibranceEffect(0.5) // 50% vibrance
result := effect.Apply(img)
```

## Effects

- Color adjustments:
  - Brightness
  - Contrast
  - Saturation
  - Vibrance
  - Colorize
  - Tint
- Geometric transformations:
  - Scale
  - Rotate
  - Translate
  - Flip
  - Crop
- Special effects:
  - Blur
  - Sharpen
  - Emboss
  - Enhance
  - Extract
- Color space effects:
  - Grayscale
  - Sepia
  - Invert
  - Threshold

## Implementation Details

Each effect:
1. Implements the Effect interface
2. Uses the meta package for parameter validation
3. Preserves image quality during processing
4. Handles edge cases appropriately

## Testing

Run tests with:
```bash
go test ./...
``` 