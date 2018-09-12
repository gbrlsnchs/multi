package internal

import (
	"sync"

	"github.com/fatih/color"
)

// ColorPicker is a hash map that loops through it
// sequentially to pick a color.
type ColorPicker struct {
	colors []*color.Color
	mu     *sync.Mutex
	i      int
}

// NewColorPicker creates a filled color picker.
func NewColorPicker() *ColorPicker {
	return &ColorPicker{
		colors: []*color.Color{
			color.New(color.FgRed),
			color.New(color.FgGreen),
			color.New(color.FgYellow),
			color.New(color.FgBlue),
			color.New(color.FgMagenta),
			color.New(color.FgCyan),
		},
		mu: &sync.Mutex{},
	}
}

// Pick retrieves a color and increments the counter
// until it reaches the edge and then loops through the map again.
func (cp *ColorPicker) Pick() *color.Color {
	defer cp.mu.Unlock()
	cp.mu.Lock()
	c := cp.colors[cp.index()]
	return c
}

func (cp *ColorPicker) index() int {
	defer func() { cp.i++ }()
	if cp.i >= len(cp.colors) {
		cp.i = 0
	}
	return cp.i
}
