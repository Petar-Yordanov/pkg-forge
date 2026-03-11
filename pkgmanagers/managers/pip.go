package managers

import (
	"time"

	"github.com/Petar-Yordanov/pkg-forge/common"
)

type Pip struct{}

func (*Pip) ID() string          { return "pip" }
func (*Pip) DisplayName() string { return "pip" }

func (*Pip) Platforms() []common.Platform {
	return []common.Platform{common.PlatformWindows, common.PlatformLinux, common.PlatformMacOS}
}

func (m *Pip) Detect() (DetectResult, error) {
	cur := common.CurrentPlatform()

	cmd := common.Command("python.exe", "python", "python3", "py")
	if err := cmd.Exists(); err != nil {
		return DetectResult{Available: false, Platform: cur}, err
	}

	return DetectResult{Available: true, Path: cmd.Path(), Platform: cur}, nil
}

func (m *Pip) GetVersion() (string, error) {
	out, err := common.Command("python", "python3", "py").
		Args("-c", "import pip; print(pip.__version__)").
		Timeout(2 * time.Second).
		RunTrimOutput()

	if err != nil {
		return "", err
	}

	return out, nil
}

func (*Pip) Install(name string, version string) error { return nil }
func (*Pip) InstallLatest(name string) error           { return nil }
func (*Pip) Uninstall(name string) error               { return nil }
