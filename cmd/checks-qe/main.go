package main

import (
	"fmt"
	"os"

	checksall "github.com/redhat-best-practices-for-k8s/checks/all"
	"github.com/spf13/cobra"

	"github.com/redhat-best-practices-for-k8s/checks-qe/scenarios/accesscontrol"
	"github.com/redhat-best-practices-for-k8s/checks-qe/scenarios/certification"
	"github.com/redhat-best-practices-for-k8s/checks-qe/scenarios/lifecycle"
	"github.com/redhat-best-practices-for-k8s/checks-qe/scenarios/manageability"
	"github.com/redhat-best-practices-for-k8s/checks-qe/scenarios/networking"
	"github.com/redhat-best-practices-for-k8s/checks-qe/scenarios/observability"
	"github.com/redhat-best-practices-for-k8s/checks-qe/scenarios/operator"
	"github.com/redhat-best-practices-for-k8s/checks-qe/scenarios/performance"
	"github.com/redhat-best-practices-for-k8s/checks-qe/scenarios/platform"
)

var version = "dev"

func main() {
	checksall.Register()

	accesscontrol.Register()
	certification.Register()
	lifecycle.Register()
	manageability.Register()
	networking.Register()
	observability.Register()
	operator.Register()
	performance.Register()
	platform.Register()

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
