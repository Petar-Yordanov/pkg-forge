package pkgmanagers

import (
	"context"

	"github.com/Petar-Yordanov/pkg-forge/common"
	"github.com/samber/lo"
)

type ManagerStatus struct {
	Cmd       string   `json:"cmd"`
	Name      string   `json:"name"`
	Platforms []string `json:"platforms"`
	Available bool     `json:"available"`
	Path      string   `json:"path,omitempty"`
	Version   string   `json:"version,omitempty"`
	Err       string   `json:"error,omitempty"`
}

func DetectAll(ctx context.Context) []ManagerStatus {
	_ = ctx

	ms := DefaultManagers()
	out := make([]ManagerStatus, 0, len(ms))

	for _, m := range ms {
		st := ManagerStatus{
			Cmd:       m.ID(),
			Name:      m.DisplayName(),
			Platforms: lo.Map(m.Platforms(), func(p common.Platform, _ int) string { return string(p) }),
		}

		dr, err := m.Detect()
		if err != nil {
			st.Available = false
			st.Err = err.Error()
			out = append(out, st)
			continue
		}

		st.Available = dr.Available
		st.Path = dr.Path

		if dr.Available {
			if ver, verr := m.GetVersion(); verr == nil {
				st.Version = ver
			}
		}

		out = append(out, st)
	}

	return out
}
