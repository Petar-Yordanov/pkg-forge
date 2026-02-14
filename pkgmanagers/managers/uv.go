package managers

import (
	"errors"

	"github.com/Petar-Yordanov/pkg-forge/common"
)

type Uv struct{}

func (Uv) ID() string { return "uv" }
func (Uv) DisplayName() string { return "uv" }
func (Uv) Platforms() []common.Platform { return []common.Platform{common.PlatformWindows, common.PlatformLinux, common.PlatformMacOS} }

func (Uv) Detect() (DetectResult, error) {
	return DetectResult{Available: false, Platform: common.CurrentPlatform()}, errors.New("not found in PATH")
}

func (Uv) Install(name string, version string) error { return nil }
func (Uv) InstallLatest(name string) error { return nil }
func (Uv) Uninstall(name string) error { return nil }
