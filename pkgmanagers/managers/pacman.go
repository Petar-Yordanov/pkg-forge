package managers

import (
	"errors"

	"github.com/Petar-Yordanov/pkg-forge/common"
)

type Pacman struct{}

func (Pacman) ID() string { return "pacman" }
func (Pacman) DisplayName() string { return "pacman" }
func (Pacman) Platforms() []common.Platform { return []common.Platform{common.PlatformLinux} }

func (Pacman) Detect() (DetectResult, error) {
	return DetectResult{Available: false, Platform: common.CurrentPlatform()}, errors.New("not applicable on this platform")
}

func (Pacman) GetVersion() (string, error) {
	return "", errors.New("not implemented")
}

func (Pacman) Install(name string, version string) error { return nil }
func (Pacman) InstallLatest(name string) error { return nil }
func (Pacman) Uninstall(name string) error { return nil }
