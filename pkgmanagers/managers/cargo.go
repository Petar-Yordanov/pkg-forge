package managers

import (
	"errors"

	"github.com/Petar-Yordanov/pkg-forge/common"
)

type Cargo struct{}

func (Cargo) ID() string { return "cargo" }
func (Cargo) DisplayName() string { return "Cargo" }
func (Cargo) Platforms() []common.Platform { return []common.Platform{common.PlatformWindows, common.PlatformLinux, common.PlatformMacOS} }

func (Cargo) Detect() (DetectResult, error) {
	return DetectResult{Available: false, Platform: common.CurrentPlatform()},
		errors.New("not applicable on this platform")
}

func (Cargo) GetVersion() (string, error) {
	return "", errors.New("not implemented")
}

func (Cargo) Install(name string, version string) error { return nil }
func (Cargo) InstallLatest(name string) error { return nil }
func (Cargo) Uninstall(name string) error { return nil }
