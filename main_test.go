package main

import (
	_ "embed"
	"fmt"
	"testing"

	"github.com/toxyl/gfx/color/blend"
	"github.com/toxyl/gfx/color/filter"
	"github.com/toxyl/gfx/color/hsla"
	"github.com/toxyl/gfx/color/rgba"
	"github.com/toxyl/gfx/composition"
	"github.com/toxyl/gfx/coordinates"
	"github.com/toxyl/gfx/filters"
	"github.com/toxyl/gfx/filters/extract"
	"github.com/toxyl/gfx/image"
	"github.com/toxyl/gfx/math"
)

func MakeTestImage(size int,
	hue, hueTolerance, hueFeather,
	sat, satTolerance, satFeather,
	lum, lumTolerance, lumFeather,
	hueShift float64,
) *image.Image {
	var (
		cf              = filter.ToColorFilter(hue, hueTolerance, hueFeather, sat, satTolerance, satFeather, lum, lumTolerance, lumFeather)
		d               = float64(size)
		rc              = float64(size) * 0.85
		w               = float64(size)
		wh              = w / 2
		h               = 20.0
		img             = image.NewWithColor(size, size, *rgba.New(0, 0, 0, 0xFF))
		colMarker       = hsla.New(cf.Col.H()+hueShift, 0.5, 0.75, 1)
		colMarkerMedium = hsla.New(cf.Col.H()+hueShift, 0.5, 0.50, 1)
		colMarkerDark   = hsla.New(cf.Col.H()+hueShift, 0.5, 0.25, 1)
		matchesHue      = func(a, b float64) bool {
			x := int(math.Round(a))
			y := int(math.Round(b))
			return x == y || x+360 == y || x-360 == y
		}
		drawMarker = func(dt, min, max float64) {
			if max <= min {
				max += 360
			}
			for y := h - dt; y <= h+dt; y++ {
				col := colMarker
				if y == h-dt || y == h+dt {
					col = colMarkerDark
				}
				for x := min; x <= max; x++ {
					c := col
					if x == min || x == max {
						c = colMarkerDark
					}
					img.SetHSLA(int((x/360)*(w/2)), int(y), c)
				}
				for x := (min + 360); x <= (max + 360); x++ {
					c := col
					if x == (min+360) || x == (max+360) {
						c = colMarkerDark
					}
					img.SetHSLA(int((x/360)*(w/2)), int(y), c)
				}
				for x := (min - 360); x <= (max - 360); x++ {
					c := col
					if x == (min-360) || x == (max-360) {
						c = colMarkerDark
					}
					img.SetHSLA(int((x/360)*(w/2)), int(y), c)
				}
			}
		}
		drawHueMarker = func(x, y, a, b float64) {
			if matchesHue(a, b) {
				img.SetHSLA(int(x), int(y), colMarker)
				img.SetHSLA(int(x)-1, int(y), colMarkerDark)
				img.SetHSLA(int(x)+1, int(y), colMarkerDark)
			}
		}
	)

	// render hue circle

	for r := 0.0; r < float64(rc)/2; r += 0.1 {
		for hue := 0.0; hue <= 360.0; hue += 0.1 {
			x, y := coordinates.PolarToCartesian(r, hue)
			x += float64(size) / 2
			y += float64(size)/2 + (h * 1.5)
			markerDrawn := false
			if matchesHue(hue, cf.MinThres.H()) {
				img.SetHSLA(int(x), int(y), colMarkerDark)
				markerDrawn = true
			}
			if matchesHue(hue, cf.Min.H()) {
				img.SetHSLA(int(x), int(y), colMarkerMedium)
				markerDrawn = true
			}
			if matchesHue(hue, cf.Col.H()) {
				img.SetHSLA(int(x), int(y), colMarker)
				markerDrawn = true
			}
			if matchesHue(hue, cf.Max.H()) {
				img.SetHSLA(int(x), int(y), colMarkerMedium)
				markerDrawn = true
			}
			if matchesHue(hue, cf.MaxThres.H()) {
				img.SetHSLA(int(x), int(y), colMarkerDark)
				markerDrawn = true
			}
			if !markerDrawn {
				img.SetHSLA(int(x), int(y), hsla.New(hue+hueShift, 1.0, 1-r/(rc/2), 1))
			}
		}
	}

	// render hue bar

	for y := 0.0; y < h; y++ {
		for x := 0.0; x <= w; x++ {
			hue := (x / wh) * 360
			if x > d {
				hue = 360 - hue
			}
			img.SetHSLA(int(x), int(y), hsla.New(hue+hueShift, 1.0, 0.5, 1))
		}
	}

	// render saturation bar

	for y := h; y < 2*h; y++ {
		for x := 0.0; x <= w; x++ {
			sat := (x / w) * 2
			if x >= wh {
				sat = 1 - (sat - 1)
			}
			img.SetHSLA(int(x), int(y), hsla.New(cf.Col.H()+hueShift, sat, 0.5, 1))
		}
	}

	// render luminance  bar

	for y := 2 * h; y < 3*h; y++ {
		for x := 0.0; x <= w; x++ {
			lum := (x / w) * 2
			if x >= wh {
				lum = 1 - (lum - 1)
			}
			img.SetHSLA(int(x), int(y), hsla.New(cf.Col.H()+hueShift, 1.0, lum, 1))
		}
	}

	// render range markers

	for y := h - 10; y < h+10; y++ {
		for x := 0.0; x <= w; x++ {
			hue := (x / wh) * 360
			drawHueMarker(x, y, hue, cf.Col.H())
			drawHueMarker(x, y, hue, cf.MinThres.H())
			drawHueMarker(x, y, hue, cf.Min.H())
			drawHueMarker(x, y, hue, cf.Max.H())
			drawHueMarker(x, y, hue, cf.MaxThres.H())
		}
	}

	drawMarker(3, cf.MinThres.H(), cf.Min.H()) // draw min fade marker
	drawMarker(3, cf.Max.H(), cf.MaxThres.H()) // draw max fade marker
	drawMarker(2, cf.Min.H(), cf.Max.H())      // draw min-max range marker

	// draw text

	img.DrawText(
		fmt.Sprintf(
			"min (thres): %6.2f\nmin:         %6.2f\nmax:         %6.2f\nmax (thres): %6.2f",
			cf.MinThres.H(),
			cf.Min.H(),
			cf.Max.H(),
			cf.MaxThres.H(),
		),
		0, int(h*4), *colMarker, true, blend.NORMAL,
	)

	return img
}

