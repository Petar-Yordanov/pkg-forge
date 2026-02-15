package managers

import "github.com/Petar-Yordanov/pkg-forge/common"

type Manager interface {
	ID() string
	DisplayName() string
	Platforms() []common.Platform

	Detect() (DetectResult, error)
	GetVersion() (string, error)

	Install(name string, version string) error
	InstallLatest(name string) error
	Uninstall(name string) error
}

type DetectResult struct {
	Available bool            `json:"available"`
	Path      string          `json:"path,omitempty"`
	Platform  common.Platform `json:"platform"`
}
