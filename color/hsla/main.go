package hsla

import (
	"fmt"

	"github.com/toxyl/gfx/math"
)

// HSL represents a color in the HSL color space.
type HSLA struct {
	h float64 // Hue, in degrees [0, 360)
	s float64 // Saturation, as a fraction [0, 1]
	l float64 // Lightness, as a fraction [0, 1]
	a float64 // Alpha, as a fraction [0, 1]
}

func New[N math.Number](h, s, l, a N) *HSLA {
	return &HSLA{h: math.Wrap(float64(h), 0.0, 360.0), s: float64(s), l: float64(l), a: float64(a)}
}

func (hsla *HSLA) String() string {
	return fmt.Sprintf("h: %f, s: %f, s: %f, a: %f", hsla.h, hsla.s, hsla.l, hsla.a)
}

func (hsla *HSLA) H() float64 { return hsla.h }
func (hsla *HSLA) S() float64 { return hsla.s }
func (hsla *HSLA) L() float64 { return hsla.l }
func (hsla *HSLA) A() float64 { return hsla.a }

func (hsla *HSLA) SetH(v float64) *HSLA { hsla.h = math.Wrap(v, 0.0, 360.0); return hsla }
func (hsla *HSLA) SetS(v float64) *HSLA { hsla.s = math.Clamp(v, 0.0, 1.0); return hsla }
func (hsla *HSLA) SetL(v float64) *HSLA { hsla.l = math.Clamp(v, 0.0, 1.0); return hsla }
func (hsla *HSLA) SetA(v float64) *HSLA { hsla.a = math.Clamp(v, 0.0, 1.0); return hsla }

func (hsla *HSLA) ShiftH(v float64) *HSLA { hsla.h = math.Wrap(hsla.h+v, 0.0, 360.0); return hsla }
func (hsla *HSLA) ShiftS(v float64) *HSLA { hsla.s = math.Clamp(hsla.s+v, 0.0, 1.0); return hsla }
func (hsla *HSLA) ShiftL(v float64) *HSLA { hsla.l = math.Clamp(hsla.l+v, 0.0, 1.0); return hsla }
func (hsla *HSLA) ShiftA(v float64) *HSLA {
	hsla.a = math.Clamp(hsla.a+v, 0.0, 1.0)
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
	return fmt.Sprintf("%d:%.2f:%.2f", int(hsla.h), hsla.s, hsla.l)
}
