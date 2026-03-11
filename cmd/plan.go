package cmd

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"sort"

	"github.com/Petar-Yordanov/pkg-forge/manifest/engine"
	"github.com/Petar-Yordanov/pkg-forge/manifest/parser"
	validate "github.com/Petar-Yordanov/pkg-forge/manifest/validator"
	"github.com/spf13/cobra"

	"github.com/charmbracelet/lipgloss"
	"github.com/charmbracelet/lipgloss/table"
)

var planJSON bool

var planCmd = &cobra.Command{
	Use:   "plan <file.yml>",
	Short: "Parse, validate, and show the execution plan without applying it",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		path := args[0]

		docs, err := parser.ParseFile(path)
		if err != nil {
			return fmt.Errorf("parse %s: %w", path, err)
		}

		v := validate.New(
			validate.RuleBasicShape{},
		)

		for i, d := range docs {
			if err := v.Validate(d); err != nil {
				return fmt.Errorf("validate doc %d: %w", i, err)
			}
		}

		statePath := filepath.Join(filepath.Dir(path), ".pkg-forge-state.sqlite")
		store, err := engine.OpenStateStore(statePath)
		if err != nil {
			return fmt.Errorf("open state store %s: %w", statePath, err)
		}
		defer store.Close()

		ctx := engine.NewDefaultContext(&engine.LogEvents{})
		ctx.ManifestPath = path
		ctx.ManifestName = filepath.Base(path)
		ctx.State = store

		r := engine.NewRunner()

		plan, err := r.BuildPlan(ctx, docs)
		if err != nil {
			return err
		}

		if planJSON {
			enc := json.NewEncoder(os.Stdout)
			enc.SetIndent("", "  ")
			return enc.Encode(plan)
		}

		renderPlanTable(plan)
		return nil
	},
}

func init() {
	rootCmd.AddCommand(planCmd)
	planCmd.Flags().BoolVar(&planJSON, "json", false, "output JSON")
}

func renderPlanTable(plan *engine.Plan) {
	if plan == nil {
		return
	}

	headerStyle := lipgloss.NewStyle().Bold(true)
	okStyle := lipgloss.NewStyle().Foreground(lipgloss.Color("42"))
	changeStyle := lipgloss.NewStyle().Foreground(lipgloss.Color("214"))
	badStyle := lipgloss.NewStyle().Foreground(lipgloss.Color("196"))
	muted := lipgloss.NewStyle().Foreground(lipgloss.Color("244"))

	items := append([]engine.PlanItem(nil), plan.Items...)
	sort.SliceStable(items, func(i, j int) bool { return items[i].Ordinal < items[j].Ordinal })

	rows := make([][]string, 0, len(items))
	for _, item := range items {
		action := string(item.Action)
		status := ""
		reason := item.Reason
		if reason == "" {
			reason = "-"
		}

		switch item.Action {
		case engine.PlanActionCreate:
			action = okStyle.Render(action)
			status = okStyle.Render("planned")
		case engine.PlanActionReplace:
			action = changeStyle.Render(action)
			status = changeStyle.Render("planned")
		case engine.PlanActionDestroy:
			action = badStyle.Render(action)
			status = badStyle.Render("planned")
		case engine.PlanActionNoOp:
			action = muted.Render(action)
			status = muted.Render("unchanged")
		case engine.PlanActionSkip:
			action = muted.Render(action)
			status = muted.Render("skipped")
		default:
			status = muted.Render("-")
		}

		rows = append(rows, []string{
			fmt.Sprintf("%d", item.Ordinal),
			action,
			item.Kind,
			item.Name,
			emptyDash(item.PackageManager),
			emptyDash(item.PackageVersion),
			status,
			reason,
		})
	}

	t := table.New().
		Border(lipgloss.RoundedBorder()).
		BorderStyle(lipgloss.NewStyle().Foreground(lipgloss.Color("238"))).
		StyleFunc(func(row, col int) lipgloss.Style {
			if row == 0 {
				return headerStyle
			}
			if col == 7 {
				return lipgloss.NewStyle().Foreground(lipgloss.Color("244"))
			}
			return lipgloss.NewStyle()
		}).
		Headers("Order", "Action", "Kind", "Name", "Manager", "PackageVersion", "Status", "Reason").
		Rows(rows...)

	fmt.Printf("[PLAN] manifest=%s\n", plan.ManifestName)
	fmt.Println(t.Render())

	fmt.Printf("\nSummary:\n")
	fmt.Printf("  create: %d\n", plan.Summary.Create)
	fmt.Printf("  replace: %d\n", plan.Summary.Replace)
	fmt.Printf("  destroy: %d\n", plan.Summary.Destroy)
	fmt.Printf("  no-op: %d\n", plan.Summary.NoOp)
	fmt.Printf("  skip: %d\n", plan.Summary.Skip)
}

func emptyDash(v string) string {
	if v == "" {
		return "-"
	}
	return v
}
