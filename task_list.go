package multi

import (
	"io"
	"sync"

	"github.com/gbrlsnchs/multi/internal"
)

// TaskList is a list of tasks.
//
// If a task doesn't have a command set,
// it runs the task list default command.
type TaskList struct {
	Name   string      `json:"name" toml:"name"`
	Desc   string      `json:"description" toml:"description"`
	Cmd    string      `json:"command" toml:"command"`
	Tasks  []*Task     `json:"tasks" toml:"tasks"`
	Stderr []io.Writer `json:"-" toml:"-"`
	Stdout []io.Writer `json:"-" toml:"-"`
}

// Start starts the task list.
func (tl *TaskList) Start(colored bool) error {
	// Calculate largest padding for logging.
	var padding int
	for _, t := range tl.Tasks {
		if nlen := len(t.Name); nlen > padding {
			padding = nlen
		}
	}

	var wg sync.WaitGroup
	wg.Add(len(tl.Tasks))
	cp := internal.NewColorPicker()
	for _, t := range tl.Tasks {
		go func(t *Task) {
			defer wg.Done()
			t.ResolveCmd(tl.Cmd, padding)
			t.color = cp.Pick()
			if !colored {
				t.color.DisableColor()
			}
			t.Run(tl.Stderr, tl.Stdout)
		}(t)
	}
	wg.Wait()
	return nil
}
