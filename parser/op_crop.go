package parser

import (
	"fmt"
	"strconv"
	"strings"
)

type Crop struct {
	X int `yaml:"x,omitempty"`
	Y int `yaml:"y,omitempty"`
	W int `yaml:"w,omitempty"`
	H int `yaml:"h,omitempty"`
}

func (c *Crop) String() string {
	return fmt.Sprintf("%s %4d %4d %4d %4d", LAYER_CROP, c.X, c.Y, c.W, c.H)
}

func parseCrop(value string) Crop {
	parts := strings.Split(value, STR_SPACE)
	x, _ := strconv.Atoi(parts[0])
	y, _ := strconv.Atoi(parts[1])
	w, _ := strconv.Atoi(parts[2])
	h, _ := strconv.Atoi(parts[3])
	return Crop{X: x, Y: y, W: w, H: h}
}
