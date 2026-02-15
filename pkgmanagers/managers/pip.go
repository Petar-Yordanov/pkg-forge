package managers

import (
	"time"

	"github.com/Petar-Yordanov/pkg-forge/common"
)

type Pip struct {
	locator ToolLocator
}

func (*Pip) ID() string          { return "pip" }
func (*Pip) DisplayName() string { return "pip" }
func (*Pip) Platforms() []common.Platform {
	return []common.Platform{common.PlatformWindows, common.PlatformLinux, common.PlatformMacOS}
}

func (m *Pip) Detect() (DetectResult, error) {
	cur := common.CurrentPlatform()

	py, err := m.locator.Resolve("python", "python3")
	if err != nil {
		return DetectResult{Available: false, Platform: cur}, err
	}

	return DetectResult{Available: true, Path: py, Platform: cur}, nil
}

func (m *Pip) GetVersion() (string, error) {
	py, err := m.locator.Resolve("python", "python3")
	if err != nil {
		return "", err
	}

	out, err := Command(py).Args("-c", "import pip; print(pip.__version__)").Timeout(2 * time.Second).RunTrimOutput()
	if err != nil {
		return "", err
	}
	return out, nil
}

func (*Pip) Install(name string, version string) error   { return nil }
func (*Pip) InstallLatest(name string) error             { return nil }
func (*Pip) Uninstall(name string) error                 { return nil }
