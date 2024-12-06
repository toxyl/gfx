package filters

import (
	"fmt"
	"strings"

	"github.com/toxyl/gfx/color/filter"
	"github.com/toxyl/gfx/filters/alphamap"
	"github.com/toxyl/gfx/filters/blur"
	"github.com/toxyl/gfx/filters/brightness"
	"github.com/toxyl/gfx/filters/colorshift"
	"github.com/toxyl/gfx/filters/contrast"
	"github.com/toxyl/gfx/filters/convolution"
	"github.com/toxyl/gfx/filters/edgedetect"
	"github.com/toxyl/gfx/filters/emboss"
	"github.com/toxyl/gfx/filters/enhance"
	"github.com/toxyl/gfx/filters/extract"
	"github.com/toxyl/gfx/filters/gamma"
	"github.com/toxyl/gfx/filters/gray"
	"github.com/toxyl/gfx/filters/hue"
	"github.com/toxyl/gfx/filters/huecontrast"
	"github.com/toxyl/gfx/filters/invert"
	"github.com/toxyl/gfx/filters/lum"
	"github.com/toxyl/gfx/filters/lumcontrast"
	"github.com/toxyl/gfx/filters/pastelize"
	"github.com/toxyl/gfx/filters/sat"
	"github.com/toxyl/gfx/filters/satcontrast"
	"github.com/toxyl/gfx/filters/sepia"
	"github.com/toxyl/gfx/filters/sharpen"
	"github.com/toxyl/gfx/filters/threshold"
	"github.com/toxyl/gfx/filters/vibrance"
	"github.com/toxyl/gfx/image"
	"github.com/toxyl/gfx/math"
)

const (
	Gray        = "gray"
	Invert      = "invert"
	Pastelize   = "pastelize"
	Sepia       = "sepia"
	Hue         = "hue"
	Sat         = "sat"
	Lum         = "lum"
	HueContrast = "hue-contrast"
	SatContrast = "sat-contrast"
	LumContrast = "lum-contrast"
	ColorShift  = "color-shift"
	Brightness  = "brightness"
	Contrast    = "contrast"
	Gamma       = "gamma"
	Vibrance    = "vibrance"
	Enhance     = "enhance"
	Sharpen     = "sharpen"
	Blur        = "blur"
	EdgeDetect  = "edge-detect"
	Emboss      = "emboss"
	Threshold   = "threshold"
	AlphaMap    = "alpha-map"
	Extract     = "extract"
	Convolution = "convolution"
	optAmt      = "amount"
	optFct      = "factor"
	optBis      = "bias"
	optHue      = "hue"
	optSat      = "sat"
	optLum      = "lum"
	optHueTlr   = "hue-tolerance"
	optSatTlr   = "sat-tolerance"
	optLumTlr   = "lum-tolerance"
	optHueFth   = "hue-feather"
	optSatFth   = "sat-feather"
	optLumFth   = "lum-feather"
	optLwr      = "lower"
	optUpp      = "upper"
	optSrc      = "source"
	optMtr      = "matrix"
)

