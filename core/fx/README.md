# FX Package

The FX package provides a comprehensive set of image processing effects and filters. It is designed to be extensible and maintainable, with a clear separation between different types of effects.

## Structure

- `color/` - Color manipulation effects (brightness, contrast, hue, etc.)
- `effects/` - Special effects (convolution, noise, etc.)
- `geometric/` - Geometric transformations (crop, rotate, scale, etc.)
- `spatial/` - Spatial filters (blur, sharpen, edge detection, etc.)
- `models/` - Effect model definitions and implementations

## Usage

```go
// Example: Applying a brightness effect
img := image.NewRGBA(...)
fx := fx.NewBrightness(0.5) // 50% brightness increase
result := fx.Apply(img)
```

## Effect Types

### Color Effects
- Brightness
- Contrast
- Gamma
- Gray
- Hue
- Invert
- Luminance
- Saturation
- Sepia
- Threshold
- Vibrance

### Geometric Effects
- Crop
- CropCircle
- FlipH
- FlipV
- Rotate
- Scale
- ToPolar
- Transform
- Translate
- TranslateWrap

### Spatial Effects
- Blur
- EdgeDetect
- Emboss
- Enhance
- Sharpen

### Special Effects
- ColorShift
- Convolution
- Extract
- Noise

## Extending

New effects can be added by implementing the `iEffect` interface and registering them with the effect registry.

## Testing

Run tests with:
```bash
go test ./...
``` 