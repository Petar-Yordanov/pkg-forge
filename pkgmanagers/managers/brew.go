package managers

import (
	"errors"
	"os"
	"os/exec"
	"strings"
	"time"

	"github.com/Petar-Yordanov/pkg-forge/common"
	"github.com/samber/lo"
)

type Brew struct{}

func (Brew) ID() string          { return "brew" }
func (Brew) DisplayName() string { return "Homebrew" }

func (Brew) Platforms() []common.Platform {
	return []common.Platform{common.PlatformMacOS}
}

func (Brew) Detect() (DetectResult, error) {
	cur := common.CurrentPlatform()

	bin, ok := lo.Find([]string{"brew"}, func(c string) bool {
		_, err := exec.LookPath(c)
		return err == nil
	})

	if !ok {
		return DetectResult{Available: false, Platform: cur}, errors.New("not found in PATH")
	}

	path, _ := exec.LookPath("brew")
	if path == "" {
		path = bin
	}

	out, err := Command(bin).Args("--version").Timeout(2 * time.Second).RunTrimOutput()
	if err != nil {
		return DetectResult{Available: true, Path: path, Platform: cur}, nil
	}

	ver := strings.TrimSpace(out)
	return DetectResult{Available: true, Path: path, Version: ver, Platform: cur}, nil
}

func (Brew) Install(name string, version string) error {
	bin, ok := lo.Find([]string{"brew"}, func(c string) bool {
		_, err := exec.LookPath(c)
		return err == nil
	})

	pkg := name
	if version != "" && !strings.Contains(name, "@") {
		pkg = name + "@" + version
	}

	_, err := Command(bin).
		Args("install", pkg).
		Timeout(15 * time.Minute).
		RunTrimOutput()
	return err
}

func (Brew) InstallLatest(name string) error {
	return (Brew{}).Install(name, "")
}

func (Brew) Uninstall(name string) error {
	bin, ok := lo.Find([]string{"brew"}, func(c string) bool {
		_, err := exec.LookPath(c)
		return err == nil
	})

	_, err := Command(bin).
		Args("uninstall", "--force", name).
		Timeout(15 * time.Minute).
		RunTrimOutput()
	return err
}
