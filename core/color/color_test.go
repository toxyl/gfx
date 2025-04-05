package color

import (
	"image/color"
	"math"
	"testing"
)

// TestColor represents a test case with RGB values and their expected conversions
type TestColor struct {
	name     string
	rgb      *RGBA64
	expected map[string]struct {
		channels  []float64
		tolerance float64
	}
}

// Test colors with known values
var testColors = []TestColor{
	{
		name: "Black",
		rgb:  &RGBA64{R: 0, G: 0, B: 0, A: 1},
		expected: map[string]struct {
			channels  []float64
			tolerance float64
		}{
			"CMY":       {[]float64{1, 1, 1, 1}, 0.001},
			"CMYK":      {[]float64{0, 0, 0, 1, 1}, 0.001},
			"Grayscale": {[]float64{0, 1}, 0.001},
			"Hex":       {[]float64{1}, 0.001},
			"HSL":       {[]float64{0, 0, 0, 1}, 0.001},
			"HSB":       {[]float64{0, 0, 0, 1}, 0.001},
			"LAB":       {[]float64{0, 0, 0, 1}, 0.001},
			"LCH":       {[]float64{0, 0, 0, 1}, 0.001},
			"HCL":       {[]float64{0, 0, 0, 1}, 0.001},
			"LSB":       {[]float64{380, 0, 0, 1}, 0.001},
			"LSL":       {[]float64{380, 0, 0, 1}, 0.001},
			"RGB8":      {[]float64{0, 0, 0, 1}, 0.001},
			"XYZ":       {[]float64{0, 0, 0, 1}, 0.001},
			"YIQ":       {[]float64{0, 0, 0, 1}, 0.001},
			"YUV":       {[]float64{0, 0, 0, 1}, 0.001},
			"YCbCr":     {[]float64{0, 0, 0, 1}, 0.001},
		},
	},
	{
		name: "White",
		rgb:  &RGBA64{R: 1, G: 1, B: 1, A: 1},
		expected: map[string]struct {
			channels  []float64
			tolerance float64
		}{
			"CMY":       {[]float64{0, 0, 0, 1}, 0.001},
			"CMYK":      {[]float64{0, 0, 0, 0, 1}, 0.001},
			"Grayscale": {[]float64{1, 1}, 0.001},
			"Hex":       {[]float64{1}, 0.001},
			"HSL":       {[]float64{0, 0, 1, 1}, 0.001},
			"HSB":       {[]float64{0, 0, 1, 1}, 0.001},
			"LAB":       {[]float64{100, 0, 0, 1}, 0.001},
			"LCH":       {[]float64{100, 0, 0, 1}, 0.001},
			"HCL":       {[]float64{0, 0, 100, 1}, 0.001},
			"LSB":       {[]float64{380, 0, 1, 1}, 0.001},
			"LSL":       {[]float64{380, 0, 1, 1}, 0.001},
			"RGB8":      {[]float64{255, 255, 255, 1}, 0.001},
			"XYZ":       {[]float64{0.95047, 1.00000, 1.08883, 1}, 0.001},
			"YIQ":       {[]float64{1, 0, 0, 1}, 0.001},
			"YUV":       {[]float64{1, 0, 0, 1}, 0.001},
			"YCbCr":     {[]float64{1, 0, 0, 1}, 0.001},
		},
	},
	{
		name: "Red",
		rgb:  &RGBA64{R: 1, G: 0, B: 0, A: 1},
		expected: map[string]struct {
			channels  []float64
			tolerance float64
		}{
			"CMY":       {[]float64{0, 1, 1, 1}, 0.001},
			"CMYK":      {[]float64{0, 1, 1, 0, 1}, 0.001},
			"Grayscale": {[]float64{0.299, 1}, 0.001},
			"Hex":       {[]float64{1}, 0.001},
			"HSL":       {[]float64{0, 1, 0.5, 1}, 0.001},
			"HSB":       {[]float64{0, 1, 1, 1}, 0.001},
			"LAB":       {[]float64{53.2408, 80.0925, 67.2032, 1}, 0.001},
			"LCH":       {[]float64{53.2408, 104.5518, 39.9990, 1}, 0.001},
			"HCL":       {[]float64{39.9990, 104.5518, 53.2408, 1}, 0.001},
			"LSB":       {[]float64{700, 1, 1, 1}, 0.001},
			"LSL":       {[]float64{700, 1, 0.5, 1}, 0.001},
			"RGB8":      {[]float64{255, 0, 0, 1}, 0.001},
			"XYZ":       {[]float64{0.4124564, 0.2126729, 0.0193339, 1}, 0.001},
			"YIQ":       {[]float64{0.299, 0.596, 0.212, 1}, 0.001},
			"YUV":       {[]float64{0.299, -0.14713, 0.615, 1}, 0.001},
			"YCbCr":     {[]float64{0.299, -0.168736, 0.5, 1}, 0.001},
		},
	},
	{
		name: "Green",
		rgb:  &RGBA64{R: 0, G: 1, B: 0, A: 1},
		expected: map[string]struct {
			channels  []float64
			tolerance float64
		}{
			"CMY":       {[]float64{1, 0, 1, 1}, 0.001},
			"CMYK":      {[]float64{1, 0, 1, 0, 1}, 0.001},
			"Grayscale": {[]float64{0.587, 1}, 0.001},
			"Hex":       {[]float64{1}, 0.001},
			"HSL":       {[]float64{120, 1, 0.5, 1}, 0.001},
			"HSB":       {[]float64{120, 1, 1, 1}, 0.001},
			"LAB":       {[]float64{87.7347, -86.1827, 83.1793, 1}, 0.001},
			"LCH":       {[]float64{87.7347, 119.7759, 136.0160, 1}, 0.001},
			"HCL":       {[]float64{136.0160, 119.7759, 87.7347, 1}, 0.001},
			"LSB":       {[]float64{520, 1, 1, 1}, 0.001},
			"LSL":       {[]float64{520, 1, 0.5, 1}, 0.001},
			"RGB8":      {[]float64{0, 255, 0, 1}, 0.001},
			"XYZ":       {[]float64{0.3575761, 0.7151522, 0.1191920, 1}, 0.001},
			"YIQ":       {[]float64{0.587, -0.275, -0.523, 1}, 0.001},
			"YUV":       {[]float64{0.587, -0.331, -0.419, 1}, 0.001},
			"YCbCr":     {[]float64{0.587, -0.331264, -0.418688, 1}, 0.001},
		},
	},
	{
		name: "Blue",
		rgb:  &RGBA64{R: 0, G: 0, B: 1, A: 1},
		expected: map[string]struct {
			channels  []float64
			tolerance float64
		}{
			"CMY":       {[]float64{1, 1, 0, 1}, 0.001},
			"CMYK":      {[]float64{1, 1, 0, 0, 1}, 0.001},
			"Grayscale": {[]float64{0.114, 1}, 0.001},
			"Hex":       {[]float64{1}, 0.001},
			"HSL":       {[]float64{240, 1, 0.5, 1}, 0.001},
			"HSB":       {[]float64{240, 1, 1, 1}, 0.001},
			"LAB":       {[]float64{32.2970, 79.1875, -107.8602, 1}, 0.001},
			"LCH":       {[]float64{32.2970, 133.8076, 306.2849, 1}, 0.001},
			"HCL":       {[]float64{306.2849, 133.8076, 32.2970, 1}, 0.001},
			"LSB":       {[]float64{450, 1, 1, 1}, 0.001},
			"LSL":       {[]float64{450, 1, 0.5, 1}, 0.001},
			"RGB8":      {[]float64{0, 0, 255, 1}, 0.001},
			"XYZ":       {[]float64{0.1804375, 0.0721750, 0.9503041, 1}, 0.001},
			"YIQ":       {[]float64{0.114, -0.321, 0.311, 1}, 0.001},
			"YUV":       {[]float64{0.114, 0.5, -0.081, 1}, 0.001},
			"YCbCr":     {[]float64{0.114, 0.5, -0.081312, 1}, 0.001},
		},
	},
}

