package managers

import (
	"errors"

	"github.com/Petar-Yordanov/pkg-forge/common"
)

type AptGet struct{}

func (AptGet) ID() string { return "apt-get" }
func (AptGet) DisplayName() string { return "APT" }
func (AptGet) Platforms() []common.Platform { return []common.Platform{common.PlatformLinux} }

func (AptGet) Detect() (DetectResult, error) {
	return DetectResult{Available: false, Platform: common.CurrentPlatform()}, errors.New("not applicable on this platform")
}

func (AptGet) Install(name string, version string) error { return nil }
func (AptGet) InstallLatest(name string) error { return nil }
func (AptGet) Uninstall(name string) error { return nil }
