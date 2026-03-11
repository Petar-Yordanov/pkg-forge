package engine

import (
	"fmt"

	"github.com/Petar-Yordanov/pkg-forge/common"
	"github.com/Petar-Yordanov/pkg-forge/manifest/parser"
)

type Runner struct{}

type plannedEntry struct {
	SourceOrdinal int
	Entry         Entry
	Row           StateRow
}

type PlanAction string

const (
	PlanActionNoOp    PlanAction = "no-op"
	PlanActionCreate  PlanAction = "create"
	PlanActionReplace PlanAction = "replace"
	PlanActionDestroy PlanAction = "destroy"
	PlanActionSkip    PlanAction = "skip"
)

type PlanItem struct {
	Ordinal        int        `json:"ordinal"`
	Action         PlanAction `json:"action"`
	Reason         string     `json:"reason,omitempty"`
	Kind           string     `json:"kind"`
	Name           string     `json:"name"`
	PackageManager string     `json:"packageManager,omitempty"`
	PackageVersion string     `json:"packageVersion,omitempty"`

	Desired *StateRow `json:"desired,omitempty"`
	Current *StateRow `json:"current,omitempty"`

	Entry Entry `json:"-"`
}

type PlanSummary struct {
	Create  int `json:"create"`
	Replace int `json:"replace"`
	Destroy int `json:"destroy"`
	NoOp    int `json:"noOp"`
	Skip    int `json:"skip"`
}

type Plan struct {
	ManifestName string      `json:"manifestName"`
	Items        []PlanItem  `json:"items"`
	Summary      PlanSummary `json:"summary"`
}

func NewRunner() *Runner { return &Runner{} }

func (r *Runner) RunDocs(ctx *Context, docs []parser.Document) error {
	plan, err := r.BuildPlan(ctx, docs)
	if err != nil {
		return err
	}
	r.PrintPlan(plan)
	return r.ApplyPlan(ctx, plan)
}

func (r *Runner) BuildPlan(ctx *Context, docs []parser.Document) (*Plan, error) {
	if ctx == nil {
		return nil, fmt.Errorf("nil context")
	}
	if ctx.State == nil {
		return nil, fmt.Errorf("nil state store")
	}
	if ctx.ManifestName == "" {
		return nil, fmt.Errorf("empty manifest name")
	}

	allEntries, desiredEntries, err := r.buildDesiredEntries(ctx, docs)
	if err != nil {
		return nil, err
	}

	currentRows, err := ctx.State.LoadManifest(ctx.ManifestName)
	if err != nil {
		return nil, fmt.Errorf("load state: %w", err)
	}

	itemsByOrdinal := make(map[int][]PlanItem, len(allEntries)+len(currentRows))

	for _, pe := range allEntries {
		ok, reason := pe.Entry.Applies(ctx.Platform)
		if ok {
			continue
		}

		raw := pe.Entry.Raw()
		item := PlanItem{
			Ordinal:        pe.SourceOrdinal,
			Action:         PlanActionSkip,
			Reason:         reason,
			Kind:           raw.Kind,
			Name:           raw.Name,
			PackageManager: raw.PackageManager,
			PackageVersion: raw.Version,
			Entry:          pe.Entry,
		}
		itemsByOrdinal[item.Ordinal] = append(itemsByOrdinal[item.Ordinal], item)
	}

	maxLen := len(currentRows)
	if len(desiredEntries) > maxLen {
		maxLen = len(desiredEntries)
	}

	for i := 0; i < maxLen; i++ {
		switch {
		case i >= len(currentRows):
			rowCopy := desiredEntries[i].Row
			raw := desiredEntries[i].Entry.Raw()
			item := PlanItem{
				Ordinal:        desiredEntries[i].SourceOrdinal,
				Action:         PlanActionCreate,
				Kind:           raw.Kind,
				Name:           raw.Name,
				PackageManager: raw.PackageManager,
				PackageVersion: raw.Version,
				Desired:        &rowCopy,
				Entry:          desiredEntries[i].Entry,
			}
			itemsByOrdinal[item.Ordinal] = append(itemsByOrdinal[item.Ordinal], item)

		case i >= len(desiredEntries):
			rowCopy := currentRows[i]
			item := PlanItem{
				Ordinal:        rowCopy.Ordinal,
				Action:         PlanActionDestroy,
				Kind:           rowCopy.EntryKind,
				Name:           rowCopy.PackageName,
				PackageManager: rowCopy.PackageManager,
				PackageVersion: rowCopy.Version,
				Current:        &rowCopy,
			}
			itemsByOrdinal[item.Ordinal] = append(itemsByOrdinal[item.Ordinal], item)

		default:
			curCopy := currentRows[i]
			newCopy := desiredEntries[i].Row
			raw := desiredEntries[i].Entry.Raw()

			action := PlanActionReplace
			if stateRowsEqual(currentRows[i], desiredEntries[i].Row) {
				action = PlanActionNoOp
			}

			item := PlanItem{
				Ordinal:        desiredEntries[i].SourceOrdinal,
				Action:         action,
				Kind:           raw.Kind,
				Name:           raw.Name,
				PackageManager: raw.PackageManager,
				PackageVersion: raw.Version,
				Current:        &curCopy,
				Desired:        &newCopy,
				Entry:          desiredEntries[i].Entry,
			}
			itemsByOrdinal[item.Ordinal] = append(itemsByOrdinal[item.Ordinal], item)
		}
	}

	items := make([]PlanItem, 0, len(allEntries)+len(currentRows))
	summary := PlanSummary{}

	for ordinal := 0; ordinal <= maxOrdinal(itemsByOrdinal); ordinal++ {
		group := itemsByOrdinal[ordinal]
		for _, item := range group {
			items = append(items, item)
			switch item.Action {
			case PlanActionCreate:
				summary.Create++
			case PlanActionReplace:
				summary.Replace++
			case PlanActionDestroy:
				summary.Destroy++
			case PlanActionNoOp:
				summary.NoOp++
			case PlanActionSkip:
				summary.Skip++
			}
		}
	}

	return &Plan{
		ManifestName: ctx.ManifestName,
		Items:        items,
		Summary:      summary,
	}, nil
}

