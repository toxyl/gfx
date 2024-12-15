package parser

import (
	"fmt"
	"strconv"
	"strings"
)

type Resize struct {
	W int `yaml:"w,omitempty"`
	H int `yaml:"h,omitempty"`
}

func (r *Resize) String() string {
	return fmt.Sprintf("%s %4d %4d", LAYER_RESIZE, r.W, r.H)
}

func parseResize(value string) Resize {
	parts := strings.Split(value, STR_SPACE)
	w, _ := strconv.Atoi(parts[0])
	h, _ := strconv.Atoi(parts[1])
	return Resize{W: w, H: h}
}
