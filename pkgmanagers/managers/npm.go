package managers

import (
	"errors"

	"github.com/Petar-Yordanov/pkg-forge/common"
)

type Npm struct{}

func (Npm) ID() string { return "npm" }
func (Npm) DisplayName() string { return "npm" }
func (Npm) Platforms() []common.Platform { return []common.Platform{common.PlatformWindows, common.PlatformLinux, common.PlatformMacOS} }

func (Npm) Detect() (DetectResult, error) {
	return DetectResult{Available: false, Platform: common.CurrentPlatform()}, errors.New("not found in PATH")
}

func (Npm) GetVersion() (string, error) {
	return "", errors.New("not implemented")
}

func (Npm) Install(name string, version string) error { return nil }
func (Npm) InstallLatest(name string) error { return nil }
func (Npm) Uninstall(name string) error { return nil }
