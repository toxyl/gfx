package projections

import (
	"github.com/toxyl/math"
)

// BaseProjection provides common functionality for all projections
type BaseProjection struct {
	// Add any common fields here if needed
}

// clampCoordinates ensures coordinates are within valid ranges
func (p *BaseProjection) clampCoordinates(x, y, width, height float64) (float64, float64) {
	return math.Clamp(x, 0, width-1), math.Clamp(y, 0, height-1)
}

// clampLatitude ensures latitude is within [-90, 90]
func (p *BaseProjection) clampLatitude(lat float64) float64 {
	return math.Clamp(lat, -90, 90)
}

// clampLongitude ensures longitude is within [-180, 180]
func (p *BaseProjection) clampLongitude(lon float64) float64 {
	return math.Clamp(lon, -180, 180)
}
