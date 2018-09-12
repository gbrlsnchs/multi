package multi

import (
	"bufio"
	"context"
	"io"
	"log"
	"os"
	"os/exec"
	"sync"
	"time"

	"github.com/fatih/color"
)

// Task is a subprocess.
type Task struct {
	Name    string        `json:"name" toml:"name"`
	Cmd     string        `json:"command" toml:"command"`
	Args    []string      `json:"args" toml:"args"`
	Env     Env           `json:"env" toml:"env"`
	TTL     time.Duration `json:"ttl" toml:"ttl"`
	Flags   int           `json:"-" toml:"-"`
	padding int           // used for task lists
	stderr  chan string
	stdout  chan string
	color   *color.Color
}

// ResolveCmd resolves the task's name and printing padding.
func (t *Task) ResolveCmd(cmd string, padding int) {
	if t.Cmd == "" {
		t.Cmd = cmd
	}
	t.padding = padding - len(t.Name)
}

// Run starts a task using a background context.
// It writes the task stdout and stderr to a list of respective writers.
func (t *Task) Run(stdoutWriters, stderrWriters []io.Writer) (*exec.Cmd, error) {
	return t.RunContext(context.Background(), stdoutWriters, stderrWriters)
}

// RunContext starts a task using a existent context.
// It writes the task stdout and stderr to a list of respective writers.
func (t *Task) RunContext(ctx context.Context, errWriters, outWriters []io.Writer) (*exec.Cmd, error) {
	if t.TTL > 0 {
		now := time.Now()
		var cancel context.CancelFunc
		ctx, cancel = context.WithDeadline(ctx, now.Add(t.TTL*time.Second))
		defer cancel()
	}

	cmd := exec.CommandContext(ctx, t.Cmd, t.Args...)
	if len(t.Env) > 0 {
		cmd.Env = append(os.Environ(), t.Env.Raw()...)
	}
	cmderr, err := cmd.StderrPipe()
	if err != nil {
		return nil, err
	}
	cmdout, err := cmd.StdoutPipe()
	if err != nil {
		return nil, err
	}
	if err = cmd.Start(); err != nil {
		return nil, err
	}

	stderr := io.MultiWriter(errWriters...)
	stdout := io.MultiWriter(outWriters...)

	var wg sync.WaitGroup
	wg.Add(2)
	go t.writeOutput(stderr, cmderr, cmd.Process.Pid, &wg)
	go t.writeOutput(stdout, cmdout, cmd.Process.Pid, &wg)
	wg.Wait()

	return cmd, nil
}

func (t *Task) writeOutput(w io.Writer, rc io.ReadCloser, pid int, wg *sync.WaitGroup) {
	prefix := []byte("%s%*s | ")
	args := []interface{}{
		t.Name,
		t.padding,
		"",
	}
	if t.Flags&Mpid > 0 {
		prefix = append(prefix, "[PID: %d] "...)
		args = append(args, pid)
	}
	l := log.New(w, t.color.Sprintf(string(prefix), args...), t.Flags)
	s := bufio.NewScanner(rc)
	for s.Scan() {
		l.Printf("%s", s.Bytes())
	}
	wg.Done()
}
