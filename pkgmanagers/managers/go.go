package managers

import (
	"strings"
	"time"

	"github.com/Petar-Yordanov/pkg-forge/common"
)

type Go struct{}

func (*Go) ID() string          { return "go" }
func (*Go) DisplayName() string { return "Go" }

func (*Go) Platforms() []common.Platform {
	return []common.Platform{common.PlatformWindows, common.PlatformLinux, common.PlatformMacOS}
}

func (m *Go) Detect() (DetectResult, error) {
	cur := common.CurrentPlatform()
	cmd := common.Command("go.exe", "go")

	if err := cmd.Exists(); err != nil {
		return DetectResult{Available: false, Platform: cur}, err
	}

	return DetectResult{Available: true, Path: cmd.Path(), Platform: cur}, nil
}

func (m *Go) GetVersion() (string, error) {
	out, err := common.Command("go.exe", "go").
		Args("env", "GOVERSION").
		Timeout(2 * time.Second).
		RunTrimOutput()

	if err != nil {
		return "", err
	}

	out = strings.TrimSpace(strings.TrimPrefix(out, "go"))
	return out, nil
}

func (*Go) Install(name string, version string) error { return nil }
func (*Go) InstallLatest(name string) error           { return nil }
func (*Go) Uninstall(name string) error               { return nil }
