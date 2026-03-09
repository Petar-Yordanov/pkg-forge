package validate

import (
	"fmt"

	"github.com/Petar-Yordanov/pkg-forge/manifest/parser"
)

type Rule interface {
	ID() string
	Check(doc parser.Document) []error
}

type Validator struct {
	rules []Rule
}

func New(rules ...Rule) *Validator {
	return &Validator{rules: rules}
}

func (v *Validator) Add(r Rule) { v.rules = append(v.rules, r) }

func (v *Validator) Validate(doc parser.Document) error {
	var errs []error
	for _, r := range v.rules {
		errs = append(errs, r.Check(doc)...)
	}
	if len(errs) == 0 {
		return nil
	}
	return fmt.Errorf("%d validation error(s): %w", len(errs), join(errs))
}

func join(errs []error) error {
	if len(errs) == 0 {
		return nil
	}
	out := errs[0]
	for i := 1; i < len(errs); i++ {
		out = fmt.Errorf("%w; %v", out, errs[i])
	}
	return out
}
