package managers

import (
	"errors"
	"strings"
	"time"

	"github.com/Petar-Yordanov/pkg-forge/common"
)

type Brew struct {
	locator ToolLocator
}

func (*Brew) ID() string          { return "brew" }
func (*Brew) DisplayName() string { return "Homebrew" }

func (*Brew) Platforms() []common.Platform {
	return []common.Platform{common.PlatformMacOS}
}

func (m *Brew) Detect() (DetectResult, error) {
	cur := common.CurrentPlatform()

	path, err := m.locator.Resolve("brew")
	if err != nil {
		return DetectResult{Available: false, Platform: cur}, err
	}

	return DetectResult{Available: true, Path: path, Platform: cur}, nil
}

func (m *Brew) GetVersion() (string, error) {
	path, err := m.locator.Resolve("brew")
	if err != nil {
		return "", err
	}

	out, err := Command(path).Args("--version").Timeout(2 * time.Second).RunTrimOutput()
	if err != nil {
		return "", err
	}
	return strings.TrimSpace(out), nil
}

func (m *Brew) Install(name string, version string) error {
	path, err := m.locator.Resolve("brew")
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

func (m *Brew) InstallLatest(name string) error { return m.Install(name, "") }

func (m *Brew) Uninstall(name string) error {
	path, err := m.locator.Resolve("brew")
	if err != nil {
		return errors.New("not found in PATH")
	}

	_, err = Command(path).
		Args("uninstall", "--force", name).
		Timeout(15 * time.Minute).
		RunTrimOutput()
	return err
}
