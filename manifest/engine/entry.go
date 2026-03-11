package engine

import (
	"github.com/Petar-Yordanov/pkg-forge/common"
	"github.com/Petar-Yordanov/pkg-forge/manifest/parser"
)

type Entry interface {
	Raw() parser.Entry
	Applies(platform common.Platform) (bool, string)
	Run(ctx *Context) error
	Uninstall(ctx *Context) error
}
