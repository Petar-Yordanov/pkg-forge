package managers

import (
	"strings"
	"time"

	"github.com/Petar-Yordanov/pkg-forge/common"
)

type Brew struct{}

func (*Brew) ID() string          { return "brew" }
func (*Brew) DisplayName() string { return "Homebrew" }

func (*Brew) Platforms() []common.Platform {
	return []common.Platform{common.PlatformMacOS}
}

func (m *Brew) Detect() (DetectResult, error) {
	cur := common.CurrentPlatform()
	cmd := common.Command("brew")

	if err := cmd.Exists(); err != nil {
		return DetectResult{Available: false, Platform: cur}, err
	}

	return DetectResult{Available: true, Path: cmd.Path(), Platform: cur}, nil
}

func (m *Brew) GetVersion() (string, error) {
	out, err := common.Command("brew").
		Args("--version").
		Timeout(2 * time.Second).
		RunTrimOutput()

	if err != nil {
		return "", err
	}
	return strings.TrimSpace(out), nil
}

func (m *Brew) Install(name string, version string) error {
	pkg := name
	if version != "" && !strings.Contains(name, "@") {
		pkg = name + "@" + version
	}

	_, err := common.Command("brew").
		Args("install", pkg).
		Timeout(15 * time.Minute).
		RunTrimOutput()

	if err != nil {
		return err
	}
	return nil
}

func (m *Brew) InstallLatest(name string) error { return m.Install(name, "") }

func (m *Brew) Uninstall(name string) error {
	_, err := common.Command("brew").
		Args("uninstall", "--force", name).
		Timeout(15 * time.Minute).
		RunTrimOutput()

	if err != nil {
		return err
	}
	return nil
}
