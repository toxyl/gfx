package color

import (
	"fmt"
	"strings"
)

// Channel holds metadata for a single channel in a color model.
type Channel struct {
	name string
	min  float64
	max  float64
	unit string
	desc string
}

// NewChannel creates a new Channel for a color model channel.
func NewChannel(name string, min, max float64, unit, desc string) *Channel {
	return &Channel{
		name: name,
		min:  min,
		max:  max,
		unit: unit,
		desc: desc,
	}
}

func (c *Channel) Name() string { return c.name }
func (c *Channel) Min() float64 { return c.min }
func (c *Channel) Max() float64 { return c.max }
func (c *Channel) Unit() string { return c.unit }
func (c *Channel) Desc() string { return c.desc }

// ColorModelMeta holds metadata for an entire color model.
type ColorModelMeta struct {
	name     string
	desc     string
	channels []*Channel
}

// NewColorModelMeta creates a new ColorModelMeta with the given name, description, and channels.
func NewColorModelMeta(name, desc string, channels ...*Channel) *ColorModelMeta {
	return &ColorModelMeta{
		name:     name,
		desc:     desc,
		channels: channels,
	}
}

// Name returns the name of the color model.
func (m *ColorModelMeta) Name() string {
	return m.name
}

// Desc returns the description of the color model.
func (m *ColorModelMeta) Desc() string {
	return m.desc
}

// Channels returns the metadata for all channels in the color model.
func (m *ColorModelMeta) Channels() []*Channel {
	return m.channels
}

// DocMarkdown returns a Markdown formatted documentation string for the color model.
func (m *ColorModelMeta) Doc() string {
	var sb strings.Builder
	sb.WriteString("# " + m.name + "\n")
	sb.WriteString("_" + m.desc + "_\n\n")
	sb.WriteString("| Channel | Min | Max | Unit | Desc |\n")
	sb.WriteString("| --- | --- | --- | --- | --- |\n")
	for _, ch := range m.channels {
		sb.WriteString(fmt.Sprintf("| %s | %f | %f | %s | %s |\n", ch.Name(), ch.Min(), ch.Max(), ch.Unit(), ch.Desc()))
	}
	return sb.String()
}
