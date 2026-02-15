//go:build linux || windows || darwin

package managers_test

import (
	"regexp"
	"strings"
	"testing"

	"github.com/Petar-Yordanov/pkg-forge/pkgmanagers/managers"
)

func TestGo_GetVersion(t *testing.T) {
	m := &managers.Go{}

	r, _ := m.Detect()
	if !r.Available {
		t.Skipf("go not available")
	}

	v, err := m.GetVersion()
	if err != nil {
		t.Fatalf("get version failed: %v", err)
	}

	v = strings.TrimSpace(v)
	if v == "" {
		t.Fatalf("expected non-empty version")
	}

	if !regexp.MustCompile(`^\d+\.\d+(\.\d+)?$`).MatchString(v) {
		t.Fatalf("unexpected version format: %q", v)
	}
}
