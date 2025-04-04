# Image Package

A modern, type-safe image processing and manipulation library that integrates with the `core/color` and `core/blendmodes` packages.

## Features

- **Thread-safe operations**: All image manipulation functions are thread-safe through mutex locks
- **High-precision color handling**: Uses 64-bit float precision for color channels
- **Composable transformations**: Chain operations to create complex image manipulations
- **Full alpha channel support**: Proper handling of transparency across all operations
- **Extensive drawing capabilities**: Draw lines, rectangles, circles, and other shapes
- **Blending operations**: Apply blend modes from the `core/blendmodes` package
- **Parallel processing**: Leverages parallel processing for performance-critical operations

## Basic Usage

### Creating Images

```go
// Create a new blank image
img, err := image.New(800, 600)
if err != nil {
    // Handle error
}

// Load an image from a file
img, err := image.LoadImage("path/to/image.jpg")
if err != nil {
    // Handle error
}
```

### Manipulating Images

```go
// Resize the image using bilinear interpolation
resized, err := img.Resize(400, 300, image.ResizeBilinear)
if err != nil {
    // Handle error
}

// Crop a region from the image
cropped, err := img.Crop(100, 100, 200, 200)
if err != nil {
    // Handle error
}

// Rotate the image by 45 degrees
rotated, err := img.Rotate(45)
if err != nil {
    // Handle error
}

// Flip the image horizontally
flipped, err := img.FlipHorizontal()
if err != nil {
    // Handle error
}

// Translate the image
translated, err := img.Translate(10, 20)
if err != nil {
    // Handle error
}
```

### Drawing Operations

```go
// Create line style
lineStyle := image.NewLineStyle()
lineStyle.Width = 5
lineStyle.Color = &color.RGBA64{R: 1.0, G: 0.0, B: 0.0, A: 1.0} // Red

// Create fill style
fillStyle := image.NewFillStyle()
fillStyle.Color = &color.RGBA64{R: 0.0, G: 0.0, B: 1.0, A: 0.5} // Semi-transparent blue

// Get a blend mode
multiply, err := blendmodes.Get("multiply")
if err != nil {
    // Handle error
}
lineStyle.BlendMode = multiply

// Draw a line
err = img.DrawLine(10, 10, 100, 100, lineStyle)
if err != nil {
    // Handle error
}

// Draw a rectangle
err = img.DrawRectangle(150, 150, 200, 100, lineStyle, fillStyle)
if err != nil {
    // Handle error
}

// Draw a circle
err = img.DrawCircle(300, 300, 50, lineStyle, fillStyle)
if err != nil {
    // Handle error
}

// Draw another image onto this one
err = img.DrawImage(otherImg, 400, 200, multiply, 0.8)
if err != nil {
    // Handle error
}
```

### Processing Images

```go
// Process each pixel with a custom function
err = img.Process(func(x, y int, c *color.RGBA64) (*color.RGBA64, error) {
    // Invert the color
    return &color.RGBA64{
        R: 1.0 - c.R,
        G: 1.0 - c.G,
        B: 1.0 - c.B,
        A: c.A,
    }, nil
})
if err != nil {
    // Handle error
}

// Process in parallel for better performance
err = img.ProcessParallel(func(x, y int, c *color.RGBA64) (*color.RGBA64, error) {
    // Grayscale conversion
    gray := 0.299*c.R + 0.587*c.G + 0.114*c.B
    return &color.RGBA64{
        R: gray,
        G: gray,
        B: gray,
        A: c.A,
    }, nil
})
if err != nil {
    // Handle error
}
```

### Saving Images

```go
// Save with default options (PNG)
err = img.Save("output.png", nil)
if err != nil {
    // Handle error
}

// Save as JPEG with 90% quality
err = img.Save("output.jpg", &image.SaveOptions{
    Format: image.FormatJPEG,
    Quality: 90,
})
if err != nil {
    // Handle error
}

// Convenience methods
err = img.SavePNG("output.png")
if err != nil {
    // Handle error
}

err = img.SaveJPEG("output.jpg", 90)
if err != nil {
    // Handle error
}
```

## Design Considerations

1. **Thread Safety**: All methods acquire appropriate locks to ensure thread safety
2. **Error Handling**: Comprehensive error messages and proper error propagation
3. **Memory Efficiency**: Image operations create new instances to maintain immutability
4. **Performance**: Critical operations use parallelization for better performance
5. **Type Safety**: Strong typing throughout the API to prevent errors

## Implementation Notes

- Image data is represented using the standard library's `image.RGBA` type internally
- All color operations use the `core/color` package's `RGBA64` type for high precision
- Blend operations use the `core/blendmodes` package's blend modes
- Transformations implement industry-standard algorithms (Bresenham, bilinear interpolation, etc.)
- Image metadata is preserved across operations when appropriate 