package image

import (
	"fmt"
	"math"

	"github.com/toxyl/gfx/core/blendmodes"
	"github.com/toxyl/gfx/core/color"
)

// LineStyle defines the style of a line.
type LineStyle struct {
	// Width of the line in pixels
	Width int

	// Color of the line
	Color *color.RGBA64

	// Blend mode to use when drawing the line
	BlendMode *blendmodes.IBlendMode

	// Alpha value (0-1) for transparency
	Alpha float64
}

// NewLineStyle creates a new line style with default values.
func NewLineStyle() *LineStyle {
	return &LineStyle{
		Width: 1,
		Color: &color.RGBA64{R: 0, G: 0, B: 0, A: 1}, // Black
		Alpha: 1.0,                                   // Fully opaque
	}
}

// DrawLine draws a line from (x1, y1) to (x2, y2) using the specified style.
func (i *Image) DrawLine(x1, y1, x2, y2 int, style *LineStyle) error {
	if style == nil {
		return fmt.Errorf("line style cannot be nil")
	}

	if style.BlendMode == nil {
		// Default to normal blend mode if none specified
		normalMode, err := blendmodes.Get("normal")
		if err != nil {
			return fmt.Errorf("failed to get default blend mode: %w", err)
		}
		style.BlendMode = normalMode
	}

	// Use Bresenham's line algorithm for drawing
	// This is a classic algorithm for drawing lines on a pixel grid
	dx := int(math.Abs(float64(x2 - x1)))
	dy := int(math.Abs(float64(y2 - y1)))

	var sx, sy int
	if x1 < x2 {
		sx = 1
	} else {
		sx = -1
	}

	if y1 < y2 {
		sy = 1
	} else {
		sy = -1
	}

	err := dx - dy

	// Apply line width
	halfWidth := style.Width / 2

	// Lock the image for writing
	i.mu.Lock()
	defer i.mu.Unlock()

	// Draw the line using Bresenham's algorithm
	x, y := x1, y1
	for {
		// Draw a square of pixels around the current point for line width
		for offsetY := -halfWidth; offsetY <= halfWidth; offsetY++ {
			for offsetX := -halfWidth; offsetX <= halfWidth; offsetX++ {
				px, py := x+offsetX, y+offsetY

				// Check if the pixel is within image bounds
				if px >= 0 && px < i.width && py >= 0 && py < i.height {
					// Get the current pixel color
					r, g, b, a := i.data.RGBA64At(px, py).RGBA()
					dstColor := &color.RGBA64{
						R: float64(r) / 65535.0,
						G: float64(g) / 65535.0,
						B: float64(b) / 65535.0,
						A: float64(a) / 65535.0,
					}

					// Blend with the line color
					result, err := style.BlendMode.Blend(dstColor, style.Color, style.Alpha)
					if err != nil {
						continue // Skip this pixel on error
					}

					// Set the pixel
					i.data.SetRGBA64(px, py, result.To16bit())
				}
			}
		}

		// Check if we've reached the end point
		if x == x2 && y == y2 {
			break
		}

		// Calculate the next point
		e2 := 2 * err
		if e2 > -dy {
			err -= dy
			x += sx
		}
		if e2 < dx {
			err += dx
			y += sy
		}
	}

	return nil
}

// FillStyle defines the style for filling shapes.
type FillStyle struct {
	// Color to fill with
	Color *color.RGBA64

	// Blend mode to use
	BlendMode *blendmodes.IBlendMode

	// Alpha value (0-1) for transparency
	Alpha float64
}

// NewFillStyle creates a new fill style with default values.
func NewFillStyle() *FillStyle {
	return &FillStyle{
		Color: &color.RGBA64{R: 0, G: 0, B: 0, A: 1}, // Black
		Alpha: 1.0,                                   // Fully opaque
	}
}

