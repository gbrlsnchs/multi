package multi

import (
	"log"
	"sync"

	"github.com/gbrlsnchs/multi/internal"
)

// TaskList is a list of tasks.
//
// If a task doesn't have a command set,
// it runs the task list default command.
type TaskList struct {
	Name  string  `json:"name" toml:"name"`
	Desc  string  `json:"description" toml:"description"`
	Cmd   string  `json:"command" toml:"command"`
	Tasks []*Task `json:"tasks" toml:"tasks"`
}

// Start starts the task list.
func (tl *TaskList) Start(noColor bool) error {
	var padding int
	for _, t := range tl.Tasks {
		if nlen := len(t.Name); nlen > padding {
			padding = nlen
		}
	}
	cp := internal.NewColorPicker()
	var wg sync.WaitGroup
	wg.Add(len(tl.Tasks))
	for _, t := range tl.Tasks {
		go func(t *Task) {
			defer wg.Done()

			if t.c = cp.Pick(); noColor {
				t.c.DisableColor()
			}

			if t.Cmd == "" {
				t.Cmd = tl.Cmd
			}

			if err := t.Run(padding); err != nil {
				log.Fatal(err)
			}
		}(t)
	}
	wg.Wait()
	return nil
}
