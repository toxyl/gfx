package parser

import (
	"fmt"
	stdcolor "image/color"
	"path/filepath"
	"sort"
	"strings"

	"github.com/toxyl/gfx/config"
	"github.com/toxyl/gfx/config/constants"
	"github.com/toxyl/gfx/core/blendmodes"
	"github.com/toxyl/gfx/core/color"
	"github.com/toxyl/gfx/core/image"
	"github.com/toxyl/gfx/fs"
	"github.com/toxyl/math"
)

type Composition struct {
	Width  int
	Height int
	Color  *color.HSL
	Crop   *Crop
	Resize *Resize
	Filter *FX
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
			spfPad(maxLenName, constants.COMP_WIDTH), constants.ASSIGN, c.Width))
	}
	if c.Height != 0 {
		comp = append(comp, spf("%s %s %d",
			spfPad(maxLenName, constants.COMP_HEIGHT), constants.ASSIGN, c.Height))
	}
	if c.Color != nil {
		comp = append(comp, spf("%s %s hsla%s%f %f %f %f%s",
			spfPad(maxLenName, constants.COMP_COLOR), constants.ASSIGN, constants.LPAREN, c.Color.H, c.Color.S, c.Color.L, c.Color.Alpha, constants.RPAREN))
	}
	if c.Filter != nil {
		comp = append(comp, spf("%s %s %s",
			spfPad(maxLenName, constants.COMP_FILTER), constants.ASSIGN, c.Filter.Name))
	}
	if c.Crop != nil {
		comp = append(comp, spf("%s %s %d %d %d %d",
			spfPad(maxLenName, constants.COMP_CROP), constants.ASSIGN, c.Crop.X, c.Crop.Y, c.Crop.W, c.Crop.H))
	}
	if c.Resize != nil {
		comp = append(comp, spf("%s %s %d %d",
			spfPad(maxLenName, constants.COMP_RESIZE), constants.ASSIGN, c.Resize.W, c.Resize.H))
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
		constants.LBRACKET, strings.ToUpper(constants.SECTION_COMPOSITION), constants.RBRACKET,
		strings.Join(comp, "\n"),
		constants.LBRACKET, strings.ToUpper(constants.SECTION_LAYERS), constants.RBRACKET,
		strings.Join(layers, "\n"),
		constants.LBRACKET, strings.ToUpper(constants.SECTION_VARS), constants.RBRACKET,
		strings.Join(vars, "\n"),
		constants.LBRACKET, strings.ToUpper(constants.SECTION_FILTERS), constants.RBRACKET,
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

// RGBA64 represents an RGBA color with 64-bit precision.
type RGBA64 struct {
	R, G, B, A float64
}

// To16bit returns the color as a standard library 16-bit RGBA.
func (c *RGBA64) To16bit() stdcolor.RGBA64 {
	return stdcolor.RGBA64{
		R: uint16(c.R * 65535),
		G: uint16(c.G * 65535),
		B: uint16(c.B * 65535),
		A: uint16(c.A * 65535),
	}
}

func (c *Composition) Render(images ...string) (*image.Image, error) {
	w, h := c.Width, c.Height
	res, err := image.New(w, h)
	if err != nil {
		return nil, fmt.Errorf("failed to create image: %w", err)
	}

	if c.Color != nil {
		// Use Process to fill with color
		col := c.Color.ToRGBA64()

		// Instead of using Process with custom types, iterate manually
		for y := 0; y < h; y++ {
			for x := 0; x < w; x++ {
				err := res.SetPixel(x, y, col.ToRGBA64())
				if err != nil {
					return nil, fmt.Errorf("failed to set pixel: %w", err)
				}
			}
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

		// Get blend mode
		blendMode, err := blendmodes.Get(l.BlendMode)
		if err != nil {
			blendMode, _ = blendmodes.Get("normal") // Use normal mode as fallback
		}

		// Draw the scaled image onto the result
		err = res.DrawImage(scaled, 0, 0, blendMode, l.Alpha)
		if err != nil {
			return nil, fmt.Errorf("failed to blend layer: %w", err)
		}
	}

	if c.Filter != nil {
		for _, filter := range c.Filter.Get() {
			if filter != nil {
				var filterErr error
				stdImg := res.ToStandard()
				newImg, filterErr := filter.Apply(stdImg)
				if filterErr != nil {
					return nil, fmt.Errorf("failed to apply filter: %w", filterErr)
				}
				res, err = image.FromImage(newImg)
				if err != nil {
					return nil, fmt.Errorf("failed to convert image: %w", err)
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
