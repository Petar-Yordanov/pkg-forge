package steps

import (
	"fmt"

	"github.com/Petar-Yordanov/pkg-forge/common"
	"github.com/Petar-Yordanov/pkg-forge/manifest/parser"
)

type StepResult struct {
	Skipped    bool
	SkipReason string
}

func RunStepSkeleton(platform common.Platform, e parser.Entry, s parser.Step) (StepResult, error) {
	if s.When != nil && len(s.When.Platform) > 0 {
		ok := false
		for _, p := range s.When.Platform {
			if p == platform {
				ok = true
				break
			}
		}
		if !ok {
			return StepResult{
			  Skipped: true,
			  SkipReason: fmt.Sprintf("platform mismatch (need %v, have %s)", s.When.Platform, platform),
			}, nil
		}
	}

	if s.Cmd != "" && s.CmdFile != "" {
		return StepResult{}, fmt.Errorf("step cannot have both cmd and cmdFile")
	}
	if s.Cmd == "" && s.CmdFile == "" {
		return StepResult{}, fmt.Errorf("step must have either cmd or cmdFile")
	}

	return StepResult{}, nil
}
