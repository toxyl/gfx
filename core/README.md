# Core Package

The core package provides the foundation for image processing and manipulation in the GFX library.

## Structure

- `blendmodes/` - Image blending operations
- `color/` - Color space conversions and utilities
- `fx/` - Image effects and filters
- `image/` - Image manipulation utilities
- `meta/` - Metadata and documentation utilities
- `projections/` - Coordinate system projections

## Usage

```go
// Example: Applying an effect with color conversion
img := image.NewRGBA(...)
effect := fx.NewVibranceEffect(0.5)
result := effect.Apply(img)
```

## Features

- Comprehensive image processing capabilities
- High-precision color handling
- Extensive effect library
- Flexible blend modes
- Coordinate system projections
- Metadata and documentation support

## Implementation Details

The package:
1. Uses consistent interfaces across components
2. Implements proper error handling
3. Preserves image quality
4. Supports extensibility

## Testing

Run tests with:
```bash
go test ./...
``` 