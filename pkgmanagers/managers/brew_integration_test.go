//go:build darwin || linux

package managers_test

import (
	"testing"

	"github.com/Petar-Yordanov/pkg-forge/pkgmanagers/managers"
)

const brewPkg = "git"

func TestBrew_Detect(t *testing.T) {
	r, _ := (managers.Brew{}).Detect()
	if !r.Available {
		t.Fatalf("expected brew available")
	}
}

func TestBrew_GetVersion(t *testing.T) {
	m := managers.Brew{}

	r, _ := m.Detect()
	if !r.Available {
		t.Skipf("brew not available")
	}

	v, err := m.GetVersion()
	if err != nil {
		t.Fatalf("get version failed: %v", err)
	}
	v = strings.TrimSpace(v)
	if v == "" {
		t.Fatalf("expected non-empty version")
	}

	if !strings.HasPrefix(v, "Homebrew") {
		t.Fatalf("unexpected version output: %q", v)
	}
}

func TestBrew_Install_Uninstall(t *testing.T) {
	m := managers.Brew{}

	if err := m.Install(brewPkg, ""); err != nil {
		t.Fatalf("install failed: %v", err)
	}
	if err := m.Uninstall(brewPkg); err != nil {
		t.Fatalf("uninstall failed: %v", err)
	}
}

func TestBrew_InstallLatest_Uninstall(t *testing.T) {
	m := managers.Brew{}

	if err := m.InstallLatest(brewPkg); err != nil {
		t.Fatalf("install latest failed: %v", err)
	}
	if err := m.Uninstall(brewPkg); err != nil {
		t.Fatalf("uninstall failed: %v", err)
	}
}
