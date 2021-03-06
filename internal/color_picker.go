package internal

import (
	"sync"

	"github.com/fatih/color"
)

// ColorPicker is a hash map that loops through it
// sequentially to pick a color.
type ColorPicker struct {
	attrs []color.Attribute
	mu    *sync.Mutex
	i     int
}

// NewColorPicker creates a filled color picker.
func NewColorPicker() *ColorPicker {
	return &ColorPicker{
		attrs: []color.Attribute{
			color.FgRed,
			color.FgGreen,
			color.FgYellow,
			color.FgBlue,
			color.FgMagenta,
			color.FgCyan,
		},
		mu: &sync.Mutex{},
	}
}

// Pick retrieves a color and increments the counter
// until it reaches the edge and then loops through the map again.
func (cp *ColorPicker) Pick() *color.Color {
	defer cp.mu.Unlock()
	cp.mu.Lock()
	return color.New(cp.attrs[cp.index()])
}

func (cp *ColorPicker) index() int {
	defer func() { cp.i++ }()
	if cp.i >= len(cp.attrs) {
		cp.i = 0
	}
	return cp.i
}
