package common

import (
	"bytes"
	"context"
	"errors"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"time"
)

var (
	ErrNilCmdFile        = errors.New("nil cmd file")
	ErrFileExistsHasArgs = errors.New("Exists requires no args")
	ErrEmptyRunner       = errors.New("missing runner")
	ErrRunnerHasSpace    = errors.New("runner must not contain whitespace")
	ErrRunnerNotFound    = errors.New("runner not found in PATH")

	ErrEmptyScriptPath   = errors.New("missing script path")
	ErrScriptNotFound    = errors.New("script not found")
	ErrScriptPathIsDir   = errors.New("script path is a directory")
	ErrScriptStatFailed  = errors.New("script stat failed")
)

type CmdFile struct {
	runner     string
	rArgs      []string
	path       string
	args       []string
	timeout    time.Duration
	retries    int
	retryDelay time.Duration
}

func CommandFile(runner string, path string) *CmdFile {
	return &CmdFile{runner: runner, path: path}
}

func (c *CmdFile) RunnerArgs(args ...string) *CmdFile {
	c.rArgs = append(c.rArgs[:0], args...)
	return c
}

func (c *CmdFile) Args(args ...string) *CmdFile {
	c.args = append(c.args[:0], args...)
	return c
}

func (c *CmdFile) Timeout(d time.Duration) *CmdFile {
	c.timeout = d
	return c
}

func (c *CmdFile) Retries(n int) *CmdFile {
	if n < 0 {
		n = 0
	}

	c.retries = n
	return c
}

func (c *CmdFile) RetryDelay(d time.Duration) *CmdFile {
	if d < 0 {
		d = 0
	}

	c.retryDelay = d
	return c
}

func (c *CmdFile) Exists() error {
	if c == nil {
		return ErrNilCmdFile
	}

	if len(c.args) != 0 {
		return ErrFileExistsHasArgs
	}

	r := strings.TrimSpace(c.runner)
	if r == "" {
		return ErrEmptyRunner
	}

	if strings.ContainsAny(r, " \t\r\n") {
		return fmt.Errorf("%w: %q", ErrRunnerHasSpace, r)
	}

	if _, err := exec.LookPath(r); err != nil {
		return fmt.Errorf("%w: %q (%v)", ErrRunnerNotFound, r, err)
	}

	p := strings.TrimSpace(c.path)
	if p == "" {
		return ErrEmptyScriptPath
	}

	p = filepath.Clean(p)

	info, err := os.Stat(p)
	if err != nil {
		if os.IsNotExist(err) {
			return fmt.Errorf("%w: %s", ErrScriptNotFound, p)
		}

		return fmt.Errorf("%w: %s (%v)", ErrScriptStatFailed, p, err)
	}

	if info.IsDir() {
		return fmt.Errorf("%w: %s", ErrScriptPathIsDir, p)
	}

	return nil
}

func (c *CmdFile) Run() (string, error) {
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

func (c *CmdFile) runOnce() (string, error) {
	ctx := context.Background()
	if c.timeout > 0 {
		var cancel context.CancelFunc
		ctx, cancel = context.WithTimeout(ctx, c.timeout)
		defer cancel()
	}

	allArgs := make([]string, 0, len(c.rArgs)+1+len(c.args))
	allArgs = append(allArgs, c.rArgs...)
	allArgs = append(allArgs, c.path)
	allArgs = append(allArgs, c.args...)

	cmd := exec.CommandContext(ctx, c.runner, allArgs...)
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

	out := strings.TrimSpace(stdout.String())
	if out == "" {
		out = strings.TrimSpace(stderr.String())
	}

	return out, nil
}

func (c *CmdFile) RunTrimOutput() (string, error) {
	out, err := c.Run()
	if err != nil {
		return "", err
	}

	return strings.TrimSpace(out), nil
}