var (
	defMtr = [][]float64{
		{1.0, 1.0, 1.0},
		{1.0, 8.0, 1.0},
		{1.0, 1.0, 1.0},
	}
	Examples = []string{
		Gray,
		Invert,
		Pastelize,
		Sepia,
		Hue + "::" + optAmt + "=0.5",
		Sat + "::" + optAmt + "=1.0",
		Lum + "::" + optAmt + "=1.0",
		HueContrast + "::" + optAmt + "=1.0",
		SatContrast + "::" + optAmt + "=1.0",
		LumContrast + "::" + optAmt + "=1.0",
		ColorShift + "::" + optHue + "=180.0::" + optSat + "=0.1::" + optLum + "=0.7",
		Brightness + "::" + optAmt + "=1.0",
		Contrast + "::" + optAmt + "=1.0",
		Gamma + "::" + optAmt + "=1.0",
		Vibrance + "::" + optAmt + "=1.0",
		Enhance + "::" + optAmt + "=1.0",
		Sharpen + "::" + optAmt + "=1.0",
		Blur + "::" + optAmt + "=1.0",
		EdgeDetect + "::" + optAmt + "=1.0",
		Emboss + "::" + optAmt + "=1.0",
		Threshold + "::" + optAmt + "=1.0",
		AlphaMap + "::" + optSrc + "=s*l::" + optLwr + "=0.1::" + optUpp + "=0.7",
		Extract + "::" +
			optHue + "=180.0::" + optHueTlr + "=90.0::" + optHueFth + "=90.0::" +
			optSat + "=0.50::" + optSatTlr + "=0.25::" + optSatFth + "=0.25::" +
			optLum + "=0.50::" + optLumTlr + "=0.25::" + optLumFth + "=0.25",
		Convolution + "::" + optMtr + "=1.0,1.0,1.0,1.0,8.0,1.0,1.0,1.0,1.0",
	}
	filterMap = map[string]func(s *ImageFilter, i *image.Image){
		Gray:        func(s *ImageFilter, i *image.Image) { gray.Apply(i) },
		Invert:      func(s *ImageFilter, i *image.Image) { invert.Apply(i) },
		Pastelize:   func(s *ImageFilter, i *image.Image) { pastelize.Apply(i) },
		Sepia:       func(s *ImageFilter, i *image.Image) { sepia.Apply(i) },
		Hue:         func(s *ImageFilter, i *image.Image) { hue.Apply(i, s.getAmt()) },
		Sat:         func(s *ImageFilter, i *image.Image) { sat.Apply(i, s.getAmt()) },
		Lum:         func(s *ImageFilter, i *image.Image) { lum.Apply(i, s.getAmt()) },
		HueContrast: func(s *ImageFilter, i *image.Image) { huecontrast.Apply(i, s.getAmt()) },
		SatContrast: func(s *ImageFilter, i *image.Image) { satcontrast.Apply(i, s.getAmt()) },
		LumContrast: func(s *ImageFilter, i *image.Image) { lumcontrast.Apply(i, s.getAmt()) },
		ColorShift:  func(s *ImageFilter, i *image.Image) { colorshift.Apply(i, s.getH(), s.getS(), s.getL()) },
		Brightness:  func(s *ImageFilter, i *image.Image) { brightness.Apply(i, s.getAmt()) },
		Contrast:    func(s *ImageFilter, i *image.Image) { contrast.Apply(i, s.getAmt()) },
		Gamma:       func(s *ImageFilter, i *image.Image) { gamma.Apply(i, s.getAmt()) },
		Vibrance:    func(s *ImageFilter, i *image.Image) { vibrance.Apply(i, s.getAmt()) },
		Enhance:     func(s *ImageFilter, i *image.Image) { enhance.Apply(i, s.getAmt()) },
		Sharpen:     func(s *ImageFilter, i *image.Image) { sharpen.Apply(i, s.getAmt()) },
		Blur:        func(s *ImageFilter, i *image.Image) { blur.Apply(i, s.getAmt()) },
		EdgeDetect:  func(s *ImageFilter, i *image.Image) { edgedetect.Apply(i, s.getAmt()) },
		Emboss:      func(s *ImageFilter, i *image.Image) { emboss.Apply(i, s.getAmt()) },
		Threshold:   func(s *ImageFilter, i *image.Image) { threshold.Apply(i, s.getAmt()) },
		AlphaMap:    func(s *ImageFilter, i *image.Image) { alphamap.Apply(i, s.getSrc(), s.getLwr(), s.getUpp()) },
		Extract: func(s *ImageFilter, i *image.Image) {
			extract.Extract(i, filter.ToColorFilter(
				s.getH(), s.getHTlr(), s.getHFth(),
				s.getS(), s.getSTlr(), s.getSFth(),
				s.getL(), s.getLTlr(), s.getLFth(),
			))
		},
		Convolution: func(s *ImageFilter, i *image.Image) {
			convolution.NewCustomFilter(
				s.getAmt(),
				1.0+s.getFct(),
				s.getBis(),
				func(a float64) (m [][]float64) { return s.getMtr() },
			).Apply(i)
		},
	}
)

