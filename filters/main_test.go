package filters

import (
	"fmt"
	"testing"

	"github.com/toxyl/gfx/image"
)

func TestFilters(t *testing.T) {
	testImage := image.NewFromURL("https://sdo.gsfc.nasa.gov/assets/img/latest/f_211_193_171pfss_512.jpg")
	tests := []struct {
		name    string
		filter  string
		options map[string]any
	}{
		{"convolution-1", CONVOLUTION, map[string]any{"matrix": []float64{
			0.025, 0.100, 0.025,
			0.050, 0.150, 0.050,
			0.075, 0.100, 0.075,
		}}},

		{"blur--1.00", BLUR, map[string]any{"amount": -1.00}},
		{"blur--0.75", BLUR, map[string]any{"amount": -0.75}},
		{"blur--0.50", BLUR, map[string]any{"amount": -0.50}},
		{"blur--0.25", BLUR, map[string]any{"amount": -0.25}},
		{"blur-0.00", BLUR, map[string]any{"amount": 0.00}},
		{"blur-0.25", BLUR, map[string]any{"amount": 0.25}},
		{"blur-0.50", BLUR, map[string]any{"amount": 0.50}},
		{"blur-0.75", BLUR, map[string]any{"amount": 0.75}},
		{"blur-1.00", BLUR, map[string]any{"amount": 1.00}},
		{"blur-1.50", BLUR, map[string]any{"amount": 1.50}},
		{"blur-2.00", BLUR, map[string]any{"amount": 2.00}},

		{"sharpen--1.00", SHARPEN, map[string]any{"amount": -1.00}},
		{"sharpen--0.75", SHARPEN, map[string]any{"amount": -0.75}},
		{"sharpen--0.50", SHARPEN, map[string]any{"amount": -0.50}},
		{"sharpen--0.25", SHARPEN, map[string]any{"amount": -0.25}},
		{"sharpen-0.00", SHARPEN, map[string]any{"amount": 0.00}},
		{"sharpen-0.25", SHARPEN, map[string]any{"amount": 0.25}},
		{"sharpen-0.50", SHARPEN, map[string]any{"amount": 0.50}},
		{"sharpen-0.75", SHARPEN, map[string]any{"amount": 0.75}},
		{"sharpen-1.00", SHARPEN, map[string]any{"amount": 1.00}},
		{"sharpen-1.50", SHARPEN, map[string]any{"amount": 1.50}},
		{"sharpen-2.00", SHARPEN, map[string]any{"amount": 2.00}},

		{"edge-detection--1.00", EDGE_DETECT, map[string]any{"amount": -1.00}},
		{"edge-detection--0.75", EDGE_DETECT, map[string]any{"amount": -0.75}},
		{"edge-detection--0.50", EDGE_DETECT, map[string]any{"amount": -0.50}},
		{"edge-detection--0.25", EDGE_DETECT, map[string]any{"amount": -0.25}},
		{"edge-detection-0.00", EDGE_DETECT, map[string]any{"amount": 0.00}},
		{"edge-detection-0.25", EDGE_DETECT, map[string]any{"amount": 0.25}},
		{"edge-detection-0.50", EDGE_DETECT, map[string]any{"amount": 0.50}},
		{"edge-detection-0.75", EDGE_DETECT, map[string]any{"amount": 0.75}},
		{"edge-detection-1.00", EDGE_DETECT, map[string]any{"amount": 1.00}},
		{"edge-detection-1.50", EDGE_DETECT, map[string]any{"amount": 1.50}},
		{"edge-detection-2.00", EDGE_DETECT, map[string]any{"amount": 2.00}},

		{"emboss--1.00", EMBOSS, map[string]any{"amount": -1.00}},
		{"emboss--0.75", EMBOSS, map[string]any{"amount": -0.75}},
		{"emboss--0.50", EMBOSS, map[string]any{"amount": -0.50}},
		{"emboss--0.25", EMBOSS, map[string]any{"amount": -0.25}},
		{"emboss-0.00", EMBOSS, map[string]any{"amount": 0.00}},
		{"emboss-0.25", EMBOSS, map[string]any{"amount": 0.25}},
		{"emboss-0.50", EMBOSS, map[string]any{"amount": 0.50}},
		{"emboss-0.75", EMBOSS, map[string]any{"amount": 0.75}},
		{"emboss-1.00", EMBOSS, map[string]any{"amount": 1.00}},
		{"emboss-1.50", EMBOSS, map[string]any{"amount": 1.50}},
		{"emboss-2.00", EMBOSS, map[string]any{"amount": 2.00}},

		{"enhance--1.00", ENHANCE, map[string]any{"amount": -1.00}},
		{"enhance--0.75", ENHANCE, map[string]any{"amount": -0.75}},
		{"enhance--0.50", ENHANCE, map[string]any{"amount": -0.50}},
		{"enhance--0.25", ENHANCE, map[string]any{"amount": -0.25}},
		{"enhance-0.00", ENHANCE, map[string]any{"amount": 0.00}},
		{"enhance-0.25", ENHANCE, map[string]any{"amount": 0.25}},
		{"enhance-0.50", ENHANCE, map[string]any{"amount": 0.50}},
		{"enhance-0.75", ENHANCE, map[string]any{"amount": 0.75}},
		{"enhance-1.00", ENHANCE, map[string]any{"amount": 1.00}},
		{"enhance-1.50", ENHANCE, map[string]any{"amount": 1.50}},
		{"enhance-2.00", ENHANCE, map[string]any{"amount": 2.00}},

		{"brightness--1.00", BRIGHTNESS, map[string]any{"amount": -1.00}},
		{"brightness--0.75", BRIGHTNESS, map[string]any{"amount": -0.75}},
		{"brightness--0.50", BRIGHTNESS, map[string]any{"amount": -0.50}},
		{"brightness--0.25", BRIGHTNESS, map[string]any{"amount": -0.25}},
		{"brightness-0.00", BRIGHTNESS, map[string]any{"amount": 0.00}},
		{"brightness-0.25", BRIGHTNESS, map[string]any{"amount": 0.25}},
		{"brightness-0.50", BRIGHTNESS, map[string]any{"amount": 0.50}},
		{"brightness-0.75", BRIGHTNESS, map[string]any{"amount": 0.75}},
		{"brightness-1.00", BRIGHTNESS, map[string]any{"amount": 1.00}},
		{"brightness-1.50", BRIGHTNESS, map[string]any{"amount": 1.50}},
		{"brightness-2.00", BRIGHTNESS, map[string]any{"amount": 2.00}},

		{"contrast--1.00", CONTRAST, map[string]any{"amount": -1.00}},
		{"contrast--0.75", CONTRAST, map[string]any{"amount": -0.75}},
		{"contrast--0.50", CONTRAST, map[string]any{"amount": -0.50}},
		{"contrast--0.25", CONTRAST, map[string]any{"amount": -0.25}},
		{"contrast-0.00", CONTRAST, map[string]any{"amount": 0.00}},
		{"contrast-0.25", CONTRAST, map[string]any{"amount": 0.25}},
		{"contrast-0.50", CONTRAST, map[string]any{"amount": 0.50}},
		{"contrast-0.75", CONTRAST, map[string]any{"amount": 0.75}},
		{"contrast-1.00", CONTRAST, map[string]any{"amount": 1.00}},
		{"contrast-1.50", CONTRAST, map[string]any{"amount": 1.50}},
		{"contrast-2.00", CONTRAST, map[string]any{"amount": 2.00}},

		{"gamma--1.00", GAMMA, map[string]any{"amount": -1.00}},
		{"gamma--0.75", GAMMA, map[string]any{"amount": -0.75}},
		{"gamma--0.50", GAMMA, map[string]any{"amount": -0.50}},
		{"gamma--0.25", GAMMA, map[string]any{"amount": -0.25}},
		{"gamma-0.00", GAMMA, map[string]any{"amount": 0.00}},
		{"gamma-0.25", GAMMA, map[string]any{"amount": 0.25}},
		{"gamma-0.50", GAMMA, map[string]any{"amount": 0.50}},
		{"gamma-0.75", GAMMA, map[string]any{"amount": 0.75}},
		{"gamma-1.00", GAMMA, map[string]any{"amount": 1.00}},
		{"gamma-1.50", GAMMA, map[string]any{"amount": 1.50}},
		{"gamma-2.00", GAMMA, map[string]any{"amount": 2.00}},

		{"lightness-contrast--1.00", LIGHTNESS_CONTRAST, map[string]any{"amount": -1.00}},
		{"lightness-contrast--0.75", LIGHTNESS_CONTRAST, map[string]any{"amount": -0.75}},
		{"lightness-contrast--0.50", LIGHTNESS_CONTRAST, map[string]any{"amount": -0.50}},
		{"lightness-contrast--0.25", LIGHTNESS_CONTRAST, map[string]any{"amount": -0.25}},
		{"lightness-contrast-0.00", LIGHTNESS_CONTRAST, map[string]any{"amount": 0.00}},
		{"lightness-contrast-0.25", LIGHTNESS_CONTRAST, map[string]any{"amount": 0.25}},
		{"lightness-contrast-0.50", LIGHTNESS_CONTRAST, map[string]any{"amount": 0.50}},
		{"lightness-contrast-0.75", LIGHTNESS_CONTRAST, map[string]any{"amount": 0.75}},
		{"lightness-contrast-1.00", LIGHTNESS_CONTRAST, map[string]any{"amount": 1.00}},
		{"lightness-contrast-1.50", LIGHTNESS_CONTRAST, map[string]any{"amount": 1.50}},
		{"lightness-contrast-2.00", LIGHTNESS_CONTRAST, map[string]any{"amount": 2.00}},

		{"saturation--1.00", SATURATION, map[string]any{"amount": -1.00}},
		{"saturation--0.75", SATURATION, map[string]any{"amount": -0.75}},
		{"saturation--0.50", SATURATION, map[string]any{"amount": -0.50}},
		{"saturation--0.25", SATURATION, map[string]any{"amount": -0.25}},
		{"saturation-0.00", SATURATION, map[string]any{"amount": 0.00}},
		{"saturation-0.25", SATURATION, map[string]any{"amount": 0.25}},
		{"saturation-0.50", SATURATION, map[string]any{"amount": 0.50}},
		{"saturation-0.75", SATURATION, map[string]any{"amount": 0.75}},
		{"saturation-1.00", SATURATION, map[string]any{"amount": 1.00}},
		{"saturation-1.50", SATURATION, map[string]any{"amount": 1.50}},
		{"saturation-2.00", SATURATION, map[string]any{"amount": 2.00}},

		{"threshold--1.00", THRESHOLD, map[string]any{"amount": -1.00}},
		{"threshold--0.75", THRESHOLD, map[string]any{"amount": -0.75}},
		{"threshold--0.50", THRESHOLD, map[string]any{"amount": -0.50}},
		{"threshold--0.25", THRESHOLD, map[string]any{"amount": -0.25}},
		{"threshold-0.00", THRESHOLD, map[string]any{"amount": 0.00}},
		{"threshold-0.25", THRESHOLD, map[string]any{"amount": 0.25}},
		{"threshold-0.50", THRESHOLD, map[string]any{"amount": 0.50}},
		{"threshold-0.75", THRESHOLD, map[string]any{"amount": 0.75}},
		{"threshold-1.00", THRESHOLD, map[string]any{"amount": 1.00}},
		{"threshold-1.50", THRESHOLD, map[string]any{"amount": 1.50}},
		{"threshold-2.00", THRESHOLD, map[string]any{"amount": 2.00}},

		{"vibrance--1.00", VIBRANCE, map[string]any{"amount": -1.00}},
		{"vibrance--0.75", VIBRANCE, map[string]any{"amount": -0.75}},
		{"vibrance--0.50", VIBRANCE, map[string]any{"amount": -0.50}},
		{"vibrance--0.25", VIBRANCE, map[string]any{"amount": -0.25}},
		{"vibrance-0.00", VIBRANCE, map[string]any{"amount": 0.00}},
		{"vibrance-0.25", VIBRANCE, map[string]any{"amount": 0.25}},
		{"vibrance-0.50", VIBRANCE, map[string]any{"amount": 0.50}},
		{"vibrance-0.75", VIBRANCE, map[string]any{"amount": 0.75}},
		{"vibrance-1.00", VIBRANCE, map[string]any{"amount": 1.00}},
		{"vibrance-1.50", VIBRANCE, map[string]any{"amount": 1.50}},
		{"vibrance-2.00", VIBRANCE, map[string]any{"amount": 2.00}},

		{"grayscale", GRAYSCALE, map[string]any{}},
		{"sepia", SEPIA, map[string]any{}},
		{"invert", INVERT, map[string]any{}},
		{"pastelize", PASTELIZE, map[string]any{}},
	}
	for i, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			NewImageFilter(tt.filter, tt.options).Apply(testImage.Clone()).SaveAsPNG(fmt.Sprintf("../test_data/filters/%03d-%s.png", i, tt.name))
		})
	}
}