// Helper function to compare float64 slices with tolerance
func compareChannels(actual, expected []float64, tolerance float64) bool {
	if len(actual) != len(expected) {
		return false
	}
	for i := range actual {
		if math.Abs(actual[i]-expected[i]) > tolerance {
			return false
		}
	}
	return true
}

// Helper function to get channels from a color model
func getChannels(color iColor) []float64 {
	switch c := color.(type) {
	case *CMY:
		return []float64{c.C, c.M, c.Y, c.Alpha}
	case *CMYK:
		return []float64{c.C, c.M, c.Y, c.K, c.Alpha}
	case *Grayscale:
		return []float64{c.Gray, c.Alpha}
	case *Hex:
		return []float64{c.Alpha}
	case *HSL:
		return []float64{c.H, c.S, c.L, c.Alpha}
	case *HSB:
		return []float64{c.Hue, c.Saturation, c.Brightness, c.Alpha}
	case *LAB:
		return []float64{c.L, c.A, c.B, c.Alpha}
	case *LCH:
		return []float64{c.L, c.C, c.H, c.Alpha}
	case *HCL:
		return []float64{c.H, c.C, c.L, c.Alpha}
	case *LSB:
		return []float64{c.Lightness, c.Saturation, c.Brightness, c.Alpha}
	case *LSL:
		return []float64{c.Wavelength, c.Saturation, c.Lightness, c.Alpha}
	case *RGB8:
		return []float64{c.R, c.G, c.B, c.Alpha}
	case *XYZ:
		return []float64{c.X, c.Y, c.Z, c.Alpha}
	case *YIQ:
		return []float64{c.Y, c.I, c.Q, c.Alpha}
	case *YUV:
		return []float64{c.Y, c.U, c.V, c.Alpha}
	case *YCbCr:
		return []float64{c.Y, c.Cb, c.Cr, c.Alpha}
	default:
		return nil
	}
}

