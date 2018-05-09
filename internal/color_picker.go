package internal

import (
	"sync"

	"github.com/fatih/color"
)

// ColorPicker is a hash map that loops through it
// sequentially to pick a color.
type ColorPicker struct {
	m     []*color.Color
	mu    *sync.Mutex
	count int
}

// NewColorPicker creates a filled color picker.
func NewColorPicker() *ColorPicker {
	m := []*color.Color{
		color.New(color.FgRed),
		color.New(color.FgGreen),
		color.New(color.FgYellow),
		color.New(color.FgBlue),
		color.New(color.FgMagenta),
		color.New(color.FgCyan),
	}

	return &ColorPicker{m: m, mu: &sync.Mutex{}}
}

// Pick retrieves a color and increments the counter
// until it reaches the edge and then loops through the map again.
func (cp *ColorPicker) Pick() *color.Color {
	cp.mu.Lock()

	c := cp.m[cp.count]

	cp.incr()
	cp.mu.Unlock()

	return c
}

func (cp *ColorPicker) incr() {
	if cp.count++; cp.count >= len(cp.m) {
		cp.count = 0
	}
}
