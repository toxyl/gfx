package convolution

import (
	"testing"

	"github.com/toxyl/gfx/image"
	"github.com/toxyl/gfx/math"
)

func TestConvolutionMatrix_Apply(t *testing.T) {
	testImage := image.NewFromURL("https://sdo.gsfc.nasa.gov/assets/img/latest/f_211_193_171pfss_512.jpg")
	fnCustom1 := func(amount float64) (matrix [][]float64) {
		return [][]float64{
			{-0.5 / amount / 4.0, 1 / amount / 6.0, -0.5 / amount / 4.0},
			{1 / amount / 8.0, math.Clamp(amount*amount, 0.0, 1.5), 1 / amount / 8.0},
			{-0.5 / amount / 4.0, 1 / amount / 6.0, -0.5 / amount / 4.0},
		}
	}
	fnCustom2 := func(amount float64) (matrix [][]float64) {
		return [][]float64{
			{0.5 / amount / 4.0, 1 / amount / 6.0, 0.5 / amount / 4.0},
			{1 / amount / 8.0, math.Clamp(amount*amount, 0.0, 1.5), 1 / amount / 8.0},
			{0.5 / amount / 4.0, 1 / amount / 6.0, 0.5 / amount / 4.0},
		}
	}
	fnCustomIdentity := func(amount float64) (matrix [][]float64) {
		return [][]float64{
			{0, 0, 0},
			{0, 1, 0},
			{0, 0, 0},
		}
	}
	tests := []struct {
		name   string
		Matrix *ConvolutionMatrix
		src    *image.Image
	}{
		{"blur-0.25", NewBlurFilter(0.25), testImage.Clone()},
		{"blur-0.50", NewBlurFilter(0.50), testImage.Clone()},
		{"blur-0.75", NewBlurFilter(0.75), testImage.Clone()},
		{"blur-1.00", NewBlurFilter(1.00), testImage.Clone()},
		{"blur-2.00", NewBlurFilter(2.00), testImage.Clone()},
		{"sharpen-0.25", NewSharpenFilter(0.25), testImage.Clone()},
		{"sharpen-0.50", NewSharpenFilter(0.50), testImage.Clone()},
		{"sharpen-0.75", NewSharpenFilter(0.75), testImage.Clone()},
		{"sharpen-1.00", NewSharpenFilter(1.00), testImage.Clone()},
		{"sharpen-2.00", NewSharpenFilter(2.00), testImage.Clone()},
		{"edge-detection-0.25", NewEdgeDetectFilter(0.25), testImage.Clone()},
		{"edge-detection-0.50", NewEdgeDetectFilter(0.50), testImage.Clone()},
		{"edge-detection-0.75", NewEdgeDetectFilter(0.75), testImage.Clone()},
		{"edge-detection-1.00", NewEdgeDetectFilter(1.00), testImage.Clone()},
		{"edge-detection-2.00", NewEdgeDetectFilter(2.00), testImage.Clone()},
		{"emboss-0.25", NewEmbossFilter(0.25), testImage.Clone()},
		{"emboss-0.50", NewEmbossFilter(0.50), testImage.Clone()},
		{"emboss-0.75", NewEmbossFilter(0.75), testImage.Clone()},
		{"emboss-1.00", NewEmbossFilter(1.00), testImage.Clone()},
		{"emboss-2.00", NewEmbossFilter(2.00), testImage.Clone()},
		{"custom1-0.25", NewCustomFilter(0.25, 1.1, -0.125, fnCustom1), testImage.Clone()},
		{"custom1-0.50", NewCustomFilter(0.50, 1.1, -0.125, fnCustom1), testImage.Clone()},
		{"custom1-0.75", NewCustomFilter(0.75, 1.1, -0.125, fnCustom1), testImage.Clone()},
		{"custom1-1.00", NewCustomFilter(1.00, 1.1, -0.125, fnCustom1), testImage.Clone()},
		{"custom1-2.00", NewCustomFilter(2.00, 1.1, -0.125, fnCustom1), testImage.Clone()},
		{"custom2-0.25", NewCustomFilter(0.25, 1.1, -0.125, fnCustom2), testImage.Clone()},
		{"custom2-0.50", NewCustomFilter(0.50, 1.1, -0.125, fnCustom2), testImage.Clone()},
		{"custom2-0.75", NewCustomFilter(0.75, 1.1, -0.125, fnCustom2), testImage.Clone()},
		{"custom2-1.00", NewCustomFilter(1.00, 1.1, -0.125, fnCustom2), testImage.Clone()},
		{"custom2-2.00", NewCustomFilter(2.00, 1.1, -0.125, fnCustom2), testImage.Clone()},
		{"custom-identity", NewCustomFilter(1.00, 1.0, 0.0, fnCustomIdentity), testImage.Clone()},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.Matrix.Apply(tt.src).SaveAsPNG("../../test_data/filters/convolution/" + tt.name + ".png")
		})
	}
}
