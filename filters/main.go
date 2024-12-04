package filters

import (
	"fmt"
	"strings"

	"github.com/toxyl/gfx/color/filter"
	"github.com/toxyl/gfx/filters/alphamap"
	"github.com/toxyl/gfx/filters/brightness"
	"github.com/toxyl/gfx/filters/colorshift"
	"github.com/toxyl/gfx/filters/contrast"
	"github.com/toxyl/gfx/filters/convolution"
	"github.com/toxyl/gfx/filters/extract"
	"github.com/toxyl/gfx/filters/gamma"
	"github.com/toxyl/gfx/filters/grayscale"
	"github.com/toxyl/gfx/filters/huerotate"
	"github.com/toxyl/gfx/filters/invert"
	"github.com/toxyl/gfx/filters/lightnesscontrast"
	"github.com/toxyl/gfx/filters/pastelize"
	"github.com/toxyl/gfx/filters/saturation"
	"github.com/toxyl/gfx/filters/sepia"
	"github.com/toxyl/gfx/filters/threshold"
	"github.com/toxyl/gfx/filters/vibrance"
	"github.com/toxyl/gfx/image"
)

const (
	ALPHAMAP           = "alpha-map"
	BLUR               = "blur"
	BRIGHTNESS         = "brightness"
	COLOR_SHIFT        = "color-shift"
	CONTRAST           = "contrast"
	EDGE_DETECT        = "edge-detect"
	EMBOSS             = "emboss"
	ENHANCE            = "enhance"
	EXTRACT            = "extract"
	GAMMA              = "gamma"
	GRAYSCALE          = "grayscale"
	HUE_ROTATE         = "hue-rotate"
	INVERT             = "invert"
	LIGHTNESS_CONTRAST = "lightness-contrast"
	PASTELIZE          = "pastelize"
	SATURATION         = "saturation"
	SEPIA              = "sepia"
	SHARPEN            = "sharpen"
	THRESHOLD          = "threshold"
	VIBRANCE           = "vibrance"

	OPTION_AMOUNT        = "amount"
	OPTION_FACTOR        = "factor"
	OPTION_BIAS          = "bias"
	OPTION_HUE           = "hue"
	OPTION_HUE_TOLERANCE = "hue-tolerance"
	OPTION_HUE_FEATHER   = "hue-feather"
	OPTION_SAT           = "sat"
	OPTION_SAT_TOLERANCE = "sat-tolerance"
	OPTION_SAT_FEATHER   = "sat-feather"
	OPTION_LUM           = "lum"
	OPTION_LUM_TOLERANCE = "lum-tolerance"
	OPTION_LUM_FEATHER   = "lum-feather"
	OPTION_FALLOFF       = "falloff"
	OPTION_LOWER         = "lower"
	OPTION_UPPER         = "upper"
	OPTION_SOURCE        = "source"
)

var (
	EXAMPLES = []string{
		ALPHAMAP + "::" + OPTION_SOURCE + "=s*l::" + OPTION_LOWER + "=0.1::" + OPTION_UPPER + "=0.7",
		BLUR + "::" + OPTION_AMOUNT + "=1.0",
		BRIGHTNESS + "::" + OPTION_AMOUNT + "=1.0",
		COLOR_SHIFT + "::" + OPTION_HUE + "=180.0::" + OPTION_SAT + "=0.1::" + OPTION_LUM + "=0.7",
		CONTRAST + "::" + OPTION_AMOUNT + "=1.0",
		EDGE_DETECT + "::" + OPTION_AMOUNT + "=1.0",
		EMBOSS + "::" + OPTION_AMOUNT + "=1.0",
		ENHANCE + "::" + OPTION_AMOUNT + "=1.0",
		EXTRACT + "::" + OPTION_HUE + "=180.0::" + OPTION_HUE_TOLERANCE + "=90.0::" + OPTION_HUE_FEATHER + "=90.0::" +
			OPTION_SAT + "=0.50::" + OPTION_SAT_TOLERANCE + "=0.25::" + OPTION_SAT_FEATHER + "=0.25::" +
			OPTION_LUM + "=0.50::" + OPTION_LUM_TOLERANCE + "=0.25::" + OPTION_LUM_FEATHER + "=0.25",
		GAMMA + "::" + OPTION_AMOUNT + "=1.0",
		GRAYSCALE,
		HUE_ROTATE + "::" + OPTION_AMOUNT + "=180.0",
		INVERT,
		LIGHTNESS_CONTRAST + "::" + OPTION_AMOUNT + "=1.0",
		PASTELIZE,
		SATURATION + "::" + OPTION_AMOUNT + "=1.0",
		SEPIA,
		SHARPEN + "::" + OPTION_AMOUNT + "=1.0",
		THRESHOLD + "::" + OPTION_AMOUNT + "=1.0",
		VIBRANCE + "::" + OPTION_AMOUNT + "=1.0",
	}
)

