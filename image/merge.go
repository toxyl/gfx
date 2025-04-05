package image

import (
	"runtime"
	"sync"

	"github.com/toxyl/gfx/core/color"
)

func (i *Image) mergeHSLA(src *Image, startX, startY, endX, endY int, fn func(x, y int, col *color.HSL) (x2, y2 int, col2 *color.HSL)) *Image {
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
				dst.SetHSLA(fn(x, row, color.HSLFromRGB(src.GetRGBA64(x, row))))
			}
		}(y)
	}
	wg.Wait()
	i.Set(dst.raw)
	return i
}
