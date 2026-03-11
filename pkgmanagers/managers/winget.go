package managers

import (
	"strings"
	"time"

	"github.com/Petar-Yordanov/pkg-forge/common"
)

type Winget struct{}

func (*Winget) ID() string          { return "winget" }
func (*Winget) DisplayName() string { return "WinGet" }

func (*Winget) Platforms() []common.Platform {
	return []common.Platform{common.PlatformWindows}
}

func (m *Winget) Detect() (DetectResult, error) {
	cur := common.CurrentPlatform()
	cmd := common.Command("winget.exe", "winget")

	if err := cmd.Exists(); err != nil {
		return DetectResult{Available: false, Platform: cur}, err
	}

	return DetectResult{Available: true, Path: cmd.Path(), Platform: cur}, nil
}

func (m *Winget) GetVersion() (string, error) {
	out, err := common.Command("winget.exe", "winget").
		Args("--version").
		Timeout(2 * time.Second).
		RunTrimOutput()

	if err != nil {
		return "", err
	}

	return strings.TrimSpace(strings.TrimPrefix(out, "v")), nil
}

func (m *Winget) Install(name string, version string) error {
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

	out, err := common.Command("winget.exe", "winget").
		Args(args...).
		Timeout(15 * time.Minute).
		RunTrimOutput()

	if err == nil {
		return nil
	}

	// TODO: This needs a more robust solution
	// Treat "already installed/no upgrade" as success for InstallLatest semantics.
	l := strings.ToLower(out)
	if strings.Contains(l, "found an existing package already installed") ||
		strings.Contains(l, "no available upgrade found") ||
		strings.Contains(l, "no newer package versions are available") {
		return nil
	}

	return err
}

func (m *Winget) InstallLatest(name string) error {
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

	out, err := common.Command("winget.exe", "winget").
		Args(args...).
		Timeout(15 * time.Minute).
		RunTrimOutput()

	if err == nil {
		return nil
	}

	// TODO: This needs a more robust solution
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
	args := []string{
		"uninstall",
		"--exact",
		"--id", name,
		"--source", "winget",
		"--silent",
		"--accept-source-agreements",
		"--disable-interactivity",
	}

	_, err := common.Command("winget.exe", "winget").
		Args(args...).
		Timeout(15 * time.Minute).
		RunTrimOutput()

	return err
}
