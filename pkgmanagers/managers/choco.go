package managers

import (
	"errors"
	"os/exec"
	"time"

	"github.com/Petar-Yordanov/pkg-forge/common"
	"github.com/samber/lo"
)

type Choco struct{}

func (Choco) ID() string { return "choco" }
func (Choco) DisplayName() string { return "Chocolatey" }
func (Choco) Platforms() []common.Platform { return []common.Platform{common.PlatformWindows} }

func (Choco) Detect() (DetectResult, error) {
	cur := common.CurrentPlatform()

	bin, ok := lo.Find([]string{"choco", "choco.exe"}, func(c string) bool {
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

	return DetectResult{Available: true, Path: path, Version: out, Platform: cur}, nil
}

func (Choco) Install(name string, version string) error {
	bin, ok := lo.Find([]string{"choco", "choco.exe"}, func(c string) bool {
		_, err := exec.LookPath(c)
		return err == nil
	})
	if !ok {
		return nil
	}

	args := []string{"install", name, "-y"}
	if version != "" {
		args = []string{"install", name, "--version", version, "-y"}
	}
	_, err := Command(bin).Args(args...).RunTrimOutput()
	return err
}

func (Choco) InstallLatest(name string) error { return (Choco{}).Install(name, "") }

func (Choco) Uninstall(name string) error {
	bin, ok := lo.Find([]string{"choco", "choco.exe"}, func(c string) bool {
		_, err := exec.LookPath(c)
		return err == nil
	})
	if !ok {
		return nil
	}

	_, err := Command(bin).Args("uninstall", name, "-y").RunTrimOutput()
	return err
}
