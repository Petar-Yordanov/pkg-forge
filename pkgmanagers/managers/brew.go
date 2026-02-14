package managers

import (
	"errors"
	"os/exec"
	"strings"
	"time"

	"github.com/Petar-Yordanov/pkg-forge/common"
)

type Brew struct{}

func (Brew) ID() string          { return "brew" }
func (Brew) DisplayName() string { return "Homebrew" }

func (Brew) Platforms() []common.Platform {
	return []common.Platform{common.PlatformMacOS}
}

func (Brew) Detect() (DetectResult, error) {
	cur := common.CurrentPlatform()

	path, err := exec.LookPath("brew")
	if err != nil {
		return DetectResult{Available: false, Platform: cur}, errors.New("not found in PATH")
	}

	out, err := Command(path).Args("--version").Timeout(2 * time.Second).RunTrimOutput()
	if err != nil {
		// Brew exists, but version command failed, so technically still "available"
		return DetectResult{Available: true, Path: path, Platform: cur}, nil
	}

	ver := strings.TrimSpace(out)
	return DetectResult{Available: true, Path: path, Version: ver, Platform: cur}, nil
}

func (Brew) Install(name string, version string) error {
	path, err := exec.LookPath("brew")
	if err != nil {
		return errors.New("not found in PATH")
	}

	pkg := name
	if version != "" && !strings.Contains(name, "@") {
		pkg = name + "@" + version
	}

	_, err = Command(path).
		Args("install", pkg).
		Timeout(15 * time.Minute).
		RunTrimOutput()
	return err
}

func (Brew) InstallLatest(name string) error {
	return (Brew{}).Install(name, "")
}

func (Brew) Uninstall(name string) error {
	path, err := exec.LookPath("brew")
	if err != nil {
		return errors.New("not found in PATH")
	}

	_, err = Command(path).
		Args("uninstall", "--force", name).
		Timeout(15 * time.Minute).
		RunTrimOutput()
	return err
}