// DrawRectangle draws a rectangle with optional border and fill.
func (i *Image) DrawRectangle(x, y, width, height int, lineStyle *LineStyle, fillStyle *FillStyle) error {
	// Validate dimensions
	if width <= 0 || height <= 0 {
		return fmt.Errorf("invalid rectangle dimensions: width=%d, height=%d", width, height)
	}

	// Lock the image for writing
	i.mu.Lock()
	defer i.mu.Unlock()

	// Get default blend mode if needed
	var normalMode *blendmodes.IBlendMode
	if (fillStyle != nil && fillStyle.BlendMode == nil) || (lineStyle != nil && lineStyle.BlendMode == nil) {
		var err error
		normalMode, err = blendmodes.Get("normal")
		if err != nil {
			return fmt.Errorf("failed to get default blend mode: %w", err)
		}
	}

	// Fill the rectangle if a fill style is provided
	if fillStyle != nil {
		if fillStyle.BlendMode == nil {
			fillStyle.BlendMode = normalMode
		}

		// Fill the rectangle
		for py := y; py < y+height; py++ {
			for px := x; px < x+width; px++ {
				// Check if the pixel is within image bounds
				if px >= 0 && px < i.width && py >= 0 && py < i.height {
					// Get the current pixel color
					r, g, b, a := i.data.RGBA64At(px, py).RGBA()
					dstColor := &color.RGBA64{
						R: float64(r) / 65535.0,
						G: float64(g) / 65535.0,
						B: float64(b) / 65535.0,
						A: float64(a) / 65535.0,
					}

					// Blend with the fill color
					result, err := fillStyle.BlendMode.Blend(dstColor, fillStyle.Color, fillStyle.Alpha)
					if err != nil {
						continue // Skip this pixel on error
					}

					// Set the pixel
					i.data.SetRGBA64(px, py, result.To16bit())
				}
			}
		}
	}

	// Draw the border if a line style is provided
	if lineStyle != nil {
		if lineStyle.BlendMode == nil {
			lineStyle.BlendMode = normalMode
		}

		// Release the lock temporarily to use the higher-level function
		i.mu.Unlock()

		// Draw the four sides of the rectangle
		err := i.DrawLine(x, y, x+width-1, y, lineStyle)
		if err != nil {
			i.mu.Lock() // Reacquire the lock before returning
			return err
		}

		err = i.DrawLine(x+width-1, y, x+width-1, y+height-1, lineStyle)
		if err != nil {
			i.mu.Lock()
			return err
		}

		err = i.DrawLine(x+width-1, y+height-1, x, y+height-1, lineStyle)
		if err != nil {
			i.mu.Lock()
			return err
		}

		err = i.DrawLine(x, y+height-1, x, y, lineStyle)
		if err != nil {
			i.mu.Lock()
			return err
		}

		// Reacquire the lock
		i.mu.Lock()
	}

	return nil
}

// DrawCircle draws a circle with optional border and fill.
func (i *Image) DrawCircle(centerX, centerY, radius int, lineStyle *LineStyle, fillStyle *FillStyle) error {
	// Validate dimensions
	if radius <= 0 {
		return fmt.Errorf("invalid circle radius: %d", radius)
	}

	// Lock the image for writing
	i.mu.Lock()
	defer i.mu.Unlock()

	// Get default blend mode if needed
	var normalMode *blendmodes.IBlendMode
	if (fillStyle != nil && fillStyle.BlendMode == nil) || (lineStyle != nil && lineStyle.BlendMode == nil) {
		var err error
		normalMode, err = blendmodes.Get("normal")
		if err != nil {
			return fmt.Errorf("failed to get default blend mode: %w", err)
		}
	}

	// Fill the circle if a fill style is provided
	if fillStyle != nil {
		if fillStyle.BlendMode == nil {
			fillStyle.BlendMode = normalMode
		}

		radiusSquared := radius * radius
		// Process all pixels in a square that contains the circle
		for py := centerY - radius; py <= centerY+radius; py++ {
			for px := centerX - radius; px <= centerX+radius; px++ {
				// Calculate the squared distance from the center
				dx := px - centerX
				dy := py - centerY
				distanceSquared := dx*dx + dy*dy

				// Check if the pixel is within the circle
				if distanceSquared <= radiusSquared {
					// Check if the pixel is within image bounds
					if px >= 0 && px < i.width && py >= 0 && py < i.height {
						// Get the current pixel color
						r, g, b, a := i.data.RGBA64At(px, py).RGBA()
						dstColor := &color.RGBA64{
							R: float64(r) / 65535.0,
							G: float64(g) / 65535.0,
							B: float64(b) / 65535.0,
							A: float64(a) / 65535.0,
						}

						// Blend with the fill color
						result, err := fillStyle.BlendMode.Blend(dstColor, fillStyle.Color, fillStyle.Alpha)
						if err != nil {
							continue // Skip this pixel on error
						}

						// Set the pixel
						i.data.SetRGBA64(px, py, result.To16bit())
					}
				}
			}
		}
	}

	// Draw the border if a line style is provided
	if lineStyle != nil {
		if lineStyle.BlendMode == nil {
			lineStyle.BlendMode = normalMode
		}

		// Use Bresenham's circle algorithm
		x, y := radius, 0
		decision := 1 - radius

		// Function to draw eight points in a circle
		drawCirclePoints := func(cx, cy, x, y, width int) {
			drawPoint := func(px, py int) {
				// Only draw if within bounds
				if px >= 0 && px < i.width && py >= 0 && py < i.height {
					// For thicker lines, draw a square around the point
					halfWidth := width / 2
					for oy := -halfWidth; oy <= halfWidth; oy++ {
						for ox := -halfWidth; ox <= halfWidth; ox++ {
							pointX, pointY := px+ox, py+oy

							// Check bounds again for the offset point
							if pointX >= 0 && pointX < i.width && pointY >= 0 && pointY < i.height {
								// Get the current pixel color
								r, g, b, a := i.data.RGBA64At(pointX, pointY).RGBA()
								dstColor := &color.RGBA64{
									R: float64(r) / 65535.0,
									G: float64(g) / 65535.0,
									B: float64(b) / 65535.0,
									A: float64(a) / 65535.0,
								}

								// Blend with the line color
								result, err := lineStyle.BlendMode.Blend(dstColor, lineStyle.Color, lineStyle.Alpha)
								if err != nil {
									return // Skip this pixel on error
								}

								// Set the pixel
								i.data.SetRGBA64(pointX, pointY, result.To16bit())
							}
						}
					}
				}
			}

			// Draw the eight octants
			drawPoint(cx+x, cy+y)
			drawPoint(cx-x, cy+y)
			drawPoint(cx+x, cy-y)
			drawPoint(cx-x, cy-y)
			drawPoint(cx+y, cy+x)
			drawPoint(cx-y, cy+x)
			drawPoint(cx+y, cy-x)
			drawPoint(cx-y, cy-x)
		}

		// Draw the initial points
		drawCirclePoints(centerX, centerY, x, y, lineStyle.Width)

		// Main loop for the Bresenham circle algorithm
		for y < x {
			y++
			if decision <= 0 {
				decision += 2*y + 1
			} else {
				x--
				decision += 2*(y-x) + 1
			}
			drawCirclePoints(centerX, centerY, x, y, lineStyle.Width)
		}
	}

	return nil
}