func extractTest(
	img *image.Image, w, h int,
	hue, hueTolerance, hueFeather,
	sat, satTolerance, satFeather,
	lum, lumTolerance, lumFeather float64,
) *image.Image {
	return extract.Extract(img, filter.ToColorFilter(hue, hueTolerance, hueFeather, sat, satTolerance, satFeather, lum, lumTolerance, lumFeather)).Resize(w, h)
}

var (
	fAIAImage = image.NewFromURL("https://sdo.gsfc.nasa.gov/assets/img/latest/f_211_193_171pfss_512.jpg")
)

func renderTestImage(name string, size int, h, hTolerance, hFeather, s, sTolerance, sFeather, l, lTolerance, lFeather float64) {
	var (
		f      = "test_data/main/" + name + ".png"
		fImage = "test_data/main/debug/" + name + ".png"
		fAIA   = "test_data/main/aia/" + name + ".png"
	)
	MakeTestImage(size, h, hTolerance, hFeather, s, sTolerance, sFeather, l, lTolerance, lFeather, 0).SaveAsPNG(fImage)
	extractTest(image.NewFromFile(fImage), size, size, h, hTolerance, hFeather, s, sTolerance, sFeather, l, lTolerance, lFeather).SaveAsPNG(f)
	extractTest(fAIAImage.Clone(), size, size, h, hTolerance, hFeather, s, sTolerance, sFeather, l, lTolerance, lFeather).SaveAsPNG(fAIA)
}

func TestImageRendering(t *testing.T) {
	var (
		SIZE          = 512
		SAT           = 0.65
		SAT_TOLERANCE = 0.35
		SAT_FEATHER   = 0.20
		LUM           = 0.65
		LUM_TOLERANCE = 0.35
		LUM_FEATHER   = 0.20
	)

	testsPerFile := map[string][]struct {
		name      string
		hue       float64
		tolerance float64
		feather   float64
	}{
		"color": {
			{"01", 0, 45, 30},
			{"02", 30, 45, 30},
			{"03", 60, 45, 30},
			{"04", 90, 45, 30},
			{"05", 120, 45, 30},
			{"06", 150, 45, 30},
			{"07", 180, 45, 30},
			{"08", 210, 45, 30},
			{"09", 240, 45, 30},
			{"10", 270, 45, 30},
			{"11", 300, 45, 30},
			{"12", 330, 45, 30},
		},
		"colors": {
			{"01-red", 0, 15, 60},
			{"02-orange", 30, 15, 60},
			{"03-yellow", 60, 15, 60},
			{"04-spring-green", 90, 15, 60},
			{"05-green", 120, 15, 60},
			{"06-turquoise", 150, 15, 60},
			{"07-cyan", 180, 15, 60},
			{"08-ocean", 210, 15, 60},
			{"09-blue", 240, 15, 60},
			{"10-violet", 270, 15, 60},
			{"11-magenta", 300, 15, 60},
			{"12-raspberry", 330, 15, 60},
		},
		"zero-crossing": {
			{"even-01", 0, 22.5, 22.5},
			{"even-02", 0, 45, 22.5},
			{"even-03", 0, 67.5, 22.5},
			{"even-04", 0, 90, 22.5},
			{"even-05", 0, 112.5, 22.5},
			{"even-06", 0, 135, 22.5},
			{"uneven-01", 35, 45, 60},
			{"uneven-02", 80, 90, 60},
			{"uneven-03", 110, 120, 60},
			{"uneven-04", 100, 90, 90},
			{"uneven-05", 130, 120, 90}, // exceeds range
			{"uneven-06", 130, 90, 89},  // in range
		},
		"full-circle": {
			{"not-over-zero", 180, 120, 60},
			{"not-over-zero-2", 180, 150, 30},
			{"not-over-zero-3", 180, 90, 90},
			{"not-over-zero-4", 180, 120, 120},
			{"not-over-zero-5", 180, 180, 120},
			{"over-zero", 0, 120, 60},
			{"over-zero-2", 0, 150, 30},
			{"over-zero-3", 0, 90, 90},
			{"over-zero-4", 0, 120, 120},
			{"over-zero-5", 0, 180, 120},
			{"fade-before-zero", 120, 120, 60},
			{"fade-before-zero-2", 120, 150, 30},
			{"fade-before-zero-3", 120, 90, 90},
			{"fade-before-zero-4", 120, 120, 120},
			{"fade-before-zero-5", 120, 180, 120},
			{"fade-after-zero", 240, 120, 60},
			{"fade-after-zero-2", 240, 150, 30},
			{"fade-after-zero-3", 240, 90, 90},
			{"fade-after-zero-4", 240, 120, 120},
			{"fade-after-zero-5", 240, 180, 120},
		},
	}
	for fileSet, tests := range testsPerFile {
		for _, tt := range tests {
			t.Run(tt.name, func(t *testing.T) {
				renderTestImage(fileSet+"-"+tt.name, SIZE,
					tt.hue, tt.tolerance, tt.feather,
					SAT, SAT_TOLERANCE, SAT_FEATHER,
					LUM, LUM_TOLERANCE, LUM_FEATHER,
				)
			})
		}
	}
}

