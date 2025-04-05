package buildlog

import (
	"fmt"
	"time"
)

func Log(action string, fn func()) {
	t := time.Now()
	fmt.Print(action + " ... ")
	defer func() { fmt.Println("done in " + time.Since(t).String() + "!") }()
	fn()
}
