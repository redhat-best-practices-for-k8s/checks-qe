package accesscontrol

import (
	checks "github.com/redhat-best-practices-for-k8s/checks"
	"github.com/redhat-best-practices-for-k8s/checks-qe/pkg/scenario"
)

func registerNamespace() {
	scenario.Register(
		scenario.Scenario{
			Name:           "accesscontrol/namespace/compliant",
			CheckName:      "access-control-namespace",
			Category:       checks.CategoryAccessControl,
			Description:    "Configured namespace with valid prefix should be compliant",
			ExpectedStatus: checks.StatusCompliant,
			Setup:          scenario.VanillaDeploymentSetup,
		},
		scenario.Scenario{
			Name:           "accesscontrol/namespace/non-compliant-openshift-prefix",
			CheckName:      "access-control-namespace",
			Category:       checks.CategoryAccessControl,
			Description:    "Configured namespace with openshift- prefix should be non-compliant",
			ExpectedStatus: checks.StatusNonCompliant,
			Setup:          scenario.VanillaDeploymentSetup,
			PostDiscovery: func(resources *checks.DiscoveredResources) {
				resources.Namespaces = append(resources.Namespaces, "openshift-test")
			},
		},
		scenario.Scenario{
			Name:           "accesscontrol/namespace/non-compliant-default",
			CheckName:      "access-control-namespace",
			Category:       checks.CategoryAccessControl,
			Description:    "Configured namespace with default prefix should be non-compliant",
			ExpectedStatus: checks.StatusNonCompliant,
			Setup:          scenario.VanillaDeploymentSetup,
			PostDiscovery: func(resources *checks.DiscoveredResources) {
				resources.Namespaces = append(resources.Namespaces, "default")
			},
		},
		scenario.Scenario{
			Name:           "accesscontrol/namespace/mixed-namespaces",
			CheckName:      "access-control-namespace",
			Category:       checks.CategoryAccessControl,
			Description:    "Mix of valid and invalid configured namespaces should be non-compliant",
			ExpectedStatus: checks.StatusNonCompliant,
			Setup:          scenario.VanillaDeploymentSetup,
			PostDiscovery: func(resources *checks.DiscoveredResources) {
				resources.Namespaces = append(resources.Namespaces, "openshift-bad")
			},
		},
		scenario.Scenario{
			Name:           "accesscontrol/namespace/compliant-cr-in-valid-ns",
			CheckName:      "access-control-namespace",
			Category:       checks.CategoryAccessControl,
			Description:    "CR in configured namespace should be compliant",
			ExpectedStatus: checks.StatusCompliant,
			Setup:          scenario.VanillaDeploymentSetup,
			PostDiscovery: func(resources *checks.DiscoveredResources) {
				ns := resources.Namespaces[0]
				resources.CRInstances = map[string]map[string][]string{
					"widgets.example.com": {ns: {"my-widget"}},
				}
			},
		},
		scenario.Scenario{
			Name:           "accesscontrol/namespace/non-compliant-cr-in-invalid-ns",
			CheckName:      "access-control-namespace",
			Category:       checks.CategoryAccessControl,
			Description:    "CR in unconfigured namespace should be non-compliant",
			ExpectedStatus: checks.StatusNonCompliant,
			Setup:          scenario.VanillaDeploymentSetup,
			PostDiscovery: func(resources *checks.DiscoveredResources) {
				resources.CRInstances = map[string]map[string][]string{
					"widgets.example.com": {"other-namespace": {"my-widget"}},
				}
			},
		},
	)
}
