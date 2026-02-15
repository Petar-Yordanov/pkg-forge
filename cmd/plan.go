package cmd

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/Petar-Yordanov/pkg-forge/manifest"
	"github.com/Petar-Yordanov/pkg-forge/manifest/parse"
	"github.com/spf13/cobra"
	"github.com/Petar-Yordanov/pkg-forge/common"
)

var planCmd = &cobra.Command{
	Use:   "plan",
	Short: "Parse config and print normalized plan",
	RunE: func(cmd *cobra.Command, args []string) error {
		cfg, _ := cmd.Flags().GetString("config")
		if cfg == "" && len(args) > 0 {
			cfg = args[0]
		}
		if cfg == "" {
			return fmt.Errorf("missing config path")
		}

		doc, err := parse.ParseFile(cfg)
		if err != nil {
			return err
		}

		baseDir := filepath.Dir(cfg)
		platform := common.CurrentPlatform()

		plans, err := manifest.Normalize(doc, baseDir, platform)
		if err != nil {
			return err
		}

		s, err := common.ToJSON(plans)
		if err != nil {
			return err
		}

		_, _ = fmt.Fprintln(os.Stdout, s)
		return nil
	},
}

func init() {
	planCmd.Flags().StringP("config", "c", "", "path to yaml config")
	rootCmd.AddCommand(planCmd)
}
