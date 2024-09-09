package sh

import (
	"context"
	"io"
	"os"
	"os/exec"

	"github.com/campbel/flow/types"
)

type SystemRunner struct {
	in  io.Reader
	out io.Writer
	err io.Writer
}

var DefaultRunner = NewSystemRunner(os.Stdin, os.Stdout, os.Stderr)

func NewSystemRunner(in io.Reader, out, err io.Writer) *SystemRunner {
	return &SystemRunner{in: in, out: out, err: err}
}

func (l *SystemRunner) Run(ctx context.Context, script string, envs []string) error {
	cmd := exec.CommandContext(ctx, "sh", "-c", script)
	cmd.Env = envs
	cmd.Stdin = l.in
	cmd.Stdout = l.out
	cmd.Stderr = l.err
	if err := cmd.Run(); err != nil {
		if exitErr, ok := err.(*exec.ExitError); ok {
			return types.NewExitError(exitErr.ExitCode(), err)
		}
		return err
	}
	return nil
}
