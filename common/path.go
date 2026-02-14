package common

import "path/filepath"

func ResolvePath(baseDir, p string) string {
	if p == "" {
		return p
	}
	if filepath.IsAbs(p) {
		return p
	}
	if baseDir == "" {
		return p
	}
	return filepath.Clean(filepath.Join(baseDir, p))
}
