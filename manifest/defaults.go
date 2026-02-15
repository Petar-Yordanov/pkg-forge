package manifest

import "maps"

type effDefaults struct {
	FailOnStderr  bool
	Env           map[string]string
	Retries       int
	RetryDelaySec int
	TimeoutSec    *int
}

func effectiveDefaults(e *Entry) effDefaults {
	d := effDefaults{
		FailOnStderr:  false,
		Env:           map[string]string{},
		Retries:       0,
		RetryDelaySec: 0,
		TimeoutSec:    nil,
	}

	if e.FailOnStderr != nil {
		d.FailOnStderr = *e.FailOnStderr
	}
	if e.Defaults == nil || e.Defaults.Step == nil {
		return d
	}

	sd := e.Defaults.Step
	if sd.FailOnStderr != nil {
		d.FailOnStderr = *sd.FailOnStderr
	}
	if sd.Env != nil {
		d.Env = maps.Clone(sd.Env)
	}
	if sd.Retries != nil {
		d.Retries = *sd.Retries
	}
	if sd.RetryDelaySec != nil {
		d.RetryDelaySec = *sd.RetryDelaySec
	}
	if sd.TimeoutSec != nil {
		d.TimeoutSec = sd.TimeoutSec
	}
	return d
}
