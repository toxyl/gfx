package models

import (
	"image"
	"image/color"

	"github.com/toxyl/gfx/core/fx"
	"github.com/toxyl/gfx/core/meta"
	"github.com/toxyl/math"
)

// Enhance represents an image enhancement effect.
type Enhance struct {
	Amount float64 // Enhancement amount (0.0 to 1.0)
	meta   *fx.EffectMeta
}

// NewEnhanceEffect creates a new enhance effect.
func NewEnhanceEffect(amount float64) *Enhance {
	e := &Enhance{
		Amount: amount,
		meta: fx.NewEffectMeta(
			"Enhance",
			"Enhances image details and contrast",
			meta.NewChannelMeta("Amount", 0.0, 1.0, "", "Enhancement amount (0.0 to 1.0)"),
		),
	}
	e.Amount = fx.ClampParameter(amount, e.meta.Parameters[0])
	return e
}

// Apply applies the enhance effect to an image.
func (e *Enhance) Apply(img image.Image) image.Image {
	bounds := img.Bounds()
	dst := image.NewRGBA(bounds)

	// Unsharp mask parameters
	radius := 1
	threshold := 0.0
	amount := e.Amount * 2.0 // Scale amount for better effect

	// Create blurred version
	blurred := image.NewRGBA(bounds)
	for y := bounds.Min.Y + radius; y < bounds.Max.Y-radius; y++ {
		for x := bounds.Min.X + radius; x < bounds.Max.X-radius; x++ {
			var r, g, b float64
			count := 0

			// Apply box blur
			for ky := -radius; ky <= radius; ky++ {
				for kx := -radius; kx <= radius; kx++ {
					px := x + kx
					py := y + ky
					pr, pg, pb, _ := img.At(px, py).RGBA()

					r += float64(pr) / 0xFFFF
					g += float64(pg) / 0xFFFF
					b += float64(pb) / 0xFFFF
					count++
				}
			}

			r /= float64(count)
			g /= float64(count)
			b /= float64(count)

			blurred.Set(x, y, color.RGBA64{
				R: uint16(r * 0xFFFF),
				G: uint16(g * 0xFFFF),
				B: uint16(b * 0xFFFF),
				A: 0xFFFF,
			})
		}
	}

	// Apply unsharp mask
	for y := bounds.Min.Y; y < bounds.Max.Y; y++ {
		for x := bounds.Min.X; x < bounds.Max.X; x++ {
			// Get original and blurred colors
			or, og, ob, oa := img.At(x, y).RGBA()
			br, bg, bb, _ := blurred.At(x, y).RGBA()

			// Convert to float64 and normalize
			orf := float64(or) / 0xFFFF
			ogf := float64(og) / 0xFFFF
			obf := float64(ob) / 0xFFFF
			brf := float64(br) / 0xFFFF
			bgf := float64(bg) / 0xFFFF
			bbf := float64(bb) / 0xFFFF

			// Calculate difference
			dr := orf - brf
			dg := ogf - bgf
			db := obf - bbf

			// Apply threshold
			if math.Abs(dr) < threshold {
				dr = 0
			}
			if math.Abs(dg) < threshold {
				dg = 0
			}
			if math.Abs(db) < threshold {
				db = 0
			}

			// Apply enhancement
			r := orf + dr*amount
			g := ogf + dg*amount
			b := obf + db*amount

			// Clamp values
			r = math.Max(0, math.Min(1, r))
			g = math.Max(0, math.Min(1, g))
			b = math.Max(0, math.Min(1, b))

			dst.Set(x, y, color.RGBA64{
				R: uint16(r * 0xFFFF),
				G: uint16(g * 0xFFFF),
				B: uint16(b * 0xFFFF),
				A: uint16(oa),
			})
		}
	}

	return dst
}

// Meta returns the effect metadata.
func (e *Enhance) Meta() *fx.EffectMeta {
	return e.meta
}
