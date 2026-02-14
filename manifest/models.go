package manifest

import (
	"time"
	"github.com/Petar-Yordanov/pkg-forge/common"
)

type Document struct {
	SchemaVersion int     `yaml:"schemaVersion" json:"schemaVersion"`
	Entries       []Entry `yaml:"entries" json:"entries"`
}

type Entry struct {
	Kind        string        `yaml:"kind" json:"kind"`
	Name        string        `yaml:"name" json:"name"`
	Version     string        `yaml:"version,omitempty" json:"version,omitempty"`
	FailOnStderr *bool        `yaml:"failOnStderr,omitempty" json:"failOnStderr,omitempty"`
	Defaults    *EntryDefaults `yaml:"defaults,omitempty" json:"defaults,omitempty"`

	PreInstall  []StepLike `yaml:"preInstall,omitempty" json:"preInstall,omitempty"`
	PostInstall []StepLike `yaml:"postInstall,omitempty" json:"postInstall,omitempty"`
	Validation  []StepLike `yaml:"validation,omitempty" json:"validation,omitempty"`

	Pre   []StepLike `yaml:"pre,omitempty" json:"pre,omitempty"`
	Post  []StepLike `yaml:"post,omitempty" json:"post,omitempty"`
	Steps []StepLike `yaml:"steps,omitempty" json:"steps,omitempty"`
	Cmd   *StepLike  `yaml:"cmd,omitempty" json:"cmd,omitempty"`
}

type EntryDefaults struct {
	Step *StepDefaults `yaml:"step,omitempty" json:"step,omitempty"`
}

type StepDefaults struct {
	FailOnStderr  *bool             `yaml:"failOnStderr,omitempty" json:"failOnStderr,omitempty"`
	Env           map[string]string `yaml:"env,omitempty" json:"env,omitempty"`
	Retries       *int              `yaml:"retries,omitempty" json:"retries,omitempty"`
	RetryDelaySec *int              `yaml:"retryDelaySec,omitempty" json:"retryDelaySec,omitempty"`
	TimeoutSec    *int              `yaml:"timeoutSec,omitempty" json:"timeoutSec,omitempty"`
}

type When struct {
	Platform common.Platform `yaml:"platform,omitempty" json:"platform,omitempty"`
}

type Step struct {
	Cmd     string `yaml:"cmd,omitempty" json:"cmd,omitempty"`
	CmdFile string `yaml:"cmdFile,omitempty" json:"cmdFile,omitempty"`

	When *When `yaml:"when,omitempty" json:"when,omitempty"`

	Shell string   `yaml:"shell,omitempty" json:"shell,omitempty"`
	Args  []string `yaml:"args,omitempty" json:"args,omitempty"`
	Cwd   string   `yaml:"cwd,omitempty" json:"cwd,omitempty"`

	Env map[string]string `yaml:"env,omitempty" json:"env,omitempty"`

	FailOnStderr  *bool `yaml:"failOnStderr,omitempty" json:"failOnStderr,omitempty"`
	Retries       *int  `yaml:"retries,omitempty" json:"retries,omitempty"`
	RetryDelaySec *int  `yaml:"retryDelaySec,omitempty" json:"retryDelaySec,omitempty"`
	TimeoutSec    *int  `yaml:"timeoutSec,omitempty" json:"timeoutSec,omitempty"`
}

type SelectNode struct {
	Default *Step            `yaml:"default,omitempty" json:"default,omitempty"`
	Items   map[string]*Step `yaml:",inline" json:"-"`
}

type StepLike struct {
	Step   *Step       `json:"step,omitempty"`
	Select *SelectNode `json:"select,omitempty"`
}

type StepPlan struct {
	Kind string `json:"kind"`
	Name string `json:"name"`

	Version string `json:"version,omitempty"`

	PreInstall  []ResolvedStep `json:"preInstall,omitempty"`
	PostInstall []ResolvedStep `json:"postInstall,omitempty"`
	Validation  []ResolvedStep `json:"validation,omitempty"`

	Pre   []ResolvedStep `json:"pre,omitempty"`
	Post  []ResolvedStep `json:"post,omitempty"`
	Steps []ResolvedStep `json:"steps,omitempty"`
	Cmd   *ResolvedStep  `json:"cmd,omitempty"`
}

type ResolvedStep struct {
	ExecKind string `json:"execKind"`
	Cmd      string `json:"cmd,omitempty"`
	CmdFile  string `json:"cmdFile,omitempty"`

	When *When `json:"when,omitempty"`

	Shell string   `json:"shell,omitempty"`
	Args  []string `json:"args,omitempty"`
	Cwd   string   `json:"cwd,omitempty"`

	Env map[string]string `json:"env,omitempty"`

	FailOnStderr bool `json:"failOnStderr"`

	Retries       int           `json:"retries"`
	RetryDelay    time.Duration `json:"retryDelay"`
	Timeout       time.Duration `json:"timeout"`
	TimeoutInfinite bool        `json:"timeoutInfinite"`
}
