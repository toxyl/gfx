package projections

// Equirectangular implements the equirectangular projection
type Equirectangular struct {
	BaseProjection
}

// NewEquirectangular creates a new equirectangular projection
func NewEquirectangular() *Equirectangular {
	return &Equirectangular{}
}

// To converts from geographic coordinates to cartesian coordinates
func (p *Equirectangular) To(latitude, longitude, width, height float64) (x, y float64) {
	// Clamp input coordinates
	latitude = p.clampLatitude(latitude)
	longitude = p.clampLongitude(longitude)

	// Convert to cartesian coordinates
	x = ((longitude + 180) / 360.0) * width
	y = ((90 - latitude) / 180.0) * height

	// Handle edge cases
	if latitude == 90 {
		y = 0
	} else if latitude == -90 {
		y = height - 1
	}
	if longitude == 180 {
		x = width - 1
	} else if longitude == -180 {
		x = 0
	}

	return x, y
}

// From converts from cartesian coordinates to geographic coordinates
func (p *Equirectangular) From(x, y, width, height float64) (latitude, longitude float64) {
	// Handle edge cases
	if y == 0 {
		latitude = 90
	} else if y == height-1 {
		latitude = -90
	} else {
		latitude = 90.0 - (y/height)*180.0
	}

	if x == 0 {
		longitude = -180
	} else if x == width-1 {
		longitude = 180
	} else {
		longitude = (x/width)*360.0 - 180.0
	}

	return latitude, longitude
}
