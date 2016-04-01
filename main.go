package main

import (
	"fmt"
	"runtime"
	"sync"
	"time"

	sdl "github.com/veandco/go-sdl2/sdl"
)

const (
	windowWidth  = 800
	windowHeight = 600
)

var (
	maxIterations = 32

	moveX = -0.5
	moveY = 0.0
	zoom  = 1.0

	running = true

	keyState = map[sdl.Keycode]bool{
		sdl.K_DOWN:  false,
		sdl.K_UP:    false,
		sdl.K_LEFT:  false,
		sdl.K_RIGHT: false,

		sdl.K_z: false,
		sdl.K_x: false,
	}

	gridX = 1
	gridY = 1

	maxIterationsMutex = sync.Mutex{}
	wg                 = sync.WaitGroup{}
)

type BrotSurface struct {
	surface       *sdl.Surface
	x, y          int32
	width, height int32
	retChannel    chan BrotSurface
}

func renderToSurface(brotSurface BrotSurface) {
	beginX := brotSurface.x
	beginY := brotSurface.y
	width := brotSurface.width
	height := brotSurface.height
	surface := brotSurface.surface

	halfWidth := float64(windowWidth) / 2.0
	halfHeight := float64(windowHeight) / 2.0
	halfZoom := zoom / 2.0
	halfZoomWidth := halfZoom * float64(windowWidth)
	halfZoomHeight := halfZoom * float64(windowHeight)
	pixels := surface.Pixels()
	pixelIndex := 0
	localMaxIterations := maxIterations

	for y := beginY; y < beginY+height; y++ {
		for x := beginX; x < beginX+width; x++ {
			pr := 1.5*(float64(x)-halfWidth)/(halfZoomWidth) + moveX
			pi := (float64(y)-halfHeight)/(halfZoomHeight) + moveY

			newRe := 0.0
			newIm := 0.0
			oldRe := 0.0
			oldIm := 0.0

			i := 0
			for ; i < localMaxIterations; i++ {
				oldRe = newRe
				oldIm = newIm

				newRe = oldRe*oldRe - oldIm*oldIm + pr
				newIm = 2*oldRe*oldIm + pi
				if (newRe*newRe + newIm*newIm) > 4 {
					break
				}
			}

			colourIndexFactor := float64(i) / float64(localMaxIterations)
			colourIndex := int(float64(len(colours)-1) * colourIndexFactor)
			colour := colours[colourIndex]
			pixels[pixelIndex+0] = byte(colour.b)
			pixels[pixelIndex+1] = byte(colour.g)
			pixels[pixelIndex+2] = byte(colour.r)
			pixelIndex += 4
		}
	}

	brotSurface.retChannel <- brotSurface
}

func eventPoll(eventChan chan sdl.Event) {
	wg.Add(1)
	defer wg.Done()

	for event := sdl.WaitEvent(); event != nil; event = sdl.WaitEvent() {
		eventChan <- event

		if !running {
			break
		}
	}
}

func handleInput() {
	movementValue := 0.1 / zoom
	if keyState[sdl.K_DOWN] {
		moveY += movementValue
	}
	if keyState[sdl.K_UP] {
		moveY -= movementValue
	}
	if keyState[sdl.K_RIGHT] {
		moveX += movementValue
	}
	if keyState[sdl.K_LEFT] {
		moveX -= movementValue
	}

	zoomValue := 0.1 * zoom
	if keyState[sdl.K_z] {
		zoom += zoomValue
	}
	if keyState[sdl.K_x] {
		zoom -= zoomValue
	}

	if keyState[sdl.K_a] {
		incMaxIterations(1)
	}
	if keyState[sdl.K_s] {
		incMaxIterations(-1)
	}

	if keyState[sdl.K_SPACE] {
		moveX = -0.5
		moveY = 0.0
		zoom = 1.0
		maxIterations = 32
	}
}

func incMaxIterations(value int) {
	maxIterations += value
	if maxIterations < 1 {
		maxIterations = 1
	}
}

func handleEvent(event *sdl.Event) {
	switch t := (*event).(type) {
	case *sdl.QuitEvent:
		running = false
		break
	case *sdl.KeyDownEvent:
		keyState[t.Keysym.Sym] = true
		break
	case *sdl.KeyUpEvent:
		keyState[t.Keysym.Sym] = false
		break
	}
}

func maxParallelism() int {
	maxProcs := runtime.GOMAXPROCS(0)
	numCPU := runtime.NumCPU()
	if maxProcs < numCPU {
		return maxProcs
	}
	return numCPU
}

func initParallelism() {
	numCpus := maxParallelism()
	runtime.GOMAXPROCS(numCpus)
	gridX = numCpus
	gridY = numCpus
}

func main() {
	initParallelism()
	sdl.Init(sdl.INIT_EVERYTHING)

	window, err := sdl.CreateWindow("mandelbrot", sdl.WINDOWPOS_UNDEFINED, sdl.WINDOWPOS_UNDEFINED,
		windowWidth, windowHeight, sdl.WINDOW_SHOWN)
	defer window.Destroy()

	windowSurface, err := window.GetSurface()
	if err != nil {
		panic(err)
	}

	brotSurfaces := []BrotSurface{}
	brotSurfaceChannel := make(chan BrotSurface)
	for y := 0; y < gridY; y++ {
		for x := 0; x < gridX; x++ {
			width := windowWidth / gridX
			beginX := x * width

			height := windowHeight / gridY
			beginY := y * height

			brotSurface := BrotSurface{}
			brotSurface.surface, _ = sdl.CreateRGBSurface(0, int32(width), int32(height), 32, 0, 0, 0, 0)
			brotSurface.x = int32(beginX)
			brotSurface.y = int32(beginY)
			brotSurface.width = int32(width)
			brotSurface.height = int32(height)
			brotSurface.retChannel = brotSurfaceChannel

			brotSurfaces = append(brotSurfaces, brotSurface)
		}
	}

	eventChannel := make(chan sdl.Event)
	go eventPoll(eventChannel)
	go func() {
		wg.Add(1)
		defer wg.Done()

		for range time.Tick(16 * time.Millisecond) {
			for i := 0; i < gridX*gridY; i++ {
				go renderToSurface(brotSurfaces[i])
			}
			for i := 0; i < gridX*gridY; i++ {
				brotSurface := <-brotSurfaceChannel
				brotSurface.surface.Blit(
					&sdl.Rect{0, 0, brotSurface.width, brotSurface.height},
					windowSurface,
					&sdl.Rect{brotSurface.x, brotSurface.y, brotSurface.width, brotSurface.height},
				)
			}
			window.UpdateSurface()

			if !running {
				return
			}
		}
	}()

	for running {
		select {
		case event := <-eventChannel:
			handleEvent(&event)
			break
		case <-time.NewTicker(16 * time.Millisecond).C:
			handleInput()
			break
		}
	}

	wg.Wait()
	sdl.Quit()
}
