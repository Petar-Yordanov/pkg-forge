package pkgmanagers

import (
	"strings"

	"github.com/Petar-Yordanov/pkg-forge/pkgmanagers/managers"
)

type Manager = managers.Manager

func DefaultManagers() []managers.Manager {
	return []managers.Manager{
		managers.AptGet{},
		managers.Dnf{},
		managers.Pacman{},
		managers.Brew{},
		managers.Cargo{},
		managers.Npm{},
		managers.Pipx{},
		managers.Uv{},
		managers.Choco{},
		managers.Go{},
		managers.Scoop{},
		managers.Winget{},
		managers.Pip{},
	}
}

var managerByID = func() map[string]Manager {
	m := make(map[string]Manager, 32)
	for _, inst := range DefaultManagers() {
		id := strings.TrimSpace(strings.ToLower(inst.ID()))
		if id == "" {
			continue
		}
		m[id] = inst
	}
	return m
}()

func InstanceFromString(id string) (Manager, bool) {
	id = strings.TrimSpace(strings.ToLower(id))
	if id == "" {
		return nil, false
	}
	m, ok := managerByID[id]
	return m, ok
}
