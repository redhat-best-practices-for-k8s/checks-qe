package observability

import (
	checks "github.com/redhat-best-practices-for-k8s/checks"
	"github.com/redhat-best-practices-for-k8s/checks-qe/pkg/builder"
	"github.com/redhat-best-practices-for-k8s/checks-qe/pkg/scenario"
)

func registerCRDStatus() {
	scenario.Register(
		scenario.Scenario{
			Name:           "observability/crd-status/compliant",
			CheckName:      "observability-crd-status",
			Category:       checks.CategoryObservability,
			Description:    "CRD with status subresource should be compliant",
			ExpectedStatus: checks.StatusCompliant,
			Setup:          scenario.VanillaDeploymentSetup,
			PostDiscovery: func(resources *checks.DiscoveredResources) {
				crd := builder.NewCRD("testresources", "example.com").
					WithStatusSubresource().
					Build()
				resources.CRDs = append(resources.CRDs, crd)
			},
		},
		scenario.Scenario{
			Name:           "observability/crd-status/non-compliant",
			CheckName:      "observability-crd-status",
			Category:       checks.CategoryObservability,
			Description:    "CRD without status subresource should be non-compliant",
			ExpectedStatus: checks.StatusNonCompliant,
			Setup:          scenario.VanillaDeploymentSetup,
			PostDiscovery: func(resources *checks.DiscoveredResources) {
				crd := builder.NewCRD("testresources", "example.com").Build()
				resources.CRDs = append(resources.CRDs, crd)
			},
		},
	)
}
