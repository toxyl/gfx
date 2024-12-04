package image

import (
	"runtime"
	"sync"

	"github.com/toxyl/gfx/color/hsla"
	"github.com/toxyl/gfx/color/rgba"
)

func (i *Image) ProcessHSLA(startX, startY, endX, endY int, fn func(x, y int, col *hsla.HSLA) (x2, y2 int, col2 *hsla.HSLA)) *Image {
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
				dst.SetHSLA(fn(x, row, dst.GetHSLA(x, row)))
			}
		}(y)
	}
	wg.Wait()
	i.Set(dst.raw)
	return i
}

// ProcessRGBA processes RGBA pixels using goroutines per row.
func (i *Image) ProcessRGBA(startX, startY, endX, endY int, fn func(x, y int, col *rgba.RGBA) (x2, y2 int, col2 *rgba.RGBA)) *Image {
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
				dst.SetRGBA(fn(x, row, dst.GetRGBA(x, row)))
			}
		}(y)
	}
	wg.Wait()
	i.Set(dst.raw)
	return i
}
