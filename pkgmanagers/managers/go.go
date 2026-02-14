package managers

import (
	"errors"
	"os/exec"
	"strings"
	"time"

	"github.com/Petar-Yordanov/pkg-forge/common"
	"github.com/samber/lo"
)

type Go struct{}

func (Go) ID() string { return "go" }
func (Go) DisplayName() string { return "Go" }
func (Go) Platforms() []common.Platform { return []common.Platform{common.PlatformWindows, common.PlatformLinux, common.PlatformMacOS} }

func (Go) Detect() (DetectResult, error) {
	cur := common.CurrentPlatform()

	bin, ok := lo.Find([]string{"go", "go.exe"}, func(c string) bool {
		_, err := exec.LookPath(c)
		return err == nil
	})
	if !ok {
		return DetectResult{Available: false, Platform: cur}, errors.New("not found in PATH")
	}

	path, _ := exec.LookPath(bin)

	out, err := Command(bin).Args("env", "GOVERSION").Timeout(2 * time.Second).RunTrimOutput()
	if err != nil {
		return DetectResult{Available: true, Path: path, Platform: cur}, nil
	}

	out = strings.TrimSpace(strings.TrimPrefix(out, "go"))
	return DetectResult{Available: true, Path: path, Version: out, Platform: cur}, nil
}

func (Go) Install(name string, version string) error { return nil }
func (Go) InstallLatest(name string) error { return nil }
func (Go) Uninstall(name string) error { return nil }
