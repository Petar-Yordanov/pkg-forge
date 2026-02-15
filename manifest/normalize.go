package manifest

import (
	"errors"
	"fmt"
	"maps"
	"time"

	"github.com/Petar-Yordanov/pkg-forge/common"
)

func Normalize(doc *Document, baseDir string, platform common.Platform) ([]StepPlan, error) {
	if doc == nil {
		return nil, errors.New("nil document")
	}

	type effDefaults struct {
		failOnStderr  bool
		env           map[string]string
		retries       int
		retryDelaySec int
		timeoutSec    *int
	}

	effectiveDefaults := func(e *Entry) effDefaults {
		d := effDefaults{
			failOnStderr:  false,
			env:           map[string]string{},
			retries:       0,
			retryDelaySec: 0,
			timeoutSec:    nil,
		}
		if e.FailOnStderr != nil {
			d.failOnStderr = *e.FailOnStderr
		}
		if e.Defaults == nil || e.Defaults.Step == nil {
			return d
		}
		sd := e.Defaults.Step
		if sd.FailOnStderr != nil {
			d.failOnStderr = *sd.FailOnStderr
		}
		if sd.Env != nil {
			d.env = maps.Clone(sd.Env)
		}
		if sd.Retries != nil {
			d.retries = *sd.Retries
		}
		if sd.RetryDelaySec != nil {
			d.retryDelaySec = *sd.RetryDelaySec
		}
		if sd.TimeoutSec != nil {
			d.timeoutSec = sd.TimeoutSec
		}
		return d
	}

	pickSelect := func(sn *SelectNode, platform common.Platform) *Step {
		if sn == nil {
			return nil
		}
		if sn.Items != nil {
			if s := sn.Items[string(platform)]; s != nil {
				return s
			}
		}
		return sn.Default
	}

	resolveStep := func(s Step, d effDefaults) ResolvedStep {
		rs := ResolvedStep{
			When:  s.When,
			Shell: s.Shell,
			Args:  append([]string(nil), s.Args...),
			Cwd:   s.Cwd,
			Env:   maps.Clone(d.env),

			FailOnStderr: d.failOnStderr,
			Retries:      d.retries,
			RetryDelay:   time.Duration(d.retryDelaySec) * time.Second,
		}

		if s.Env != nil {
			maps.Copy(rs.Env, s.Env)
		}
		if s.FailOnStderr != nil {
			rs.FailOnStderr = *s.FailOnStderr
		}
		if s.Retries != nil {
			rs.Retries = *s.Retries
		}
		if s.RetryDelaySec != nil {
			rs.RetryDelay = time.Duration(*s.RetryDelaySec) * time.Second
		}

		timeoutPtr := d.timeoutSec
		if s.TimeoutSec != nil {
			timeoutPtr = s.TimeoutSec
		}
		if timeoutPtr == nil {
			rs.TimeoutInfinite = true
		} else {
			rs.Timeout = time.Duration(*timeoutPtr) * time.Second
		}

		if s.Cmd != "" && s.CmdFile != "" {
			rs.ExecKind = "invalid"
			rs.Cmd = s.Cmd
			rs.CmdFile = common.ResolvePath(baseDir, s.CmdFile)
			return rs
		}
		if s.Cmd != "" {
			rs.ExecKind = "cmd"
			rs.Cmd = s.Cmd
			return rs
		}
		rs.ExecKind = "cmdFile"
		rs.CmdFile = common.ResolvePath(baseDir, s.CmdFile)
		return rs
	}

	resolveOne := func(sl StepLike, d effDefaults) (ResolvedStep, bool, error) {
		if sl.Select != nil {
			sel := pickSelect(sl.Select, platform)
			if sel == nil {
				return ResolvedStep{}, false, nil
			}
			rs := resolveStep(*sel, d)
			if rs.When != nil && rs.When.Platform != "" && rs.When.Platform != platform {
				return ResolvedStep{}, false, nil
			}
			return rs, true, nil
		}
		if sl.Step == nil {
			return ResolvedStep{}, false, errors.New("empty step")
		}
		rs := resolveStep(*sl.Step, d)
		if rs.When != nil && rs.When.Platform != "" && rs.When.Platform != platform {
			return ResolvedStep{}, false, nil
		}
		return rs, true, nil
	}

	resolveList := func(list []StepLike, d effDefaults) ([]ResolvedStep, error) {
		if len(list) == 0 {
			return nil, nil
		}
		out := make([]ResolvedStep, 0, len(list))
		for i := range list {
			rs, ok, err := resolveOne(list[i], d)
			if err != nil {
				return nil, fmt.Errorf("step[%d]: %w", i, err)
			}
			if ok {
				out = append(out, rs)
			}
		}
		return out, nil
	}

	normalizeEntry := func(e *Entry) (StepPlan, error) {
		if e.Kind == "" {
			return StepPlan{}, errors.New("missing kind")
		}
		if e.Name == "" {
			return StepPlan{}, errors.New("missing name")
		}
		d := effectiveDefaults(e)

		p := StepPlan{Kind: e.Kind, Name: e.Name, Version: e.Version}
		var err error

		if p.PreInstall, err = resolveList(e.PreInstall, d); err != nil {
			return StepPlan{}, fmt.Errorf("preInstall: %w", err)
		}
		if p.PostInstall, err = resolveList(e.PostInstall, d); err != nil {
			return StepPlan{}, fmt.Errorf("postInstall: %w", err)
		}
		if p.Validation, err = resolveList(e.Validation, d); err != nil {
			return StepPlan{}, fmt.Errorf("validation: %w", err)
		}
		if p.Pre, err = resolveList(e.Pre, d); err != nil {
			return StepPlan{}, fmt.Errorf("pre: %w", err)
		}
		if p.Post, err = resolveList(e.Post, d); err != nil {
			return StepPlan{}, fmt.Errorf("post: %w", err)
		}
		if p.Steps, err = resolveList(e.Steps, d); err != nil {
			return StepPlan{}, fmt.Errorf("steps: %w", err)
		}
		if e.Cmd != nil {
			rs, ok, err := resolveOne(*e.Cmd, d)
			if err != nil {
				return StepPlan{}, fmt.Errorf("cmd: %w", err)
			}
			if ok {
				p.Cmd = &rs
			}
		}

		return p, nil
	}

	out := make([]StepPlan, 0, len(doc.Entries))
	for i := range doc.Entries {
		e := doc.Entries[i]
		p, err := normalizeEntry(&e)
		if err != nil {
			return nil, fmt.Errorf("entries[%d] %s/%s: %w", i, e.Kind, e.Name, err)
		}
		out = append(out, p)
	}
	return out, nil
}
