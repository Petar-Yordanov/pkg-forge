package cmd

import (
	"fmt"
	"path/filepath"

	"github.com/Petar-Yordanov/pkg-forge/manifest/engine"
	"github.com/Petar-Yordanov/pkg-forge/manifest/parser"
	validate "github.com/Petar-Yordanov/pkg-forge/manifest/validator"
	"github.com/spf13/cobra"
)

var applyCmd = &cobra.Command{
	Use:   "apply <file.yml>",
	Short: "Parse, validate, plan, and apply a manifest",
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

		r.PrintPlan(plan)

		return r.ApplyPlan(ctx, plan)
	},
}

func init() {
	rootCmd.AddCommand(applyCmd)
}
