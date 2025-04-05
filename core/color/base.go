// Package color provides a comprehensive set of color models and conversion utilities.
// This file contains the base interface and utilities used by all color models.
package color

import (
	"fmt"

	"github.com/toxyl/math"
)

// iColor is the internal interface for color models.
type iColor interface {
	// Meta returns metadata about the color model.
	Meta() *ColorModelMeta
	// ToRGBA64 converts the color to RGBA64.
	ToRGBA64() *RGBA64
	// FromRGBA64 converts an RGBA64 color to this color model.
	FromRGBA64(rgba *RGBA64) iColor
	// FromSlice initializes the color from a slice of float64 values.
	FromSlice(values []float64) error
}

// newColor creates a new instance of a color model and initializes it with the given values.
//
// Parameters:
//   - factory: A function that creates a new instance of the color model
//   - values: The channel values to initialize the model with
//
// Returns:
//   - A new instance of the color model
//   - An error if initialization fails
func newColor[N math.Number, M iColor](factory func() M, values ...N) (M, error) {
	model := factory()
	vals := make([]float64, len(values))
	for i, v := range values {
		vals[i] = float64(v)
	}
	if err := model.FromSlice(vals); err != nil {
		return model, fmt.Errorf("failed to initialize model: %w", err)
	}
	return model, nil
}

// FromSlice initializes a color model from a slice of float64 values.
func FromSlice[M iColor](model M, values []float64) error {
	return model.FromSlice(values)
}

// NewChannelMeta creates a new ChannelMeta with the given name, min, max, unit, and description.
func NewChannelMeta[N math.Number](name string, min, max N, unit, description string) *Channel {
	return NewChannel(name, float64(min), float64(max), unit, description)
}

// NewModelMeta creates a new ColorModelMeta with the given name, description, and channels.
func NewModelMeta(name, description string, channels ...*Channel) *ColorModelMeta {
	return NewColorModelMeta(name, description, channels...)
}

// ValidateChannelValue validates that a value is within the valid range for a channel.
func ValidateChannelValue(value float64, channel *Channel) error {
	if value < channel.Min() || value > channel.Max() {
		return fmt.Errorf("value %f is outside valid range [%f, %f] for channel %s", value, channel.Min(), channel.Max(), channel.Name())
	}
	return nil
}

// ClampChannelValue restricts value to be within the channel's range.
func ClampChannelValue(value float64, channel *Channel) float64 {
	return math.Clamp(value, channel.Min(), channel.Max())
}
