package fx

import (
	"fmt"
	"strings"

	"github.com/toxyl/gfx/core/meta"
)

// EffectMeta holds metadata for an effect.
type EffectMeta struct {
	Name        string              // Name of the effect (e.g., "Brightness", "Blur")
	Description string              // Description of the effect
	Parameters  []*meta.ChannelMeta // Metadata for each parameter
}

// NewEffectMeta creates a new EffectMeta instance.
// It validates the input parameters and ensures they are valid.
//
// Parameters:
//   - name: The name of the effect
//   - description: A description of the effect
//   - parameters: Metadata for each parameter
//
// Returns:
//   - A pointer to a new EffectMeta instance
func NewEffectMeta(name, description string, parameters ...*meta.ChannelMeta) *EffectMeta {
	if name == "" {
		panic("effect name cannot be empty")
	}
	if description == "" {
		panic("effect description cannot be empty")
	}
	return &EffectMeta{
		Name:        name,
		Description: description,
		Parameters:  parameters,
	}
}

// Doc returns a Markdown formatted documentation string for the effect.
// The documentation includes the effect name, description, and a table of parameter metadata.
//
// Returns:
//   - A Markdown formatted string
func (m *EffectMeta) Doc() string {
	var sb strings.Builder
	sb.WriteString("# " + m.Name + "\n")
	sb.WriteString("_" + m.Description + "_\n\n")
	if len(m.Parameters) > 0 {
		sb.WriteString("## Parameters\n\n")
		sb.WriteString("| Parameter | Min | Max | Unit | Description |\n")
		sb.WriteString("| --- | --- | --- | --- | --- |\n")
		for _, param := range m.Parameters {
			sb.WriteString(fmt.Sprintf("| %s | %f | %f | %s | %s |\n",
				param.Name, param.Min, param.Max, param.Unit, param.Description))
		}
	}
	return sb.String()
}

// ValidateParameter checks if a value is within the valid range for a parameter.
// Returns an error if the value is outside the valid range.
//
// Parameters:
//   - value: The value to validate
//   - param: The parameter metadata
//
// Returns:
//   - An error if the value is invalid, nil otherwise
func ValidateParameter(value float64, param *meta.ChannelMeta) error {
	return meta.ValidateChannelValue(value, param)
}

// ClampParameter clamps a value to the valid range for a parameter.
// Returns the clamped value.
//
// Parameters:
//   - value: The value to clamp
//   - param: The parameter metadata
//
// Returns:
//   - The clamped value
func ClampParameter(value float64, param *meta.ChannelMeta) float64 {
	return meta.ClampChannelValue(value, param)
}
