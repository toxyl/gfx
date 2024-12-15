package parser

import (
	"fmt"
	"strconv"
	"strings"
)

type Offset struct {
	X int `yaml:"x,omitempty"`
	Y int `yaml:"y,omitempty"`
}

func (o *Offset) String() string {
	return fmt.Sprintf("%s %4d %4d", LAYER_OFFSET, o.X, o.Y)
}

func parseOffset(value string) Offset {
	parts := strings.Split(value, STR_SPACE)
	x, _ := strconv.Atoi(parts[0])
	y, _ := strconv.Atoi(parts[1])
	return Offset{X: x, Y: y}
}
