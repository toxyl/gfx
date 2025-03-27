package image

import (
	"image"
	"math"
	"sync"
)

// ToPolar converts a rectangular image into a polar image using the provided angle range,
// then rotates the result by the given rotation (in degrees). The final image dimensions are preserved.
func (i *Image) ToPolar(angleStart, angleEnd, rotation float64) *Image {
	i.Lock()
	defer i.Unlock()

	w := i.raw.Bounds().Dx()
	h := i.raw.Bounds().Dy()
	dst := image.NewRGBA(image.Rect(0, 0, w, h))
	cx := float64(w) / 2.0
	cy := float64(h) / 2.0
	maxR := math.Min(cx, cy)

	// Precompute the rotation parameters.
	angleRad := rotation * math.Pi / 180.0
	cosA := math.Cos(angleRad)
	sinA := math.Sin(angleRad)

	// For each pixel in the final output image:
	for y := range h {
		for x := range w {
			// First, compute the corresponding coordinate in the polar (pre-rotated) image
			// by applying the inverse rotation.
			xt := float64(x) - cx
			yt := float64(y) - cy
			// Inverse rotation transformation.
			px := cosA*xt + sinA*yt + cx
			py := -sinA*xt + cosA*yt + cy

			// Compute vector from the center.
			dx := px - cx
			dy := py - cy
			r := math.Sqrt(dx*dx + dy*dy)

			// Compute angle (in degrees) relative to the center and adjust by angleStart.
			theta := (math.Atan2(dy, dx) * 180.0 / math.Pi) - angleStart
			if theta < 0 {
				theta += 360
			}

			var proportion float64
			if angleEnd >= angleStart {
				// If outside the allowed angle range, leave the pixel unchanged (transparent).
				if theta < angleStart || theta > angleEnd {
					dst.Set(x, y, image.Transparent)
					continue
				}
				proportion = (theta - angleStart) / (angleEnd - angleStart)
			} else {
				// Handle wrapped range, e.g. (90, -90) or (270, 90)
				if theta < angleStart && theta > angleEnd {
					dst.Set(x, y, image.Transparent)
					continue
				}
				totalRange := (360 - angleStart) + angleEnd
				if theta >= angleStart {
					proportion = (theta - angleStart) / totalRange
				} else {
					proportion = (theta + (360 - angleStart)) / totalRange
				}
			}

			// Map radial distance to vertical coordinate in the source image.
			srcY := int(r / maxR * float64(h-1))
			srcX := max(w-int(proportion*float64(w-1)), 0)
			if srcX >= w {
				srcX = w - 1
			}
			if srcY < 0 {
				srcY = 0
			}
			if srcY >= h {
				srcY = h - 1
			}

			// Set the final pixel from the original image.
			dst.Set(x, y, i.raw.At(srcX, srcY))
		}
	}

	return &Image{raw: dst, path: i.path, mu: &sync.Mutex{}}
}
