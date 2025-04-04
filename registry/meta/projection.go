package meta

// Coordinate represents a spatial coordinate with metadata
type Coordinate struct {
	Name        string  `json:"name"`
	Description string  `json:"description"`
	Unit        string  `json:"unit"`
	Max         float64 `json:"max"`
	Min         float64 `json:"min"`
}

// NewCoordinate creates a new coordinate with the given parameters
func NewCoordinate(name, description, unit string, max, min float64) *Coordinate {
	return &Coordinate{
		Name:        name,
		Description: description,
		Unit:        unit,
		Max:         max,
		Min:         min,
	}
}

// Projection represents metadata for a map projection
type Projection struct {
	Name        string      `json:"name"`
	Description string      `json:"description"`
	FromCoord   *Coordinate `json:"from_coord"`
	ToCoord     *Coordinate `json:"to_coord"`
}

// NewProjection creates a new projection metadata object
func NewProjection(name, description string, fromCoord, toCoord *Coordinate) *Projection {
	return &Projection{
		Name:        name,
		Description: description,
		FromCoord:   fromCoord,
		ToCoord:     toCoord,
	}
}
