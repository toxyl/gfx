package image

import (
	"runtime"
	"sync"

	"github.com/toxyl/gfx/color/hsla"
)

func (i *Image) mergeHSLA(src *Image, startX, startY, endX, endY int, fn func(x, y int, col *hsla.HSLA) (x2, y2 int, col2 *hsla.HSLA)) *Image {
	numCores := runtime.NumCPU()
	sem := make(chan struct{}, numCores)
	var wg sync.WaitGroup
	dst := i.Clone()

	for y := startY; y < endY; y++ {
		wg.Add(1)
		sem <- struct{}{}
		go func(row int) {
			defer func() {
				<-sem
				wg.Done()
			}()
			for x := startX; x < endX; x++ {
				dst.SetHSLA(fn(x, row, src.GetHSLA(x, row)))
			}
		}(y)
	}
	wg.Wait()
	i.Set(dst.raw)
	return i
}