type ImageFilter struct {
	Type    string         `yaml:"type"`
	Options map[string]any `yaml:"options"`
}

func (f *ImageFilter) String() string {
	res := f.Type
	for k, v := range f.Options {
		res += "::" + k + "=" + fmt.Sprint(v)
	}
	return res
}

func NewImageFilter(typ string, options map[string]any) *ImageFilter {
	return &ImageFilter{
		Type:    typ,
		Options: options,
	}
}

func (s *ImageFilter) getFlt(option string, def float64) float64 {
	v, ok := s.Options[option]
	if ok && v != nil {
		switch val := v.(type) {
		case float64:
			return val
		case float32:
			return float64(val)
		case int:
			return float64(val)
		case int8:
			return float64(val)
		case int16:
			return float64(val)
		case int32:
			return float64(val)
		case int64:
			return float64(val)
		case uint:
			return float64(val)
		case uint8:
			return float64(val)
		case uint16:
			return float64(val)
		case uint32:
			return float64(val)
		case uint64:
			return float64(val)
		default:
			// Do nothing, fall back to default value
		}
	}
	return def
}

func (s *ImageFilter) getStr(option string, def string) string {
	v, ok := s.Options[option]
	if ok && v != nil {
		return v.(string)
	}
	return def
}

func (s *ImageFilter) getCMtr(option string, def [][]float64) [][]float64 {
	v, ok := s.Options[option]
	if ok && v != nil {
		in := v.([]float64)

		// Calculate the number of rows and columns
		rows := int(math.Sqrt(float64(len(in))))

		// Check if the length of the input slice is a perfect square
		if rows*rows != len(in) {
			fmt.Printf("input matrix is not a perfect square, falling back to default\n")
			return defMtr
		}

		// Initialize the matrix
		m := make([][]float64, rows)
		for i := range m {
			m[i] = make([]float64, rows)
		}

		// Fill the matrix with values from the input slice
		for i := 0; i < rows; i++ {
			for j := 0; j < rows; j++ {
				m[i][j] = in[i*rows+j]
			}
		}

		return m
	}
	return def
}

func (s *ImageFilter) getSrc() string      { return s.getStr(optSrc, "s*l") }
func (s *ImageFilter) getLwr() float64     { return s.getFlt(optLwr, 0.0) }
func (s *ImageFilter) getUpp() float64     { return s.getFlt(optUpp, 0.0) }
func (s *ImageFilter) getAmt() float64     { return s.getFlt(optAmt, 1.0) }
func (s *ImageFilter) getH() float64       { return s.getFlt(optHue, 0.0) }
func (s *ImageFilter) getHTlr() float64    { return s.getFlt(optHueTlr, 180.0) }
func (s *ImageFilter) getHFth() float64    { return s.getFlt(optHueFth, 0.0) }
func (s *ImageFilter) getS() float64       { return s.getFlt(optSat, 0.50) }
func (s *ImageFilter) getSTlr() float64    { return s.getFlt(optSatTlr, 0.50) }
func (s *ImageFilter) getSFth() float64    { return s.getFlt(optSatFth, 0.0) }
func (s *ImageFilter) getL() float64       { return s.getFlt(optLum, 0.50) }
func (s *ImageFilter) getLTlr() float64    { return s.getFlt(optLumTlr, 0.50) }
func (s *ImageFilter) getLFth() float64    { return s.getFlt(optLumFth, 0.0) }
func (s *ImageFilter) getFct() float64     { return s.getFlt(optFct, 1.0) }
func (s *ImageFilter) getBis() float64     { return s.getFlt(optBis, 0.0) }
func (s *ImageFilter) getMtr() [][]float64 { return s.getCMtr(optMtr, defMtr) }

func (s *ImageFilter) Apply(i *image.Image) *image.Image {
	if f, ok := filterMap[strings.ToLower(s.Type)]; ok {
		f(s, i)
		return i
	}
	fmt.Printf("unknown filter type: %s\n", s.Type)
	return i
}
