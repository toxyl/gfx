package rgba

import (
	"fmt"
	"image/color"

	"github.com/toxyl/gfx/math"
)

// RGBA is similar to color.RGBA but doesn't store RGB with premultiplied alpha
type RGBA struct {
	r uint8
	g uint8
	b uint8
	a uint8
}

func New[N math.Number](r, g, b, a N) *RGBA {
	rgba := RGBA{
		r: uint8(math.Round(r)),
		g: uint8(math.Round(g)),
		b: uint8(math.Round(b)),
		a: uint8(math.Round(a)),
	}
	return &rgba
}

func (rgba *RGBA) String() string {
	return fmt.Sprintf("r: %d, g: %d, b: %d, a: %d", rgba.r, rgba.g, rgba.b, rgba.a)
}

func (rgba *RGBA) R() uint8 { return rgba.r }
func (rgba *RGBA) G() uint8 { return rgba.g }
func (rgba *RGBA) B() uint8 { return rgba.b }
func (rgba *RGBA) A() uint8 { return rgba.a }

func (rgba *RGBA) SetR(v uint8) *RGBA { rgba.r = v; return rgba }
func (rgba *RGBA) SetG(v uint8) *RGBA { rgba.g = v; return rgba }
func (rgba *RGBA) SetB(v uint8) *RGBA { rgba.b = v; return rgba }
func (rgba *RGBA) SetA(v uint8) *RGBA { rgba.a = v; return rgba }

func (rgba *RGBA) RGB() RGBA {
	return RGBA{r: rgba.r, g: rgba.g, b: rgba.b, a: 0xFF}
}

// RGBA returns a version of the color that is alpha-premultiplied.
func (rgba *RGBA) RGBA() color.RGBA {
	c := color.RGBA{
		R: rgba.r,
		G: rgba.g,
		B: rgba.b,
		A: 0xFF,
	}
	if c.A == 0 { // Avoid division by zero if original alpha is zero
		return color.RGBA{0, 0, 0, 0}
	}
	scale := float32(rgba.a) / float32(c.A)
	return color.RGBA{
		R: uint8(float32(c.R) * scale),
		G: uint8(float32(c.G) * scale),
		B: uint8(float32(c.B) * scale),
		A: rgba.a,
	}
}
