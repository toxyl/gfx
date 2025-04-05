# Layers Package

The layers package provides a flexible system for composing images from multiple layers with blend modes, alpha channels, and stacked effects.

## Structure

- Layer management
- Composition handling
- Effect stacking
- Blend mode integration

## Usage

```go
// Example: Creating a composition with multiple layers
w, h := 640, 480
col := color.NewRGB(0, 0, 0, 1)
fxs := []fx.Effect{ fx.NewBlur(4) }
fxsl := []fx.Effect{ fx.NewSharpen(4) }
blend := blendmodes.NewNormal()
alpha := 0.75
outfile := "/tmp/test.png"

c := composition.New(w, h, col, fxs, layer.New(w, h, blend, alpha, fxsl...))
c.Render(outfile)
```

## Features

- Layer composition:
  - Multiple layers
  - Blend modes
  - Alpha channels
  - Effect stacking
- Background options:
  - Solid color
  - Transparency
- Layer operations:
  - Add/remove layers
  - Reorder layers
  - Modify layer properties
- Effect management:
  - Stack multiple effects
  - Apply effects to individual layers
  - Apply effects to composition

## Implementation Details

The package:
1. Provides a flexible layer system
2. Integrates with blend modes and effects
3. Handles alpha compositing
4. Supports effect stacking
5. Manages layer ordering

## Testing

Run tests with:
```bash
go test ./...
``` 