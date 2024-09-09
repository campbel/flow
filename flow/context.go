package flow

import (
	"bytes"
	"context"
	"errors"
	"os"
	"strings"
	"sync"
	"text/template"

	"github.com/campbel/flow/langs/sh"
	"github.com/campbel/flow/meta/logger"
	"github.com/campbel/flow/types"
)

type ShellRunner interface {
	Run(ctx context.Context, script string, envs []string) error
}

type Context struct {
	flowfile    types.Flowfile
	workdir     string
	homedir     string
	args        string
	shellRunner ShellRunner
}

func NewContext(flowfile types.Flowfile, workdir, homedir, args string) *Context {
	return &Context{
		flowfile:    flowfile,
		workdir:     workdir,
		homedir:     homedir,
		args:        args,
		shellRunner: sh.DefaultRunner,
	}
}

func (c *Context) Execute(ctx context.Context, flow types.Flow) error {
	return c.executeFlow(ctx, flow, osEnviron(), nil)
}

func (c *Context) executeFlow(ctx context.Context, flow types.Flow, envs map[string]string, vars map[string]any) error {
	if flow.If != nil {
		if err := c.executeStep(ctx, *flow.If, mergeMaps(envs, flow.Envs), mergeMaps(vars, flow.Vars)); err != nil {
			return nil
		}
	}
	return c.executeSteps(ctx, flow.Steps, mergeMaps(envs, flow.Envs), mergeMaps(vars, flow.Vars))
}

func (c *Context) executeSteps(ctx context.Context, steps []types.Step, envs map[string]string, vars map[string]any) error {
	var wg sync.WaitGroup
	for _, s := range steps {
		if s.If != nil {
			if err := c.executeStep(ctx, *s.If, envs, vars); err != nil {
				continue
			}
		}

		rng := s.Range
		if len(rng) == 0 {
			rng = []any{nil}
		}

		for i, item := range rng {
			loopVars := mergeMaps(vars, map[string]any{"index": i, "item": item})
			if s.Go != nil {
				wg.Add(1)
				go func(s types.Step, envs map[string]string, vars map[string]any) {
					defer wg.Done()
					if err := c.executeStep(ctx, s, envs, vars); err != nil {
						return
					}
				}(*s.Go, envs, loopVars)
			} else if s.Defer != nil {
				defer c.executeStep(ctx, *s.Defer, envs, loopVars)
			} else if s.Group != nil {
				if err := c.executeSteps(ctx, s.Group, envs, loopVars); err != nil {
					return err
				}
			} else {
				if err := c.executeStep(ctx, s, envs, loopVars); err != nil {
					return err
				}
			}
		}
	}
	wg.Wait()
	return nil
}

func (c *Context) executeStep(ctx context.Context, step types.Step, envs map[string]string, vars map[string]any) error {
	if step.Name != "" {
		logger.Info("executing step", "name", step.Name)
	}
	if step.Shell != "" {
		return c.shellRunner.Run(ctx, c.templ(step.Shell, vars), mapToSlice(envs))
	}
	if step.Flow != "" {
		flow, ok := c.flowfile.Flows[step.Flow]
		if ok {
			return c.executeFlow(ctx, flow, envs, vars)
		}
		return errors.New("flow not found: " + step.Flow)
	}
	return nil
}

func (c *Context) templ(s string, vars map[string]any) string {
	data := mergeMaps(map[string]any{
		"workdir": c.workdir,
		"homedir": c.homedir,
		"args":    c.args,
	}, vars)
	t, err := template.New("flow").Parse(s)
	if err != nil {
		return ""
	}
	var buf bytes.Buffer
	if err := t.ExecuteTemplate(&buf, "flow", data); err != nil {
		return ""
	}
	return buf.String()
}

func mergeMaps[M any](maps ...map[string]M) map[string]M {
	result := make(map[string]M)
	for _, m := range maps {
		for k, v := range m {
			result[k] = v
		}
	}
	return result
}

func mapToSlice(m map[string]string) []string {
	result := make([]string, 0, len(m))
	for k, v := range m {
		result = append(result, k+"="+v)
	}
	return result
}

func osEnviron() map[string]string {
	envs := make(map[string]string)
	for _, e := range os.Environ() {
		pair := strings.SplitN(e, "=", 2)
		envs[pair[0]] = pair[1]
	}
	return envs
}
