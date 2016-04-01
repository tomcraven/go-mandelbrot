package main

import (
	"math/rand"
	"os"
)

type Colour struct {
	r, g, b int
}

var (
	colours []Colour

	themes = map[string]func(){
		"fire":          fire,
		"full-spectrum": fullSpectrum,
		"leaf":          leaf,
		"water":         water,
		"beach":         beach,
		"random":        random,
	}
)

func init() {
	colourTheme := "beach"
	if len(os.Args) > 1 {
		colourTheme = os.Args[1]
	}

	if themeFunc, ok := themes[colourTheme]; ok {
		themeFunc()
	} else {
		fullSpectrum()
	}
}

func random() {
	for i := 0; i < 512; i++ {
		colours = append(colours, Colour{
			rand.Intn(255),
			rand.Intn(255),
			rand.Intn(255),
		})
	}
}

func beach() {
	addColourRange(Colour{0, 0, 0}, Colour{39, 90, 99}, 4)
	addColourRange(Colour{39, 90, 99}, Colour{239, 255, 95}, 5)
	addColourRange(Colour{239, 255, 94}, Colour{39, 90, 99}, 5)
	addColourRange(Colour{39, 90, 99}, Colour{0, 0, 0}, 10)
}

func water() {
	addColourRange(Colour{0, 0, 0}, Colour{0, 191, 255}, 4)
	addColourRange(Colour{0, 191, 255}, Colour{0, 0, 0}, 4)
}

func leaf() {
	addColourRange(Colour{0, 0, 0}, Colour{0, 255, 106}, 4)
	addColourRange(Colour{0, 255, 106}, Colour{40, 130, 17}, 4)
	addColourRange(Colour{40, 130, 17}, Colour{58, 87, 50}, 10)
	addColourRange(Colour{58, 87, 50}, Colour{0, 0, 0}, 10)
}

func fire() {
	addColourRange(Colour{0, 0, 0}, Colour{255, 0, 0}, 4)
	addColourRange(Colour{255, 0, 0}, Colour{255, 255, 0}, 4)
	addColourRange(Colour{255, 255, 0}, Colour{0, 0, 255}, 10)
	addColourRange(Colour{0, 0, 255}, Colour{0, 0, 0}, 50)
}

func fullSpectrum() {
	addColourRange(Colour{0, 0, 0}, Colour{0, 0, 255}, 1)
	addColourRange(Colour{0, 0, 255}, Colour{255, 0, 255}, 1)
	addColourRange(Colour{255, 0, 255}, Colour{255, 0, 0}, 1)
	addColourRange(Colour{255, 0, 0}, Colour{255, 255, 0}, 1)
	addColourRange(Colour{255, 255, 0}, Colour{0, 255, 0}, 1)
	addColourRange(Colour{0, 255, 0}, Colour{0, 255, 255}, 1)
	addColourRange(Colour{0, 255, 255}, Colour{0, 0, 255}, 1)
	addColourRange(Colour{0, 0, 255}, Colour{0, 0, 0}, 1)
}

func addColourRange(from, to Colour, increment int) {
	rDiff := to.r - from.r
	gDiff := to.g - from.g
	bDiff := to.b - from.b

	abs := func(n int) int {
		if n < 0 {
			return -n
		}
		return n
	}
	absMax := func(a, b int) int {
		if abs(a) > abs(b) {
			return abs(a)
		}
		return abs(b)
	}

	maxDiff := absMax(rDiff, absMax(gDiff, bDiff))
	rInc := float64(rDiff) / float64(maxDiff)
	gInc := float64(gDiff) / float64(maxDiff)
	bInc := float64(bDiff) / float64(maxDiff)

	runningR := float64(from.r)
	runningG := float64(from.g)
	runningB := float64(from.b)
	for i := 0; i < abs(maxDiff); i += increment {
		colours = append(colours, Colour{
			int(runningR),
			int(runningG),
			int(runningB),
		})

		incFactor := float64(increment)
		runningR += (rInc * incFactor)
		runningG += (gInc * incFactor)
		runningB += (bInc * incFactor)
	}
}
