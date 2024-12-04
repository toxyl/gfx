package composition

import (
	"github.com/toxyl/gfx/color/blend"
	"github.com/toxyl/gfx/filters"
	"github.com/toxyl/gfx/image"
)

type Layer struct {
	data    *image.Image
	Path    string                 `yaml:"path"`
	Alpha   float64                `yaml:"alpha"`
	Mode    blend.BlendMode        `yaml:"mode"`
	Crop    *CropConfig            `yaml:"crop"`
	Filters []*filters.ImageFilter `yaml:"filters"`
}

func NewLayer(path string, alpha float64, blendMode blend.BlendMode, cropConfig *CropConfig, filters ...*filters.ImageFilter) *Layer {
	cl := Layer{
		data:    nil,
		Path:    path,
		Alpha:   alpha,
		Mode:    blendMode,
		Crop:    cropConfig,
		Filters: filters,
	}
	cl.load()
	return &cl
}

func NewLayerFromImage(image *image.Image, alpha float64, blendMode blend.BlendMode, cropConfig *CropConfig, filters ...*filters.ImageFilter) *Layer {
	cl := Layer{
		data:    image,
		Path:    "",
		Alpha:   alpha,
		Mode:    blendMode,
		Crop:    cropConfig,
		Filters: filters,
	}
	return &cl
}

func (cl *Layer) load() {
	if cl.Path != "" {
		cl.data = image.NewFromURL(cl.Path)
		if cl.data == nil {
			// this wasn't a URL, maybe it's a file
			cl.data = image.NewFromFile(cl.Path)
		}
	}
}

func (cl *Layer) Render(w, h int) *image.Image {
	cl.load()
	res := cl.data.Resize(w, h)
	for _, filter := range cl.Filters {
		if filter != nil {
			res = filter.Apply(res)
		}
	}
	if cl.Crop != nil {
		res = res.Crop(cl.Crop.X, cl.Crop.Y, cl.Crop.W, cl.Crop.H)
	}
	return res
}