func TestImageBlending(t *testing.T) {
	tstImg1 := MakeTestImage(512, 0, 90, 0, 0.5, 0.5, 0.0, 0.5, 0.5, 0.0, 0)
	tstImg2 := MakeTestImage(512, 0, 90, 0, 0.5, 0.5, 0.0, 0.5, 0.5, 0.0, 180)

	tests := []blend.BlendMode{
		blend.NORMAL,
		//
		blend.DARKEN,
		blend.MULTIPLY,
		blend.COLOR_BURN,
		blend.LINEAR_BURN,
		//
		blend.LIGHTEN,
		blend.SCREEN,
		blend.ADD,
		//
		blend.OVERLAY,
		blend.SOFT_LIGHT,
		blend.HARD_LIGHT,
		blend.PIN_LIGHT,
		//
		blend.DIFFERENCE,
		blend.EXCLUSION,
		blend.SUBTRACT,
		blend.DIVIDE,
		//
		blend.AVERAGE,
		blend.NEGATION,
		//
		blend.ERASE,
	}

	for _, tt := range tests {
		cmpBlend := composition.New(
			512, 512,
			composition.NewLayerFromImage(tstImg1, 1.00, blend.NORMAL, nil),
			composition.NewLayerFromImage(tstImg2, 1.00, tt, nil),
		)
		cmpBlend.Render().SaveAsPNG("test_data/blendmode/" + string(tt) + ".png")
	}
}

func TestComposition(t *testing.T) {
	tests := []struct {
		name string
		path string
	}{
		{"1", "test_data/compositions/1.yaml"},
		{"2", "test_data/compositions/2.yaml"},
		{"3", "test_data/compositions/3.yaml"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			composition.Load(tt.path).Render().SaveAsPNG("test_data/compositions/render/" + tt.name + ".png")
		})
	}
}

func TestTextDrawing(t *testing.T) {
	var (
		fontColor = hsla.New(30, 1.0, 0.5, 1.0)
	)
	type test struct {
		name string
		text string
		x    int
		y    int
		col  *hsla.HSLA
		glow bool
		mode blend.BlendMode
	}
	tests := []test{
		{"Normal", "Hello World", 1, 2, fontColor, false, blend.NORMAL},
		{"Add", "Hello World", 1, 2, fontColor, false, blend.ADD},
		{"Darken", "Hello World", 1, 2, fontColor, false, blend.DARKEN},
		{"Overlay", "Hello World", 1, 2, fontColor, false, blend.OVERLAY},
		{"Screen", "Hello World", 1, 2, fontColor, false, blend.SCREEN},
		{"Exclusion", "Hello World", 1, 2, fontColor, false, blend.EXCLUSION},
		{"Lighten", "Hello World", 1, 2, fontColor, false, blend.LIGHTEN},
		{"Multiply", "Hello World", 1, 2, fontColor, false, blend.MULTIPLY},
		{"Negation", "Hello World", 1, 2, fontColor, false, blend.NEGATION},
		{"Average", "Hello World", 1, 2, fontColor, false, blend.AVERAGE},
		{"Color Burn", "Hello World", 1, 2, fontColor, false, blend.COLOR_BURN},
		{"Difference", "Hello World", 1, 2, fontColor, false, blend.DIFFERENCE},
		{"Divide", "Hello World", 1, 2, fontColor, false, blend.DIVIDE},
		{"Hard Light", "Hello World", 1, 2, fontColor, false, blend.HARD_LIGHT},
		{"Linear Burn", "Hello World", 1, 2, fontColor, false, blend.LINEAR_BURN},
		{"Pin Light", "Hello World", 1, 2, fontColor, false, blend.PIN_LIGHT},
		{"Soft Light", "Hello World", 1, 2, fontColor, false, blend.SOFT_LIGHT},
		{"Subtract", "Hello World", 1, 2, fontColor, false, blend.SUBTRACT},

		{"Normal - Glow", "Hello World", 1, 2, fontColor, true, blend.NORMAL},
		{"Add - Glow", "Hello World", 1, 2, fontColor, true, blend.ADD},
		{"Darken - Glow", "Hello World", 1, 2, fontColor, true, blend.DARKEN},
		{"Overlay - Glow", "Hello World", 1, 2, fontColor, true, blend.OVERLAY},
		{"Screen - Glow", "Hello World", 1, 2, fontColor, true, blend.SCREEN},
		{"Exclusion - Glow", "Hello World", 1, 2, fontColor, true, blend.EXCLUSION},
		{"Lighten - Glow", "Hello World", 1, 2, fontColor, true, blend.LIGHTEN},
		{"Multiply - Glow", "Hello World", 1, 2, fontColor, true, blend.MULTIPLY},
		{"Negation - Glow", "Hello World", 1, 2, fontColor, true, blend.NEGATION},
		{"Average - Glow", "Hello World", 1, 2, fontColor, true, blend.AVERAGE},
		{"Color Burn - Glow", "Hello World", 1, 2, fontColor, true, blend.COLOR_BURN},
		{"Difference - Glow", "Hello World", 1, 2, fontColor, true, blend.DIFFERENCE},
		{"Divide - Glow", "Hello World", 1, 2, fontColor, true, blend.DIVIDE},
		{"Hard Light - Glow", "Hello World", 1, 2, fontColor, true, blend.HARD_LIGHT},
		{"Linear Burn - Glow", "Hello World", 1, 2, fontColor, true, blend.LINEAR_BURN},
		{"Pin Light - Glow", "Hello World", 1, 2, fontColor, true, blend.PIN_LIGHT},
		{"Soft Light - Glow", "Hello World", 1, 2, fontColor, true, blend.SOFT_LIGHT},
		{"Subtract - Glow", "Hello World", 1, 2, fontColor, true, blend.SUBTRACT},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			i := image.NewWithColor(101, 13, *rgba.New(0x00, 0x33, 0x33, 0xFF))
			i.DrawText(tt.text, tt.x, tt.y, *tt.col, tt.glow, tt.mode).SaveAsPNG("test_data/text/" + tt.name + ".png")
		})
	}
}

