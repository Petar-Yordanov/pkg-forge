package managers

import (
	"bytes"
	"context"
	"errors"
	"os/exec"
	"strings"
	"time"
)

type Cmd struct {
	bin     string
	args    []string
	timeout time.Duration
}

func Command(bin string) *Cmd {
	return &Cmd{bin: bin}
}

func (c *Cmd) Args(args ...string) *Cmd {
	c.args = append(c.args[:0], args...)
	return c
}

func (c *Cmd) Timeout(d time.Duration) *Cmd {
	c.timeout = d
	return c
}

func (c *Cmd) Run() (string, error) {
	ctx := context.Background()
	if c.timeout > 0 {
		var cancel context.CancelFunc
		ctx, cancel = context.WithTimeout(ctx, c.timeout)
		defer cancel()
	}

	cmd := exec.CommandContext(ctx, c.bin, c.args...)
	var stdout, stderr bytes.Buffer
	cmd.Stdout = &stdout
	cmd.Stderr = &stderr

	if err := cmd.Run(); err != nil {
		if errors.Is(ctx.Err(), context.DeadlineExceeded) {
			return "", errors.New("probe timed out")
		}
		msg := strings.TrimSpace(stderr.String())
		if msg == "" {
			msg = strings.TrimSpace(stdout.String())
		}
		if msg == "" {
			msg = err.Error()
		}
		return "", errors.New(msg)
	}

	out := strings.TrimSpace(stdout.String())
	if out == "" {
		out = strings.TrimSpace(stderr.String())
	}
	return out, nil
}

func (c *Cmd) RunTrimOutput() (string, error) {
	out, err := c.Run()
	if err != nil {
		return "", err
	}
	return strings.TrimSpace(out), nil
}