func (r *Runner) ApplyPlan(ctx *Context, plan *Plan) error {
	if ctx == nil {
		return fmt.Errorf("nil context")
	}
	if ctx.State == nil {
		return fmt.Errorf("nil state store")
	}
	if plan == nil {
		return fmt.Errorf("nil plan")
	}

	for i := len(plan.Items) - 1; i >= 0; i-- {
		item := plan.Items[i]
		switch item.Action {
		case PlanActionReplace, PlanActionDestroy:
			if item.Current == nil {
				continue
			}
			if err := r.uninstallStateRow(ctx, *item.Current); err != nil {
				return err
			}
		}
	}

	finalRows := make([]StateRow, 0, len(plan.Items))

	for _, item := range plan.Items {
		switch item.Action {
		case PlanActionDestroy, PlanActionSkip:
			continue

		case PlanActionNoOp:
			if item.Desired != nil {
				finalRows = append(finalRows, *item.Desired)
			}
			continue

		case PlanActionCreate, PlanActionReplace:
			if item.Entry == nil {
				return fmt.Errorf("plan item %d has no entry", item.Ordinal)
			}

			raw := item.Entry.Raw()
			ctx.Events.OnEntryStart(raw)

			ok, reason := item.Entry.Applies(ctx.Platform)
			if !ok {
				ctx.Events.OnEntrySkip(raw, reason)
				ctx.Events.OnEntryDone(raw)
				continue
			}

			if err := item.Entry.Run(ctx); err != nil {
				ctx.Events.OnError(raw, err)
				return err
			}

			ctx.Events.OnEntryDone(raw)

			if item.Desired != nil {
				finalRows = append(finalRows, *item.Desired)
			}
		}
	}

	if err := ctx.State.ReplaceManifest(ctx.ManifestName, finalRows); err != nil {
		return fmt.Errorf("save state: %w", err)
	}

	return nil
}

