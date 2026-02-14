package traverse

import (
	"fmt"

	"github.com/Petar-Yordanov/pkg-forge/manifest"
)

type Decision int

const (
	Keep Decision = iota
	Skip
	Stop
)

type Rule interface {
	Decide(ref StepRef) (Decision, error)
}

var ErrWalkStopped = fmt.Errorf("walk stopped")

func WalkWithRules(plans []manifest.StepPlan, rules []Rule, fn WalkFn) error {
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
