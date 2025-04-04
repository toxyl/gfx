package color

import (
	"testing"

	"github.com/toxyl/gfx/core/color/constants"
)

// TestEdgeCases tests color conversions with edge cases
func TestEdgeCases(t *testing.T) {
	tests := []struct {
		name string
		rgb  *RGBA64
		want map[string][]float64
	}{
		{
			name: "Almost White",
			rgb:  &RGBA64{R: constants.AlmostWhiteR, G: constants.AlmostWhiteG, B: constants.AlmostWhiteB, A: constants.AlmostWhiteA},
			want: map[string][]float64{
				"LAB": {constants.AlmostWhiteLABL, constants.AlmostWhiteLABA, constants.AlmostWhiteLABB, constants.AlmostWhiteA},
				"LCH": {constants.AlmostWhiteLABL, constants.AlmostWhiteLABA, constants.AlmostWhiteLABB, constants.AlmostWhiteA},
				"HCL": {0, 0, constants.AlmostWhiteLABL, constants.AlmostWhiteA},
			},
		},
		{
			name: "Almost Black",
			rgb:  &RGBA64{R: constants.AlmostBlackR, G: constants.AlmostBlackG, B: constants.AlmostBlackB, A: constants.AlmostBlackA},
			want: map[string][]float64{
				"LAB": {constants.AlmostBlackLABL, constants.AlmostBlackLABA, constants.AlmostBlackLABB, constants.AlmostBlackA},
				"LCH": {constants.AlmostBlackLABL, constants.AlmostBlackLABA, constants.AlmostBlackLABB, constants.AlmostBlackA},
				"HCL": {0, 0, constants.AlmostBlackLABL, constants.AlmostBlackA},
			},
		},
		{
			name: "50% Gray",
			rgb:  &RGBA64{R: constants.Gray50R, G: constants.Gray50G, B: constants.Gray50B, A: constants.Gray50A},
			want: map[string][]float64{
				"LAB": {constants.Gray50LABL, constants.Gray50LABA, constants.Gray50LABB, constants.Gray50A},
				"LCH": {constants.Gray50LABL, constants.Gray50LABA, constants.Gray50LABB, constants.Gray50A},
				"HCL": {0, 0, constants.Gray50LABL, constants.Gray50A},
				"YUV": {constants.Gray50YUVY, constants.Gray50YUVU, constants.Gray50YUVV, constants.Gray50A},
				"YIQ": {constants.Gray50YIQY, constants.Gray50YIQI, constants.Gray50YIQQ, constants.Gray50A},
			},
		},
		{
			name: "Transparent White",
			rgb:  &RGBA64{R: constants.TransparentWhiteR, G: constants.TransparentWhiteG, B: constants.TransparentWhiteB, A: constants.TransparentWhiteA},
			want: map[string][]float64{
				"LAB": {constants.TransparentWhiteLABL, constants.TransparentWhiteLABA, constants.TransparentWhiteLABB, constants.TransparentWhiteA},
				"LCH": {constants.TransparentWhiteLABL, constants.TransparentWhiteLABA, constants.TransparentWhiteLABB, constants.TransparentWhiteA},
				"HCL": {0, 0, constants.TransparentWhiteLABL, constants.TransparentWhiteA},
			},
		},
		{
			name: "Semi-Transparent Red",
			rgb:  &RGBA64{R: constants.SemiTransparentRedR, G: constants.SemiTransparentRedG, B: constants.SemiTransparentRedB, A: constants.SemiTransparentRedA},
			want: map[string][]float64{
				"LAB": {constants.SemiTransparentRedLABL, constants.SemiTransparentRedLABA, constants.SemiTransparentRedLABB, constants.SemiTransparentRedA},
				"LCH": {constants.SemiTransparentRedLABL, 104.55, 40.00, constants.SemiTransparentRedA},
				"HCL": {40.00, 104.55, constants.SemiTransparentRedLABL, constants.SemiTransparentRedA},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Test LAB conversion
			lab := LABFromRGB(tt.rgb)
			if !compareChannels(getChannels(lab), tt.want["LAB"], 0.01) {
				t.Errorf("LAB conversion failed for %s: got %v, want %v", tt.name, getChannels(lab), tt.want["LAB"])
			}

			// Test LCH conversion
			lch := LCHFromRGB(tt.rgb)
			if !compareChannels(getChannels(lch), tt.want["LCH"], 0.01) {
				t.Errorf("LCH conversion failed for %s: got %v, want %v", tt.name, getChannels(lch), tt.want["LCH"])
			}

			// Test HCL conversion
			hcl := HCLFromRGB(tt.rgb)
			if !compareChannels(getChannels(hcl), tt.want["HCL"], 0.01) {
				t.Errorf("HCL conversion failed for %s: got %v, want %v", tt.name, getChannels(hcl), tt.want["HCL"])
			}

			if yuvVals, ok := tt.want["YUV"]; ok {
				yuv := YUVFromRGB(tt.rgb)
				if !compareChannels(getChannels(yuv), yuvVals, 0.01) {
					t.Errorf("YUV conversion failed for %s: got %v, want %v", tt.name, getChannels(yuv), yuvVals)
				}
			}

			if yiqVals, ok := tt.want["YIQ"]; ok {
				yiq := YIQFromRGB(tt.rgb)
				if !compareChannels(getChannels(yiq), yiqVals, 0.01) {
					t.Errorf("YIQ conversion failed for %s: got %v, want %v", tt.name, getChannels(yiq), yiqVals)
				}
			}
		})
	}
}

