//go:build linux || windows || darwin

package managers_test

import (
	"regexp"
	"strings"
	"testing"

	"github.com/Petar-Yordanov/pkg-forge/pkgmanagers/managers"
)

func TestPip_GetVersion(t *testing.T) {
	m := managers.Pip{}

	r, _ := m.Detect()
	if !r.Available {
		t.Skipf("python/pip not available")
	}

	v, err := m.GetVersion()
	if err != nil {
		t.Fatalf("get version failed: %v", err)
	}

	v = strings.TrimSpace(v)
	if v == "" {
		t.Fatalf("expected non-empty version")
	}

	if !regexp.MustCompile(`^\d+(\.\d+){1,3}([a-zA-Z0-9\.\-\+]+)?$`).MatchString(v) {
		t.Fatalf("unexpected version format: %q", v)
	}
}
