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

func Test_Rendering(t *testing.T) {
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

func Test_Blending(t *testing.T) {
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
