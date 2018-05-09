package multi

import (
	"bufio"
	"io"
	"log"
	"os/exec"

	"github.com/fatih/color"
)

// Task is a task to be run by multi.
type Task struct {
	Name string   `json:"name,omitempty"`
	Cmd  string   `json:"command,omitempty"`
	Args []string `json:"args,omitemptu"`
	c    *color.Color
}

// Run starts a task and streams its stdout
// and stderr to the process stdout.
func (t *Task) Run(padding int) error {
	padding -= len(t.Name)
	var (
		err error
		outRd,
		errRd io.ReadCloser
		cmd = exec.Command(t.Cmd, t.Args...)
	)

	// Stdout.
	if outRd, err = cmd.StdoutPipe(); err != nil {
		return err
	}

	// Stderr.
	if errRd, err = cmd.StderrPipe(); err != nil {
		return err
	}

	go func() {
		s := bufio.NewScanner(outRd)

		for s.Scan() {
			label := t.c.Sprintf("%s%*s (PID: %d) |", t.Name, padding, "", cmd.Process.Pid)

			log.Printf("%s %s\n", label, s.Text())
		}
	}()
	go func() {
		s := bufio.NewScanner(errRd)

		for s.Scan() {
			label := t.c.Sprintf("%s%*s (PID: %d) |", t.Name, padding, "", cmd.Process.Pid)

			log.Printf("%s %s\n", label, s.Text())
		}
	}()

	if err = cmd.Start(); err != nil {
		return err
	}

	msg := t.c.Sprintf("Starting task %s with PID %d", t.Name, cmd.Process.Pid)

	log.Println(msg)

	return cmd.Wait()
}
