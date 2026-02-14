//go:build darwin
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