// DrawImage draws another image onto this image at the specified position.
func (i *Image) DrawImage(other *Image, x, y int, mode *blendmodes.IBlendMode, alpha float64) error {
	if other == nil {
		return fmt.Errorf("source image cannot be nil")
	}

	// Get default blend mode if needed
	if mode == nil {
		var err error
		mode, err = blendmodes.Get("normal")
		if err != nil {
			return fmt.Errorf("failed to get default blend mode: %w", err)
		}
	}

	// Lock both images
	i.mu.Lock()
	defer i.mu.Unlock()

	other.mu.RLock()
	defer other.mu.RUnlock()

	// Get dimensions of both images
	srcW, srcH := other.width, other.height

	// Calculate the intersection of the draw area and the destination image
	startX := max(0, x)
	startY := max(0, y)
	endX := min(i.width, x+srcW)
	endY := min(i.height, y+srcH)

	// Check if there's any overlap
	if startX >= endX || startY >= endY {
		return nil // No overlap, nothing to draw
	}

	// Draw the overlapping region
	for dy := startY; dy < endY; dy++ {
		for dx := startX; dx < endX; dx++ {
			// Calculate the position in the source image
			srcX := dx - x
			srcY := dy - y

			// Get the source pixel color
			r, g, b, a := other.data.RGBA64At(srcX, srcY).RGBA()
			srcColor := &color.RGBA64{
				R: float64(r) / 65535.0,
				G: float64(g) / 65535.0,
				B: float64(b) / 65535.0,
				A: float64(a) / 65535.0,
			}

			// If the source pixel is fully transparent, skip it
			if srcColor.A == 0 {
				continue
			}

			// Get the destination pixel color
			r, g, b, a = i.data.RGBA64At(dx, dy).RGBA()
			dstColor := &color.RGBA64{
				R: float64(r) / 65535.0,
				G: float64(g) / 65535.0,
				B: float64(b) / 65535.0,
				A: float64(a) / 65535.0,
			}

			// Apply the blend mode
			result, err := mode.Blend(dstColor, srcColor, alpha)
			if err != nil {
				continue // Skip this pixel on error
			}

			// Set the pixel
			i.data.SetRGBA64(dx, dy, result.To16bit())
		}
	}

	return nil
}

// max returns the larger of two integers
func max(a, b int) int {
	if a > b {
		return a
	}
	return b
}

// min returns the smaller of two integers
func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}
