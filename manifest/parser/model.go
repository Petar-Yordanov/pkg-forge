package parser

import "github.com/Petar-Yordanov/pkg-forge/common"

type Document struct {
	Entries []Entry `yaml:"entries"`
}

type When struct {
	Platform []common.Platform `yaml:"platform"`
}

type Entry struct {
	Kind string `yaml:"kind"`
	Name string `yaml:"name"`

	Version        string `yaml:"version,omitempty"`
	PackageManager string `yaml:"packageManager,omitempty"`

	Steps []Step `yaml:"steps,omitempty"`

	PreInstall  []Step `yaml:"preInstall,omitempty"`
	PostInstall []Step `yaml:"postInstall,omitempty"`
	Validation  []Step `yaml:"validation,omitempty"`

	When         *When             `yaml:"when,omitempty"`
	FailOnStderr *bool             `yaml:"failOnStderr,omitempty"`
	Env          map[string]string `yaml:"env,omitempty"`
}

type Step struct {
	Cmd     string   `yaml:"cmd,omitempty"`
	CmdFile string   `yaml:"cmdFile,omitempty"`
	Shell   string   `yaml:"shell,omitempty"`
	Args    []string `yaml:"args,omitempty"`

	When *When `yaml:"when,omitempty"`

	FailOnStderr  *bool             `yaml:"failOnStderr,omitempty"`
	Env           map[string]string `yaml:"env,omitempty"`
	TimeoutSec    int               `yaml:"timeoutSec,omitempty"`
	Retries       int               `yaml:"retries,omitempty"`
	RetryDelaySec int               `yaml:"retryDelaySec,omitempty"`
}
