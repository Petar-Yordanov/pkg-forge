package steps

import (
	"fmt"
	"path/filepath"
	"strings"
	"time"

	"github.com/Petar-Yordanov/pkg-forge/common"
	"github.com/Petar-Yordanov/pkg-forge/manifest/parser"
)

type StepResult struct {
	Skipped    bool
	SkipReason string
	Output     string
}

func RunStep(platform common.Platform, e parser.Entry, s parser.Step) (StepResult, error) {
	if s.When != nil && len(s.When.Platform) > 0 {
		ok := false
		for _, p := range s.When.Platform {
			if p == platform {
				ok = true
				break
			}
		}
		if !ok {
			return StepResult{
				Skipped:    true,
				SkipReason: fmt.Sprintf("platform mismatch (need %v, have %s)", s.When.Platform, platform),
			}, nil
		}
	}

	if s.Cmd != "" && s.CmdFile != "" {
		return StepResult{}, fmt.Errorf("step cannot have both cmd and cmdFile")
	}
	if s.Cmd == "" && s.CmdFile == "" {
		return StepResult{}, fmt.Errorf("step must have either cmd or cmdFile")
	}

	env := mergeEnv(e.Env, s.Env)
	failOnStderr := mergedFailOnStderr(e.FailOnStderr, s.FailOnStderr)
	timeout := stepTimeout(s)
	retries := stepRetries(s)
	retryDelay := stepRetryDelay(s)

	if s.Cmd != "" {
		runner, runnerArgs, err := resolveInlineShell(platform, s.Shell)
		if err != nil {
			return StepResult{}, err
		}

		out, err := common.Command(runner).
			Args(append(runnerArgs, s.Cmd)...).
			Env(env).
			Timeout(timeout).
			Retries(retries).
			RetryDelay(retryDelay).
			FailOnStderr(failOnStderr).
			RunTrimOutput()
		if err != nil {
			return StepResult{}, err
		}

		return StepResult{Output: out}, nil
	}

	runner, runnerArgs, err := resolveFileShell(platform, s.Shell, s.CmdFile)
	if err != nil {
		return StepResult{}, err
	}

	out, err := common.CommandFile(runner, filepath.Clean(s.CmdFile)).
		RunnerArgs(runnerArgs...).
		Args(s.Args...).
		Env(env).
		Timeout(timeout).
		Retries(retries).
		RetryDelay(retryDelay).
		FailOnStderr(failOnStderr).
		RunTrimOutput()
	if err != nil {
		return StepResult{}, err
	}

	return StepResult{Output: out}, nil
}

func mergeEnv(entryEnv map[string]string, stepEnv map[string]string) map[string]string {
	if len(entryEnv) == 0 && len(stepEnv) == 0 {
		return nil
	}

	out := make(map[string]string, len(entryEnv)+len(stepEnv))
	for k, v := range entryEnv {
		out[k] = v
	}
	for k, v := range stepEnv {
		out[k] = v
	}
	return out
}

func mergedFailOnStderr(entryValue, stepValue *bool) bool {
	if stepValue != nil {
		return *stepValue
	}
	if entryValue != nil {
		return *entryValue
	}
	return false
}

func stepTimeout(s parser.Step) time.Duration {
	if s.TimeoutSec <= 0 {
		return 0
	}
	return time.Duration(s.TimeoutSec) * time.Second
}

func stepRetries(s parser.Step) int {
	if s.Retries < 0 {
		return 0
	}
	return s.Retries
}

func stepRetryDelay(s parser.Step) time.Duration {
	if s.RetryDelaySec <= 0 {
		return 0
	}
	return time.Duration(s.RetryDelaySec) * time.Second
}

func resolveInlineShell(platform common.Platform, shell string) (string, []string, error) {
	switch normalizeShell(shell, platform) {
	case "powershell":
		if platform == common.PlatformWindows {
			return "powershell.exe", []string{"-NoProfile", "-NonInteractive", "-ExecutionPolicy", "Bypass", "-Command"}, nil
		}
		return "pwsh", []string{"-NoProfile", "-NonInteractive", "-Command"}, nil
	case "pwsh":
		return "pwsh", []string{"-NoProfile", "-NonInteractive", "-Command"}, nil
	case "cmd":
		return "cmd.exe", []string{"/C"}, nil
	case "bash":
		return "bash", []string{"-c"}, nil
	case "zsh":
		return "zsh", []string{"-c"}, nil
	case "sh":
		return "sh", []string{"-c"}, nil
	default:
		return "", nil, fmt.Errorf("unsupported inline shell %q", shell)
	}
}

func resolveFileShell(platform common.Platform, shell string, scriptPath string) (string, []string, error) {
	switch normalizeFileShell(shell, platform, scriptPath) {
	case "powershell":
		if platform == common.PlatformWindows {
			return "powershell.exe", nil, nil
		}
		return "pwsh", nil, nil
	case "pwsh":
		return "pwsh", nil, nil
	case "cmd":
		return "cmd.exe", []string{"/C"}, nil
	case "bash":
		return "bash", nil, nil
	case "zsh":
		return "zsh", nil, nil
	case "sh":
		return "sh", nil, nil
	default:
		return "", nil, fmt.Errorf("unsupported file shell %q", shell)
	}
}

func normalizeShell(shell string, platform common.Platform) string {
	s := strings.ToLower(strings.TrimSpace(shell))
	if s != "" {
		return s
	}

	if platform == common.PlatformWindows {
		return "powershell"
	}
	return "sh"
}

func normalizeFileShell(shell string, platform common.Platform, scriptPath string) string {
	s := strings.ToLower(strings.TrimSpace(shell))
	if s != "" {
		return s
	}

	ext := strings.ToLower(filepath.Ext(scriptPath))
	switch ext {
	case ".ps1":
		return "powershell"
	case ".bat", ".cmd":
		return "cmd"
	case ".sh":
		return "sh"
	case ".zsh":
		return "zsh"
	}

	if platform == common.PlatformWindows {
		return "powershell"
	}
	return "sh"
}
