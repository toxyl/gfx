package text

import (
	"fmt"

	"github.com/toxyl/gfx/core/blendmodes"
	"github.com/toxyl/gfx/core/image"
)

func Draw(img *image.Image, x, y int, text string, mode *blendmodes.IBlendMode) error {
	// The DrawText method is not implemented in core/image yet
	// This is a temporary implementation that returns an error
	return fmt.Errorf("DrawText not implemented in core/image package yet")
}
