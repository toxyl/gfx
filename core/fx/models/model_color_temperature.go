package models

import (
	"image"
	"image/color"

	"github.com/toxyl/gfx/core/fx"
	"github.com/toxyl/gfx/core/meta"
	"github.com/toxyl/math"
)

// ColorTemperature represents a color temperature adjustment effect.
type ColorTemperature struct {
	Temperature float64 // Color temperature in Kelvin (1000 to 40000)
	meta        *fx.EffectMeta
}

// NewColorTemperatureEffect creates a new color temperature adjustment effect.
func NewColorTemperatureEffect(temperature float64) *ColorTemperature {
	ct := &ColorTemperature{
		Temperature: temperature,
		meta: fx.NewEffectMeta(
			"Color Temperature",
			"Adjusts the color temperature of an image",
			meta.NewChannelMeta("Temperature", 1000.0, 40000.0, "K", "Color temperature in Kelvin (1000 to 40000)"),
		),
	}
	ct.Temperature = fx.ClampParameter(temperature, ct.meta.Parameters[0])
	return ct
}

// kelvinToRGB converts a color temperature in Kelvin to RGB values.
func (ct *ColorTemperature) kelvinToRGB() (r, g, b float64) {
	// Convert temperature to RGB using the black body radiation formula
	temp := ct.Temperature / 100.0

	if temp <= 66.0 {
		r = 255.0
		g = 99.4708025861*math.Log(temp) - 161.1195681661
		if temp <= 19.0 {
			b = 0.0
		} else {
			b = 138.5177312231*math.Log(temp-10.0) - 305.0447927307
		}
	} else {
		r = 329.698727446 * math.Pow(temp-60.0, -0.1332047592)
		g = 288.1221695283 * math.Pow(temp-60.0, -0.0755148492)
		b = 255.0
	}

	// Normalize to 0-1 range
	r = math.Max(0, math.Min(1, r/255.0))
	g = math.Max(0, math.Min(1, g/255.0))
	b = math.Max(0, math.Min(1, b/255.0))

	return r, g, b
}

// Apply applies the color temperature adjustment effect to an image.
func (ct *ColorTemperature) Apply(img image.Image) image.Image {
	bounds := img.Bounds()
	dst := image.NewRGBA(bounds)

	// Get the target RGB values for the temperature
	targetR, targetG, targetB := ct.kelvinToRGB()

	// Calculate the average RGB values of the image
	var avgR, avgG, avgB float64
	pixelCount := float64(bounds.Dx() * bounds.Dy())

	for y := bounds.Min.Y; y < bounds.Max.Y; y++ {
		for x := bounds.Min.X; x < bounds.Max.X; x++ {
			r, g, b, _ := img.At(x, y).RGBA()
			avgR += float64(r) / 0xFFFF
			avgG += float64(g) / 0xFFFF
			avgB += float64(b) / 0xFFFF
		}
	}

	avgR /= pixelCount
	avgG /= pixelCount
	avgB /= pixelCount

	// Calculate adjustment factors
	adjR := targetR / avgR
	adjG := targetG / avgG
	adjB := targetB / avgB

	for y := bounds.Min.Y; y < bounds.Max.Y; y++ {
		for x := bounds.Min.X; x < bounds.Max.X; x++ {
			r, g, b, a := img.At(x, y).RGBA()

			// Convert to float64 for calculations
			rf := float64(r) / 0xFFFF
			gf := float64(g) / 0xFFFF
			bf := float64(b) / 0xFFFF

			// Apply color temperature adjustment
			rf = math.Max(0, math.Min(1, rf*adjR))
			gf = math.Max(0, math.Min(1, gf*adjG))
			bf = math.Max(0, math.Min(1, bf*adjB))

			// Convert back to uint32
			r = uint32(math.Max(0, math.Min(0xFFFF, rf*0xFFFF)))
			g = uint32(math.Max(0, math.Min(0xFFFF, gf*0xFFFF)))
			b = uint32(math.Max(0, math.Min(0xFFFF, bf*0xFFFF)))

			dst.Set(x, y, color.RGBA64{
				R: uint16(r),
				G: uint16(g),
				B: uint16(b),
				A: uint16(a),
			})
		}
	}

	return dst
}

// Meta returns the effect metadata.
func (ct *ColorTemperature) Meta() *fx.EffectMeta {
	return ct.meta
}
