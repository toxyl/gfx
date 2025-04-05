package image

import (
	"runtime"
	"sync"

	"github.com/toxyl/gfx/core/color"
)

func (i *Image) ProcessHSLA(startX, startY, endX, endY int, fn func(x, y int, col *color.HSL) (x2, y2 int, col2 *color.HSL)) *Image {
	numCores := runtime.NumCPU()
	sem := make(chan struct{}, numCores) // Limit goroutines to the number of cores.
	var wg sync.WaitGroup
	dst := i.Clone()
	for y := startY; y < endY; y++ {
		wg.Add(1)
		sem <- struct{}{} // Acquire a slot.
		go func(row int) {
			defer func() {
				<-sem // Release the slot.
				wg.Done()
			}()
			for x := startX; x < endX; x++ {
				dst.SetHSLA(fn(x, row, color.HSLFromRGB(dst.GetRGBA64(x, row))))
			}
		}(y)
	}
	wg.Wait()
	i.Set(dst.raw)
	return i
}

// ProcessRGBA processes RGBA pixels using goroutines per row.
func (i *Image) ProcessRGBA(startX, startY, endX, endY int, fn func(x, y int, col *color.RGBA64) (x2, y2 int, col2 *color.RGBA64)) *Image {
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
				x2, y2, c := fn(x, row, dst.GetRGBA64(x, row))
				dst.SetRGBA(x2, y2, c)
			}
		}(y)
	}
	wg.Wait()
	i.Set(dst.raw)
	return i
}
