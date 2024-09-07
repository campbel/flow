package sh

import (
	"context"
	"io"
	"os"
	"os/exec"
)

type Runner struct {
	in  io.Reader
	out io.Writer
	err io.Writer
}

var DefaultRunner = NewRunner(os.Stdin, os.Stdout, os.Stderr)

func NewRunner(in io.Reader, out, err io.Writer) *Runner {
	return &Runner{in: in, out: out, err: err}
}

func (l *Runner) Run(ctx context.Context, script string, envs []string) error {
	cmd := exec.CommandContext(ctx, "sh", "-c", script)
	cmd.Env = envs
	cmd.Stdin = l.in
	cmd.Stdout = l.out
	cmd.Stderr = l.err
	return cmd.Run()
}