func TestFilters(t *testing.T) {
	var (
		testImage    = image.NewFromURL("https://sdo.gsfc.nasa.gov/assets/img/latest/f_211_193_171pfss_512.jpg")
		customFilter = []float64{
			0.025, 0.100, 0.025,
			0.050, 0.150, 0.050,
			0.075, 0.100, 0.075,
		}
		noArgs    = map[string]any{}
		amtNeg100 = map[string]any{"amount": -1.00}
	)
	type filterConf struct {
		name    string
		filter  string
		options map[string]any
	}
	testGroups := map[string][]filterConf{
		"colorize": {
			{"grayscale", filters.Gray, noArgs},
			{"invert", filters.Invert, noArgs},
			{"pastelize", filters.Pastelize, noArgs},
			{"sepia", filters.Sepia, noArgs},
		},
		"hue": {
			{"-1.00", filters.Hue, amtNeg100},
			{"-0.75", filters.Hue, map[string]any{"amount": -0.75}},
			{"-0.50", filters.Hue, map[string]any{"amount": -0.50}},
			{"-0.25", filters.Hue, map[string]any{"amount": -0.25}},
			{"0.00", filters.Hue, map[string]any{"amount": 0.00}},
			{"0.25", filters.Hue, map[string]any{"amount": 0.25}},
			{"0.50", filters.Hue, map[string]any{"amount": 0.50}},
			{"0.75", filters.Hue, map[string]any{"amount": 0.75}},
			{"1.00", filters.Hue, map[string]any{"amount": 1.00}},
			{"1.50", filters.Hue, map[string]any{"amount": 1.50}},
			{"2.00", filters.Hue, map[string]any{"amount": 2.00}},
		},
		"sat": {
			{"-1.00", filters.Sat, amtNeg100},
			{"-0.75", filters.Sat, map[string]any{"amount": -0.75}},
			{"-0.50", filters.Sat, map[string]any{"amount": -0.50}},
			{"-0.25", filters.Sat, map[string]any{"amount": -0.25}},
			{"0.00", filters.Sat, map[string]any{"amount": 0.00}},
			{"0.25", filters.Sat, map[string]any{"amount": 0.25}},
			{"0.50", filters.Sat, map[string]any{"amount": 0.50}},
			{"0.75", filters.Sat, map[string]any{"amount": 0.75}},
			{"1.00", filters.Sat, map[string]any{"amount": 1.00}},
			{"1.50", filters.Sat, map[string]any{"amount": 1.50}},
			{"2.00", filters.Sat, map[string]any{"amount": 2.00}},
		},
		"lum": {
			{"-1.00", filters.Lum, amtNeg100},
			{"-0.75", filters.Lum, map[string]any{"amount": -0.75}},
			{"-0.50", filters.Lum, map[string]any{"amount": -0.50}},
			{"-0.25", filters.Lum, map[string]any{"amount": -0.25}},
			{"0.00", filters.Lum, map[string]any{"amount": 0.00}},
			{"0.25", filters.Lum, map[string]any{"amount": 0.25}},
			{"0.50", filters.Lum, map[string]any{"amount": 0.50}},
			{"0.75", filters.Lum, map[string]any{"amount": 0.75}},
			{"1.00", filters.Lum, map[string]any{"amount": 1.00}},
			{"1.50", filters.Lum, map[string]any{"amount": 1.50}},
			{"2.00", filters.Lum, map[string]any{"amount": 2.00}},
		},
		"hue-contrast": {
			{"-1.00", filters.HueContrast, amtNeg100},
			{"-0.75", filters.HueContrast, map[string]any{"amount": -0.75}},
			{"-0.50", filters.HueContrast, map[string]any{"amount": -0.50}},
			{"-0.25", filters.HueContrast, map[string]any{"amount": -0.25}},
			{"0.00", filters.HueContrast, map[string]any{"amount": 0.00}},
			{"0.25", filters.HueContrast, map[string]any{"amount": 0.25}},
			{"0.50", filters.HueContrast, map[string]any{"amount": 0.50}},
			{"0.75", filters.HueContrast, map[string]any{"amount": 0.75}},
			{"1.00", filters.HueContrast, map[string]any{"amount": 1.00}},
			{"1.50", filters.HueContrast, map[string]any{"amount": 1.50}},
			{"2.00", filters.HueContrast, map[string]any{"amount": 2.00}},
		},
		"sat-contrast": {
			{"-1.00", filters.SatContrast, amtNeg100},
			{"-0.75", filters.SatContrast, map[string]any{"amount": -0.75}},
			{"-0.50", filters.SatContrast, map[string]any{"amount": -0.50}},
			{"-0.25", filters.SatContrast, map[string]any{"amount": -0.25}},
			{"0.00", filters.SatContrast, map[string]any{"amount": 0.00}},
			{"0.25", filters.SatContrast, map[string]any{"amount": 0.25}},
			{"0.50", filters.SatContrast, map[string]any{"amount": 0.50}},
			{"0.75", filters.SatContrast, map[string]any{"amount": 0.75}},
			{"1.00", filters.SatContrast, map[string]any{"amount": 1.00}},
			{"1.50", filters.SatContrast, map[string]any{"amount": 1.50}},
			{"2.00", filters.SatContrast, map[string]any{"amount": 2.00}},
		},
		"lum-contrast": {
			{"-1.00", filters.LumContrast, amtNeg100},
			{"-0.75", filters.LumContrast, map[string]any{"amount": -0.75}},
			{"-0.50", filters.LumContrast, map[string]any{"amount": -0.50}},
			{"-0.25", filters.LumContrast, map[string]any{"amount": -0.25}},
			{"0.00", filters.LumContrast, map[string]any{"amount": 0.00}},
			{"0.25", filters.LumContrast, map[string]any{"amount": 0.25}},
			{"0.50", filters.LumContrast, map[string]any{"amount": 0.50}},
			{"0.75", filters.LumContrast, map[string]any{"amount": 0.75}},
			{"1.00", filters.LumContrast, map[string]any{"amount": 1.00}},
			{"1.50", filters.LumContrast, map[string]any{"amount": 1.50}},
			{"2.00", filters.LumContrast, map[string]any{"amount": 2.00}},
		},
		"color-shift": {
			{"-1.00", filters.ColorShift, map[string]any{"hue": -1.00, "sat": -1.00, "lum": -1.00}},
			{"-0.75", filters.ColorShift, map[string]any{"hue": -0.75, "sat": -0.75, "lum": -0.75}},
			{"-0.50", filters.ColorShift, map[string]any{"hue": -0.50, "sat": -0.50, "lum": -0.50}},
			{"-0.25", filters.ColorShift, map[string]any{"hue": -0.25, "sat": -0.25, "lum": -0.25}},
			{"0.00", filters.ColorShift, map[string]any{"hue": 0.00, "sat": 0.00, "lum": 0.00}},
			{"0.25", filters.ColorShift, map[string]any{"hue": 0.25, "sat": 0.25, "lum": 0.25}},
			{"0.50", filters.ColorShift, map[string]any{"hue": 0.50, "sat": 0.50, "lum": 0.50}},
			{"0.75", filters.ColorShift, map[string]any{"hue": 0.75, "sat": 0.75, "lum": 0.75}},
			{"1.00", filters.ColorShift, map[string]any{"hue": 1.00, "sat": 1.00, "lum": 1.00}},
			{"1.50", filters.ColorShift, map[string]any{"hue": 1.50, "sat": 1.50, "lum": 1.50}},
			{"2.00", filters.ColorShift, map[string]any{"hue": 2.00, "sat": 2.00, "lum": 2.00}},
		},
		"brightness": {
			{"-1.00", filters.Brightness, amtNeg100},
			{"-0.75", filters.Brightness, map[string]any{"amount": -0.75}},
			{"-0.50", filters.Brightness, map[string]any{"amount": -0.50}},
			{"-0.25", filters.Brightness, map[string]any{"amount": -0.25}},
			{"0.00", filters.Brightness, map[string]any{"amount": 0.00}},
			{"0.25", filters.Brightness, map[string]any{"amount": 0.25}},
			{"0.50", filters.Brightness, map[string]any{"amount": 0.50}},
			{"0.75", filters.Brightness, map[string]any{"amount": 0.75}},
			{"1.00", filters.Brightness, map[string]any{"amount": 1.00}},
			{"1.50", filters.Brightness, map[string]any{"amount": 1.50}},
			{"2.00", filters.Brightness, map[string]any{"amount": 2.00}},
		},
		"contrast": {
			{"-1.00", filters.Contrast, amtNeg100},
			{"-0.75", filters.Contrast, map[string]any{"amount": -0.75}},
			{"-0.50", filters.Contrast, map[string]any{"amount": -0.50}},
			{"-0.25", filters.Contrast, map[string]any{"amount": -0.25}},
			{"0.00", filters.Contrast, map[string]any{"amount": 0.00}},
			{"0.25", filters.Contrast, map[string]any{"amount": 0.25}},
			{"0.50", filters.Contrast, map[string]any{"amount": 0.50}},
			{"0.75", filters.Contrast, map[string]any{"amount": 0.75}},
			{"1.00", filters.Contrast, map[string]any{"amount": 1.00}},
			{"1.25", filters.Contrast, map[string]any{"amount": 1.25}},
			{"1.50", filters.Contrast, map[string]any{"amount": 1.50}},
			{"2.00", filters.Contrast, map[string]any{"amount": 2.00}},
			{"4.00", filters.Contrast, map[string]any{"amount": 4.00}},
			{"8.00", filters.Contrast, map[string]any{"amount": 8.00}},
			{"16.00", filters.Contrast, map[string]any{"amount": 16.00}},
			{"32.00", filters.Contrast, map[string]any{"amount": 32.00}},
			{"64.00", filters.Contrast, map[string]any{"amount": 64.00}},
			{"128.00", filters.Contrast, map[string]any{"amount": 128.00}},
		},
		"gamma": {
			{"-1.00", filters.Gamma, amtNeg100},
			{"-0.75", filters.Gamma, map[string]any{"amount": -0.75}},
			{"-0.50", filters.Gamma, map[string]any{"amount": -0.50}},
			{"-0.25", filters.Gamma, map[string]any{"amount": -0.25}},
			{"0.00", filters.Gamma, map[string]any{"amount": 0.00}},
			{"0.25", filters.Gamma, map[string]any{"amount": 0.25}},
			{"0.50", filters.Gamma, map[string]any{"amount": 0.50}},
			{"0.75", filters.Gamma, map[string]any{"amount": 0.75}},
			{"1.00", filters.Gamma, map[string]any{"amount": 1.00}},
			{"1.50", filters.Gamma, map[string]any{"amount": 1.50}},
			{"2.00", filters.Gamma, map[string]any{"amount": 2.00}},
		},
		"vibrance": {
			{"-1.00", filters.Vibrance, amtNeg100},
			{"-0.75", filters.Vibrance, map[string]any{"amount": -0.75}},
			{"-0.50", filters.Vibrance, map[string]any{"amount": -0.50}},
			{"-0.25", filters.Vibrance, map[string]any{"amount": -0.25}},
			{"0.00", filters.Vibrance, map[string]any{"amount": 0.00}},
			{"0.25", filters.Vibrance, map[string]any{"amount": 0.25}},
			{"0.50", filters.Vibrance, map[string]any{"amount": 0.50}},
			{"0.75", filters.Vibrance, map[string]any{"amount": 0.75}},
			{"1.00", filters.Vibrance, map[string]any{"amount": 1.00}},
			{"1.50", filters.Vibrance, map[string]any{"amount": 1.50}},
			{"2.00", filters.Vibrance, map[string]any{"amount": 2.00}},
		},
		"enhance": {
			{"-1.00", filters.Enhance, amtNeg100},
			{"-0.75", filters.Enhance, map[string]any{"amount": -0.75}},
			{"-0.50", filters.Enhance, map[string]any{"amount": -0.50}},
			{"-0.25", filters.Enhance, map[string]any{"amount": -0.25}},
			{"0.00", filters.Enhance, map[string]any{"amount": 0.00}},
			{"0.25", filters.Enhance, map[string]any{"amount": 0.25}},
			{"0.50", filters.Enhance, map[string]any{"amount": 0.50}},
			{"0.75", filters.Enhance, map[string]any{"amount": 0.75}},
			{"1.00", filters.Enhance, map[string]any{"amount": 1.00}},
			{"1.50", filters.Enhance, map[string]any{"amount": 1.50}},
			{"2.00", filters.Enhance, map[string]any{"amount": 2.00}},
		},
		"sharpen": {
			{"-1.00", filters.Sharpen, amtNeg100},
			{"-0.75", filters.Sharpen, map[string]any{"amount": -0.75}},
			{"-0.50", filters.Sharpen, map[string]any{"amount": -0.50}},
			{"-0.25", filters.Sharpen, map[string]any{"amount": -0.25}},
			{"0.00", filters.Sharpen, map[string]any{"amount": 0.00}},
			{"0.25", filters.Sharpen, map[string]any{"amount": 0.25}},
			{"0.50", filters.Sharpen, map[string]any{"amount": 0.50}},
			{"0.75", filters.Sharpen, map[string]any{"amount": 0.75}},
			{"1.00", filters.Sharpen, map[string]any{"amount": 1.00}},
			{"1.50", filters.Sharpen, map[string]any{"amount": 1.50}},
			{"2.00", filters.Sharpen, map[string]any{"amount": 2.00}},
		},
		"blur": {
			{"-1.00", filters.Blur, amtNeg100},
			{"-0.75", filters.Blur, map[string]any{"amount": -0.75}},
			{"-0.50", filters.Blur, map[string]any{"amount": -0.50}},
			{"-0.25", filters.Blur, map[string]any{"amount": -0.25}},
			{"0.00", filters.Blur, map[string]any{"amount": 0.00}},
			{"0.25", filters.Blur, map[string]any{"amount": 0.25}},
			{"0.50", filters.Blur, map[string]any{"amount": 0.50}},
			{"0.75", filters.Blur, map[string]any{"amount": 0.75}},
			{"1.00", filters.Blur, map[string]any{"amount": 1.00}},
			{"1.50", filters.Blur, map[string]any{"amount": 1.50}},
			{"2.00", filters.Blur, map[string]any{"amount": 2.00}},
		},
		"edge-detection": {
			{"-1.00", filters.EdgeDetect, amtNeg100},
			{"-0.75", filters.EdgeDetect, map[string]any{"amount": -0.75}},
			{"-0.50", filters.EdgeDetect, map[string]any{"amount": -0.50}},
			{"-0.25", filters.EdgeDetect, map[string]any{"amount": -0.25}},
			{"0.00", filters.EdgeDetect, map[string]any{"amount": 0.00}},
			{"0.25", filters.EdgeDetect, map[string]any{"amount": 0.25}},
			{"0.50", filters.EdgeDetect, map[string]any{"amount": 0.50}},
			{"0.75", filters.EdgeDetect, map[string]any{"amount": 0.75}},
			{"1.00", filters.EdgeDetect, map[string]any{"amount": 1.00}},
			{"1.50", filters.EdgeDetect, map[string]any{"amount": 1.50}},
			{"2.00", filters.EdgeDetect, map[string]any{"amount": 2.00}},
		},
		"emboss": {
			{"-1.00", filters.Emboss, amtNeg100},
			{"-0.75", filters.Emboss, map[string]any{"amount": -0.75}},
			{"-0.50", filters.Emboss, map[string]any{"amount": -0.50}},
			{"-0.25", filters.Emboss, map[string]any{"amount": -0.25}},
			{"0.00", filters.Emboss, map[string]any{"amount": 0.00}},
			{"0.25", filters.Emboss, map[string]any{"amount": 0.25}},
			{"0.50", filters.Emboss, map[string]any{"amount": 0.50}},
			{"0.75", filters.Emboss, map[string]any{"amount": 0.75}},
			{"1.00", filters.Emboss, map[string]any{"amount": 1.00}},
			{"1.50", filters.Emboss, map[string]any{"amount": 1.50}},
			{"2.00", filters.Emboss, map[string]any{"amount": 2.00}},
		},
		"threshold": {
			{"-1.00", filters.Threshold, amtNeg100},
			{"-0.75", filters.Threshold, map[string]any{"amount": -0.75}},
			{"-0.50", filters.Threshold, map[string]any{"amount": -0.50}},
			{"-0.25", filters.Threshold, map[string]any{"amount": -0.25}},
			{"0.00", filters.Threshold, map[string]any{"amount": 0.00}},
			{"0.25", filters.Threshold, map[string]any{"amount": 0.25}},
			{"0.50", filters.Threshold, map[string]any{"amount": 0.50}},
			{"0.75", filters.Threshold, map[string]any{"amount": 0.75}},
			{"1.00", filters.Threshold, map[string]any{"amount": 1.00}},
			{"1.50", filters.Threshold, map[string]any{"amount": 1.50}},
			{"2.00", filters.Threshold, map[string]any{"amount": 2.00}},
		},
		"alpha-map-s": {
			{"0.00", filters.AlphaMap, map[string]any{"source": "s", "lower": 0.00, "upper": 1.0}},
			{"0.25", filters.AlphaMap, map[string]any{"source": "s", "lower": 0.25, "upper": 1.0}},
			{"0.50", filters.AlphaMap, map[string]any{"source": "s", "lower": 0.50, "upper": 1.0}},
			{"0.75", filters.AlphaMap, map[string]any{"source": "s", "lower": 0.75, "upper": 1.0}},
			{"1.00", filters.AlphaMap, map[string]any{"source": "s", "lower": 1.00, "upper": 1.0}},
		},
		"alpha-map-l": {
			{"0.00", filters.AlphaMap, map[string]any{"source": "l", "lower": 0.00, "upper": 1.0}},
			{"0.25", filters.AlphaMap, map[string]any{"source": "l", "lower": 0.25, "upper": 1.0}},
			{"0.50", filters.AlphaMap, map[string]any{"source": "l", "lower": 0.50, "upper": 1.0}},
			{"0.75", filters.AlphaMap, map[string]any{"source": "l", "lower": 0.75, "upper": 1.0}},
			{"1.00", filters.AlphaMap, map[string]any{"source": "l", "lower": 1.00, "upper": 1.0}},
		},
		"alpha-map-sl": {
			{"0.00", filters.AlphaMap, map[string]any{"source": "s*l", "lower": 0.00, "upper": 1.00}},
			{"0.25", filters.AlphaMap, map[string]any{"source": "s*l", "lower": 0.25, "upper": 1.00}},
			{"0.50", filters.AlphaMap, map[string]any{"source": "s*l", "lower": 0.50, "upper": 1.00}},
			{"0.75", filters.AlphaMap, map[string]any{"source": "s*l", "lower": 0.75, "upper": 1.00}},
			{"1.00", filters.AlphaMap, map[string]any{"source": "s*l", "lower": 1.00, "upper": 1.00}},
		},
		"alpha-map-s-inv": {
			{"0.00", filters.AlphaMap, map[string]any{"source": "s", "lower": 1.00, "upper": 0.00}},
			{"0.25", filters.AlphaMap, map[string]any{"source": "s", "lower": 1.00, "upper": 0.25}},
			{"0.50", filters.AlphaMap, map[string]any{"source": "s", "lower": 1.00, "upper": 0.50}},
			{"0.75", filters.AlphaMap, map[string]any{"source": "s", "lower": 1.00, "upper": 0.75}},
			{"1.00", filters.AlphaMap, map[string]any{"source": "s", "lower": 1.00, "upper": 1.00}},
		},
		"alpha-map-l-inv": {
			{"0.00", filters.AlphaMap, map[string]any{"source": "l", "lower": 1.00, "upper": 0.00}},
			{"0.25", filters.AlphaMap, map[string]any{"source": "l", "lower": 1.00, "upper": 0.25}},
			{"0.50", filters.AlphaMap, map[string]any{"source": "l", "lower": 1.00, "upper": 0.50}},
			{"0.75", filters.AlphaMap, map[string]any{"source": "l", "lower": 1.00, "upper": 0.75}},
			{"1.00", filters.AlphaMap, map[string]any{"source": "l", "lower": 1.00, "upper": 1.00}},
		},
		"alpha-map-sl-inv": {
			{"0.00", filters.AlphaMap, map[string]any{"source": "s*l", "lower": 1.00, "upper": 0.00}},
			{"0.25", filters.AlphaMap, map[string]any{"source": "s*l", "lower": 1.00, "upper": 0.25}},
			{"0.50", filters.AlphaMap, map[string]any{"source": "s*l", "lower": 1.00, "upper": 0.50}},
			{"0.75", filters.AlphaMap, map[string]any{"source": "s*l", "lower": 1.00, "upper": 0.75}},
			{"1.00", filters.AlphaMap, map[string]any{"source": "s*l", "lower": 1.00, "upper": 1.00}},
		},
		"convolution": {
			{"-1.00", filters.Convolution, map[string]any{"matrix": customFilter, "bias": 0.00, "factor": -1.00}},
			{"-0.75", filters.Convolution, map[string]any{"matrix": customFilter, "bias": 0.00, "factor": -0.75}},
			{"-0.50", filters.Convolution, map[string]any{"matrix": customFilter, "bias": 0.00, "factor": -0.50}},
			{"-0.25", filters.Convolution, map[string]any{"matrix": customFilter, "bias": 0.00, "factor": -0.25}},
			{"0.00", filters.Convolution, map[string]any{"matrix": customFilter, "bias": 0.00, "factor": 0.00}},
			{"0.25", filters.Convolution, map[string]any{"matrix": customFilter, "bias": 0.00, "factor": 0.25}},
			{"0.50", filters.Convolution, map[string]any{"matrix": customFilter, "bias": 0.00, "factor": 0.50}},
			{"0.75", filters.Convolution, map[string]any{"matrix": customFilter, "bias": 0.00, "factor": 0.75}},
			{"1.00", filters.Convolution, map[string]any{"matrix": customFilter, "bias": 0.00, "factor": 1.00}},
		},
	}
	for k, tests := range testGroups {
		for i, tt := range tests {
			t.Run(tt.name, func(t *testing.T) {
				res := filters.NewImageFilter(tt.filter, tt.options).Apply(testImage.Clone())
				res.SaveAsPNG(fmt.Sprintf("test_data/filters/%s-%03d-%s.png", k, i, tt.name))
			})
		}
	}
}