// TestColorConversions tests the conversion between different color models
func TestColorConversions(t *testing.T) {
	for _, tc := range testColors {
		t.Run(tc.name, func(t *testing.T) {
			// Test CMY conversion
			cmy := CMYFromRGB(tc.rgb)
			if !compareChannels(getChannels(cmy), tc.expected["CMY"].channels, tc.expected["CMY"].tolerance) {
				t.Errorf("CMY conversion failed for %s: got %v, want %v", tc.name, getChannels(cmy), tc.expected["CMY"].channels)
			}

			// Test CMYK conversion
			cmyk := CMYKFromRGB(tc.rgb)
			if !compareChannels(getChannels(cmyk), tc.expected["CMYK"].channels, tc.expected["CMYK"].tolerance) {
				t.Errorf("CMYK conversion failed for %s: got %v, want %v", tc.name, getChannels(cmyk), tc.expected["CMYK"].channels)
			}

			// Test Grayscale conversion
			gray := GrayscaleFromRGB(tc.rgb)
			if !compareChannels(getChannels(gray), tc.expected["Grayscale"].channels, tc.expected["Grayscale"].tolerance) {
				t.Errorf("Grayscale conversion failed for %s: got %v, want %v", tc.name, getChannels(gray), tc.expected["Grayscale"].channels)
			}

			// Test Hex conversion
			hex := HexFromRGB(tc.rgb)
			if !compareChannels(getChannels(hex), tc.expected["Hex"].channels, tc.expected["Hex"].tolerance) {
				t.Errorf("Hex conversion failed for %s: got %v, want %v", tc.name, getChannels(hex), tc.expected["Hex"].channels)
			}

			// Test HSL conversion
			hsl := HSLFromRGB(tc.rgb)
			if !compareChannels(getChannels(hsl), tc.expected["HSL"].channels, tc.expected["HSL"].tolerance) {
				t.Errorf("HSL conversion failed for %s: got %v, want %v", tc.name, getChannels(hsl), tc.expected["HSL"].channels)
			}

			// Test HSB conversion
			hsb := HSBFromRGB(tc.rgb)
			if !compareChannels(getChannels(hsb), tc.expected["HSB"].channels, tc.expected["HSB"].tolerance) {
				t.Errorf("HSB conversion failed for %s: got %v, want %v", tc.name, getChannels(hsb), tc.expected["HSB"].channels)
			}

			// Test LAB conversion
			lab := LABFromRGB(tc.rgb)
			if !compareChannels(getChannels(lab), tc.expected["LAB"].channels, tc.expected["LAB"].tolerance) {
				t.Errorf("LAB conversion failed for %s: got %v, want %v", tc.name, getChannels(lab), tc.expected["LAB"].channels)
			}

			// Test LCH conversion
			lch := LCHFromRGB(tc.rgb)
			if !compareChannels(getChannels(lch), tc.expected["LCH"].channels, tc.expected["LCH"].tolerance) {
				t.Errorf("LCH conversion failed for %s: got %v, want %v", tc.name, getChannels(lch), tc.expected["LCH"].channels)
			}

			// Test HCL conversion
			hcl := HCLFromRGB(tc.rgb)
			if !compareChannels(getChannels(hcl), tc.expected["HCL"].channels, tc.expected["HCL"].tolerance) {
				t.Errorf("HCL conversion failed for %s: got %v, want %v", tc.name, getChannels(hcl), tc.expected["HCL"].channels)
			}

			// Test LSB conversion
			lsb := LSBFromRGB(tc.rgb)
			if !compareChannels(getChannels(lsb), tc.expected["LSB"].channels, tc.expected["LSB"].tolerance) {
				t.Errorf("LSB conversion failed for %s: got %v, want %v", tc.name, getChannels(lsb), tc.expected["LSB"].channels)
			}

			// Test LSL conversion
			lsl := LSLFromRGB(tc.rgb)
			if !compareChannels(getChannels(lsl), tc.expected["LSL"].channels, tc.expected["LSL"].tolerance) {
				t.Errorf("LSL conversion failed for %s: got %v, want %v", tc.name, getChannels(lsl), tc.expected["LSL"].channels)
			}

			// Test RGB8 conversion
			rgb8 := RGB8FromRGB(tc.rgb)
			if !compareChannels(getChannels(rgb8), tc.expected["RGB8"].channels, tc.expected["RGB8"].tolerance) {
				t.Errorf("RGB8 conversion failed for %s: got %v, want %v", tc.name, getChannels(rgb8), tc.expected["RGB8"].channels)
			}

			// Test XYZ conversion
			xyz := XYZFromRGB(tc.rgb)
			if !compareChannels(getChannels(xyz), tc.expected["XYZ"].channels, tc.expected["XYZ"].tolerance) {
				t.Errorf("XYZ conversion failed for %s: got %v, want %v", tc.name, getChannels(xyz), tc.expected["XYZ"].channels)
			}

			// Test YIQ conversion
			yiq := YIQFromRGB(tc.rgb)
			if !compareChannels(getChannels(yiq), tc.expected["YIQ"].channels, tc.expected["YIQ"].tolerance) {
				t.Errorf("YIQ conversion failed for %s: got %v, want %v", tc.name, getChannels(yiq), tc.expected["YIQ"].channels)
			}

			// Test YUV conversion
			yuv := YUVFromRGB(tc.rgb)
			if !compareChannels(getChannels(yuv), tc.expected["YUV"].channels, tc.expected["YUV"].tolerance) {
				t.Errorf("YUV conversion failed for %s: got %v, want %v", tc.name, getChannels(yuv), tc.expected["YUV"].channels)
			}

			// Test YCbCr conversion
			ycbcr := YCbCrFromRGB(tc.rgb)
			if !compareChannels(getChannels(ycbcr), tc.expected["YCbCr"].channels, tc.expected["YCbCr"].tolerance) {
				t.Errorf("YCbCr conversion failed for %s: got %v, want %v", tc.name, getChannels(ycbcr), tc.expected["YCbCr"].channels)
			}
		})
	}
}

