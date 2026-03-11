package engine

import "github.com/Petar-Yordanov/pkg-forge/common"

type Context struct {
	Platform     common.Platform
	ManifestPath string
	ManifestName string
	State        *StateStore
	Events       Events
}
