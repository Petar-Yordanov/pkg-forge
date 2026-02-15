//go:build windows

package managers_test

import (
	"regexp"
	"strings"
	"testing"

	"github.com/Petar-Yordanov/pkg-forge/pkgmanagers/managers"
)

const scoopPkg = "git"

func TestScoop_Detect(t *testing.T) {
	m := &managers.Scoop{}

	r, _ := m.Detect()
	if !r.Available {
		t.Fatalf("expected scoop available")
	}
}

func TestScoop_GetVersion(t *testing.T) {
	m := &managers.Scoop{}

	r, _ := m.Detect()
	if !r.Available {
		t.Skipf("scoop not available")
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

func TestScoop_Install_Uninstall(t *testing.T) {
	m := &managers.Scoop{}

	if err := m.Install(scoopPkg, ""); err != nil {
		t.Fatalf("install failed: %v", err)
	}
	if err := m.Uninstall(scoopPkg); err != nil {
		t.Fatalf("uninstall failed: %v", err)
	}
}

func TestScoop_InstallLatest_Uninstall(t *testing.T) {
	m := &managers.Scoop{}

	if err := m.InstallLatest(scoopPkg); err != nil {
		t.Fatalf("install latest failed: %v", err)
	}
	if err := m.Uninstall(scoopPkg); err != nil {
		t.Fatalf("uninstall failed: %v", err)
	}
}