type ImageFilter struct {
	Type    string         `yaml:"type"`
	Options map[string]any `yaml:"options"`
}

func NewImageFilter(typ string, options map[string]any) *ImageFilter {
	return &ImageFilter{
		Type:    typ,
		Options: options,
	}
}

func (s *ImageFilter) getFloat(option string, def float64) float64 {
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

func (s *ImageFilter) getString(option string, def string) string {
	v, ok := s.Options[option]
	if ok && v != nil {
		return v.(string)
	}
	return def
}

func (s *ImageFilter) getSource() string          { return s.getString(OPTION_SOURCE, "s*l") }
func (s *ImageFilter) getLowerThreshold() float64 { return s.getFloat(OPTION_LOWER, 0.0) }
func (s *ImageFilter) getUpperThreshold() float64 { return s.getFloat(OPTION_UPPER, 0.0) }
func (s *ImageFilter) getAmount() float64         { return s.getFloat(OPTION_AMOUNT, 1.0) }
func (s *ImageFilter) getHue() float64            { return s.getFloat(OPTION_HUE, 0.0) }
func (s *ImageFilter) getHueTolerance() float64   { return s.getFloat(OPTION_HUE_TOLERANCE, 180.0) }
func (s *ImageFilter) getHueFeather() float64     { return s.getFloat(OPTION_HUE_FEATHER, 0.0) }
func (s *ImageFilter) getSat() float64            { return s.getFloat(OPTION_SAT, 0.50) }
func (s *ImageFilter) getSatTolerance() float64   { return s.getFloat(OPTION_SAT_TOLERANCE, 0.50) }
func (s *ImageFilter) getSatFeather() float64     { return s.getFloat(OPTION_SAT_FEATHER, 0.0) }
func (s *ImageFilter) getLum() float64            { return s.getFloat(OPTION_LUM, 0.50) }
func (s *ImageFilter) getLumTolerance() float64   { return s.getFloat(OPTION_LUM_TOLERANCE, 0.50) }
func (s *ImageFilter) getLumFeather() float64     { return s.getFloat(OPTION_LUM_FEATHER, 0.0) }

func (s *ImageFilter) Apply(i *image.Image) *image.Image {
	switch strings.ToLower(s.Type) {
	case GRAYSCALE:
		return grayscale.Apply(i)
	case INVERT:
		return invert.Apply(i)
	case SEPIA:
		return sepia.Apply(i)
	case PASTELIZE:
		return pastelize.Apply(i)
	case SHARPEN:
		return convolution.NewSharpenFilter(s.getAmount()).Apply(i)
	case BLUR:
		return convolution.NewBlurFilter(s.getAmount()).Apply(i)
	case EMBOSS:
		return convolution.NewEmbossFilter(s.getAmount()).Apply(i)
	case EDGE_DETECT:
		return convolution.NewEdgeDetectFilter(s.getAmount()).Apply(i)
	case ENHANCE:
		return convolution.NewEnhanceFilter(s.getAmount()).Apply(i)
	case CONTRAST:
		return contrast.Apply(i, s.getAmount())
	case LIGHTNESS_CONTRAST:
		return lightnesscontrast.Apply(i, s.getAmount())
	case BRIGHTNESS:
		return brightness.Apply(i, s.getAmount())
	case THRESHOLD:
		return threshold.Apply(i, uint8(s.getAmount()*255.0))
	case SATURATION:
		return saturation.Apply(i, s.getAmount())
	case HUE_ROTATE:
		return huerotate.Apply(i, s.getAmount())
	case VIBRANCE:
		return vibrance.Apply(i, s.getAmount())
	case GAMMA:
		return gamma.Apply(i, s.getAmount())
	case ALPHAMAP:
		return alphamap.Apply(i, s.getSource(), s.getLowerThreshold(), s.getUpperThreshold())
	case COLOR_SHIFT:
		return colorshift.Apply(i, s.getHue(), s.getSat(), s.getLum())
	case EXTRACT:
		return extract.Extract(i, filter.ToColorFilter(
			s.getHue(), s.getHueTolerance(), s.getHueFeather(),
			s.getSat(), s.getSatTolerance(), s.getSatFeather(),
			s.getLum(), s.getLumTolerance(), s.getLumFeather(),
		))
	default:
		fmt.Printf("unknown filter type: %s\n", s.Type)
		return i
	}
}
