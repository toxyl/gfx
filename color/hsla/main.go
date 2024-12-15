package hsla

import (
	"fmt"

	"github.com/toxyl/gfx/math"
)

// HSL represents a color in the HSL color space.
type HSLA struct {
	Hue   float64 `yaml:"h"` // Hue, in degrees [0, 360)
	Sat   float64 `yaml:"s"` // Saturation, as a fraction [0, 1]
	Lum   float64 `yaml:"l"` // Lightness, as a fraction [0, 1]
	Alpha float64 `yaml:"a"` // Alpha, as a fraction [0, 1]
}

func New[N math.Number](h, s, l, a N) *HSLA {
	return &HSLA{Hue: math.Wrap(float64(h), 0.0, 360.0), Sat: float64(s), Lum: float64(l), Alpha: float64(a)}
}

func (hsla *HSLA) String() string {
	return fmt.Sprintf("h: %f, s: %f, s: %f, a: %f", hsla.Hue, hsla.Sat, hsla.Lum, hsla.Alpha)
}

func (hsla *HSLA) H() float64 { return hsla.Hue }
func (hsla *HSLA) S() float64 { return hsla.Sat }
func (hsla *HSLA) L() float64 { return hsla.Lum }
func (hsla *HSLA) A() float64 { return hsla.Alpha }

func (hsla *HSLA) SetH(v float64) *HSLA { hsla.Hue = math.Wrap(v, 0.0, 360.0); return hsla }
func (hsla *HSLA) SetS(v float64) *HSLA { hsla.Sat = math.Clamp(v, 0.0, 1.0); return hsla }
func (hsla *HSLA) SetL(v float64) *HSLA { hsla.Lum = math.Clamp(v, 0.0, 1.0); return hsla }
func (hsla *HSLA) SetA(v float64) *HSLA { hsla.Alpha = math.Clamp(v, 0.0, 1.0); return hsla }

func (hsla *HSLA) ShiftH(v float64) *HSLA { hsla.Hue = math.Wrap(hsla.Hue+v, 0.0, 360.0); return hsla }
func (hsla *HSLA) ShiftS(v float64) *HSLA { hsla.Sat = math.Clamp(hsla.Sat+v, 0.0, 1.0); return hsla }
func (hsla *HSLA) ShiftL(v float64) *HSLA { hsla.Lum = math.Clamp(hsla.Lum+v, 0.0, 1.0); return hsla }
func (hsla *HSLA) ShiftA(v float64) *HSLA {
	hsla.Alpha = math.Clamp(hsla.Alpha+v, 0.0, 1.0)
	return hsla
}
func (hsla *HSLA) Shift(h, s, l, a float64) *HSLA {
	if h != 0 {
		hsla.ShiftH(h)
	}
	if s != 0 {
		hsla.ShiftS(s)
	}
	if l != 0 {
		hsla.ShiftL(l)
	}
	if a != 0 {
		hsla.ShiftA(a)
	}
	return hsla
}

func (hsla *HSLA) HSLString() string {
	return fmt.Sprintf("%d:%.2f:%.2f", int(hsla.Hue), hsla.Sat, hsla.Lum)
}
