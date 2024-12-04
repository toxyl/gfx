package convert

import (
	"image/color"

	"github.com/toxyl/gfx/color/hsla"
	"github.com/toxyl/gfx/color/rgba"
	"github.com/toxyl/gfx/math"
)

// RGBAToHSLA convert a RGBA color to a HSLA color.
func RGBAToHSLA(col *rgba.RGBA) *hsla.HSLA {
	r := float64(col.R()) / 255.0
	g := float64(col.G()) / 255.0
	b := float64(col.B()) / 255.0
	a := float64(col.A()) / 255.0

	max := math.Max(r, math.Max(g, b))
	min := math.Min(r, math.Min(g, b))
	l := (max + min) / 2

	var h, s float64
	delta := max - min

	if delta == 0 {
		h = 0
		s = 0
	} else {
		if max == r {
			h = math.Mod((g-b)/delta, 6)
		} else if max == g {
			h = (b-r)/delta + 2
		} else {
			h = (r-g)/delta + 4
		}

		h *= 60
		if h < 0 {
			h += 360
		}

		s = delta / (1 - math.Abs(2*l-1))
	}

	return hsla.New(h, s, l, a)
}

// HSLAToRGBA converts a HSLA color to a RGBA color.
func HSLAToRGBA(col *hsla.HSLA) *rgba.RGBA {
	// Normalize the hue value to the range [0, 360)
	h := col.H() // math.Wrap(, 0, 360)
	s := col.S()
	l := col.L()

	var r, g, b, a float64

	c := (1 - math.Abs(2*l-1)) * s
	x := c * (1 - math.Abs(math.Mod(h/60, 2)-1))
	m := l - c/2

	switch {
	case 0 <= h && h < 60:
		r, g, b = c, x, 0
	case 60 <= h && h < 120:
		r, g, b = x, c, 0
	case 120 <= h && h < 180:
		r, g, b = 0, c, x
	case 180 <= h && h < 240:
		r, g, b = 0, x, c
	case 240 <= h && h < 300:
		r, g, b = x, 0, c
	case 300 <= h && h < 360:
		r, g, b = c, 0, x
	}

	r = math.Round((r + m) * 0xFF)
	g = math.Round((g + m) * 0xFF)
	b = math.Round((b + m) * 0xFF)
	a = math.Round(col.A() * 0xFF)

	return rgba.New(r, g, b, a)
}

// RGBAToRGBAPremul converts a RGBA color to a premultiplied RGBA color.
func RGBAToRGBAPremul(col *rgba.RGBA) color.RGBA {
	res := color.RGBA{
		R: col.R(),
		G: col.G(),
		B: col.B(),
		A: 0xFF,
	}
	if res.A == 0 { // Avoid division by zero if original alpha is zero
		return color.RGBA{0, 0, 0, 0}
	}
	scale := float32(col.A()) / float32(res.A)
	return color.RGBA{
		uint8(float32(res.R) * scale),
		uint8(float32(res.G) * scale),
		uint8(float32(res.B) * scale),
		col.A(),
	}
}

// RGBAPremulToRGBA converts a premultiplied RGBA color to a RGBA color.
func RGBAPremulToRGBA(col color.RGBA) *rgba.RGBA {
	if col.A == 0 {
		// If alpha is zero, return fully transparent black (0, 0, 0, 0).
		return rgba.New(0, 0, 0, 0)
	}
	scale := 255.0 / float32(col.A)
	return rgba.New(
		float32(col.R)*scale,
		float32(col.G)*scale,
		float32(col.B)*scale,
		float32(col.A),
	)
}

// RGBAPremulToHSLA converts a premultiplied RGBA color to a HSLA color.
func RGBAPremulToHSLA(col color.RGBA) *hsla.HSLA {
	cc := RGBAPremulToRGBA(col)
	rf := float64(cc.R()) / 255.0
	gf := float64(cc.G()) / 255.0
	bf := float64(cc.B()) / 255.0
	af := float64(cc.A()) / 255.0

	max := math.Max(rf, math.Max(gf, bf))
	min := math.Min(rf, math.Min(gf, bf))
	delta := max - min

	h := 0.0
	s := 0.0
	l := (max + min) / 2

	if delta != 0 {
		s = delta / (1 - math.Abs(2*l-1))
	}

	switch {
	case delta == 0:
		h = 0
	case max == rf:
		h = math.Mod((gf-bf)/delta, 6)
	case max == gf:
		h = (bf-rf)/delta + 2
	case max == bf:
		h = (rf-gf)/delta + 4
	}
	h *= 60
	return hsla.New(h, s, l, af)
}

// HSLAToRGBAPremul convert a HSLA color to a premultiplied RGBA color.
func HSLAToRGBAPremul(col *hsla.HSLA) color.RGBA {
	return HSLAToRGBA(col).RGBA()
}
