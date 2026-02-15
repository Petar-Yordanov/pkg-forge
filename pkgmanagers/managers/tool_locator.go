package managers

import (
	"errors"
	"os/exec"
	"sync"
)

var ErrToolNotFound = errors.New("not found in PATH")

type ToolLocator struct {
	once sync.Once
	path string
	err  error
}

func (l *ToolLocator) Resolve(candidates ...string) (string, error) {
	l.once.Do(func() {
		for _, c := range candidates {
			if p, err := exec.LookPath(c); err == nil {
				l.path = p
				l.err = nil
				return
			}
		}
		l.err = ErrToolNotFound
	})

	if l.err != nil {
		return "", l.err
	}
	return l.path, nil
}
