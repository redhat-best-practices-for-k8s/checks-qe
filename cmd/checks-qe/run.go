package main

import (
	"context"
	"fmt"
	"os"

	"github.com/redhat-best-practices-for-k8s/checks-qe/pkg/cluster"
	"github.com/redhat-best-practices-for-k8s/checks-qe/pkg/runner"
	"github.com/redhat-best-practices-for-k8s/checks-qe/pkg/scenario"
	"github.com/spf13/cobra"
)

func runCmd() *cobra.Command {
	var (
		kubeconfig    string
		category      string
		scenarioName  string
		parallel      int
		verbose       bool
		skipProbe     bool
		skipOCP       bool
		mode          string
		operatorNS    string
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
			fmt.Fprintf(os.Stdout, "checks-qe %s — cluster: %s, mode: %s\n\n", version, clusterType, mode)

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

			opts := runner.Options{
				Parallel:  parallel,
				Verbose:   verbose,
				SkipProbe: skipProbe,
				SkipOCP:   skipOCP,
			}
			printFn := func(msg string) {
				fmt.Fprint(os.Stdout, msg)
			}

			var summary runner.Summary
			switch mode {
			case "operator":
				ctrlClient, err := cluster.NewOperatorClient(cl.Config)
				if err != nil {
					return fmt.Errorf("creating operator client: %w", err)
				}
				if err := cluster.VerifyOperatorRunning(context.Background(), ctrlClient, operatorNS); err != nil {
					return fmt.Errorf("operator preflight check: %w", err)
				}
				r := runner.NewOperator(cl.Interface, ctrlClient, cl.Config, cl.IsOCP, opts, printFn)
				summary = r.Run(scenarios)
			default:
				r := runner.New(cl.Interface, cl.Config, cl.IsOCP, opts, printFn)
				summary = r.Run(scenarios)
			}

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
	cmd.Flags().StringVar(&mode, "mode", "direct", "execution mode: direct (call checks) or operator (via bps-operator)")
	cmd.Flags().StringVar(&operatorNS, "operator-namespace", "bps-operator-system", "namespace where the bps-operator is deployed (operator mode only)")

	return cmd
}
