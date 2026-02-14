//go:build windows

package managers_test

import (
	"testing"

	"github.com/Petar-Yordanov/pkg-forge/pkgmanagers/managers"
)

const chocoPkg = "git"

func TestChoco_Detect(t *testing.T) {
	r, _ := (managers.Choco{}).Detect()
	if !r.Available {
		t.Fatalf("expected choco available")
	}
}

func TestChoco_Install_Uninstall(t *testing.T) {
	m := managers.Choco{}

	if err := m.Install(chocoPkg, ""); err != nil {
		t.Fatalf("install failed: %v", err)
	}
	if err := m.Uninstall(chocoPkg); err != nil {
		t.Fatalf("uninstall failed: %v", err)
	}
}

func TestChoco_InstallLatest_Uninstall(t *testing.T) {
	m := managers.Choco{}

	if err := m.InstallLatest(chocoPkg); err != nil {
		t.Fatalf("install latest failed: %v", err)
	}
	if err := m.Uninstall(chocoPkg); err != nil {
		t.Fatalf("uninstall failed: %v", err)
	}
}
