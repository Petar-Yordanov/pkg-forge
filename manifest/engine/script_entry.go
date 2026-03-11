package engine

import (
	"github.com/Petar-Yordanov/pkg-forge/common"
	"github.com/Petar-Yordanov/pkg-forge/manifest/parser"
	"github.com/Petar-Yordanov/pkg-forge/manifest/steps"
)

type ScriptEntry struct {
	e parser.Entry
}

func NewScriptEntry(e parser.Entry) *ScriptEntry { return &ScriptEntry{e: e} }

func (s *ScriptEntry) Raw() parser.Entry { return s.e }

func (s *ScriptEntry) Applies(platform common.Platform) (bool, string) {
	return AppliesWhen(platform, s.e.When)
}

func (s *ScriptEntry) Run(ctx *Context) error {
	for _, st := range s.e.Steps {
		ctx.Events.OnStep(s.e, st)

		res, err := steps.RunStep(ctx.Platform, s.e, st)
		if err != nil {
			return err
		}
		if res.Skipped {
			ctx.Events.OnStepSkip(s.e, st, res.SkipReason)
		}
	}

	for _, st := range s.e.Validation {
		ctx.Events.OnValidation(s.e, st)

		res, err := steps.RunStep(ctx.Platform, s.e, st)
		if err != nil {
			return err
		}
		if res.Skipped {
			ctx.Events.OnStepSkip(s.e, st, res.SkipReason)
		}
	}

	return nil
}

func (s *ScriptEntry) Uninstall(ctx *Context) error {
	_ = ctx
	return nil
}
