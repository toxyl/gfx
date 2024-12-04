package composition

import (
	"github.com/toxyl/gfx/filters"
	img "github.com/toxyl/gfx/image"

	"github.com/toxyl/flo"
)

type CropConfig struct {
	X int `yaml:"x"`
	Y int `yaml:"y"`
	W int `yaml:"w"`
	H int `yaml:"h"`
}

type Composition struct {
	Enabled bool                   `yaml:"enabled"`
	Name    string                 `yaml:"name"`
	Layers  []*Layer               `yaml:"layers"`
	Crop    *CropConfig            `yaml:"crop"`
	Filters []*filters.ImageFilter `yaml:"filters"`
	Width   int                    `yaml:"width"`
	Height  int                    `yaml:"height"`
}

func New(width, height int, layers ...*Layer) *Composition {
	c := Composition{
		Enabled: true,
		Name:    "",
		Layers:  layers,
		Crop:    nil,
		Width:   width,
		Height:  height,
	}
	return &c
}

func (c *Composition) Load(path string) *Composition {
	if err := flo.File(path).LoadYAML(&c); err != nil {
		panic("failed to load composition config file: " + err.Error())
	}
	return c
}

func (c *Composition) Save(path string) *Composition {
	if err := flo.File(path).StoreYAML(&c); err != nil {
		panic("failed to save composition config file: " + err.Error())
	}
	return c
}

func NewFromFile(path string) *Composition {
	c := &Composition{
		Enabled: true,
		Name:    "",
		Layers:  []*Layer{},
		Crop:    nil,
		Width:   0,
		Height:  0,
	}
	return c.Load(path)
}

func (c *Composition) Render() *img.Image {
	w, h := c.Width, c.Height
	res := img.New(w, h)
	ll := len(c.Layers)
	for _, l := range c.Layers {
		scaled := l.Render(w, h)
		if ll == 1 {
			res = scaled
			continue
		}
		res.Draw(scaled, 0, 0, w, h, 0, 0, w, h, l.Mode, l.Alpha)
	}
	for _, filter := range c.Filters {
		if filter != nil {
			res = filter.Apply(res)
		}
	}
	if c.Crop != nil {
		res = res.Crop(c.Crop.X, c.Crop.Y, c.Crop.W, c.Crop.H)
	}
	return res
}
