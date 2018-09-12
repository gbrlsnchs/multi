package multi

import (
	"context"
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
	Flags  int         `json:"-" toml:"-"`
}

// Start starts the task list using a background context.
func (tl *TaskList) Start() error {
	return tl.StartContext(context.Background())
}

// StartContext starts the task list using a existent content.
func (tl *TaskList) StartContext(ctx context.Context) error {
	// Calculate largest padding for logging.
	var padding int
	for _, t := range tl.Tasks {
		if nlen := len(t.Name); nlen > padding {
			padding = nlen
		}
	}

	done := make(chan struct{})
	errCh := make(chan error)

	var wg sync.WaitGroup
	wg.Add(len(tl.Tasks))

	cp := internal.NewColorPicker()
	for _, t := range tl.Tasks {
		go func(t *Task) {
			defer wg.Done()
			t.ResolveCmd(tl.Cmd, padding)
			t.color = cp.Pick()
			t.Flags = tl.Flags
			if tl.Flags&Mcolor == 0 {
				t.color.DisableColor()
			}
			cmd, err := t.RunContext(ctx, tl.Stderr, tl.Stdout)
			if err != nil {
				errCh <- err
				return
			}
			// TODO: implement detached command
			if err = cmd.Wait(); err != nil {
				errCh <- err
			}
		}(t)
	}
	go func() {
		wg.Wait()
		close(done)
	}()

	select {
	case err := <-errCh:
		return err
	case <-done:
		return nil
	}
}
