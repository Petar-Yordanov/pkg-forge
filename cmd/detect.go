package cmd

import (
	"context"
	"encoding/json"
	"fmt"
	"os"
	"sort"

	detect "github.com/Petar-Yordanov/pkg-forge/pkgmanagers"
	"github.com/spf13/cobra"
)

var detectJSON bool

var detectCmd = &cobra.Command{
	Use:   "detect",
	Short: "Detect available package managers",
	RunE: func(cmd *cobra.Command, args []string) error {
		statuses := detect.DetectAll(context.Background())
		sort.Slice(statuses, func(i, j int) bool { return statuses[i].Name < statuses[j].Name })

		if detectJSON {
			enc := json.NewEncoder(os.Stdout)
			enc.SetIndent("", "  ")
			return enc.Encode(statuses)
		}

		for _, s := range statuses {
			if s.Available {
				v := s.Version
				if v == "" {
					v = "(unknown)"
				}
				fmt.Printf("%-8s  available  path=%s  version=%s\n", s.Name, s.Path, v)
			} else {
				fmt.Printf("%-8s  missing    %s\n", s.Name, s.Err)
			}
		}
		return nil
	},
}

func init() {
	rootCmd.AddCommand(detectCmd)
	detectCmd.Flags().BoolVar(&detectJSON, "json", false, "output JSON")
}
