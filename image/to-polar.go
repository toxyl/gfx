package image

import (
	"image"
	"math"
	"sync"
)

// ToPolar converts a rectangular image into a polar image using the provided angle range,
// applies a configurable fisheye effect, then rotates the result by the given rotation (in degrees).
// The final image dimensions are preserved.
func (i *Image) ToPolar(angleStart, angleEnd, rotation, fisheye float64) *Image {
	i.Lock()
	defer i.Unlock()

	w := i.raw.Bounds().Dx()
	h := i.raw.Bounds().Dy()
	dst := image.NewRGBA(image.Rect(0, 0, w, h))
	cx := float64(w) / 2.0
	cy := float64(h) / 2.0
	maxR := math.Min(cx, cy)

	// Precompute rotation parameters.
	angleRad := rotation * math.Pi / 180.0
	cosA := math.Cos(angleRad)
	sinA := math.Sin(angleRad)

	// For each pixel in the final output image:
	for y := 0; y < h; y++ {
		for x := 0; x < w; x++ {
			// Apply inverse rotation to get the corresponding coordinate in the pre-rotated polar image.
			xt := float64(x) - cx
			yt := float64(y) - cy
			px := cosA*xt + sinA*yt + cx
			py := -sinA*xt + cosA*yt + cy

			// Compute vector from the center.
			dx := px - cx
			dy := py - cy
			r := math.Sqrt(dx*dx + dy*dy)

			// Apply fisheye effect by remapping the normalized radius.
			// When fisheye == 0, no change occurs. For fisheye > 0 the mapping produces a barrel distortion.
			norm := r / maxR
			if fisheye != 0 {
				// Adjust the exponent based on fisheye strength.
				// A higher fisheye value compresses distances further out.
				norm = math.Pow(norm, 1.0/(1.0+fisheye))
				r = norm * maxR
			}

			// Compute angle (in degrees) relative to the center and adjust by angleStart.
			theta := (math.Atan2(dy, dx) * 180.0 / math.Pi) - angleStart
			if theta < 0 {
				theta += 360
			}

			// Determine angular proportion.
			var proportion float64
			if angleEnd >= angleStart {
				// If outside the allowed angle range, set pixel transparent.
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

			// Map the radial distance to the vertical coordinate in the source image.
			srcY := int(r / maxR * float64(h-1))
			// Map the angular proportion to the horizontal coordinate in the source image.
			srcX := max(int(proportion*float64(w-1)), 0)
			if srcX >= w {
				srcX = w - 1
			}
			if srcY < 0 {
				srcY = 0
			}
			if srcY >= h {
				srcY = h - 1
			}

			dst.Set(x, y, i.raw.At(srcX, srcY))
		}
	}

	return &Image{raw: dst, path: i.path, mu: &sync.Mutex{}}
}

// A helper function to mimic math.Max for ints.
func max(a, b int) int {
	if a > b {
		return a
	}
	return b
}
