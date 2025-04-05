package models

import (
	"image"
	"image/color"

	"github.com/toxyl/gfx/core/fx"
	"github.com/toxyl/gfx/core/meta"
	"github.com/toxyl/math"
)

// Point represents a control point on the curve.
type Point struct {
	X float64 // Input value (0.0 to 1.0)
	Y float64 // Output value (0.0 to 1.0)
}

// Curves represents a curves adjustment effect.
type Curves struct {
	Points []Point // Control points for the curve
	meta   *fx.EffectMeta
}

// NewCurvesEffect creates a new curves adjustment effect.
func NewCurvesEffect(points []Point) *Curves {
	c := &Curves{
		Points: points,
		meta: fx.NewEffectMeta(
			"Curves",
			"Adjusts the tonal range of an image using a custom curve",
			meta.NewChannelMeta("Points", 0.0, 1.0, "", "Control points for the curve"),
		),
	}

	// Sort points by X value
	for i := 0; i < len(c.Points); i++ {
		for j := i + 1; j < len(c.Points); j++ {
			if c.Points[i].X > c.Points[j].X {
				c.Points[i], c.Points[j] = c.Points[j], c.Points[i]
			}
		}
	}

	// Ensure first point is at (0,0) and last point is at (1,1)
	if len(c.Points) == 0 || c.Points[0].X > 0 {
		c.Points = append([]Point{{X: 0, Y: 0}}, c.Points...)
	}
	if c.Points[len(c.Points)-1].X < 1 {
		c.Points = append(c.Points, Point{X: 1, Y: 1})
	}

	return c
}

// interpolate returns the interpolated Y value for a given X using cubic spline interpolation.
func (c *Curves) interpolate(x float64) float64 {
	if len(c.Points) < 2 {
		return x
	}

	// Find the segment containing x
	var i int
	for i = 0; i < len(c.Points)-1; i++ {
		if x >= c.Points[i].X && x <= c.Points[i+1].X {
			break
		}
	}

	if i >= len(c.Points)-1 {
		return c.Points[len(c.Points)-1].Y
	}

	// Cubic spline interpolation
	t := (x - c.Points[i].X) / (c.Points[i+1].X - c.Points[i].X)
	t2 := t * t
	t3 := t2 * t

	h00 := 2*t3 - 3*t2 + 1
	h10 := t3 - 2*t2 + t
	h01 := -2*t3 + 3*t2
	h11 := t3 - t2

	// Calculate tangents
	var m0, m1 float64
	if i == 0 {
		m0 = (c.Points[i+1].Y - c.Points[i].Y) / (c.Points[i+1].X - c.Points[i].X)
	} else {
		m0 = (c.Points[i+1].Y - c.Points[i-1].Y) / (c.Points[i+1].X - c.Points[i-1].X)
	}

	if i == len(c.Points)-2 {
		m1 = (c.Points[i+1].Y - c.Points[i].Y) / (c.Points[i+1].X - c.Points[i].X)
	} else {
		m1 = (c.Points[i+2].Y - c.Points[i].Y) / (c.Points[i+2].X - c.Points[i].X)
	}

	return h00*c.Points[i].Y + h10*(c.Points[i+1].X-c.Points[i].X)*m0 +
		h01*c.Points[i+1].Y + h11*(c.Points[i+1].X-c.Points[i].X)*m1
}

// Apply applies the curves adjustment effect to an image.
func (c *Curves) Apply(img image.Image) image.Image {
	bounds := img.Bounds()
	dst := image.NewRGBA(bounds)

	for y := bounds.Min.Y; y < bounds.Max.Y; y++ {
		for x := bounds.Min.X; x < bounds.Max.X; x++ {
			r, g, b, a := img.At(x, y).RGBA()

			// Convert to float64 for calculations
			rf := float64(r) / 0xFFFF
			gf := float64(g) / 0xFFFF
			bf := float64(b) / 0xFFFF

			// Apply curve adjustment
			rf = c.interpolate(rf)
			gf = c.interpolate(gf)
			bf = c.interpolate(bf)

			// Convert back to uint32
			r = uint32(math.Max(0, math.Min(0xFFFF, rf*0xFFFF)))
			g = uint32(math.Max(0, math.Min(0xFFFF, gf*0xFFFF)))
			b = uint32(math.Max(0, math.Min(0xFFFF, bf*0xFFFF)))

			dst.Set(x, y, color.RGBA64{
				R: uint16(r),
				G: uint16(g),
				B: uint16(b),
				A: uint16(a),
			})
		}
	}

	return dst
}

// Meta returns the effect metadata.
func (c *Curves) Meta() *fx.EffectMeta {
	return c.meta
}
