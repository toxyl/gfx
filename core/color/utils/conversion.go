package utils

import "github.com/toxyl/math"

// RGBToHSL converts RGB values to HSL color space
func RGBToHSL(r, g, b float64) (h, s, l float64) {
	max := math.Max(r, math.Max(g, b))
	min := math.Min(r, math.Min(g, b))

	// Calculate lightness
	l = (max + min) / 2

	// If max and min are equal, it's a shade of gray
	if max == min {
		return 0, 0, l
	}

	// Calculate saturation
	if l > 0.5 {
		s = (max - min) / (2 - max - min)
	} else {
		s = (max - min) / (max + min)
	}

	// Calculate hue in degrees
	switch max {
	case r:
		h = (g - b) / (max - min)
		if g < b {
			h += 6
		}
	case g:
		h = (b-r)/(max-min) + 2
	case b:
		h = (r-g)/(max-min) + 4
	}
	h *= 60 // Convert to degrees
	if h < 0 {
		h += 360
	}
	h = math.Mod(h, 360) // Ensure hue is in [0,360] range

	return h, s, l
}

// HSLToRGB converts HSL values to RGB color space
func HSLToRGB(h, s, l float64) (r, g, b float64) {
	if s == 0 {
		return l, l, l
	}

	// Convert hue to normalized value
	h = h / 360.0

	var q float64
	if l < 0.5 {
		q = l * (1 + s)
	} else {
		q = l + s - l*s
	}
	p := 2*l - q

	r = HueToRGB(p, q, h+1.0/3.0)
	g = HueToRGB(p, q, h)
	b = HueToRGB(p, q, h-1.0/3.0)

	return r, g, b
}

// HueToRGB is a helper function for HSLToRGB
func HueToRGB(p, q, t float64) float64 {
	if t < 0 {
		t += 1
	}
	if t > 1 {
		t -= 1
	}
	if t < 1.0/6.0 {
		return p + (q-p)*6*t
	}
	if t < 1.0/2.0 {
		return q
	}
	if t < 2.0/3.0 {
		return p + (q-p)*(2.0/3.0-t)*6
	}
	return p
}

// RGBToHSB converts RGB values to HSB color space
func RGBToHSB(r, g, b float64) (h, s, v float64) {
	max := math.Max(r, math.Max(g, b))
	min := math.Min(r, math.Min(g, b))

	// Calculate value (brightness)
	v = max

	// If max is 0, it's black
	if max == 0 {
		return 0, 0, 0
	}

	// Calculate saturation
	s = (max - min) / max

	// If max and min are equal, it's a shade of gray
	if max == min {
		return 0, 0, v
	}

	// Calculate hue in degrees
	switch max {
	case r:
		h = (g - b) / (max - min)
		if g < b {
			h += 6
		}
	case g:
		h = (b-r)/(max-min) + 2
	case b:
		h = (r-g)/(max-min) + 4
	}
	h *= 60 // Convert to degrees

	return h, s, v
}

// HSBToRGB converts HSB values to RGB color space
func HSBToRGB(h, s, v float64) (r, g, b float64) {
	if s == 0 {
		return v, v, v
	}

	// Convert hue to normalized value
	h = h / 360.0

	i := math.Floor(h * 6)
	f := h*6 - i
	p := v * (1 - s)
	q := v * (1 - s*f)
	t := v * (1 - s*(1-f))

	switch int(i) % 6 {
	case 0:
		r, g, b = v, t, p
	case 1:
		r, g, b = q, v, p
	case 2:
		r, g, b = p, v, t
	case 3:
		r, g, b = p, q, v
	case 4:
		r, g, b = t, p, v
	case 5:
		r, g, b = v, p, q
	}

	return r, g, b
}
