package validate

import (
	"github.com/Petar-Yordanov/pkg-forge/manifest"
	"github.com/Petar-Yordanov/pkg-forge/manifest/traverse"
)

type Validator interface {
	ValidateStep(ref traverse.StepRef) error
	// ValidateEntry or ValidateDocument eventually
}

func ValidatePlans(plans []manifest.StepPlan, validators ...Validator) error {
	return traverse.WalkPlans(plans, func(ref traverse.StepRef) error {
		for _, v := range validators {
			if err := v.ValidateStep(ref); err != nil {
				return err
			}
		}
		return nil
	})
}
