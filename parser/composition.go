package parser

import (
	"fmt"
	"image/color"
	"path/filepath"
	"sort"
	"strings"

	"github.com/toxyl/gfx/config"
	"github.com/toxyl/gfx/core/blendmodes/registry"
	"github.com/toxyl/gfx/core/colormodels/hsl"
	"github.com/toxyl/gfx/core/image"
	"github.com/toxyl/gfx/fs"
	"github.com/toxyl/math"
)

type Composition struct {
	Width  int
	Height int
	Color  *hsl.HSL
	Crop   *Crop
	Resize *Resize
	Filter *Filter
	Vars   []*Var
	Layers []*Layer
}

func (c *Composition) String() string {
	maxLenName := math.MaxLenStr(config.COMPOSITION...)
	spf := fmt.Sprintf
	spfPad := func(l int, s string) string { return spf("%-*s", l, s) }
	comp := []string{}
	if c.Width != 0 {
		comp = append(comp, spf("%s %s %d",
			spfPad(maxLenName, config.COMP_WIDTH), config.ASSIGN, c.Width))
	}
	if c.Height != 0 {
		comp = append(comp, spf("%s %s %d",
			spfPad(maxLenName, config.COMP_HEIGHT), config.ASSIGN, c.Height))
	}
	if c.Color != nil {
		comp = append(comp, spf("%s %s hsla%s%f %f %f %f%s",
			spfPad(maxLenName, config.COMP_COLOR), config.ASSIGN, config.LPAREN, c.Color.H(), c.Color.S(), c.Color.L(), c.Color.A(), config.RPAREN))
	}
	if c.Filter != nil {
		comp = append(comp, spf("%s %s %s",
			spfPad(maxLenName, config.COMP_FILTER), config.ASSIGN, c.Filter.Name))
	}
	if c.Crop != nil {
		comp = append(comp, spf("%s %s %d %d %d %d",
			spfPad(maxLenName, config.COMP_CROP), config.ASSIGN, c.Crop.X, c.Crop.Y, c.Crop.W, c.Crop.H))
	}
	if c.Resize != nil {
		comp = append(comp, spf("%s %s %d %d",
			spfPad(maxLenName, config.COMP_RESIZE), config.ASSIGN, c.Resize.W, c.Resize.H))
	}
	vars := []string{}
	filters := []string{}
	layers := []string{}
	collectedFilters := map[string]struct{}{}
	if c.Layers != nil {
		for _, l := range c.Layers {
			if l == nil {
				continue
			}
		}
		for _, l := range c.Layers {
			if l == nil {
				continue
			}
			layers = append(layers, l.String())
			if l.Filter != nil {
				if _, ok := collectedFilters[l.Filter.Name]; !ok {
					filters = append(filters, l.Filter.String())
					collectedFilters[l.Filter.Name] = struct{}{}
				}
			}
		}
	}
	sort.Slice(c.Vars, func(i, j int) bool {
		return c.Vars[i].Name < c.Vars[j].Name
	})
	for _, v := range c.Vars {
		vars = append(vars, v.String())
	}
	if c.Filter != nil {
		filters = append(filters, c.Filter.String())
	}
	return spf(
		`%s%s%s
%s

%s%s%s
%s

%s%s%s
%s

%s%s%s
%s
`,
		config.LBRACKET, strings.ToUpper(config.SECTION_COMPOSITION), config.RBRACKET,
		strings.Join(comp, "\n"),
		config.LBRACKET, strings.ToUpper(config.SECTION_LAYERS), config.RBRACKET,
		strings.Join(layers, "\n"),
		config.LBRACKET, strings.ToUpper(config.SECTION_VARS), config.RBRACKET,
		strings.Join(vars, "\n"),
		config.LBRACKET, strings.ToUpper(config.SECTION_FILTERS), config.RBRACKET,
		strings.Join(filters, "\n"),
	)
}

func (c *Composition) Load(path string) *Composition {
	str := ""
	if err := fs.LoadStringInto(path, &str); err != nil {
		panic("failed to load composition: " + err.Error())
	}
	comp, err := ParseComposition(str)
	if err != nil {
		panic("failed to parse composition: " + err.Error())
	}
	c = comp
	return c
}

func (c *Composition) Save(path string) *Composition {
	if err := fs.SaveString(path, c.String()); err != nil {
		panic("failed to save composition: " + err.Error())
	}
	return c
}

func (c *Composition) Render(images ...string) (*image.Image, error) {
	w, h := c.Width, c.Height
	res, err := image.New(w, h)
	if err != nil {
		return nil, fmt.Errorf("failed to create image: %w", err)
	}

	if c.Color != nil {
		// Use Process to fill with color
		err = res.Process(func(x, y int, _ *color.RGBA64) (*color.RGBA64, error) {
			return &color.RGBA64{
				R: c.Color.R(),
				G: c.Color.G(),
				B: c.Color.B(),
				A: c.Color.A(),
			}, nil
		})
		if err != nil {
			return nil, fmt.Errorf("failed to fill image with color: %w", err)
		}
	}

	numLayers := len(c.Layers)
	for i := numLayers - 1; i >= 0; i-- {
		l := c.Layers[i]
		l.Source = strings.TrimSpace(l.Source)

		for i, img := range images {
			placeholder := "$IMG" + fmt.Sprint(i+1)
			if l.Source == placeholder {
				absImagePath, err := filepath.Abs(img)
				if err != nil {
					absImagePath = img
				}
				l.Source = strings.ReplaceAll(l.Source, placeholder, absImagePath)
			}
		}

		scaled, err := l.Render(w, h)
		if err != nil || scaled == nil {
			fmt.Printf("WARN: failed to render layer source %s, ignoring layer: %v\n", l.Source, err)
			continue // rendering failed, maybe URL or file wasn't available
		}

		// Use Blend method for drawing one image onto another
		blendMode := registry.Get(l.BlendMode)
		err = res.Blend(scaled, 0, 0, w, h, 0, 0, w, h, blendMode, l.Alpha)
		if err != nil {
			return nil, fmt.Errorf("failed to blend layer: %w", err)
		}
	}

	if c.Filter != nil {
		for _, filter := range c.Filter.Get() {
			if filter != nil {
				// Apply filter will need to be updated to return error
				var filterErr error
				res, filterErr = filter.Apply(res)
				if filterErr != nil {
					return nil, fmt.Errorf("failed to apply filter: %w", filterErr)
				}
			}
		}
	}

	if c.Crop != nil && c.Crop.W > 0 && c.Crop.H > 0 {
		res, err = res.Crop(c.Crop.X, c.Crop.Y, c.Crop.W, c.Crop.H)
		if err != nil {
			return nil, fmt.Errorf("failed to crop image: %w", err)
		}
	}

	if c.Resize != nil && c.Resize.W > 0 && c.Resize.H > 0 {
		res, err = res.Resize(c.Resize.W, c.Resize.H, image.ResizeBilinear)
		if err != nil {
			return nil, fmt.Errorf("failed to resize image: %w", err)
		}
		return res, nil
	}

	resized, err := res.Resize(w, h, image.ResizeBilinear)
	if err != nil {
		return nil, fmt.Errorf("failed to resize image: %w", err)
	}
	return resized, nil
}

func NewComposition(w, h int) *Composition {
	c := Composition{
		Width:  w,
		Height: h,
		Layers: []*Layer{},
		Color:  nil,
		Crop:   nil,
		Resize: nil,
		Filter: nil,
	}
	return &c
}
