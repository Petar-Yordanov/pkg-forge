//go:build windows

package managers_test

import (
	"testing"

	"github.com/Petar-Yordanov/pkg-forge/pkgmanagers/managers"
)

const wingetPkgID = "Git.Git"

func TestWinget_Detect(t *testing.T) {
	r, _ := (managers.Winget{}).Detect()
	if !r.Available {
		t.Fatalf("expected winget available")
	}
}

func TestWinget_Install_Uninstall(t *testing.T) {
	m := managers.Winget{}

	if err := m.Install(wingetPkgID, ""); err != nil {
		t.Fatalf("install failed: %v", err)
	}
	if err := m.Uninstall(wingetPkgID); err != nil {
		t.Fatalf("uninstall failed: %v", err)
	}
}

func TestWinget_InstallLatest_Uninstall(t *testing.T) {
	m := managers.Winget{}

	if err := m.InstallLatest(wingetPkgID); err != nil {
		t.Fatalf("install latest failed: %v", err)
	}
	if err := m.Uninstall(wingetPkgID); err != nil {
		t.Fatalf("uninstall failed: %v", err)
	}
}
