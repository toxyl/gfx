package projections

import (
	"math"
)

func init() {
	// Register Equirectangular projection
	Default.Register(
		&RegistryProjection{
			Meta: NewProjectionMeta(
				"equirectangular",
				"A linear mapping between cartesian x/y and lat/lon.",
				NewCoordinateMeta("lat", "The latitude where 90 is North and -90 is South.", "°", 90, -90),
				NewCoordinateMeta("lon", "The longitude where -180 is West and 180 is East.", "°", -180, 180),
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
			Meta: NewProjectionMeta(
				"mercator",
				"A conformal cylindrical map projection.",
				NewCoordinateMeta("lat", "The latitude where 90 is North and -90 is South.", "°", 90, -90),
				NewCoordinateMeta("lon", "The longitude where -180 is West and 180 is East.", "°", -180, 180),
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
			Meta: NewProjectionMeta(
				"sinusoidal",
				"An equal-area pseudocylindrical projection.",
				NewCoordinateMeta("lat", "The latitude where 90 is North and -90 is South.", "°", 90, -90),
				NewCoordinateMeta("lon", "The longitude where -180 is West and 180 is East.", "°", -180, 180),
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
			Meta: NewProjectionMeta(
				"stereographic",
				"A conformal map projection.",
				NewCoordinateMeta("lat", "The latitude where 90 is North and -90 is South.", "°", 90, -90),
				NewCoordinateMeta("lon", "The longitude where -180 is West and 180 is East.", "°", -180, 180),
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
			Meta: NewProjectionMeta(
				"rectpolar",
				"Converts from rectangular to polar coordinates.",
				NewCoordinateMeta("radius", "The distance from the center.", "", 1, 0),
				NewCoordinateMeta("angle", "The angle in radians.", "rad", 2*math.Pi, 0),
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
			Meta: NewProjectionMeta(
				"polarrect",
				"Converts from polar to rectangular coordinates.",
				NewCoordinateMeta("radius", "The distance from the center.", "", 1, 0),
				NewCoordinateMeta("angle", "The angle in radians.", "rad", 2*math.Pi, 0),
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
