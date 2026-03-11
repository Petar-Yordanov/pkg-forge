package pkgmanagers

import (
	"fmt"
	"strings"

	"github.com/Petar-Yordanov/pkg-forge/common"
	"github.com/Petar-Yordanov/pkg-forge/pkgmanagers/managers"
)

type Manager = managers.Manager

func DefaultManagers() []managers.Manager {
	return []managers.Manager{
		&managers.AptGet{},
		&managers.Dnf{},
		&managers.Pacman{},
		&managers.Brew{},
		&managers.Cargo{},
		&managers.Npm{},
		&managers.Pipx{},
		&managers.Uv{},
		&managers.Choco{},
		&managers.Go{},
		&managers.Scoop{},
		&managers.Winget{},
		&managers.Pip{},
	}
}

var managerByDisplayName = func() map[string]Manager {
	m := make(map[string]Manager, 32)
	for _, inst := range DefaultManagers() {
		key := normalizeManagerDisplayName(inst.DisplayName())
		if key == "" {
			continue
		}
		m[key] = inst
	}

	if inst, ok := findManagerByID("apt-get"); ok {
		m["apt"] = inst
	}
	if inst, ok := findManagerByID("dnf"); ok {
		m["dnf"] = inst
	}
	if inst, ok := findManagerByID("pacman"); ok {
		m["pacman"] = inst
	}
	if inst, ok := findManagerByID("brew"); ok {
		m["homebrew"] = inst
	}
	if inst, ok := findManagerByID("cargo"); ok {
		m["cargo"] = inst
	}
	if inst, ok := findManagerByID("npm"); ok {
		m["npm"] = inst
	}
	if inst, ok := findManagerByID("pipx"); ok {
		m["pipx"] = inst
	}
	if inst, ok := findManagerByID("uv"); ok {
		m["uv"] = inst
	}
	if inst, ok := findManagerByID("choco"); ok {
		m["chocolatey"] = inst
	}
	if inst, ok := findManagerByID("go"); ok {
		m["go"] = inst
	}
	if inst, ok := findManagerByID("scoop"); ok {
		m["scoop"] = inst
	}
	if inst, ok := findManagerByID("winget"); ok {
		m["winget"] = inst
	}
	if inst, ok := findManagerByID("pip"); ok {
		m["pip"] = inst
	}

	return m
}()

func InstanceFromString(displayName string) (Manager, bool) {
	key := normalizeManagerDisplayName(displayName)
	if key == "" {
		return nil, false
	}
	m, ok := managerByDisplayName[key]
	return m, ok
}

func ResolveManagerForEntry(platform common.Platform, displayName string) (Manager, error) {
	displayName = strings.TrimSpace(displayName)
	if displayName != "" {
		m, ok := InstanceFromString(displayName)
		if !ok {
			return nil, fmt.Errorf("unknown package manager display name %q", displayName)
		}
		if !supportsPlatform(m, platform) {
			return nil, fmt.Errorf("package manager %q does not support platform %s", displayName, platform)
		}
		if _, err := m.Detect(); err != nil {
			return nil, fmt.Errorf("package manager %q is not available: %w", displayName, err)
		}
		return m, nil
	}

	for _, m := range DefaultManagers() {
		if !supportsPlatform(m, platform) {
			continue
		}
		dr, err := m.Detect()
		if err != nil {
			continue
		}
		if dr.Available {
			return m, nil
		}
	}

	return nil, fmt.Errorf("no available package manager detected for platform %s", platform)
}

func supportsPlatform(m Manager, platform common.Platform) bool {
	for _, p := range m.Platforms() {
		if p == platform {
			return true
		}
	}
	return false
}

func normalizeManagerDisplayName(s string) string {
	s = strings.TrimSpace(strings.ToLower(s))
	s = strings.ReplaceAll(s, "_", "")
	s = strings.ReplaceAll(s, "-", "")
	s = strings.ReplaceAll(s, " ", "")
	return s
}

func findManagerByID(id string) (Manager, bool) {
	id = strings.TrimSpace(strings.ToLower(id))
	if id == "" {
		return nil, false
	}

	for _, inst := range DefaultManagers() {
		if strings.EqualFold(inst.ID(), id) {
			return inst, true
		}
	}

	return nil, false
}
