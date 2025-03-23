package parser

import (
	"fmt"
)

type Offset struct {
	X int `yaml:"x,omitempty"`
	Y int `yaml:"y,omitempty"`
}

func (o *Offset) String() string {
	return fmt.Sprintf("%s %4d %4d", LAYER_OFFSET, o.X, o.Y)
}
