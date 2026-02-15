package managers

import (
	"time"

	"github.com/Petar-Yordanov/pkg-forge/common"
)

type Choco struct {
	locator ToolLocator
}

func (*Choco) ID() string          { return "choco" }
func (*Choco) DisplayName() string { return "Chocolatey" }
func (*Choco) Platforms() []common.Platform {
	return []common.Platform{common.PlatformWindows}
}

func (m *Choco) Detect() (DetectResult, error) {
	cur := common.CurrentPlatform()

	path, err := m.locator.Resolve("choco.exe", "choco")
	if err != nil {
		return DetectResult{Available: false, Platform: cur}, err
	}

	return DetectResult{Available: true, Path: path, Platform: cur}, nil
}

func (m *Choco) GetVersion() (string, error) {
	path, err := m.locator.Resolve("choco.exe", "choco")
	if err != nil {
		return "", err
	}

	out, err := Command(path).Args("--version").Timeout(2 * time.Second).RunTrimOutput()
	if err != nil {
		return "", err
	}
	return out, nil
}

func (m *Choco) Install(name string, version string) error {
	path, err := m.locator.Resolve("choco.exe", "choco")
	if err != nil {
		return nil
	}

	args := []string{"install", name, "-y"}
	if version != "" {
		args = []string{"install", name, "--version", version, "-y"}
	}
	_, err = Command(path).Args(args...).RunTrimOutput()
	return err
}

func (m *Choco) InstallLatest(name string) error { return m.Install(name, "") }

func (m *Choco) Uninstall(name string) error {
	path, err := m.locator.Resolve("choco.exe", "choco")
	if err != nil {
		return nil
	}

	_, err = Command(path).Args("uninstall", name, "-y").RunTrimOutput()
	return err
}
