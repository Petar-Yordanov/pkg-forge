package managers

import (
	"errors"
	"os/exec"
	"regexp"
	"time"

	"github.com/Petar-Yordanov/pkg-forge/common"
	"github.com/samber/lo"
)

type Scoop struct{}

func (Scoop) ID() string { return "scoop" }
func (Scoop) DisplayName() string { return "Scoop" }
func (Scoop) Platforms() []common.Platform { return []common.Platform{common.PlatformWindows} }

func (Scoop) Detect() (DetectResult, error) {
	cur := common.CurrentPlatform()

	bin, ok := lo.Find([]string{"scoop", "scoop.cmd"}, func(c string) bool {
		_, err := exec.LookPath(c)
		return err == nil
	})
	if !ok {
		return DetectResult{Available: false, Platform: cur}, errors.New("not found in PATH")
	}

	path, _ := exec.LookPath(bin)

	out, err := Command(bin).Args("--version").Timeout(2 * time.Second).RunTrimOutput()
	if err != nil {
		return DetectResult{Available: true, Path: path, Platform: cur}, nil
	}

	return DetectResult{Available: true, Path: path, Version: parseScoopVersion(out), Platform: cur}, nil
}

func (Scoop) Install(name string, version string) error {
	bin, ok := lo.Find([]string{"scoop", "scoop.cmd"}, func(c string) bool {
		_, err := exec.LookPath(c)
		return err == nil
	})
	if !ok {
		return nil
	}

	pkg := name
	if version != "" {
		pkg = name + "@" + version
	}

	_, err := Command(bin).Args("install", pkg).RunTrimOutput()
	return err
}

func (Scoop) InstallLatest(name string) error { return (Scoop{}).Install(name, "") }

func (Scoop) Uninstall(name string) error {
	bin, ok := lo.Find([]string{"scoop", "scoop.cmd"}, func(c string) bool {
		_, err := exec.LookPath(c)
		return err == nil
	})
	if !ok {
		return nil
	}

	_, err := Command(bin).Args("uninstall", name).RunTrimOutput()
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
