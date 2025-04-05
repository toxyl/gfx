package dialogbox

import (
	"fmt"

	"github.com/toxyl/gfx/core/blendmodes"
	"github.com/toxyl/gfx/core/image"
)

// TODO: This is a placeholder implementation until the core/image package
// provides similar drawing functionality as the legacy package.
func Draw(img *image.Image, x, y, w, h int, title string, content *image.Image, mode *blendmodes.IBlendMode) error {
	// For now, just blend the content image onto the main image

	// Calculate content position
	contentX := x + 2
	contentY := y + (h - content.Height()) - 2

	// Draw the content onto the main image
	err := img.DrawImage(content, contentX, contentY, mode, 1.0)
	if err != nil {
		return fmt.Errorf("failed to draw content: %w", err)
	}

	return nil
}
