package traverse

import (
	"fmt"
	"github.com/Petar-Yordanov/pkg-forge/manifest"
)

type WalkFn func(ref StepRef) error

func WalkPlans(plans []manifest.StepPlan, fn WalkFn) error {
	for ei := range plans {
		p := plans[ei]

		emitList := func(phase Phase, list []manifest.ResolvedStep) error {
			for i := range list {
				rs := &list[i]
				if err := fn(StepRef{
					EntryIndex: ei,
					Kind:       p.Kind,
					Name:       p.Name,
					Version:    p.Version,
					Phase:      phase,
					Index:      i,
					Step:       rs,
				}); err != nil {
					return err
				}
			}
			return nil
		}

		if err := emitList(PhasePreInstall, p.PreInstall); err != nil {
			return fmt.Errorf("%s/%s %s: %w", p.Kind, p.Name, PhasePreInstall, err)
		}
		if err := emitList(PhasePostInstall, p.PostInstall); err != nil {
			return fmt.Errorf("%s/%s %s: %w", p.Kind, p.Name, PhasePostInstall, err)
		}
		if err := emitList(PhaseValidation, p.Validation); err != nil {
			return fmt.Errorf("%s/%s %s: %w", p.Kind, p.Name, PhaseValidation, err)
		}
		if err := emitList(PhasePre, p.Pre); err != nil {
			return fmt.Errorf("%s/%s %s: %w", p.Kind, p.Name, PhasePre, err)
		}
		if err := emitList(PhaseSteps, p.Steps); err != nil {
			return fmt.Errorf("%s/%s %s: %w", p.Kind, p.Name, PhaseSteps, err)
		}
		if err := emitList(PhasePost, p.Post); err != nil {
			return fmt.Errorf("%s/%s %s: %w", p.Kind, p.Name, PhasePost, err)
		}

		if p.Cmd != nil {
			rs := p.Cmd
			if err := fn(StepRef{
				EntryIndex: ei,
				Kind:       p.Kind,
				Name:       p.Name,
				Version:    p.Version,
				Phase:      PhaseCmd,
				Index:      -1,
				Step:       rs,
			}); err != nil {
				return fmt.Errorf("%s/%s %s: %w", p.Kind, p.Name, PhaseCmd, err)
			}
		}
	}
	return nil
}
