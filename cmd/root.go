package cmd

import (
	"os"

	"github.com/spf13/cobra"
)

var (
	cfgPath string
	tagsCSV string
)

var rootCmd = &cobra.Command{
	Use:   "pkg-forge",
	Short: "pkg-forge",
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		os.Exit(1)
	}
}

func init() {
	rootCmd.PersistentFlags().StringVarP(&cfgPath, "config", "c", "", "path to YAML config")
	rootCmd.PersistentFlags().StringVarP(&tagsCSV, "tags", "t", "", "comma-separated tags to include (optional)")
}