// TestRoundTripConversions tests that converting from RGB to another color model and back
// produces the same RGB values (within tolerance)
func TestRoundTripConversions(t *testing.T) {
	for _, tc := range testColors {
		t.Run(tc.name, func(t *testing.T) {
			// Test CMY round trip
			cmy := CMYFromRGB(tc.rgb)
			rgb := cmy.ToRGB()
			if !compareChannels([]float64{rgb.R, rgb.G, rgb.B, rgb.A}, []float64{tc.rgb.R, tc.rgb.G, tc.rgb.B, tc.rgb.A}, 0.001) {
				t.Errorf("CMY round trip failed for %s", tc.name)
			}

			// Test CMYK round trip
			cmyk := CMYKFromRGB(tc.rgb)
			rgb = cmyk.ToRGB()
			if !compareChannels([]float64{rgb.R, rgb.G, rgb.B, rgb.A}, []float64{tc.rgb.R, tc.rgb.G, tc.rgb.B, tc.rgb.A}, 0.001) {
				t.Errorf("CMYK round trip failed for %s", tc.name)
			}

			// Test Grayscale round trip
			gray := GrayscaleFromRGB(tc.rgb)
			rgb = gray.ToRGB()
			if !compareChannels([]float64{rgb.R, rgb.G, rgb.B, rgb.A}, []float64{tc.rgb.R, tc.rgb.G, tc.rgb.B, tc.rgb.A}, 0.001) {
				t.Errorf("Grayscale round trip failed for %s", tc.name)
			}

			// Test Hex round trip
			hex := HexFromRGB(tc.rgb)
			rgb = hex.ToRGB()
			if !compareChannels([]float64{rgb.R, rgb.G, rgb.B, rgb.A}, []float64{tc.rgb.R, tc.rgb.G, tc.rgb.B, tc.rgb.A}, 0.001) {
				t.Errorf("Hex round trip failed for %s", tc.name)
			}

			// Test HSL round trip
			hsl := HSLFromRGB(tc.rgb)
			rgb = hsl.ToRGB()
			if !compareChannels([]float64{rgb.R, rgb.G, rgb.B, rgb.A}, []float64{tc.rgb.R, tc.rgb.G, tc.rgb.B, tc.rgb.A}, 0.001) {
				t.Errorf("HSL round trip failed for %s", tc.name)
			}

			// Test HSB round trip
			hsb := HSBFromRGB(tc.rgb)
			rgb = hsb.ToRGB()
			if !compareChannels([]float64{rgb.R, rgb.G, rgb.B, rgb.A}, []float64{tc.rgb.R, tc.rgb.G, tc.rgb.B, tc.rgb.A}, 0.001) {
				t.Errorf("HSB round trip failed for %s", tc.name)
			}

			// Test LAB round trip
			lab := LABFromRGB(tc.rgb)
			rgb = lab.ToRGB()
			if !compareChannels([]float64{rgb.R, rgb.G, rgb.B, rgb.A}, []float64{tc.rgb.R, tc.rgb.G, tc.rgb.B, tc.rgb.A}, 0.001) {
				t.Errorf("LAB round trip failed for %s", tc.name)
			}

			// Test LCH round trip
			lch := LCHFromRGB(tc.rgb)
			rgb = lch.ToRGB()
			if !compareChannels([]float64{rgb.R, rgb.G, rgb.B, rgb.A}, []float64{tc.rgb.R, tc.rgb.G, tc.rgb.B, tc.rgb.A}, 0.001) {
				t.Errorf("LCH round trip failed for %s", tc.name)
			}

			// Test HCL round trip
			hcl := HCLFromRGB(tc.rgb)
			rgb = hcl.ToRGB()
			if !compareChannels([]float64{rgb.R, rgb.G, rgb.B, rgb.A}, []float64{tc.rgb.R, tc.rgb.G, tc.rgb.B, tc.rgb.A}, 0.001) {
				t.Errorf("HCL round trip failed for %s", tc.name)
			}

			// Test LSB round trip
			lsb := LSBFromRGB(tc.rgb)
			rgb = lsb.ToRGB()
			if !compareChannels([]float64{rgb.R, rgb.G, rgb.B, rgb.A}, []float64{tc.rgb.R, tc.rgb.G, tc.rgb.B, tc.rgb.A}, 0.001) {
				t.Errorf("LSB round trip failed for %s", tc.name)
			}

			// Test LSL round trip
			lsl := LSLFromRGB(tc.rgb)
			rgb = lsl.ToRGB()
			if !compareChannels([]float64{rgb.R, rgb.G, rgb.B, rgb.A}, []float64{tc.rgb.R, tc.rgb.G, tc.rgb.B, tc.rgb.A}, 0.001) {
				t.Errorf("LSL round trip failed for %s", tc.name)
			}

			// Test RGB8 round trip
			rgb8 := RGB8FromRGB(tc.rgb)
			rgb = rgb8.ToRGB()
			if !compareChannels([]float64{rgb.R, rgb.G, rgb.B, rgb.A}, []float64{tc.rgb.R, tc.rgb.G, tc.rgb.B, tc.rgb.A}, 0.001) {
				t.Errorf("RGB8 round trip failed for %s", tc.name)
			}

			// Test XYZ round trip
			xyz := XYZFromRGB(tc.rgb)
			rgb = xyz.ToRGB()
			if !compareChannels([]float64{rgb.R, rgb.G, rgb.B, rgb.A}, []float64{tc.rgb.R, tc.rgb.G, tc.rgb.B, tc.rgb.A}, 0.001) {
				t.Errorf("XYZ round trip failed for %s", tc.name)
			}

			// Test YIQ round trip
			yiq := YIQFromRGB(tc.rgb)
			rgb = yiq.ToRGB()
			if !compareChannels([]float64{rgb.R, rgb.G, rgb.B, rgb.A}, []float64{tc.rgb.R, tc.rgb.G, tc.rgb.B, tc.rgb.A}, 0.001) {
				t.Errorf("YIQ round trip failed for %s", tc.name)
			}

			// Test YUV round trip
			yuv := YUVFromRGB(tc.rgb)
			rgb = yuv.ToRGB()
			if !compareChannels([]float64{rgb.R, rgb.G, rgb.B, rgb.A}, []float64{tc.rgb.R, tc.rgb.G, tc.rgb.B, tc.rgb.A}, 0.001) {
				t.Errorf("YUV round trip failed for %s", tc.name)
			}

			// Test YCbCr round trip
			ycbcr := YCbCrFromRGB(tc.rgb)
			rgb = ycbcr.ToRGB()
			if !compareChannels([]float64{rgb.R, rgb.G, rgb.B, rgb.A}, []float64{tc.rgb.R, tc.rgb.G, tc.rgb.B, tc.rgb.A}, 0.001) {
				t.Errorf("YCbCr round trip failed for %s", tc.name)
			}
		})
	}
}

