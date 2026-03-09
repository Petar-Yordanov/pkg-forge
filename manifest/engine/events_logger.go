package engine

import (
	"fmt"
	"strings"

	"github.com/Petar-Yordanov/pkg-forge/common"
	"github.com/Petar-Yordanov/pkg-forge/manifest/parser"
)

type LogEvents struct {
	doc       int
	entryStep int
}

func (l *LogEvents) OnDocStart(i int) {
	l.doc = i
	fmt.Printf("[DOC %d] start\n", i)
}

func (l *LogEvents) OnDocDone(i int) {
	fmt.Printf("[DOC %d] done\n", i)
}

func (l *LogEvents) OnEntryStart(e parser.Entry) {
	l.entryStep = 0
	fmt.Printf("[ENTRY] start kind=%s name=%s%s%s\n",
		e.Kind,
		e.Name,
		fmtOpt(" pm=", e.PackageManager),
		fmtOpt(" version=", e.Version),
	)
}

func (l *LogEvents) OnEntrySkip(e parser.Entry, reason string) {
	fmt.Printf("[SKIP]  kind=%s name=%s reason=%s\n", e.Kind, e.Name, reason)
}

func (l *LogEvents) OnEntryDone(e parser.Entry) {
	fmt.Printf("[ENTRY] done  kind=%s name=%s\n", e.Kind, e.Name)
}

func (l *LogEvents) OnPreInstall(e parser.Entry, s parser.Step)  { l.logStep("preInstall", e, s) }
func (l *LogEvents) OnPostInstall(e parser.Entry, s parser.Step) { l.logStep("postInstall", e, s) }
func (l *LogEvents) OnValidation(e parser.Entry, s parser.Step)  { l.logStep("validation", e, s) }
func (l *LogEvents) OnStep(e parser.Entry, s parser.Step)        { l.logStep("step", e, s) }

func (l *LogEvents) OnInstall(e parser.Entry) {
	fmt.Printf("[DO]    install name=%s%s%s\n",
		e.Name,
		fmtOpt(" pm=", e.PackageManager),
		fmtOpt(" version=", e.Version),
	)
}

func (l *LogEvents) OnError(e parser.Entry, err error) {
	fmt.Printf("[ERR]   kind=%s name=%s error=%v\n", e.Kind, e.Name, err)
}

func (l *LogEvents) logStep(phase string, e parser.Entry, s parser.Step) {
	idx := l.entryStep
	l.entryStep++

	src := stepSource(s)
	prev := stepPreview(s, 90)

	when := ""
	if s.When != nil && len(s.When.Platform) > 0 {
		when = " when=" + strings.Join(platformsToStrings(s.When.Platform), ",")
	}

	fmt.Printf("[STEP]  %s #%d entry=%s src=%s%s cmd=%s\n",
		phase,
		idx,
		e.Name,
		src,
		when,
		prev,
	)
}

func stepSource(s parser.Step) string {
	if s.Cmd != "" {
		return "cmd"
	}
	if s.CmdFile != "" {
		return "cmdFile"
	}
	return "none"
}

func stepPreview(s parser.Step, max int) string {
	var raw string
	if s.Cmd != "" {
		raw = s.Cmd
	} else {
		raw = s.CmdFile
		if s.Shell != "" {
			raw = s.Shell + " " + raw
		}
		if len(s.Args) > 0 {
			raw += " " + strings.Join(s.Args, " ")
		}
	}

	raw = strings.TrimSpace(raw)
	raw = strings.Join(strings.Fields(raw), " ")

	if len(raw) <= max {
		return raw
	}
	return raw[:max-3] + "..."
}

func fmtOpt(prefix, v string) string {
	if strings.TrimSpace(v) == "" {
		return ""
	}
	return prefix + v
}

func platformsToStrings(ps []common.Platform) []string {
	out := make([]string, 0, len(ps))
	for _, p := range ps {
		out = append(out, string(p))
	}
	return out
}

func (l *LogEvents) OnStepSkip(e parser.Entry, s parser.Step, reason string) {
	fmt.Printf("[SKIP]  step entry=%s src=%s reason=%s cmd=%s\n",
		e.Name,
		stepSource(s),
		reason,
		stepPreview(s, 90),
	)
}
