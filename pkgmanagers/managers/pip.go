package managers

import (
	"errors"
	"os/exec"
	"time"

	"github.com/Petar-Yordanov/pkg-forge/common"
	"github.com/samber/lo"
)

type Pip struct{}

func (Pip) ID() string { return "pip" }
func (Pip) DisplayName() string { return "pip" }
func (Pip) Platforms() []common.Platform { return []common.Platform{common.PlatformWindows, common.PlatformLinux, common.PlatformMacOS} }

func (Pip) Detect() (DetectResult, error) {
	cur := common.CurrentPlatform()

	bin, ok := lo.Find([]string{"python", "python3"}, func(c string) bool {
		_, err := exec.LookPath(c)
		return err == nil
	})
	if !ok {
		return DetectResult{Available: false, Platform: cur}, errors.New("not found in PATH")
	}

	path, _ := exec.LookPath(bin)

	out, err := Command(bin).Args("-c", "import pip; print(pip.__version__)").Timeout(2 * time.Second).RunTrimOutput()
	if err != nil {
		return DetectResult{Available: true, Path: path, Platform: cur}, nil
	}

	return DetectResult{Available: true, Path: path, Version: out, Platform: cur}, nil
}

func (Pip) Install(name string, version string) error { return nil }
func (Pip) InstallLatest(name string) error { return nil }
func (Pip) Uninstall(name string) error { return nil }