func (r *Runner) PrintPlan(plan *Plan) {
	if plan == nil {
		return
	}

	fmt.Printf("[PLAN] manifest=%s\n", plan.ManifestName)
	for _, item := range plan.Items {
		switch item.Action {
		case PlanActionCreate:
			fmt.Printf("  + create   %s\n", formatPlanItem(item))
		case PlanActionReplace:
			fmt.Printf("  ~ replace  %s\n", formatReplace(item.Current, item.Desired))
		case PlanActionDestroy:
			fmt.Printf("  - destroy  %s\n", formatPlanItem(item))
		case PlanActionNoOp:
			fmt.Printf("  = no-op    %s\n", formatPlanItem(item))
		case PlanActionSkip:
			fmt.Printf("  · skip     %s reason=%s\n", formatPlanItem(item), item.Reason)
		}
	}

	fmt.Printf("\nSummary:\n")
	fmt.Printf("  create: %d\n", plan.Summary.Create)
	fmt.Printf("  replace: %d\n", plan.Summary.Replace)
	fmt.Printf("  destroy: %d\n", plan.Summary.Destroy)
	fmt.Printf("  no-op: %d\n", plan.Summary.NoOp)
	fmt.Printf("  skip: %d\n", plan.Summary.Skip)
}

func (r *Runner) buildDesiredEntries(ctx *Context, docs []parser.Document) ([]plannedEntry, []plannedEntry, error) {
	all := make([]plannedEntry, 0)
	desired := make([]plannedEntry, 0)

	manifestOrdinal := 0
	desiredOrdinal := 0

	for _, d := range docs {
		list, err := BuildEntries(d)
		if err != nil {
			return nil, nil, err
		}

		for _, ent := range list {
			row, err := BuildStateRow(ctx.ManifestName, desiredOrdinal, ent.Raw())
			if err != nil {
				return nil, nil, err
			}

			pe := plannedEntry{
				SourceOrdinal: manifestOrdinal,
				Entry:         ent,
				Row:           row,
			}
			all = append(all, pe)

			ok, _ := ent.Applies(ctx.Platform)
			if ok {
				desired = append(desired, pe)
				desiredOrdinal++
			}

			manifestOrdinal++
		}
	}

	return all, desired, nil
}

func (r *Runner) uninstallStateRow(ctx *Context, row StateRow) error {
	if row.EntryKind != "package" {
		return nil
	}

	raw := parser.Entry{
		Kind:           row.EntryKind,
		Name:           row.PackageName,
		Version:        row.Version,
		PackageManager: row.PackageManager,
	}

	ent := NewPackageEntry(raw)
	if err := ent.Uninstall(ctx); err != nil {
		ctx.Events.OnError(raw, err)
		return err
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

func stateRowsEqual(a, b StateRow) bool {
	return a.ManifestName == b.ManifestName &&
		a.Ordinal == b.Ordinal &&
		a.EntryKind == b.EntryKind &&
		a.PackageName == b.PackageName &&
		a.Version == b.Version &&
		a.PackageManager == b.PackageManager &&
		a.StepHash == b.StepHash
}

func formatPlanItem(item PlanItem) string {
	switch item.Kind {
	case "package":
		return fmt.Sprintf("package %s (manager=%s packageVersion=%s)", item.Name, emptyAsDash(item.PackageManager), emptyAsDash(item.PackageVersion))
	case "script":
		return fmt.Sprintf("script %s", item.Name)
	default:
		return fmt.Sprintf("%s %s", item.Kind, item.Name)
	}
}

func formatReplace(current, desired *StateRow) string {
	if current == nil {
		return formatPlanRow(desired)
	}
	if desired == nil {
		return formatPlanRow(current)
	}

	if current.EntryKind == "package" && desired.EntryKind == "package" && current.PackageName == desired.PackageName {
		return fmt.Sprintf(
			"package %s (manager=%s packageVersion=%s -> manager=%s packageVersion=%s)",
			desired.PackageName,
			emptyAsDash(current.PackageManager),
			emptyAsDash(current.Version),
			emptyAsDash(desired.PackageManager),
			emptyAsDash(desired.Version),
		)
	}

	return fmt.Sprintf("%s -> %s", formatPlanRow(current), formatPlanRow(desired))
}

func formatPlanRow(row *StateRow) string {
	if row == nil {
		return "<nil>"
	}

	switch row.EntryKind {
	case "package":
		return fmt.Sprintf("package %s (manager=%s packageVersion=%s)", row.PackageName, emptyAsDash(row.PackageManager), emptyAsDash(row.Version))
	case "script":
		return fmt.Sprintf("script %s", row.PackageName)
	default:
		return fmt.Sprintf("%s %s", row.EntryKind, row.PackageName)
	}
}

func emptyAsDash(v string) string {
	if v == "" {
		return "-"
	}
	return v
}

func maxOrdinal(m map[int][]PlanItem) int {
	max := -1
	for k := range m {
		if k > max {
			max = k
		}
	}
	return max
}
