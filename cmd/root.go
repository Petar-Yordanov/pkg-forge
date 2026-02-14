package cmd

import (
	"os"

	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:   "pkg-forge",
	Short: "pkg-forge",
}

func Execute() {
	_ = rootCmd.Execute()
	os.Exit(0)
}
