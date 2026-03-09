package validate

import (
	"fmt"

	"github.com/Petar-Yordanov/pkg-forge/manifest/parser"
)

type RuleBasicShape struct{}
func (RuleBasicShape) ID() string { return "rule.basicShape" }

func (RuleBasicShape) Check(doc parser.Document) []error {
	var errs []error
	for i, e := range doc.Entries {
		if e.Kind == "" {
			errs = append(errs, fmt.Errorf("entries[%d].kind is required", i))
		}
		if e.Name == "" {
			errs = append(errs, fmt.Errorf("entries[%d].name is required", i))
		}
	}
	return errs
}
