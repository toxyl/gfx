package blendmodes

import (
	"strings"

	"github.com/toxyl/gfx/core/blendmodes/constants"
)

// IMeta contains metadata about a blend mode.
type IMeta struct {
	name        string
	description string
	category    string
}

// NewIMeta creates a new blend mode metadata instance.
func NewIMeta(name, description, category string) (*IMeta, error) {
	if err := validateBlendModeName(name); err != nil {
		return nil, err
	}
	if err := validateCategory(category); err != nil {
		return nil, err
	}
	return &IMeta{
		name:        name,
		description: description,
		category:    category,
	}, nil
}

// Name returns the name of the blend mode.
func (m *IMeta) Name() string {
	return m.name
}

// Description returns the description of the blend mode.
func (m *IMeta) Description() string {
	return m.description
}

// Category returns the category of the blend mode.
func (m *IMeta) Category() string {
	return m.category
}

// Doc returns a markdown-formatted documentation string for the blend mode
func (b *IMeta) Doc() string {
	sb := strings.Builder{}
	sb.WriteString("| ")
	sb.WriteString(b.name)
	sb.WriteString(" | ")
	sb.WriteString(b.category)
	sb.WriteString(" | ")
	if strings.TrimSpace(b.description) != "" {
		sb.WriteString(b.description)
	}
	sb.WriteString(" |")
	return sb.String()
}

// IsComponentMode returns true if the blend mode is in the component category
func (b *IMeta) IsComponentMode() bool {
	return b.category == constants.CategoryComponent
}
