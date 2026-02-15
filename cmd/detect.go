package cmd

import (
	"context"
	"encoding/json"
	"fmt"
	"os"
	"sort"

	"github.com/Petar-Yordanov/pkg-forge/pkgmanagers"
	"github.com/spf13/cobra"

	"github.com/charmbracelet/lipgloss"
	"github.com/charmbracelet/lipgloss/table"
)

var detectJSON bool

var detectCmd = &cobra.Command{
	Use:   "detect",
	Short: "Detect available package managers",
	RunE: func(cmd *cobra.Command, args []string) error {
		statuses := pkgmanagers.DetectAll(context.Background())
		sort.Slice(statuses, func(i, j int) bool { return statuses[i].Name < statuses[j].Name })

		if detectJSON {
			enc := json.NewEncoder(os.Stdout)
			enc.SetIndent("", "  ")
			return enc.Encode(statuses)
		}

		headerStyle := lipgloss.NewStyle().Bold(true)
		okStyle := lipgloss.NewStyle().Foreground(lipgloss.Color("42"))   // green
		badStyle := lipgloss.NewStyle().Foreground(lipgloss.Color("196")) // red
		muted := lipgloss.NewStyle().Foreground(lipgloss.Color("244"))

		rows := make([][]string, 0, len(statuses))
		for _, s := range statuses {
			status := badStyle.Render("missing")
			path := muted.Render("-")
			ver := muted.Render("-")
			errStr := s.Err

			if s.Available {
				status = okStyle.Render("available")
				if s.Path != "" {
					path = s.Path
				}
				if s.Version != "" {
					ver = s.Version
				} else {
					ver = muted.Render("(unknown)")
				}
				errStr = ""
			}

			platforms := "-"
			if len(s.Platforms) > 0 {
				platforms = fmt.Sprintf("%v", s.Platforms)
			}

			rows = append(rows, []string{
				s.Name,
				s.Cmd,
				platforms,
				status,
				path,
				ver,
				errStr,
			})
		}

		t := table.New().
			Border(lipgloss.RoundedBorder()).
			BorderStyle(lipgloss.NewStyle().Foreground(lipgloss.Color("238"))).
			StyleFunc(func(row, col int) lipgloss.Style {
				if row == 0 {
					return headerStyle
				}

				if col == 6 {
					return lipgloss.NewStyle().Foreground(lipgloss.Color("244"))
				}
				return lipgloss.NewStyle()
			}).
			Headers("Name", "Cmd", "Platforms", "Status", "Path", "Version", "Error").
			Rows(rows...)

		fmt.Println(t.Render())
		return nil
	},
}

func init() {
	rootCmd.AddCommand(detectCmd)
	detectCmd.Flags().BoolVar(&detectJSON, "json", false, "output JSON")
}
