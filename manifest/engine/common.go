package engine

import (
	"strings"

	"github.com/Petar-Yordanov/pkg-forge/common"
	"github.com/Petar-Yordanov/pkg-forge/manifest/parser"
)

func AppliesWhen(cur common.Platform, w *parser.When) (bool, string) {
	if w == nil || len(w.Platform) == 0 {
		return true, ""
	}
	for _, p := range w.Platform {
		if strings.EqualFold(string(p), string(cur)) {
			return true, ""
		}
	}
	return false, "platform mismatch"
}
