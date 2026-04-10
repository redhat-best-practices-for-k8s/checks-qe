package main

import (
	"fmt"
	"os"
	"text/tabwriter"

	"github.com/redhat-best-practices-for-k8s/checks-qe/pkg/scenario"
	"github.com/spf13/cobra"
)

func listCmd() *cobra.Command {
	var category string

	cmd := &cobra.Command{
		Use:   "list",
		Short: "List available QE scenarios",
		Run: func(_ *cobra.Command, _ []string) {
			var scenarios []scenario.Scenario
			if category != "" {
				scenarios = scenario.ByCategory(category)
			} else {
				scenarios = scenario.All()
			}

			w := tabwriter.NewWriter(os.Stdout, 0, 4, 2, ' ', 0)
			fmt.Fprintln(w, "SCENARIO\tCHECK\tEXPECTED\tFLAGS")
			for _, s := range scenarios {
				flags := ""
				if s.Privileged {
					flags += "priv "
				}
				if s.RequiresOCP {
					flags += "ocp "
				}
				if s.RequiresProbe {
					flags += "probe "
				}
				fmt.Fprintf(w, "%s\t%s\t%s\t%s\n", s.Name, s.CheckName, s.ExpectedStatus, flags)
			}
			_ = w.Flush()
			fmt.Fprintf(os.Stdout, "\n%d scenario(s)\n", len(scenarios))
		},
	}

	cmd.Flags().StringVar(&category, "category", "", "filter by category")
	return cmd
}