// TestSecondaryColors tests color conversions for secondary colors
func TestSecondaryColors(t *testing.T) {
	tests := []struct {
		name string
		rgb  *RGBA64
		want map[string][]float64
	}{
		{
			name: "Yellow",
			rgb:  &RGBA64{R: constants.YellowR, G: constants.YellowG, B: constants.YellowB, A: constants.YellowA},
			want: map[string][]float64{
				"HSL": {constants.YellowHSLH, constants.YellowHSLS, constants.YellowHSLL, constants.YellowA},
				"HSB": {constants.YellowHSLH, constants.YellowHSBA, constants.YellowHSBB, constants.YellowA},
				"LAB": {constants.YellowLABL, constants.YellowLABA, constants.YellowLABB, constants.YellowA},
				"LSB": {constants.YellowLSBW, constants.YellowLSBS, constants.YellowLSBB, constants.YellowA},
				"LSL": {constants.YellowLSBW, constants.YellowLSBS, constants.YellowLSLL, constants.YellowA},
			},
		},
		{
			name: "Cyan",
			rgb:  &RGBA64{R: constants.CyanR, G: constants.CyanG, B: constants.CyanB, A: constants.CyanA},
			want: map[string][]float64{
				"HSL": {constants.CyanHSLH, constants.CyanHSLS, constants.CyanHSLL, constants.CyanA},
				"HSB": {constants.CyanHSLH, constants.CyanHSBA, constants.CyanHSBB, constants.CyanA},
				"LAB": {constants.CyanLABL, constants.CyanLABA, constants.CyanLABB, constants.CyanA},
				"LSB": {constants.CyanLSBW, constants.CyanLSBS, constants.CyanLSBB, constants.CyanA},
				"LSL": {constants.CyanLSBW, constants.CyanLSBS, constants.CyanLSLL, constants.CyanA},
			},
		},
		{
			name: "Magenta",
			rgb:  &RGBA64{R: constants.MagentaR, G: constants.MagentaG, B: constants.MagentaB, A: constants.MagentaA},
			want: map[string][]float64{
				"HSL": {constants.MagentaHSLH, constants.MagentaHSLS, constants.MagentaHSLL, constants.MagentaA},
				"HSB": {constants.MagentaHSLH, constants.MagentaHSBA, constants.MagentaHSBB, constants.MagentaA},
				"LAB": {constants.MagentaLABL, constants.MagentaLABA, constants.MagentaLABB, constants.MagentaA},
				"LSB": {constants.MagentaLSBW, constants.MagentaLSBS, constants.MagentaLSBB, constants.MagentaA},
				"LSL": {constants.MagentaLSBW, constants.MagentaLSBS, constants.MagentaLSLL, constants.MagentaA},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Test HSL conversion
			hsl := HSLFromRGB(tt.rgb)
			if !compareChannels(getChannels(hsl), tt.want["HSL"], 0.01) {
				t.Errorf("HSL conversion failed for %s: got %v, want %v", tt.name, getChannels(hsl), tt.want["HSL"])
			}

			// Test HSB conversion
			hsb := HSBFromRGB(tt.rgb)
			if !compareChannels(getChannels(hsb), tt.want["HSB"], 0.01) {
				t.Errorf("HSB conversion failed for %s: got %v, want %v", tt.name, getChannels(hsb), tt.want["HSB"])
			}

			// Test LAB conversion
			lab := LABFromRGB(tt.rgb)
			if !compareChannels(getChannels(lab), tt.want["LAB"], 0.01) {
				t.Errorf("LAB conversion failed for %s: got %v, want %v", tt.name, getChannels(lab), tt.want["LAB"])
			}

			// Test LSB conversion
			lsb := LSBFromRGB(tt.rgb)
			if !compareChannels(getChannels(lsb), tt.want["LSB"], 0.01) {
				t.Errorf("LSB conversion failed for %s: got %v, want %v", tt.name, getChannels(lsb), tt.want["LSB"])
			}

			// Test LSL conversion
			lsl := LSLFromRGB(tt.rgb)
			if !compareChannels(getChannels(lsl), tt.want["LSL"], 0.01) {
				t.Errorf("LSL conversion failed for %s: got %v, want %v", tt.name, getChannels(lsl), tt.want["LSL"])
			}
		})
	}
}

