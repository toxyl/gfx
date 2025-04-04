package projections

import (
	"math"

	"github.com/toxyl/gfx/registry/meta"
)

func init() {
	// Register Equirectangular projection
	Default.Register(
		&RegistryProjection{
			Meta: meta.NewProjection(
				"equirectangular",
				"A linear mapping between cartesian x/y and lat/lon.",
				meta.NewCoordinate("lat", "The latitude where 90 is North and -90 is South.", "°", 90, -90),
				meta.NewCoordinate("lon", "The longitude where -180 is West and 180 is East.", "°", -180, 180),
			),
			To: func(latitude, longitude, w, h float64) (x, y float64) {
				return NewEquirectangular().To(latitude, longitude, w, h)
			},
			From: func(x, y, w, h float64) (latitude, longitude float64) {
				return NewEquirectangular().From(x, y, w, h)
			},
		},
	)

	// Register Mercator projection
	Default.Register(
		&RegistryProjection{
			Meta: meta.NewProjection(
				"mercator",
				"A conformal cylindrical map projection.",
				meta.NewCoordinate("lat", "The latitude where 90 is North and -90 is South.", "°", 90, -90),
				meta.NewCoordinate("lon", "The longitude where -180 is West and 180 is East.", "°", -180, 180),
			),
			To: func(latitude, longitude, w, h float64) (x, y float64) {
				return NewMercator().To(latitude, longitude, w, h)
			},
			From: func(x, y, w, h float64) (latitude, longitude float64) {
				return NewMercator().From(x, y, w, h)
			},
		},
	)

	// Register Sinusoidal projection
	Default.Register(
		&RegistryProjection{
			Meta: meta.NewProjection(
				"sinusoidal",
				"An equal-area pseudocylindrical projection.",
				meta.NewCoordinate("lat", "The latitude where 90 is North and -90 is South.", "°", 90, -90),
				meta.NewCoordinate("lon", "The longitude where -180 is West and 180 is East.", "°", -180, 180),
			),
			To: func(latitude, longitude, w, h float64) (x, y float64) {
				return NewSinusoidal().To(latitude, longitude, w, h)
			},
			From: func(x, y, w, h float64) (latitude, longitude float64) {
				return NewSinusoidal().From(x, y, w, h)
			},
		},
	)

	// Register Stereographic projection
	Default.Register(
		&RegistryProjection{
			Meta: meta.NewProjection(
				"stereographic",
				"A conformal map projection.",
				meta.NewCoordinate("lat", "The latitude where 90 is North and -90 is South.", "°", 90, -90),
				meta.NewCoordinate("lon", "The longitude where -180 is West and 180 is East.", "°", -180, 180),
			),
			To: func(latitude, longitude, w, h float64) (x, y float64) {
				return NewStereographic().To(latitude, longitude, w, h)
			},
			From: func(x, y, w, h float64) (latitude, longitude float64) {
				return NewStereographic().From(x, y, w, h)
			},
		},
	)

	// Register RectPolar projection
	Default.Register(
		&RegistryProjection{
			Meta: meta.NewProjection(
				"rectpolar",
				"Converts from rectangular to polar coordinates.",
				meta.NewCoordinate("radius", "The distance from the center.", "", 1, 0),
				meta.NewCoordinate("angle", "The angle in radians.", "rad", 2*math.Pi, 0),
			),
			To: func(x, y, w, h float64) (radius, angle float64) {
				return NewRectPolar().To(x, y, w, h)
			},
			From: func(radius, angle, w, h float64) (x, y float64) {
				return NewRectPolar().From(radius, angle, w, h)
			},
		},
	)

	// Register PolarRect projection
	Default.Register(
		&RegistryProjection{
			Meta: meta.NewProjection(
				"polarrect",
				"Converts from polar to rectangular coordinates.",
				meta.NewCoordinate("radius", "The distance from the center.", "", 1, 0),
				meta.NewCoordinate("angle", "The angle in radians.", "rad", 2*math.Pi, 0),
			),
			To: func(radius, angle, w, h float64) (x, y float64) {
				return NewPolarRect().To(radius, angle, w, h)
			},
			From: func(x, y, w, h float64) (radius, angle float64) {
				return NewPolarRect().From(x, y, w, h)
			},
		},
	)
}
