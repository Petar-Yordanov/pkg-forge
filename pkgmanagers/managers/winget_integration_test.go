//go:build windows

package managers_test

import (
	"regexp"
	"strings"
	"testing"

	"github.com/Petar-Yordanov/pkg-forge/pkgmanagers/managers"
)

const wingetPkgID = "7zip.7zip"

func TestWinget_Detect(t *testing.T) {
	m := &managers.Winget{}

	r, _ := m.Detect()
	if !r.Available {
		t.Fatalf("expected winget available")
	}
}

func TestWinget_GetVersion(t *testing.T) {
	m := &managers.Winget{}

	r, _ := m.Detect()
	if !r.Available {
		t.Skipf("winget not available")
	}

	v, err := m.GetVersion()
	if err != nil {
		t.Fatalf("get version failed: %v", err)
	}

	v = strings.TrimSpace(v)
	if v == "" {
		t.Fatalf("expected non-empty version")
	}

	if !regexp.MustCompile(`^\d+(\.\d+){1,3}$`).MatchString(v) {
		t.Fatalf("unexpected version format: %q", v)
	}
}

func TestWinget_Install_Uninstall(t *testing.T) {
	m := &managers.Winget{}

	if err := m.Install(wingetPkgID, ""); err != nil {
		t.Fatalf("install failed: %v", err)
	}
	if err := m.Uninstall(wingetPkgID); err != nil {
		t.Fatalf("uninstall failed: %v", err)
	}
}

func TestWinget_InstallLatest_Uninstall(t *testing.T) {
	m := &managers.Winget{}

	if err := m.InstallLatest(wingetPkgID); err != nil {
		t.Fatalf("install latest failed: %v", err)
	}
	if err := m.Uninstall(wingetPkgID); err != nil {
		t.Fatalf("uninstall failed: %v", err)
	}
}
