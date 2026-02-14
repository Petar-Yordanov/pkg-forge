package managers

import (
	"errors"

	"github.com/Petar-Yordanov/pkg-forge/common"
)

type Brew struct{}

func (Brew) ID() string { return "brew" }
func (Brew) DisplayName() string { return "Homebrew" }
func (Brew) Platforms() []common.Platform { return []common.Platform{common.PlatformMacOS, common.PlatformLinux} }

func (Brew) Detect() (DetectResult, error) {
	return DetectResult{Available: false, Platform: common.CurrentPlatform()}, errors.New("not applicable on this platform")
}

func (Brew) Install(name string, version string) error { return nil }
func (Brew) InstallLatest(name string) error { return nil }
func (Brew) Uninstall(name string) error { return nil }
