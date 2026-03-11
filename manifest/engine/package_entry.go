package engine

import (
	"fmt"
	"strings"

	"github.com/Petar-Yordanov/pkg-forge/common"
	"github.com/Petar-Yordanov/pkg-forge/manifest/parser"
	"github.com/Petar-Yordanov/pkg-forge/manifest/steps"
	"github.com/Petar-Yordanov/pkg-forge/pkgmanagers"
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

		res, err := steps.RunStep(ctx.Platform, p.e, s)
		if err != nil {
			return err
		}
		if res.Skipped {
			ctx.Events.OnStepSkip(p.e, s, res.SkipReason)
		}
	}

	ctx.Events.OnInstall(p.e)

	m, err := pkgmanagers.ResolveManagerForEntry(ctx.Platform, p.e.PackageManager)
	if err != nil {
		return err
	}

	version := strings.TrimSpace(p.e.Version)
	if strings.EqualFold(version, "latest") || version == "" {
		if err := m.InstallLatest(p.e.Name); err != nil {
			return fmt.Errorf("install latest %s with %s: %w", p.e.Name, m.ID(), err)
		}
	} else {
		if err := m.Install(p.e.Name, version); err != nil {
			return fmt.Errorf("install %s version %s with %s: %w", p.e.Name, version, m.ID(), err)
		}
	}

	for _, s := range p.e.PostInstall {
		ctx.Events.OnPostInstall(p.e, s)

		res, err := steps.RunStep(ctx.Platform, p.e, s)
		if err != nil {
			return err
		}
		if res.Skipped {
			ctx.Events.OnStepSkip(p.e, s, res.SkipReason)
		}
	}

	for _, s := range p.e.Validation {
		ctx.Events.OnValidation(p.e, s)

		res, err := steps.RunStep(ctx.Platform, p.e, s)
		if err != nil {
			return err
		}
		if res.Skipped {
			ctx.Events.OnStepSkip(p.e, s, res.SkipReason)
		}
	}

	return nil
}

func (p *PackageEntry) Uninstall(ctx *Context) error {
	ctx.Events.OnUninstall(p.e)

	m, err := pkgmanagers.ResolveManagerForEntry(ctx.Platform, p.e.PackageManager)
	if err != nil {
		return err
	}

	if err := m.Uninstall(p.e.Name); err != nil {
		return fmt.Errorf("uninstall %s with %s: %w", p.e.Name, m.ID(), err)
	}
	return nil
}
