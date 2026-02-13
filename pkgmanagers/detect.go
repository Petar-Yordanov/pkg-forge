package detect

import (
	"bytes"
	"context"
	"errors"
	"os/exec"
	"runtime"
	"strings"
	"time"
)

type ManagerStatus struct {
	Cmd       string   `json:"cmd"`
	Name      string   `json:"name"`
	Platforms []string `json:"platforms"`
	Available bool     `json:"available"`
	Path      string   `json:"path,omitempty"`
	Version   string   `json:"version,omitempty"`
	Err       string   `json:"error,omitempty"`
}

type Spec struct {
	Cmd          string
	Name         string
	Platforms    []string
	Candidates   []string
	VersionArgs  []string
	SkipRunProbe bool
}

func DetectAll(ctx context.Context) []ManagerStatus {
	specs := defaultSpecs()
	out := make([]ManagerStatus, 0, len(specs))
	for _, s := range specs {
		out = append(out, detectOne(ctx, s))
	}
	return out
}

func defaultSpecs() []Spec {
	all := []string{"windows", "linux", "macos"}

	return []Spec{
		{Cmd: "apt-get", Name: "APT", Platforms: []string{"linux"}, Candidates: []string{"apt-get"}, VersionArgs: []string{"--version"}},
		{Cmd: "dnf", Name: "DNF", Platforms: []string{"linux"}, Candidates: []string{"dnf"}, VersionArgs: []string{"--version"}},
		{Cmd: "pacman", Name: "pacman", Platforms: []string{"linux"}, Candidates: []string{"pacman"}, VersionArgs: []string{"--version"}},
		{Cmd: "brew", Name: "Homebrew", Platforms: []string{"macos", "linux"}, Candidates: []string{"brew"}, VersionArgs: []string{"--version"}},
		{Cmd: "winget", Name: "WinGet", Platforms: []string{"windows"}, Candidates: []string{"winget"}, VersionArgs: []string{"--version"}},
		{Cmd: "scoop", Name: "Scoop", Platforms: []string{"windows"}, Candidates: []string{"scoop", "scoop.cmd"}, VersionArgs: []string{"--version"}},
		{Cmd: "choco", Name: "Chocolatey", Platforms: []string{"windows"}, Candidates: []string{"choco", "choco.exe"}, VersionArgs: []string{"--version"}},
		{Cmd: "pip", Name: "pip", Platforms: all, Candidates: []string{"pip", "pip3"}, VersionArgs: []string{"--version"}},
		{Cmd: "npm", Name: "npm", Platforms: all, Candidates: []string{"npm", "npm.cmd"}, VersionArgs: []string{"--version"}},
		{Cmd: "cargo", Name: "Cargo", Platforms: all, Candidates: []string{"cargo", "cargo.exe"}, VersionArgs: []string{"--version"}},
		{Cmd: "go", Name: "Go", Platforms: all, Candidates: []string{"go", "go.exe"}, VersionArgs: []string{"version"}},
		{Cmd: "pipx", Name: "pipx", Platforms: all, Candidates: []string{"pipx", "pipx.exe"}, VersionArgs: []string{"--version"}},
		{Cmd: "uv", Name: "uv", Platforms: all, Candidates: []string{"uv", "uv.exe"}, VersionArgs: []string{"--version"}},
	}
}

func detectOne(ctx context.Context, s Spec) ManagerStatus {
	curPlatform := normalizePlatform(runtime.GOOS)
	if !listContains(s.Platforms, curPlatform) {
		return ManagerStatus{
			Cmd:       s.Cmd,
			Name:      s.Name,
			Platforms: append([]string(nil), s.Platforms...),
			Available: false,
			Err:       "not applicable on this platform",
		}
	}

	path, bin := lookPathAny(s.Candidates)
	if path == "" {
		return ManagerStatus{
			Cmd:       s.Cmd,
			Name:      s.Name,
			Platforms: append([]string(nil), s.Platforms...),
			Available: false,
			Err:       "not found in PATH",
		}
	}

	st := ManagerStatus{
		Cmd:       s.Cmd,
		Name:      s.Name,
		Platforms: append([]string(nil), s.Platforms...),
		Available: true,
		Path:      path,
	}

	if s.SkipRunProbe {
		return st
	}

	ver, err := runVersion(ctx, bin, s.VersionArgs)
	if err != nil {
		st.Err = err.Error()
		return st
	}
	st.Version = ver
	return st
}

func normalizePlatform(goos string) string {
	switch strings.ToLower(goos) {
	case "windows":
		return "windows"
	case "darwin":
		return "macos"
	case "linux":
		return "linux"
	default:
		return strings.ToLower(goos)
	}
}

func listContains(xs []string, v string) bool {
	v = strings.ToLower(v)
	for _, x := range xs {
		if strings.ToLower(x) == v {
			return true
		}
	}
	return false
}

func lookPathAny(candidates []string) (path string, bin string) {
	for _, c := range candidates {
		p, err := exec.LookPath(c)
		if err == nil && p != "" {
			return p, c
		}
	}
	return "", ""
}

func runVersion(ctx context.Context, bin string, args []string) (string, error) {
	cctx, cancel := context.WithTimeout(ctx, 2*time.Second)
	defer cancel()

	cmd := exec.CommandContext(cctx, bin, args...)
	var stdout, stderr bytes.Buffer
	cmd.Stdout = &stdout
	cmd.Stderr = &stderr

	if err := cmd.Run(); err != nil {
		msg := strings.TrimSpace(stderr.String())
		if msg == "" {
			msg = strings.TrimSpace(stdout.String())
		}
		if msg == "" {
			msg = err.Error()
		}
		if errors.Is(cctx.Err(), context.DeadlineExceeded) {
			return "", errors.New("probe timed out")
		}
		return "", errors.New(msg)
	}

	combined := strings.TrimSpace(stdout.String())
	if combined == "" {
		combined = strings.TrimSpace(stderr.String())
	}
	return firstLine(combined), nil
}

func firstLine(s string) string {
	if s == "" {
		return ""
	}
	if i := strings.IndexByte(s, '\n'); i >= 0 {
		return strings.TrimSpace(s[:i])
	}
	return strings.TrimSpace(s)
}
