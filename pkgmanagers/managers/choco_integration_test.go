//go:build windows

package managers_test

import (
	"testing"
	"strings"
	"regexp"
	"github.com/Petar-Yordanov/pkg-forge/pkgmanagers/managers"
)

const chocoPkg = "git"

func TestChoco_Detect(t *testing.T) {
	r, _ := (managers.Choco{}).Detect()
	if !r.Available {
		t.Fatalf("expected choco available")
	}
}

func TestChoco_GetVersion(t *testing.T) {
	m := managers.Choco{}

	r, _ := m.Detect()
	if !r.Available {
		t.Skipf("choco not available")
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
