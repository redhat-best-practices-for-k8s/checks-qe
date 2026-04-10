package main

import (
	"fmt"
	"os"

	"github.com/redhat-best-practices-for-k8s/checks-qe/pkg/cluster"
	"github.com/redhat-best-practices-for-k8s/checks-qe/pkg/runner"
	"github.com/redhat-best-practices-for-k8s/checks-qe/pkg/scenario"
	"github.com/spf13/cobra"
)

func runCmd() *cobra.Command {
	var (
		kubeconfig   string
		category     string
		scenarioName string
		parallel     int
		verbose      bool
		skipProbe    bool
		skipOCP      bool
	)

	cmd := &cobra.Command{
		Use:   "run",
		Short: "Run QE scenarios against a live cluster",
		RunE: func(_ *cobra.Command, _ []string) error {
			cl, err := cluster.NewClient(kubeconfig)
			if err != nil {
				return fmt.Errorf("connecting to cluster: %w", err)
			}

			clusterType := "Kubernetes"
			if cl.IsOCP {
				clusterType = "OpenShift"
			}
			fmt.Fprintf(os.Stdout, "checks-qe %s — cluster: %s\n\n", version, clusterType)

			var scenarios []scenario.Scenario
			switch {
			case scenarioName != "":
				scenarios = scenario.Filtered(scenarioName)
			case category != "":
				scenarios = scenario.ByCategory(category)
			default:
				scenarios = scenario.All()
			}

			if len(scenarios) == 0 {
				fmt.Fprintln(os.Stderr, "no scenarios matched the given filters")
				return nil
			}

			fmt.Fprintf(os.Stdout, "Running %d scenario(s) with parallelism %d\n\n", len(scenarios), parallel)

			r := runner.New(cl.Interface, cl.Config, cl.IsOCP, runner.Options{
				Parallel:  parallel,
				Verbose:   verbose,
				SkipProbe: skipProbe,
				SkipOCP:   skipOCP,
			}, func(msg string) {
				fmt.Fprint(os.Stdout, msg)
			})

			summary := r.Run(scenarios)
			fmt.Fprintf(os.Stdout, "\n%s\n", summary)

			if summary.Failed > 0 || summary.Errors > 0 {
				return fmt.Errorf("%d failed, %d errors", summary.Failed, summary.Errors)
			}
			return nil
		},
	}

	cmd.Flags().StringVar(&kubeconfig, "kubeconfig", "", "path to kubeconfig (default: $KUBECONFIG or ~/.kube/config)")
	cmd.Flags().StringVar(&category, "category", "", "run only scenarios in this category")
	cmd.Flags().StringVar(&scenarioName, "scenario", "", "run only scenarios matching this name pattern")
	cmd.Flags().IntVar(&parallel, "parallel", 4, "number of scenarios to run concurrently")
	cmd.Flags().BoolVar(&verbose, "verbose", false, "show resource details on all results")
	cmd.Flags().BoolVar(&skipProbe, "skip-probe", false, "skip scenarios requiring probe pods")
	cmd.Flags().BoolVar(&skipOCP, "skip-ocp", false, "skip scenarios requiring OpenShift")

	return cmd
}
