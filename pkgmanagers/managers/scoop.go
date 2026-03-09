package managers

import (
	"regexp"
	"time"

	"github.com/Petar-Yordanov/pkg-forge/common"
)

type Scoop struct{}

func (*Scoop) ID() string          { return "scoop" }
func (*Scoop) DisplayName() string { return "Scoop" }

func (*Scoop) Platforms() []common.Platform {
	return []common.Platform{common.PlatformWindows}
}

func (m *Scoop) Detect() (DetectResult, error) {
	cur := common.CurrentPlatform()
	cmd := common.Command("scoop.cmd", "scoop")

	if err := cmd.Exists(); err != nil {
		return DetectResult{Available: false, Platform: cur}, err
	}

	return DetectResult{Available: true, Path: cmd.Path(), Platform: cur}, nil
}

func (m *Scoop) GetVersion() (string, error) {
	out, err := common.Command("scoop.cmd", "scoop").
		Args("--version").
		Timeout(2 * time.Second).
		RunTrimOutput()

	if err != nil {
		return "", err
	}

	return parseScoopVersion(out), nil
}

func (m *Scoop) Install(name string, version string) error {
	pkg := name
	if version != "" {
		pkg = name + "@" + version
	}

	_, err := common.Command("scoop.cmd", "scoop").
		Args("install", pkg).
		Timeout(15 * time.Minute).
		RunTrimOutput()

	return err
}

func (m *Scoop) InstallLatest(name string) error { return m.Install(name, "") }

func (m *Scoop) Uninstall(name string) error {
	_, err := common.Command("scoop.cmd", "scoop").
		Args("uninstall", name).
		Timeout(15 * time.Minute).
		RunTrimOutput()

	return err
}

func parseScoopVersion(raw string) string {
	reTag := regexp.MustCompile(`\btag:\s*v(\d+(?:\.\d+){1,3})\b`)
	reBump := regexp.MustCompile(`\bBump to version\s+(\d+(?:\.\d+){1,3})\b`)

	if m := reTag.FindStringSubmatch(raw); len(m) == 2 {
		return m[1]
	}

	if m := reBump.FindStringSubmatch(raw); len(m) == 2 {
		return m[1]
	}

	out := raw
	if m := regexp.MustCompile(`\r?\n`).Split(out, 2); len(m) > 0 {
		out = m[0]
	}

	return out
}