// TestPastelColors tests color conversions for pastel colors
func TestPastelColors(t *testing.T) {
	tests := []struct {
		name string
		rgb  *RGBA64
		want map[string][]float64
	}{
		{
			name: "Pastel Pink",
			rgb:  &RGBA64{R: constants.PastelPinkR, G: constants.PastelPinkG, B: constants.PastelPinkB, A: constants.PastelPinkA},
			want: map[string][]float64{
				"HSL": {constants.PastelPinkHSLH, constants.PastelPinkHSLS, constants.PastelPinkHSLL, constants.PastelPinkA},
				"HSB": {constants.PastelPinkHSLH, constants.PastelPinkHSBA, constants.PastelPinkHSBB, constants.PastelPinkA},
				"LAB": {constants.PastelPinkLABL, constants.PastelPinkLABA, constants.PastelPinkLABB, constants.PastelPinkA},
			},
		},
		{
			name: "Pastel Blue",
			rgb:  &RGBA64{R: constants.PastelBlueR, G: constants.PastelBlueG, B: constants.PastelBlueB, A: constants.PastelBlueA},
			want: map[string][]float64{
				"HSL": {constants.PastelBlueHSLH, constants.PastelBlueHSLS, constants.PastelBlueHSLL, constants.PastelBlueA},
				"HSB": {constants.PastelBlueHSLH, constants.PastelBlueHSBA, constants.PastelBlueHSBB, constants.PastelBlueA},
				"LAB": {constants.PastelBlueLABL, constants.PastelBlueLABA, constants.PastelBlueLABB, constants.PastelBlueA},
			},
		},
		{
			name: "Pastel Green",
			rgb:  &RGBA64{R: constants.PastelGreenR, G: constants.PastelGreenG, B: constants.PastelGreenB, A: constants.PastelGreenA},
			want: map[string][]float64{
				"HSL": {constants.PastelGreenHSLH, constants.PastelGreenHSLS, constants.PastelGreenHSLL, constants.PastelGreenA},
				"HSB": {constants.PastelGreenHSLH, constants.PastelGreenHSBA, constants.PastelGreenHSBB, constants.PastelGreenA},
				"LAB": {constants.PastelGreenLABL, constants.PastelGreenLABA, constants.PastelGreenLABB, constants.PastelGreenA},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Test HSL conversion
			hsl := HSLFromRGB(tt.rgb)
			if !compareChannels(getChannels(hsl), tt.want["HSL"], 0.01) {
				t.Errorf("HSL conversion failed for %s: got %v, want %v", tt.name, getChannels(hsl), tt.want["HSL"])
			}

			// Test HSB conversion
			hsb := HSBFromRGB(tt.rgb)
			if !compareChannels(getChannels(hsb), tt.want["HSB"], 0.01) {
				t.Errorf("HSB conversion failed for %s: got %v, want %v", tt.name, getChannels(hsb), tt.want["HSB"])
			}

			// Test LAB conversion
			lab := LABFromRGB(tt.rgb)
			if !compareChannels(getChannels(lab), tt.want["LAB"], 0.01) {
				t.Errorf("LAB conversion failed for %s: got %v, want %v", tt.name, getChannels(lab), tt.want["LAB"])
			}
		})
	}
}
