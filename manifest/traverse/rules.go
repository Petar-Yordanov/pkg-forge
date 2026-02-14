package traverse

type Decision int

const (
	Keep Decision = iota
	Skip
	Stop
)

type Rule interface {
	Decide(ref StepRef) (Decision, error)
}

func WalkWithRules(plans []StepPlan, rules []Rule, fn WalkFn) error {
	return WalkPlans(plans, func(ref StepRef) error {
		for _, r := range rules {
			dec, err := r.Decide(ref)
			if err != nil {
				return err
			}
			switch dec {
			case Skip:
				return nil
			case Stop:
				return ErrWalkStopped
			}
		}
		return fn(ref)
	})
}

var ErrWalkStopped = fmt.Errorf("walk stopped")