func TestColor64Basics(t *testing.T) {
	// Test New function
	col := New(0.5, 0.6, 0.7, 0.8)
	if col.R != 0.5 || col.G != 0.6 || col.B != 0.7 || col.A != 0.8 {
		t.Errorf("New created incorrect color: got %+v", col)
	}

	// Test clamping in New
	col = New(1.5, -0.1, 2.0, -0.5)
	if col.R != 1.0 || col.G != 0.0 || col.B != 1.0 || col.A != 0.0 {
		t.Errorf("New failed to clamp values: got %+v", col)
	}

	// Test Set
	col = &Color64{}
	col.Set(0.2, 0.3, 0.4, 0.5)
	if col.R != 0.2 || col.G != 0.3 || col.B != 0.4 || col.A != 0.5 {
		t.Errorf("Set failed: got %+v", col)
	}

	// Test Get
	copy := col.Get()
	if copy.R != 0.2 || copy.G != 0.3 || copy.B != 0.4 || copy.A != 0.5 {
		t.Errorf("Get failed: got %+v", copy)
	}

	// Test SetUint8
	col = &Color64{}
	col.SetUint8(color.RGBA{R: 51, G: 102, B: 153, A: 204}) // ~20%, 40%, 60%, 80%
	if !(0.19 < col.R && col.R < 0.21) ||
		!(0.39 < col.G && col.G < 0.41) ||
		!(0.59 < col.B && col.B < 0.61) ||
		!(0.79 < col.A && col.A < 0.81) {
		t.Errorf("SetUint8 failed: got %+v", col)
	}

	// Test SetUint16
	col = &Color64{}
	col.SetUint16(color.RGBA64{R: 13107, G: 26214, B: 39321, A: 52428}) // ~20%, 40%, 60%, 80%
	if !(0.19 < col.R && col.R < 0.21) ||
		!(0.39 < col.G && col.G < 0.41) ||
		!(0.59 < col.B && col.B < 0.61) ||
		!(0.79 < col.A && col.A < 0.81) {
		t.Errorf("SetUint16 failed: got %+v", col)
	}

	// Test ToUint8
	col = New(0.2, 0.4, 0.6, 0.8)
	rgba8 := col.ToUint8()
	if rgba8.R != 51 || rgba8.G != 102 || rgba8.B != 153 || rgba8.A != 204 {
		t.Errorf("ToUint8 failed: got %+v", rgba8)
	}

	// Test ToUint16
	col = New(0.2, 0.4, 0.6, 0.8)
	rgba16 := col.ToUint16()
	if rgba16.R != 13107 || rgba16.G != 26214 || rgba16.B != 39321 || rgba16.A != 52428 {
		t.Errorf("ToUint16 failed: got %+v", rgba16)
	}

	// Test String
	col = New(0.1, 0.2, 0.3, 0.4)
	str := col.String()
	if str != "RGBA(r: 0.100000, g: 0.200000, b: 0.300000, a: 0.400000)" {
		t.Errorf("String failed: got %s", str)
	}
}

