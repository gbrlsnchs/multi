package internal

import (
	"sync"

	"github.com/fatih/color"
)

// ColorPicker is a hash map that loops through it
// sequentially to pick a color.
type ColorPicker struct {
	colors []*color.Color
	count  int
	mu     *sync.Mutex
}

// NewColorPicker creates a filled color picker.
func NewColorPicker() *ColorPicker {
	colors := []*color.Color{
		color.New(color.FgRed),
		color.New(color.FgGreen),
		color.New(color.FgYellow),
		color.New(color.FgBlue),
		color.New(color.FgMagenta),
		color.New(color.FgCyan),
	}

	return &ColorPicker{colors: colors, mu: &sync.Mutex{}}
}

// Pick retrieves a color and increments the counter
// until it reaches the edge and then loops through the map again.
func (cp *ColorPicker) Pick() *color.Color {
	cp.mu.Lock()

	c := cp.colors[cp.count]

	cp.incr()
	cp.mu.Unlock()

	return c
}

func (cp *ColorPicker) incr() {
	if cp.count++; cp.count >= len(cp.colors) {
		cp.count = 0
	}
}
