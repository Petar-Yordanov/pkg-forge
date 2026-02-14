//go:build windows

package managers_test

import (
	"testing"

	"github.com/Petar-Yordanov/pkg-forge/pkgmanagers/managers"
)

const scoopPkg = "git"

func TestScoop_Detect(t *testing.T) {
	r, _ := (managers.Scoop{}).Detect()
	if !r.Available {
		t.Fatalf("expected scoop available")
	}
}

func TestScoop_Install_Uninstall(t *testing.T) {
	m := managers.Scoop{}

	if err := m.Install(scoopPkg, ""); err != nil {
		t.Fatalf("install failed: %v", err)
	}
	if err := m.Uninstall(scoopPkg); err != nil {
		t.Fatalf("uninstall failed: %v", err)
	}
}

func TestScoop_InstallLatest_Uninstall(t *testing.T) {
	m := managers.Scoop{}

	if err := m.InstallLatest(scoopPkg); err != nil {
		t.Fatalf("install latest failed: %v", err)
	}
	if err := m.Uninstall(scoopPkg); err != nil {
		t.Fatalf("uninstall failed: %v", err)
	}
}
