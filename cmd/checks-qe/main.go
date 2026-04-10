package main

import (
	"fmt"
	"os"

	checksall "github.com/redhat-best-practices-for-k8s/checks/all"
	"github.com/spf13/cobra"

	_ "github.com/redhat-best-practices-for-k8s/checks-qe/scenarios/accesscontrol"
	_ "github.com/redhat-best-practices-for-k8s/checks-qe/scenarios/lifecycle"
	_ "github.com/redhat-best-practices-for-k8s/checks-qe/scenarios/manageability"
	_ "github.com/redhat-best-practices-for-k8s/checks-qe/scenarios/networking"
	_ "github.com/redhat-best-practices-for-k8s/checks-qe/scenarios/observability"
	_ "github.com/redhat-best-practices-for-k8s/checks-qe/scenarios/operator"
	_ "github.com/redhat-best-practices-for-k8s/checks-qe/scenarios/performance"
	_ "github.com/redhat-best-practices-for-k8s/checks-qe/scenarios/platform"
)

var version = "dev"

func main() {
	checksall.Register()

	root := &cobra.Command{
		Use:   "checks-qe",
		Short: "Scenario-based QE testing for the checks library",
	}

	root.AddCommand(runCmd())
	root.AddCommand(listCmd())
	root.AddCommand(&cobra.Command{
		Use:   "version",
		Short: "Print version",
		Run:   func(_ *cobra.Command, _ []string) { fmt.Println(version) },
	})

	if err := root.Execute(); err != nil {
		os.Exit(1)
	}
}
