package managers

import (
	"errors"

	"github.com/Petar-Yordanov/pkg-forge/common"
)

type Dnf struct{}

func (Dnf) ID() string { return "dnf" }
func (Dnf) DisplayName() string { return "DNF" }
func (Dnf) Platforms() []common.Platform { return []common.Platform{common.PlatformLinux} }

func (Dnf) Detect() (DetectResult, error) {
	return DetectResult{Available: false, Platform: common.CurrentPlatform()}, errors.New("not applicable on this platform")
}

func (Dnf) GetVersion() (string, error) {
	return "", errors.New("not implemented")
}

func (Dnf) Install(name string, version string) error { return nil }
func (Dnf) InstallLatest(name string) error { return nil }
func (Dnf) Uninstall(name string) error { return nil }
