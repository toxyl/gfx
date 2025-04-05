package projections

import (
	"fmt"
	"strings"
)

type (
	// Coordinate represents a spatial coordinate with metadata
	Coordinate struct {
		name string
		desc string
		unit string
		min  float64
		max  float64
	}

	// Coordinates represents a pair of coordinates (A and B)
	Coordinates struct {
		A *Coordinate
		B *Coordinate
	}

	// ProjectionMeta represents metadata for a map projection
	ProjectionMeta struct {
		name        string
		desc        string
		coordinates *Coordinates
	}
)

// Name returns the name of the coordinate
func (c *Coordinate) Name() string { return c.name }

// Description returns the description of the coordinate
func (c *Coordinate) Description() string { return c.desc }

// Desc is an alias for Description for backward compatibility
func (c *Coordinate) Desc() string { return c.Description() }

// Unit returns the unit of measurement for the coordinate
func (c *Coordinate) Unit() string { return c.unit }

// Min returns the minimum value of the coordinate
func (c *Coordinate) Min() float64 { return c.min }

// Max returns the maximum value of the coordinate
func (c *Coordinate) Max() float64 { return c.max }

// NewCoordinateMeta creates a new coordinate with the given parameters
func NewCoordinateMeta(name, description, unit string, max, min float64) *Coordinate {
	return &Coordinate{
		name: name,
		desc: description,
		unit: unit,
		min:  min,
		max:  max,
	}
}

// Name returns the name of the projection
func (p *ProjectionMeta) Name() string { return p.name }

// Description returns the description of the projection
func (p *ProjectionMeta) Description() string { return p.desc }

// Desc is an alias for Description for backward compatibility
func (p *ProjectionMeta) Desc() string { return p.Description() }

// NameA returns the name of the first coordinate
func (p *ProjectionMeta) NameA() string { return p.coordinates.A.Name() }

// DescriptionA returns the description of the first coordinate
func (p *ProjectionMeta) DescriptionA() string { return p.coordinates.A.Description() }

// DescA is an alias for DescriptionA for backward compatibility
func (p *ProjectionMeta) DescA() string { return p.DescriptionA() }

// UnitA returns the unit of the first coordinate
func (p *ProjectionMeta) UnitA() string { return p.coordinates.A.Unit() }

// MinA returns the minimum value of the first coordinate
func (p *ProjectionMeta) MinA() float64 { return p.coordinates.A.Min() }

// MaxA returns the maximum value of the first coordinate
func (p *ProjectionMeta) MaxA() float64 { return p.coordinates.A.Max() }

// NameB returns the name of the second coordinate
func (p *ProjectionMeta) NameB() string { return p.coordinates.B.Name() }

// DescriptionB returns the description of the second coordinate
func (p *ProjectionMeta) DescriptionB() string { return p.coordinates.B.Description() }

// DescB is an alias for DescriptionB for backward compatibility
func (p *ProjectionMeta) DescB() string { return p.DescriptionB() }

// UnitB returns the unit of the second coordinate
func (p *ProjectionMeta) UnitB() string { return p.coordinates.B.Unit() }

// MinB returns the minimum value of the second coordinate
func (p *ProjectionMeta) MinB() float64 { return p.coordinates.B.Min() }

// MaxB returns the maximum value of the second coordinate
func (p *ProjectionMeta) MaxB() float64 { return p.coordinates.B.Max() }

// Doc returns a markdown-formatted documentation string for the projection
func (p *ProjectionMeta) Doc() string {
	sb := strings.Builder{}
	sb.WriteString("# " + p.name + "\n")
	sb.WriteString("_" + p.desc + "_\n\n")
	sb.WriteString("| Component | Min | Max | Unit | Description |\n")
	sb.WriteString("| --- | --- | --- | --- | --- |\n")
	sb.WriteString(fmt.Sprintf("| %s | %f | %f | %s | %s |\n", p.NameA(), p.MinA(), p.MaxA(), p.UnitA(), p.DescriptionA()))
	sb.WriteString(fmt.Sprintf("| %s | %f | %f | %s | %s |\n", p.NameB(), p.MinB(), p.MaxB(), p.UnitB(), p.DescriptionB()))
	return sb.String()
}

// NewProjectionMeta creates a new projection metadata object
func NewProjectionMeta(name, description string, coordinateA, coordinateB *Coordinate) *ProjectionMeta {
	return &ProjectionMeta{
		name: name,
		desc: description,
		coordinates: &Coordinates{
			A: coordinateA,
			B: coordinateB,
		},
	}
}
