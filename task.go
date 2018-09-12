package multi

import (
	"bufio"
	"fmt"
	"io"
	"os"
	"os/exec"
	"sync"

	"github.com/fatih/color"
)

// Task is a task to be run by multi.
type Task struct {
	Name    string   `json:"name" toml:"name"`
	Cmd     string   `json:"command" toml:"command"`
	Args    []string `json:"args" toml:"args"`
	Env     Env      `json:"env" toml:"env"`
	padding int      // used for task lists
	stderr  chan string
	stdout  chan string
	color   *color.Color
}

// ResolveCmd resolves the task's name.
func (t *Task) ResolveCmd(cmd string, padding int) {
	if t.Cmd == "" {
		t.Cmd = cmd
	}
	t.padding = padding - len(t.Name)
}

// Run starts a task and streams its stdout
// and stderr to the process stdout.
func (t *Task) Run(stdoutWriters, stderrWriters []io.Writer) error {
	cmd := exec.Command(t.Cmd, t.Args...)
	if len(t.Env) > 0 {
		cmd.Env = append(os.Environ(), t.Env.Raw()...)
	}
	cmderr, err := cmd.StderrPipe()
	if err != nil {
		return err
	}
	cmdout, err := cmd.StdoutPipe()
	if err != nil {
		return err
	}
	if err = cmd.Start(); err != nil {
		return err
	}

	stderr := io.MultiWriter(stderrWriters...)
	stdout := io.MultiWriter(stdoutWriters...)

	var wg sync.WaitGroup
	wg.Add(2)
	t.writeOutput(stderr, cmderr, cmd.Process.Pid, &wg)
	t.writeOutput(stdout, cmdout, cmd.Process.Pid, &wg)
	wg.Wait()

	return cmd.Wait()
}

func (t *Task) writeOutput(w io.Writer, rc io.ReadCloser, pid int, wg *sync.WaitGroup) {
	s := bufio.NewScanner(rc)
	for s.Scan() {
		label := t.color.Sprintf("%s%*s (PID: %d) |", t.Name, t.padding, "", pid)
		fmt.Fprintf(w, "%s %s\n", label, s.Bytes())
	}
	wg.Done()
}
