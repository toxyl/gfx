package projections

// Projection defines the interface for all projection types
type Projection interface {
	// To converts from geographic coordinates (latitude, longitude) to cartesian coordinates (x, y)
	To(latitude, longitude, width, height float64) (x, y float64)
	// From converts from cartesian coordinates (x, y) to geographic coordinates (latitude, longitude)
	From(x, y, width, height float64) (latitude, longitude float64)
}

// RegistryProjection wraps a Projection for registration in the registry
type RegistryProjection struct {
	Meta *ProjectionMeta
	To   func(latitude, longitude, w, h float64) (x, y float64)
	From func(x, y, w, h float64) (latitude, longitude float64)
}

// Name returns the name of the projection
func (p *RegistryProjection) Name() string {
	return p.Meta.Name()
}

// Convert converts coordinates from one projection to another
func Convert(x, y, w, h float64, src, dst *RegistryProjection) (float64, float64) {
	lat, lon := src.From(x, y, w, h)
	return dst.To(lat, lon, w, h)
}

// ToXY implements the registry.ToXY interface
func (p *RegistryProjection) ToXY(a, b, w, h float64) (x, y float64) {
	return p.To(a, b, w, h)
}

// FromXY implements the registry.FromXY interface
func (p *RegistryProjection) FromXY(x, y, w, h float64) (a, b float64) {
	return p.From(x, y, w, h)
}
