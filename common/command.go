package common

import (
	"bytes"
	"context"
	"errors"
	"fmt"
	"os"
	"os/exec"
	"strings"
	"sync"
	"time"
)

var (
	ErrNilCmd           = errors.New("nil cmd")
	ErrExistsHasArgs    = errors.New("Exists requires no args")
	ErrEmptyBin         = errors.New("empty command")
	ErrBinContainsSpace = errors.New("command must not contain whitespace")

	ErrToolNotFound  = errors.New("not found in PATH")
	ErrProbeTimedOut = errors.New("probe timed out")
	ErrStderrOutput  = errors.New("stderr output detected")
)

type Cmd struct {
	candidates []string
	binPath    string
	binName    string

	args         []string
	timeout      time.Duration
	retries      int
	retryDelay   time.Duration
	env          map[string]string
	failOnStderr bool

	resolveOnce sync.Once
	resolveErr  error
}

func Command(candidates ...string) *Cmd {
	c := &Cmd{}
	return c.Resolve(candidates...)
}

func (c *Cmd) Args(args ...string) *Cmd {
	c.args = append(c.args[:0], args...)
	return c
}

func (c *Cmd) Timeout(d time.Duration) *Cmd {
	c.timeout = d
	return c
}

func (c *Cmd) Retries(n int) *Cmd {
	if n < 0 {
		n = 0
	}

	c.retries = n
	return c
}

func (c *Cmd) RetryDelay(d time.Duration) *Cmd {
	if d < 0 {
		d = 0
	}

	c.retryDelay = d
	return c
}

func (c *Cmd) Env(env map[string]string) *Cmd {
	if len(env) == 0 {
		c.env = nil
		return c
	}

	c.env = make(map[string]string, len(env))
	for k, v := range env {
		c.env[k] = v
	}
	return c
}

func (c *Cmd) FailOnStderr(v bool) *Cmd {
	c.failOnStderr = v
	return c
}

func (c *Cmd) Resolve(candidates ...string) *Cmd {
	c.candidates = c.candidates[:0]
	for _, cand := range candidates {
		cand = strings.TrimSpace(cand)

		if cand != "" {
			c.candidates = append(c.candidates, cand)
		}
	}

	c.ensureResolved()
	return c
}

func (c *Cmd) Exists() error {
	if c == nil {
		return ErrNilCmd
	}

	if len(c.args) != 0 {
		return ErrExistsHasArgs
	}

	return c.ensureResolved()
}

func (c *Cmd) Run() (string, error) {
	if c == nil {
		return "", ErrNilCmd
	}
	if err := c.ensureResolved(); err != nil {
		return "", err
	}

	var lastErr error
	attempts := 1 + c.retries

	for i := 0; i < attempts; i++ {
		out, err := c.runOnce()

		if err == nil {
			return out, nil
		}

		lastErr = err
		if i < attempts-1 && c.retryDelay > 0 {
			time.Sleep(c.retryDelay)
		}
	}
	return "", lastErr
}

func (c *Cmd) ensureResolved() error {
	c.resolveOnce.Do(func() {
		if len(c.candidates) == 0 {
			c.resolveErr = ErrEmptyBin
			return
		}

		for _, cand := range c.candidates {
			if strings.ContainsAny(cand, " \t\r\n") {
				c.resolveErr = fmt.Errorf("%w: %q", ErrBinContainsSpace, cand)
				return
			}
		}

		for _, cand := range c.candidates {
			if p, err := exec.LookPath(cand); err == nil {
				c.binPath = p
				c.binName = cand
				c.resolveErr = nil
				return
			}
		}
		c.resolveErr = fmt.Errorf("%w: %q", ErrToolNotFound, c.candidates)
	})
	return c.resolveErr
}

func (c *Cmd) runOnce() (string, error) {
	ctx := context.Background()
	if c.timeout > 0 {
		var cancel context.CancelFunc
		ctx, cancel = context.WithTimeout(ctx, c.timeout)
		defer cancel()
	}

	cmd := exec.CommandContext(ctx, c.binPath, c.args...)
	if len(c.env) > 0 {
		env := os.Environ()
		for k, v := range c.env {
			env = append(env, k+"="+v)
		}
		cmd.Env = env
	}

	var stdout, stderr bytes.Buffer
	cmd.Stdout = &stdout
	cmd.Stderr = &stderr

	if err := cmd.Run(); err != nil {
		if errors.Is(ctx.Err(), context.DeadlineExceeded) {
			return "", ErrProbeTimedOut
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

	stdoutText := strings.TrimSpace(stdout.String())
	stderrText := strings.TrimSpace(stderr.String())

	if c.failOnStderr && stderrText != "" {
		return "", fmt.Errorf("%w: %s", ErrStderrOutput, stderrText)
	}

	if stdoutText != "" {
		return stdoutText, nil
	}
	return stderrText, nil
}

func (c *Cmd) RunTrimOutput() (string, error) {
	out, err := c.Run()
	if err != nil {
		return "", err
	}
	return strings.TrimSpace(out), nil
}

func (c *Cmd) Path() string { return c.binPath }