func TestRGBA64(t *testing.T) {
	// Test NewRGBA64
	col, err := NewRGBA64(0.5, 0.6, 0.7, 0.8)
	if err != nil {
		t.Errorf("NewRGBA64 returned error: %v", err)
	}
	if col.R != 0.5 || col.G != 0.6 || col.B != 0.7 || col.A != 0.8 {
		t.Errorf("NewRGBA64 created incorrect color: got %+v", col)
	}

	// Test FromSlice
	col = &RGBA64{}
	err = col.FromSlice([]float64{0.1, 0.2, 0.3, 0.4})
	if err != nil {
		t.Errorf("FromSlice returned error: %v", err)
	}
	if col.R != 0.1 || col.G != 0.2 || col.B != 0.3 || col.A != 0.4 {
		t.Errorf("FromSlice failed: got %+v", col)
	}

	// Test Copy
	col, _ = NewRGBA64(0.5, 0.6, 0.7, 0.8)
	copy := col.Copy()
	if copy.R != 0.5 || copy.G != 0.6 || copy.B != 0.7 || copy.A != 0.8 {
		t.Errorf("Copy failed: got %+v", copy)
	}

	// Test Clamp
	col, _ = NewRGBA64(1.5, -0.1, 2.0, -0.5)
	col.Clamp()
	if col.R != 1.0 || col.G != 0.0 || col.B != 1.0 || col.A != 0.0 {
		t.Errorf("Clamp failed: got %+v", col)
	}
}

