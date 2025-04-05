package parser

import (
	"fmt"

	"github.com/toxyl/gfx/config/constants"
)

type Var struct {
	Name  string
	Value any
}

func (v *Var) String() string {
	return fmt.Sprintf("%s %s %v", v.Name, constants.ASSIGN, v.Value)
}
