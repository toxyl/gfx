package coordinates

import "math"

func LatLonToCartesian(latitude, longitude float64, w int, h int) (x int, y int) {
	x = int(((longitude + 180) / 360.0) * float64(w))
	y = int(((90 - latitude) / 180.0) * float64(h))

	if x < 0 {
		x = 0
	} else if x >= w {
		x = w - 1
	}

	if y < 0 {
		y = 0
	} else if y >= h {
		y = h - 1
	}

	return
}

func CartesianToLatLon(x, y int, w int, h int) (lat float64, lon float64) {
	if x < 0 {
		x = 0
	} else if x >= w {
		x = w - 1
	}

	if y < 0 {
		y = 0
	} else if y >= h {
		y = h - 1
	}

	lat = 90.0 - (float64(y)/float64(h))*180.0
	lon = (float64(x)/float64(w))*360.0 - 180.0

	return
}

// PolarToCartesian converts polar coordinates (radius, angle in degrees) to cartesian coordinates (x, y).
func PolarToCartesian(radius, angle float64) (x, y float64) {
	// Wrap angle to be within 0 to 360 degrees
	angle = math.Mod(angle, 360)

	// Convert degrees to radians
	angleRad := angle * (math.Pi / 180)

	x = radius * math.Cos(angleRad)
	y = radius * math.Sin(angleRad)

	return x, y
}