func TestHSL(t *testing.T) {
	// Test NewHSL
	hsl, err := NewHSL(180, 0.5, 0.5, 1.0)
	if err != nil {
		t.Errorf("NewHSL returned error: %v", err)
	}
	if hsl.H != 180 || hsl.S != 0.5 || hsl.L != 0.5 || hsl.Alpha != 1.0 {
		t.Errorf("NewHSL created incorrect color: got %+v", hsl)
	}

	// Test HSLFromRGB
	rgba, _ := NewRGBA64(0.5, 0.6, 0.7, 0.8)
	hsl = HSLFromRGB(rgba)
	if hsl.Alpha != 0.8 {
		t.Errorf("HSLFromRGB failed to copy alpha: got %f", hsl.Alpha)
	}

	// Test ToRGB (conversion back to RGB)
	rgba = hsl.ToRGB()
	if rgba == nil {
		t.Errorf("ToRGB returned nil")
	}
	if rgba.A != 0.8 {
		t.Errorf("ToRGB failed to copy alpha: got %f", rgba.A)
	}
}

func TestColorParsing(t *testing.T) {
	// Test parseHex
	hexColor, err := parseHex("#ff0000")
	if err != nil {
		t.Errorf("parseHex returned error: %v", err)
	}
	rgb8, ok := hexColor.(*RGB8)
	if !ok {
		t.Errorf("parseHex didn't return RGB8")
	} else if rgb8.R != 255 || rgb8.G != 0 || rgb8.B != 0 || rgb8.Alpha != 1.0 {
		t.Errorf("parseHex parsed incorrect color: got %+v", rgb8)
	}

	// Test parseRGB
	rgbColor, err := parseRGB("rgb(100, 150, 200)")
	if err != nil {
		t.Errorf("parseRGB returned error: %v", err)
	}
	rgb8, ok = rgbColor.(*RGB8)
	if !ok {
		t.Errorf("parseRGB didn't return RGB8")
	} else if rgb8.R != 100 || rgb8.G != 150 || rgb8.B != 200 || rgb8.Alpha != 1.0 {
		t.Errorf("parseRGB parsed incorrect color: got %+v", rgb8)
	}

	// Test parseHSL
	hslColor, err := parseHSL("hsl(180, 50%, 75%)")
	if err != nil {
		t.Errorf("parseHSL returned error: %v", err)
	}
	hsl, ok := hslColor.(*HSL)
	if !ok {
		t.Errorf("parseHSL didn't return HSL")
	} else if hsl.H != 180 || hsl.S != 0.5 || hsl.L != 0.75 || hsl.Alpha != 1.0 {
		t.Errorf("parseHSL parsed incorrect color: got %+v", hsl)
	}
}

func TestColorFormatting(t *testing.T) {
	// Create a test color
	rgba, _ := NewRGBA64(1.0, 0.0, 0.0, 1.0) // Pure red

	// Test ToHex
	hex := ToHex(rgba)
	if hex != "#FF0000" {
		t.Errorf("ToHex failed: got %s, expected #FF0000", hex)
	}

	// Test ToRGBString
	rgbStr := ToRGBString(rgba)
	if rgbStr != "rgb(255, 0, 0)" {
		t.Errorf("ToRGBString failed: got %s, expected rgb(255, 0, 0)", rgbStr)
	}

	// Test Format
	hexFmt, err := Format(rgba, "hex")
	if err != nil || hexFmt != "#FF0000" {
		t.Errorf("Format with 'hex' failed: got %s, error: %v", hexFmt, err)
	}

	rgbFmt, err := Format(rgba, "rgb")
	if err != nil || rgbFmt != "rgb(255, 0, 0)" {
		t.Errorf("Format with 'rgb' failed: got %s, error: %v", rgbFmt, err)
	}
}
