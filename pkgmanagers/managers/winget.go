package managers

import (
	"errors"
	"strings"
	"time"

	"github.com/Petar-Yordanov/pkg-forge/common"
)

type Winget struct {
	locator ToolLocator
}

func (*Winget) ID() string          { return "winget" }
func (*Winget) DisplayName() string { return "WinGet" }
func (*Winget) Platforms() []common.Platform {
	return []common.Platform{common.PlatformWindows}
}

func (m *Winget) Detect() (DetectResult, error) {
	cur := common.CurrentPlatform()

	path, err := m.locator.Resolve("winget.exe", "winget")
	if err != nil {
		return DetectResult{Available: false, Platform: cur}, errors.New("not found in PATH")
	}

	return DetectResult{Available: true, Path: path, Platform: cur}, nil
}

func (m *Winget) GetVersion() (string, error) {
	path, err := m.locator.Resolve("winget.exe", "winget")
	if err != nil {
		return "", errors.New("not found in PATH")
	}

	out, err := Command(path).Args("--version").Timeout(2 * time.Second).RunTrimOutput()
	if err != nil {
		return "", err
	}

	return strings.TrimSpace(strings.TrimPrefix(out, "v")), nil
}

func (m *Winget) Install(name string, version string) error {
	path, err := m.locator.Resolve("winget.exe", "winget")
	if err != nil {
		return nil
	}

	args := []string{
		"install",
		"--exact",
		"--id", name,
		"--source", "winget",
		"--silent",
		"--accept-source-agreements",
		"--accept-package-agreements",
		"--disable-interactivity",
	}

	if version != "" {
		args = append(args, "--version", version)
	}

	_, err = Command(path).Args(args...).RunTrimOutput()
	return err
}

func (m *Winget) InstallLatest(name string) error {
	path, err := m.locator.Resolve("winget.exe", "winget")
	if err != nil {
		return nil
	}

	args := []string{
		"install",
		"--exact",
		"--id", name,
		"--source", "winget",
		"--silent",
		"--accept-source-agreements",
		"--accept-package-agreements",
		"--disable-interactivity",
	}

	out, err := Command(path).Args(args...).RunTrimOutput()
	if err == nil {
		return nil
	}

	// TODO: More robust solution for this
	// Treat "already installed/no upgrade" as success for InstallLatest semantics.
	l := strings.ToLower(out)
	if strings.Contains(l, "no available upgrade found") ||
		strings.Contains(l, "no newer package versions are available") ||
		strings.Contains(l, "found an existing package already installed") {
		return nil
	}

	return err
}

func (m *Winget) Uninstall(name string) error {
	path, err := m.locator.Resolve("winget.exe", "winget")
	if err != nil {
		return nil
	}

	args := []string{
		"uninstall",
		"--exact",
		"--id", name,
		"--source", "winget",
		"--silent",
		"--accept-source-agreements",
		"--disable-interactivity",
	}

	_, err = Command(path).Args(args...).RunTrimOutput()
	return err
}
