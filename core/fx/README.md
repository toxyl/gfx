# FX Package

The FX package provides a comprehensive set of image processing effects and filters. It is designed to be extensible and maintainable, with a consistent model-based architecture.

## Structure

- `models/` - Effect model definitions and implementations
- `utils/` - Utility functions and helpers
- `meta.go` - Metadata and parameter validation
- `registry.go` - Effect registration and management

## Usage

```go
// Example: Applying a brightness effect
img := image.NewRGBA(...)
effect := models.NewBrightnessEffect(0.5) // 50% brightness increase
result := effect.Apply(img)
```

## Effect Types

### Color Effects
- Brightness
- Color Balance
- Color Temperature
- Colorize
- Contrast
- Curves
- Gamma
- Grayscale
- Hue
- Invert
- Levels
- Luminance
- Luminance Contrast
- Pastelize
- Saturation
- Saturation Contrast
- Sepia
- Threshold
- Tint
- Vibrance

### Geometric Effects
- Crop (Rectangular)
- Crop (Circular)
- Flip (Horizontal)
- Flip (Vertical)
- Rotate
- Scale
- ToPolar/FromPolar
- Transform
- Translate

### Spatial Effects
- Convolution
- Emboss
- Enhance
- Extract

## Extending

New effects can be added by implementing the `Effect` interface in the models package:

```go
type Effect interface {
    Apply(image.Image) image.Image
    Meta() *EffectMeta
}
```

Each effect should:
1. Define its parameters using the meta package
2. Implement parameter validation
3. Include comprehensive documentation
4. Handle image bounds properly
5. Preserve alpha channel when appropriate

## Testing

Run tests with:
```bash
go test ./...
``` 