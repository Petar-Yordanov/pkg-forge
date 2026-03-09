package cmd

import (
	"fmt"

	"github.com/Petar-Yordanov/pkg-forge/manifest/engine"
	"github.com/Petar-Yordanov/pkg-forge/manifest/parser"
	"github.com/Petar-Yordanov/pkg-forge/manifest/validator"
	"github.com/spf13/cobra"
)

var applyCmd = &cobra.Command{
	Use:   "apply <file.yml>",
	Short: "Parse a YAML file and optionally apply it",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		path := args[0]
		doApply, _ := cmd.Flags().GetBool("apply")

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

		_ = doApply

		ctx := engine.NewDefaultContext(&engine.LogEvents{})
		r := engine.NewRunner()
		return r.RunDocs(ctx, docs)
	},
}

func init() {
	applyCmd.Flags().Bool("apply", false, "Apply the entries (skeleton for now)")
	rootCmd.AddCommand(applyCmd)
}
