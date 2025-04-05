# Projections Package

The projections package provides utilities for working with different coordinate systems and projections.

## Structure

- Coordinate system conversions
- Projection utilities
- Transformation functions

## Usage

```go
// Example: Converting between coordinate systems
point := projections.NewPoint(1.0, 2.0)
polar := point.ToPolar()
cartesian := polar.ToCartesian()
```

## Projections

- Cartesian
- Polar
- Spherical
- Cylindrical
- Mercator
- Stereographic
- Orthographic

## Implementation Details

Each projection:
1. Provides conversion methods to other coordinate systems
2. Implements proper transformation formulas
3. Handles coordinate system limitations
4. Preserves accuracy during conversions

## Testing

Run tests with:
```bash
go test ./...
``` 