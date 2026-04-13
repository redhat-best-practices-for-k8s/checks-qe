package observability

import (
	checks "github.com/redhat-best-practices-for-k8s/checks"
	"github.com/redhat-best-practices-for-k8s/checks-qe/pkg/scenario"
)

func registerAPICompat() {
	scenario.Register(
		scenario.Scenario{
			Name:           "observability/api-compatibility/compliant-not-ocp",
			CheckName:      "observability-compatibility-with-next-ocp-release",
			Category:       checks.CategoryObservability,
			Description:    "Non-OCP cluster should be compliant for API compatibility",
			ExpectedStatus: checks.StatusCompliant,
			Setup:          scenario.VanillaDeploymentSetup,
		},
		scenario.Scenario{
			Name:           "observability/api-compatibility/compliant-no-api-counts",
			CheckName:      "observability-compatibility-with-next-ocp-release",
			Category:       checks.CategoryObservability,
			Description:    "OCP cluster with no API request counts should be compliant",
			ExpectedStatus: checks.StatusCompliant,
			Setup:          scenario.VanillaDeploymentSetup,
			PostDiscovery: func(resources *checks.DiscoveredResources) {
				resources.OpenshiftVersion = "4.14.5"
			},
		},
	)
}
