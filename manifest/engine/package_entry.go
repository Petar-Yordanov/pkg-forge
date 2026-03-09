package engine

import (
	"github.com/Petar-Yordanov/pkg-forge/common"
	"github.com/Petar-Yordanov/pkg-forge/manifest/parser"
	"github.com/Petar-Yordanov/pkg-forge/manifest/steps"
)

type PackageEntry struct {
	e parser.Entry
}

func NewPackageEntry(e parser.Entry) *PackageEntry { return &PackageEntry{e: e} }

func (p *PackageEntry) Raw() parser.Entry { return p.e }

func (p *PackageEntry) Applies(platform common.Platform) (bool, string) {
	return AppliesWhen(platform, p.e.When)
}

func (p *PackageEntry) Run(ctx *Context) error {
	for _, s := range p.e.PreInstall {
		ctx.Events.OnPreInstall(p.e, s)

		res, err := steps.RunStepSkeleton(ctx.Platform, p.e, s)
		if err != nil {
			return err
		}
		if res.Skipped {
			ctx.Events.OnStepSkip(p.e, s, res.SkipReason)
		}
	}

	ctx.Events.OnInstall(p.e)

	for _, s := range p.e.PostInstall {
		ctx.Events.OnPostInstall(p.e, s)

		res, err := steps.RunStepSkeleton(ctx.Platform, p.e, s)
		if err != nil {
			return err
		}
		if res.Skipped {
			ctx.Events.OnStepSkip(p.e, s, res.SkipReason)
		}
	}

	for _, s := range p.e.Validation {
		ctx.Events.OnValidation(p.e, s)

		res, err := steps.RunStepSkeleton(ctx.Platform, p.e, s)
		if err != nil {
			return err
		}
		if res.Skipped {
			ctx.Events.OnStepSkip(p.e, s, res.SkipReason)
		}
	}

	return nil
}
