package models

import (
	"image"
	"image/color"

	gfxcolor "github.com/toxyl/gfx/core/color"
	"github.com/toxyl/gfx/core/fx"
	"github.com/toxyl/gfx/core/meta"
	"github.com/toxyl/math"
)

// Vibrance represents a vibrance adjustment effect.
type Vibrance struct {
	Amount float64 // Vibrance adjustment value (-1.0 to 1.0)
	meta   *fx.EffectMeta
}

// NewVibranceEffect creates a new vibrance adjustment effect.
func NewVibranceEffect(amount float64) *Vibrance {
	v := &Vibrance{
		Amount: amount,
		meta: fx.NewEffectMeta(
			"Vibrance",
			"Adjusts the saturation of an image while protecting skin tones",
			meta.NewChannelMeta("Amount", -1.0, 1.0, "", "Vibrance adjustment value (-1.0 to 1.0)"),
		),
	}
	v.Amount = fx.ClampParameter(amount, v.meta.Parameters[0])
	return v
}

// isSkinTone returns true if the color is likely a skin tone.
func (v *Vibrance) isSkinTone(r, g, b float64) bool {
	// Convert to YCbCr color space
	y := 0.299*r + 0.587*g + 0.114*b
	cb := -0.168736*r - 0.331264*g + 0.5*b
	cr := 0.5*r - 0.418688*g - 0.081312*b

	// Check if the color is within typical skin tone ranges
	return y > 0.2 && y < 0.8 && // Luminance range
		cb > -0.1 && cb < 0.1 && // Blue-yellow range
		cr > 0.1 && cr < 0.3 // Red-green range
}

// Apply applies the vibrance adjustment effect to an image.
func (v *Vibrance) Apply(img image.Image) (image.Image, error) {
	bounds := img.Bounds()
	dst := image.NewRGBA(bounds)

	for y := bounds.Min.Y; y < bounds.Max.Y; y++ {
		for x := bounds.Min.X; x < bounds.Max.X; x++ {
			r, g, b, a := img.At(x, y).RGBA()

			// Convert to float64 for calculations
			rf := float64(r) / 0xFFFF
			gf := float64(g) / 0xFFFF
			bf := float64(b) / 0xFFFF

			// Calculate saturation
			max := math.Max(math.Max(rf, gf), bf)
			min := math.Min(math.Min(rf, gf), bf)
			saturation := (max - min) / max

			// Apply vibrance adjustment with skin tone protection
			if !v.isSkinTone(rf, gf, bf) || v.Amount < 0 {
				// Calculate adjustment factor based on saturation
				adjustment := 1.0 + v.Amount*(1.0-saturation)

				// Apply adjustment to each channel
				rf = math.Max(0, math.Min(1, rf*adjustment))
				gf = math.Max(0, math.Min(1, gf*adjustment))
				bf = math.Max(0, math.Min(1, bf*adjustment))
			}

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

	return dst, nil
}

// Meta returns the effect metadata.
func (v *Vibrance) Meta() *fx.EffectMeta {
	return v.meta
}

// Description returns the description of the vibrance effect.
func (v *Vibrance) Description() string {
	return v.meta.Description
}

// GetArgs returns the effect arguments.
func (v *Vibrance) GetArgs() []fx.FunctionArg {
	return []fx.FunctionArg{
		{
			Name:        "Amount",
			Type:        "float64",
			Description: "Vibrance adjustment value (-1.0 to 1.0)",
			Value:       v.Amount,
		},
	}
}

// Name returns the name of the effect.
func (v *Vibrance) Name() string {
	return v.meta.Name
}

// GetColorModel returns the color model used by this effect.
func (v *Vibrance) GetColorModel() *gfxcolor.Color64 {
	return nil // Vibrance effect doesn't modify the color model
}

// ValidateArgs validates the effect arguments.
func (v *Vibrance) ValidateArgs(args ...any) error {
	if len(args) != 1 {
		return fx.ErrInvalidArgument
	}
	amount, ok := args[0].(float64)
	if !ok {
		return fx.ErrInvalidArgument
	}
	if amount < -1.0 || amount > 1.0 {
		return fx.ErrInvalidArgument
	}
	return nil
}
