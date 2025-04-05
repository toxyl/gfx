package models

import (
	"image"
	"image/color"
)

// PolygonPoint represents a 2D point with normalized coordinates (0-1) for polygon cropping
type PolygonPoint struct {
	X float64
	Y float64
}

// CropPolygon represents a polygon crop effect.
type CropPolygon struct {
	Points []PolygonPoint
}

// Apply applies the polygon crop effect to an image.
func (c *CropPolygon) Apply(img image.Image) (image.Image, error) {
	bounds := img.Bounds()
	width := bounds.Max.X - bounds.Min.X
	height := bounds.Max.Y - bounds.Min.Y

	// Convert points to pixel coordinates
	points := make([]image.Point, len(c.Points))
	for i, p := range c.Points {
		points[i] = image.Point{
			X: int(float64(width) * p.X),
			Y: int(float64(height) * p.Y),
		}
	}

	// Create new image with same dimensions
	dst := image.NewRGBA(bounds)

	for y := bounds.Min.Y; y < bounds.Max.Y; y++ {
		for x := bounds.Min.X; x < bounds.Max.X; x++ {
			// Check if point is inside polygon
			if isPointInPolygon(image.Point{X: x, Y: y}, points) {
				dst.Set(x, y, img.At(x, y))
			} else {
				dst.Set(x, y, color.RGBA64{})
			}
		}
	}

	return dst, nil
}

// isPointInPolygon checks if a point is inside a polygon using ray casting algorithm
func isPointInPolygon(point image.Point, polygon []image.Point) bool {
	if len(polygon) < 3 {
		return false
	}

	inside := false
	for i, j := 0, len(polygon)-1; i < len(polygon); i++ {
		if (polygon[i].Y > point.Y) != (polygon[j].Y > point.Y) {
			slope := (point.X-polygon[i].X)*(polygon[j].Y-polygon[i].Y) -
				(polygon[j].X-polygon[i].X)*(point.Y-polygon[i].Y)
			if slope == 0 {
				return true
			}
			if (slope < 0) != (polygon[j].Y < polygon[i].Y) {
				inside = !inside
			}
		}
		j = i
	}

	return inside
}
