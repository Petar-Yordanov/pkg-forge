package managers

import (
	"errors"

	"github.com/Petar-Yordanov/pkg-forge/common"
)

type Pipx struct{}

func (Pipx) ID() string { return "pipx" }
func (Pipx) DisplayName() string { return "pipx" }
func (Pipx) Platforms() []common.Platform { return []common.Platform{common.PlatformWindows, common.PlatformLinux, common.PlatformMacOS} }

func (Pipx) Detect() (DetectResult, error) {
	return DetectResult{Available: false, Platform: common.CurrentPlatform()}, errors.New("not found in PATH")
}

func (Pipx) GetVersion() (string, error) {
	return "", errors.New("not implemented")
}

func (Pipx) Install(name string, version string) error { return nil }
func (Pipx) InstallLatest(name string) error { return nil }
func (Pipx) Uninstall(name string) error { return nil }
