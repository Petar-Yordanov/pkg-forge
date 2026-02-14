package managers

import (
	"errors"
	"os/exec"
	"strings"
	"time"

	"github.com/Petar-Yordanov/pkg-forge/common"
	"github.com/samber/lo"
)

type Winget struct{}

func (Winget) ID() string { return "winget" }
func (Winget) DisplayName() string { return "WinGet" }
func (Winget) Platforms() []common.Platform { return []common.Platform{common.PlatformWindows} }

func (Winget) Detect() (DetectResult, error) {
	cur := common.CurrentPlatform()

	bin, ok := lo.Find([]string{"winget", "winget.exe"}, func(c string) bool {
		_, err := exec.LookPath(c)
		return err == nil
	})
	if !ok {
		return DetectResult{Available: false, Platform: cur}, errors.New("not found in PATH")
	}

	path, _ := exec.LookPath(bin)

	out, err := Command(bin).Args("--version").Timeout(2 * time.Second).RunTrimOutput()
	if err != nil {
		return DetectResult{Available: true, Path: path, Platform: cur}, nil
	}

	out = strings.TrimSpace(strings.TrimPrefix(out, "v"))
	return DetectResult{Available: true, Path: path, Version: out, Platform: cur}, nil
}

func (Winget) Install(name string, version string) error {
	bin, ok := lo.Find([]string{"winget", "winget.exe"}, func(c string) bool {
		_, err := exec.LookPath(c)
		return err == nil
	})
	if !ok {
		return nil
	}

	args := []string{
		"install",
		"--exact",
		"--id", name,
		"--source", "winget",
		"--silent",
		"--accept-source-agreements",
		"--accept-package-agreements",
		"--disable-interactivity",
	}

	if version != "" {
		args = append(args, "--version", version)
	}

	_, err := Command(bin).Args(args...).RunTrimOutput()
	return err
}

func (Winget) InstallLatest(name string) error { return (Winget{}).Install(name, "") }

func (Winget) Uninstall(name string) error {
	bin, ok := lo.Find([]string{"winget", "winget.exe"}, func(c string) bool {
		_, err := exec.LookPath(c)
		return err == nil
	})
	if !ok {
		return nil
	}

	args := []string{
		"uninstall",
		"--exact",
		"--id", name,
		"--silent",
		"--disable-interactivity",
	}

	_, err := Command(bin).Args(args...).RunTrimOutput()
	return err
}
