package engine

import (
	"fmt"

	"github.com/Petar-Yordanov/pkg-forge/common"
	"github.com/Petar-Yordanov/pkg-forge/manifest/parser"
)

type Runner struct{}

func NewRunner() *Runner { return &Runner{} }

func (r *Runner) RunDocs(ctx *Context, docs []parser.Document) error {
	for di, d := range docs {
		ctx.Events.OnDocStart(di)

		list, err := BuildEntries(d)
		if err != nil {
			return err
		}

		for _, ent := range list {
			raw := ent.Raw()
			ctx.Events.OnEntryStart(raw)

			ok, reason := ent.Applies(ctx.Platform)
			if !ok {
				ctx.Events.OnEntrySkip(raw, reason)
				ctx.Events.OnEntryDone(raw)
				continue
			}

			if err := ent.Run(ctx); err != nil {
				ctx.Events.OnError(raw, err)
				return err
			}

			ctx.Events.OnEntryDone(raw)
		}

		ctx.Events.OnDocDone(di)
	}
	return nil
}

func BuildEntries(doc parser.Document) ([]Entry, error) {
	out := make([]Entry, 0, len(doc.Entries))
	for i, e := range doc.Entries {
		switch e.Kind {
		case "package":
			out = append(out, NewPackageEntry(e))
		case "script":
			out = append(out, NewScriptEntry(e))
		default:
			return nil, fmt.Errorf("entries[%d].kind: unsupported %q", i, e.Kind)
		}
	}
	return out, nil
}

func NewDefaultContext(ev Events) *Context {
	return &Context{
		Platform: common.CurrentPlatform(),
		Events:   ev,
	}
}
