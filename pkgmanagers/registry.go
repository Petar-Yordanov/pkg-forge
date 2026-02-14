package pkgmanagers

import "github.com/Petar-Yordanov/pkg-forge/pkgmanagers/managers"

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
