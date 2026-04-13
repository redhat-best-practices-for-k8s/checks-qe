package operator

import (
	olmv1alpha1 "github.com/operator-framework/api/pkg/operators/v1alpha1"
	checks "github.com/redhat-best-practices-for-k8s/checks"
	"github.com/redhat-best-practices-for-k8s/checks-qe/pkg/builder"
	"github.com/redhat-best-practices-for-k8s/checks-qe/pkg/scenario"
)

func registerInstall() {
	scenario.Register(
		scenario.Scenario{
			Name:           "operator/install-status-succeeded/compliant",
			CheckName:      "operator-install-status-succeeded",
			Category:       checks.CategoryOperator,
			Description:    "CSV in Succeeded phase should be compliant",
			ExpectedStatus: checks.StatusCompliant,
			Setup:          scenario.VanillaDeploymentSetup,
			PostDiscovery: func(resources *checks.DiscoveredResources) {
				csv := builder.NewCSV("test-operator.v1.0.0", resources.Namespaces[0]).
					WithPhase(olmv1alpha1.CSVPhaseSucceeded).
					Build()
				resources.CSVs = append(resources.CSVs, csv)
			},
		},
		scenario.Scenario{
			Name:           "operator/install-status-succeeded/non-compliant",
			CheckName:      "operator-install-status-succeeded",
			Category:       checks.CategoryOperator,
			Description:    "CSV in Installing phase should be non-compliant",
			ExpectedStatus: checks.StatusNonCompliant,
			Setup:          scenario.VanillaDeploymentSetup,
			PostDiscovery: func(resources *checks.DiscoveredResources) {
				csv := builder.NewCSV("test-operator.v1.0.0", resources.Namespaces[0]).
					WithPhase(olmv1alpha1.CSVPhaseInstalling).
					Build()
				resources.CSVs = append(resources.CSVs, csv)
			},
		},
	)

	scenario.Register(
		scenario.Scenario{
			Name:           "operator/install-source/compliant",
			CheckName:      "operator-install-source",
			Category:       checks.CategoryOperator,
			Description:    "CSV with OLM namespace annotation should be compliant",
			ExpectedStatus: checks.StatusCompliant,
			Setup:          scenario.VanillaDeploymentSetup,
			PostDiscovery: func(resources *checks.DiscoveredResources) {
				csv := builder.NewCSV("test-operator.v1.0.0", resources.Namespaces[0]).
					WithOLMInstalled().
					Build()
				resources.CSVs = append(resources.CSVs, csv)
			},
		},
		scenario.Scenario{
			Name:           "operator/install-source/non-compliant",
			CheckName:      "operator-install-source",
			Category:       checks.CategoryOperator,
			Description:    "CSV without OLM namespace annotation should be non-compliant",
			ExpectedStatus: checks.StatusNonCompliant,
			Setup:          scenario.VanillaDeploymentSetup,
			PostDiscovery: func(resources *checks.DiscoveredResources) {
				csv := builder.NewCSV("test-operator.v1.0.0", resources.Namespaces[0]).Build()
				resources.CSVs = append(resources.CSVs, csv)
			},
		},
	)

	scenario.Register(
		scenario.Scenario{
			Name:           "operator/install-status-no-privileges/compliant",
			CheckName:      "operator-install-status-no-privileges",
			Category:       checks.CategoryOperator,
			Description:    "CSV without SCC cluster permissions should be compliant",
			ExpectedStatus: checks.StatusCompliant,
			Setup:          scenario.VanillaDeploymentSetup,
			PostDiscovery: func(resources *checks.DiscoveredResources) {
				csv := builder.NewCSV("test-operator.v1.0.0", resources.Namespaces[0]).Build()
				resources.CSVs = append(resources.CSVs, csv)
			},
		},
		scenario.Scenario{
			Name:           "operator/install-status-no-privileges/non-compliant",
			CheckName:      "operator-install-status-no-privileges",
			Category:       checks.CategoryOperator,
			Description:    "CSV with SCC cluster permissions should be non-compliant",
			ExpectedStatus: checks.StatusNonCompliant,
			Setup:          scenario.VanillaDeploymentSetup,
			PostDiscovery: func(resources *checks.DiscoveredResources) {
				csv := builder.NewCSV("test-operator.v1.0.0", resources.Namespaces[0]).
					WithClusterPermissionSCC().
					Build()
				resources.CSVs = append(resources.CSVs, csv)
			},
		},
	)
}
