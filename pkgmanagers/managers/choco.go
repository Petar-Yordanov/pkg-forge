package managers

import (
	"time"

	"github.com/Petar-Yordanov/pkg-forge/common"
)

type Choco struct{}

func (*Choco) ID() string          { return "choco" }
func (*Choco) DisplayName() string { return "Chocolatey" }

func (*Choco) Platforms() []common.Platform {
	return []common.Platform{common.PlatformWindows}
}

func (m *Choco) Detect() (DetectResult, error) {
	cur := common.CurrentPlatform()
	cmd := common.Command("choco.exe", "choco")

	if err := cmd.Exists(); err != nil {
		return DetectResult{Available: false, Platform: cur}, err
	}

	return DetectResult{Available: true, Path: cmd.Path(), Platform: cur}, nil
}

func (m *Choco) GetVersion() (string, error) {
	out, err := common.Command("choco.exe", "choco").
		Args("--version").
		Timeout(2 * time.Second).
		RunTrimOutput()

	if err != nil {
		return "", err
	}

	return out, nil
}

func (m *Choco) Install(name string, version string) error {
	args := []string{"install", name, "-y"}
	if version != "" {
		args = []string{"install", name, "--version", version, "-y"}
	}

	_, err := common.Command("choco.exe", "choco").
		Args(args...).
		Timeout(15 * time.Minute).
		RunTrimOutput()

	return err
}

func (m *Choco) InstallLatest(name string) error { return m.Install(name, "") }

func (m *Choco) Uninstall(name string) error {
	_, err := common.Command("choco.exe", "choco").
		Args("uninstall", name, "-y").
		Timeout(15 * time.Minute).
		RunTrimOutput()

	return err
}
