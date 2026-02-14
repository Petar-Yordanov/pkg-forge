package validate

type Validator interface {
	ValidateStep(ref StepRef) error
	// ValidateEntry or ValidateDocument eventually
}

func ValidatePlans(plans []StepPlan, validators ...Validator) error {
	return WalkPlans(plans, func(ref StepRef) error {
		for _, v := range validators {
			if err := v.ValidateStep(ref); err != nil {
				return err
			}
		}
		return nil
	})
}
