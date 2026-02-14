package traverse

type Phase string

const (
	PhasePreInstall  Phase = "preInstall"
	PhasePostInstall Phase = "postInstall"
	PhaseValidation  Phase = "validation"
	PhasePre         Phase = "pre"
	PhasePost        Phase = "post"
	PhaseSteps       Phase = "steps"
	PhaseCmd         Phase = "cmd"
)

type StepRef struct {
	EntryIndex int
	Kind       string
	Name       string
	Version    string

	Phase Phase
	Index int
	Step  *ResolvedStep
}
