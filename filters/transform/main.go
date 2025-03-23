// filters/transform/main.go
package transform

import (
	"github.com/toxyl/gfx/filters/meta"
	"github.com/toxyl/gfx/image"
)

var Meta = meta.New("transform", []*meta.FilterMetaDataArg{
	{Name: "transform-x", Default: 0.0},
	{Name: "transform-y", Default: 0.0},
	{Name: "rotate", Default: 0.0},
	{Name: "scale", Default: 0.0},
	{Name: "offset-x", Default: 0.0},
	{Name: "offset-y", Default: 0.0},
})

// Apply performs a composite transformation on the image in the following order:
//  1. Translate (with wrap-around)
//  2. Rotate (around the effective center)
//  3. Scale (around the effective center)
//
// The parameters:
//   - transform-x, transform-y: percentages (-1..1) to translate relative to the image center.
//   - rotate: rotation angle in degrees (-360..360).
//   - scale: percentage to scale the image; 0 means no change,
//     0.5 means scale up by 50% (factor 1.5), -0.5 means scale down by 50% (factor 0.5).
//   - offset-x, offset-y: offsets (-1..1) from the image center that define the rotation/scaling center.
func Apply(img *image.Image, translateX, translateY, rotate, scale, offsetX, offsetY float64) *image.Image {
	// Compute absolute translation offsets relative to the image center.
	absTx := int(translateX * float64(img.CW()))
	absTy := int(translateY * float64(img.CH()))
	// Compute effective center for rotation and scaling.
	cx := img.CW() + int(offsetX*float64(img.CW()))
	cy := img.CH() + int(offsetY*float64(img.CH()))
	// Compute scale factor: 1.0 means no change.
	factor := 1.0 + scale

	// First, apply translation with wrap-around.
	img.Set(img.Translate(absTx, absTy, true).Get())
	// Then, apply combined rotation and scaling.
	img.Set(img.TransformRotateScale(rotate, factor, cx, cy).Get())
	return img
}
