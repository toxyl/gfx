// Package meta provides metadata types and utilities for color models and blend modes.
package meta

import (
	"fmt"
	"strings"

	"github.com/toxyl/math"
)

// ChannelMeta holds metadata for a single channel in a color model.
type ChannelMeta struct {
	Name        string  // Name of the channel (e.g., "R", "H", "L")
	Min         float64 // Minimum value for the channel
	Max         float64 // Maximum value for the channel
	Unit        string  // Unit of measurement (e.g., "", "°", "%")
	Description string  // Description of the channel
}

// NewChannelMeta creates a new ChannelMeta instance.
// It validates the input parameters and ensures they are valid.
//
// Parameters:
//   - name: The name of the channel
//   - min: The minimum value for the channel
//   - max: The maximum value for the channel
//   - unit: The unit of measurement
//   - description: A description of the channel
//
// Returns:
//   - A pointer to a new ChannelMeta instance
func NewChannelMeta[N math.Number](name string, min, max N, unit, description string) *ChannelMeta {
	if name == "" {
		panic("channel name cannot be empty")
	}
	if float64(min) >= float64(max) {
		panic(fmt.Sprintf("invalid channel range: min (%v) must be less than max (%v)", min, max))
	}
	if description == "" {
		panic("channel description cannot be empty")
	}
	return &ChannelMeta{
		Name:        name,
		Min:         float64(min),
		Max:         float64(max),
		Unit:        unit,
		Description: description,
	}
}

// ColorModelMeta holds metadata for an entire color model.
type ColorModelMeta struct {
	Name        string         // Name of the color model (e.g., "RGB", "HSL")
	Description string         // Description of the color model
	Channels    []*ChannelMeta // Metadata for each channel
}

// NewModelMeta creates a new ColorModelMeta instance.
// It validates the input parameters and ensures they are valid.
//
// Parameters:
//   - name: The name of the color model
//   - description: A description of the color model
//   - channels: Metadata for each channel in the model
//
// Returns:
//   - A pointer to a new ColorModelMeta instance
func NewModelMeta(name, description string, channels ...*ChannelMeta) *ColorModelMeta {
	if name == "" {
		panic("color model name cannot be empty")
	}
	if description == "" {
		panic("color model description cannot be empty")
	}
	if len(channels) == 0 {
		panic("color model must have at least one channel")
	}
	return &ColorModelMeta{
		Name:        name,
		Description: description,
		Channels:    channels,
	}
}

// Doc returns a Markdown formatted documentation string for the color model.
// The documentation includes the model name, description, and a table of channel metadata.
//
// Returns:
//   - A Markdown formatted string
func (m *ColorModelMeta) Doc() string {
	var sb strings.Builder
	sb.WriteString("# " + m.Name + "\n")
	sb.WriteString("_" + m.Description + "_\n\n")
	sb.WriteString("| Channel | Min | Max | Unit | Description |\n")
	sb.WriteString("| --- | --- | --- | --- | --- |\n")
	for _, ch := range m.Channels {
		sb.WriteString(fmt.Sprintf("| %s | %f | %f | %s | %s |\n",
			ch.Name, ch.Min, ch.Max, ch.Unit, ch.Description))
	}
	return sb.String()
}

// ValidateChannelValue checks if a value is within the valid range for a channel.
// Returns an error if the value is outside the valid range.
//
// Parameters:
//   - value: The value to validate
//   - channel: The channel metadata
//
// Returns:
//   - An error if the value is invalid, nil otherwise
func ValidateChannelValue(value float64, channel *ChannelMeta) error {
	if value < channel.Min || value > channel.Max {
		return fmt.Errorf("value %f is outside valid range [%f, %f] for channel %s",
			value, channel.Min, channel.Max, channel.Name)
	}
	return nil
}

// ClampChannelValue restricts value to be within the channel's range.
func ClampChannelValue(value float64, channel *ChannelMeta) float64 {
	return math.Clamp(value, channel.Min, channel.Max)
}

// BlendModeMeta holds metadata for a blend mode.
type BlendModeMeta struct {
	Name        string // Name of the blend mode (e.g., "Normal", "Multiply")
	Description string // Description of the blend mode
}

// NewBlendModeMeta creates a new BlendModeMeta instance.
// It validates the input parameters and ensures they are valid.
//
// Parameters:
//   - name: The name of the blend mode
//   - description: A description of the blend mode
//
// Returns:
//   - A pointer to a new BlendModeMeta instance
func NewBlendModeMeta(name, description string) *BlendModeMeta {
	if name == "" {
		panic("blend mode name cannot be empty")
	}
	if description == "" {
		panic("blend mode description cannot be empty")
	}
	return &BlendModeMeta{
		Name:        name,
		Description: description,
	}
}

// Doc returns a Markdown formatted documentation string for the blend mode.
// The documentation includes the blend mode name and description.
//
// Returns:
//   - A Markdown formatted string
func (m *BlendModeMeta) Doc() string {
	var sb strings.Builder
	sb.WriteString("# " + m.Name + "\n")
	sb.WriteString("_" + m.Description + "_\n")
	return sb.String()
}
