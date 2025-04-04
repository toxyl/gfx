package color

import "fmt"

// Palette represents a collection of colors
type Palette struct {
	name   string
	colors []iColor
}

// NewPalette creates a new empty palette
func NewPalette(name string) *Palette {
	return &Palette{
		name: name,
	}
}

// Add adds a color to the palette
func (p *Palette) Add(color iColor) *Palette {
	p.colors = append(p.colors, color)
	return p
}

// Get returns the color at the specified index
func (p *Palette) Get(index int) (iColor, error) {
	if index < 0 || index >= len(p.colors) {
		return nil, fmt.Errorf("index out of range")
	}
	return p.colors[index], nil
}

// Name returns the palette's name
func (p *Palette) Name() string {
	return p.name
}

// Colors returns all colors in the palette
func (p *Palette) Colors() []iColor {
	return p.colors
}

// Predefined palettes
var (
	// WebSafe is a palette of web-safe colors
	WebSafe = NewPalette("Web Safe").
		Add(&RGB8{R: 0, G: 0, B: 0, Alpha: 1.0}).
		Add(&RGB8{R: 51, G: 51, B: 51, Alpha: 1.0}).
		Add(&RGB8{R: 102, G: 102, B: 102, Alpha: 1.0}).
		Add(&RGB8{R: 153, G: 153, B: 153, Alpha: 1.0}).
		Add(&RGB8{R: 204, G: 204, B: 204, Alpha: 1.0}).
		Add(&RGB8{R: 255, G: 255, B: 255, Alpha: 1.0})

	// MaterialDesign is a palette of Material Design colors
	MaterialDesign = NewPalette("Material Design").
			Add(&RGB8{R: 244, G: 67, B: 54, Alpha: 1.0}).  // Red
			Add(&RGB8{R: 233, G: 30, B: 99, Alpha: 1.0}).  // Pink
			Add(&RGB8{R: 156, G: 39, B: 176, Alpha: 1.0}). // Purple
			Add(&RGB8{R: 103, G: 58, B: 183, Alpha: 1.0}). // Deep Purple
			Add(&RGB8{R: 63, G: 81, B: 181, Alpha: 1.0}).  // Indigo
			Add(&RGB8{R: 33, G: 150, B: 243, Alpha: 1.0}). // Blue
			Add(&RGB8{R: 3, G: 169, B: 244, Alpha: 1.0}).  // Light Blue
			Add(&RGB8{R: 0, G: 188, B: 212, Alpha: 1.0}).  // Cyan
			Add(&RGB8{R: 0, G: 150, B: 136, Alpha: 1.0}).  // Teal
			Add(&RGB8{R: 76, G: 175, B: 80, Alpha: 1.0}).  // Green
			Add(&RGB8{R: 139, G: 195, B: 74, Alpha: 1.0}). // Light Green
			Add(&RGB8{R: 205, G: 220, B: 57, Alpha: 1.0}). // Lime
			Add(&RGB8{R: 255, G: 235, B: 59, Alpha: 1.0}). // Yellow
			Add(&RGB8{R: 255, G: 193, B: 7, Alpha: 1.0}).  // Amber
			Add(&RGB8{R: 255, G: 152, B: 0, Alpha: 1.0}).  // Orange
			Add(&RGB8{R: 255, G: 87, B: 34, Alpha: 1.0})   // Deep Orange
)