func TestColorBlending(t *testing.T) {
	tests := []struct {
		name string
		c1   *rgba.RGBA
		c2   *rgba.RGBA
	}{
		{"b&w1", rgba.New(0, 0, 0, 0), rgba.New(0xFF, 0xFF, 0xFF, 0xFF)},
		{"b&w2", rgba.New(0, 0, 0, 0), rgba.New(0xFF, 0xFF, 0xFF, 0x7F)},
		{"b&w3", rgba.New(0, 0, 0, 0xFF), rgba.New(0xFF, 0xFF, 0xFF, 0x00)},
		{"b&w4", rgba.New(0, 0, 0, 0xFF), rgba.New(0xFF, 0xFF, 0xFF, 0x7F)},
		{"w&w", rgba.New(0xFF, 0xFF, 0xFF, 0xFF), rgba.New(0xFF, 0xFF, 0xFF, 0x7F)},
		{"b&b", rgba.New(0x00, 0x00, 0x00, 0xFF), rgba.New(0x00, 0x00, 0x00, 0x7F)},
		{"r&g", rgba.New(0xFF, 0x00, 0x00, 0x7F), rgba.New(0x00, 0xFF, 0x00, 0x7F)},
		{"g&r", rgba.New(0x00, 0xFF, 0x00, 0x7F), rgba.New(0xFF, 0x00, 0x00, 0x7F)},
		{"w&r", rgba.New(0xFF, 0xFF, 0xFF, 0xFF), rgba.New(0xFF, 0x00, 0x00, 0x7F)},
		{"w&r2", rgba.New(0xFF, 0xFF, 0xFF, 0x01), rgba.New(0xFF, 0x00, 0x00, 0x01)},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			fmt.Printf("blend.RGBA({%s}, {%s}, normal, 1) = %s\n", tt.c1, tt.c2, blend.RGBA(tt.c1, tt.c2, blend.NORMAL, 1.0))
		})
	}
}
