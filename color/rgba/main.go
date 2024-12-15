package rgba

import (
	"fmt"
	"image/color"

	"github.com/toxyl/gfx/math"
)

// RGBA is similar to color.RGBA but doesn't store RGB with premultiplied alpha
type RGBA struct {
	Red   uint8 `yaml:"r"`
	Green uint8 `yaml:"g"`
	Blue  uint8 `yaml:"b"`
	Alpha uint8 `yaml:"a"`
}

func New[N math.Number](r, g, b, a N) *RGBA {
	rgba := RGBA{
		Red:   uint8(math.Round(r)),
		Green: uint8(math.Round(g)),
		Blue:  uint8(math.Round(b)),
		Alpha: uint8(math.Round(a)),
	}
	return &rgba
}

func (rgba *RGBA) String() string {
	return fmt.Sprintf("r: %d, g: %d, b: %d, a: %d", rgba.Red, rgba.Green, rgba.Blue, rgba.Alpha)
}

func (rgba *RGBA) R() uint8 { return rgba.Red }
func (rgba *RGBA) G() uint8 { return rgba.Green }
func (rgba *RGBA) B() uint8 { return rgba.Blue }
func (rgba *RGBA) A() uint8 { return rgba.Alpha }

func (rgba *RGBA) SetR(v uint8) *RGBA { rgba.Red = v; return rgba }
func (rgba *RGBA) SetG(v uint8) *RGBA { rgba.Green = v; return rgba }
func (rgba *RGBA) SetB(v uint8) *RGBA { rgba.Blue = v; return rgba }
func (rgba *RGBA) SetA(v uint8) *RGBA { rgba.Alpha = v; return rgba }

func (rgba *RGBA) RGB() RGBA {
	return RGBA{Red: rgba.Red, Green: rgba.Green, Blue: rgba.Blue, Alpha: 0xFF}
}

// RGBA returns a version of the color that is alpha-premultiplied.
func (rgba *RGBA) RGBA() color.RGBA {
	c := color.RGBA{
		R: rgba.Red,
		G: rgba.Green,
		B: rgba.Blue,
		A: 0xFF,
	}
	if c.A == 0 { // Avoid division by zero if original alpha is zero
		return color.RGBA{0, 0, 0, 0}
	}
	scale := float32(rgba.Alpha) / float32(c.A)
	return color.RGBA{
		R: uint8(float32(c.R) * scale),
		G: uint8(float32(c.G) * scale),
		B: uint8(float32(c.B) * scale),
		A: rgba.Alpha,
	}
}
