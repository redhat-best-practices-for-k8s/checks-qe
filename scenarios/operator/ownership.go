package operator

import (
	checks "github.com/redhat-best-practices-for-k8s/checks"
	"github.com/redhat-best-practices-for-k8s/checks-qe/pkg/builder"
	"github.com/redhat-best-practices-for-k8s/checks-qe/pkg/scenario"
)

func init() {
	scenario.Register(
		scenario.Scenario{
			Name:           "operator/single-crd-owner/compliant",
			CheckName:      "operator-single-crd-owner",
			Category:       checks.CategoryOperator,
			Description:    "CRD owned by single operator should be compliant",
			ExpectedStatus: checks.StatusCompliant,
			Setup:          scenario.VanillaDeploymentSetup,
			PostDiscovery: func(resources *checks.DiscoveredResources) {
				csv := builder.NewCSV("test-operator.v1.0.0", resources.Namespaces[0]).
					WithOwnedCRD("widgets.example.com").
					Build()
				resources.CSVs = append(resources.CSVs, csv)
			},
		},
		scenario.Scenario{
			Name:           "operator/single-crd-owner/non-compliant",
			CheckName:      "operator-single-crd-owner",
			Category:       checks.CategoryOperator,
			Description:    "CRD owned by multiple operators should be non-compliant",
			ExpectedStatus: checks.StatusNonCompliant,
			Setup:          scenario.VanillaDeploymentSetup,
			PostDiscovery: func(resources *checks.DiscoveredResources) {
				csv1 := builder.NewCSV("operator-a.v1.0.0", resources.Namespaces[0]).
					WithOwnedCRD("widgets.example.com").
					Build()
				csv2 := builder.NewCSV("operator-b.v1.0.0", resources.Namespaces[0]).
					WithOwnedCRD("widgets.example.com").
					Build()
				resources.CSVs = append(resources.CSVs, csv1, csv2)
			},
		},
	)

	scenario.Register(
		scenario.Scenario{
			Name:           "operator/multiple-same-operators/compliant",
			CheckName:      "operator-multiple-same-operators",
			Category:       checks.CategoryOperator,
			Description:    "Single operator instance should be compliant",
			ExpectedStatus: checks.StatusCompliant,
			Setup:          scenario.VanillaDeploymentSetup,
			PostDiscovery: func(resources *checks.DiscoveredResources) {
				csv := builder.NewCSV("test-operator.v1.0.0", resources.Namespaces[0]).Build()
				resources.CSVs = append(resources.CSVs, csv)
			},
		},
		scenario.Scenario{
			Name:           "operator/multiple-same-operators/non-compliant",
			CheckName:      "operator-multiple-same-operators",
			Category:       checks.CategoryOperator,
			Description:    "Same operator installed twice should be non-compliant",
			ExpectedStatus: checks.StatusNonCompliant,
			Setup:          scenario.VanillaDeploymentSetup,
			PostDiscovery: func(resources *checks.DiscoveredResources) {
				csv1 := builder.NewCSV("test-operator.v1.0.0", resources.Namespaces[0]).Build()
				csv2 := builder.NewCSV("test-operator.v2.0.0", resources.Namespaces[0]).
					WithVersion("2.0.0").
					Build()
				resources.CSVs = append(resources.CSVs, csv1, csv2)
			},
		},
	)

	scenario.Register(
		scenario.Scenario{
			Name:           "operator/pods-no-hugepages/compliant",
			CheckName:      "operator-pods-no-hugepages",
			Category:       checks.CategoryOperator,
			Description:    "Pods without hugepage requests should be compliant",
			ExpectedStatus: checks.StatusCompliant,
			Setup:          scenario.VanillaDeploymentSetup,
		},
	)

	scenario.Register(
		scenario.Scenario{
			Name:           "operator/namespace-mode/compliant",
			CheckName:      "operator-single-or-multi-namespaced-allowed-in-tenant-namespaces",
			Category:       checks.CategoryOperator,
			Description:    "Operator with OwnNamespace mode in tenant namespace should be compliant",
			ExpectedStatus: checks.StatusCompliant,
			Setup:          scenario.VanillaDeploymentSetup,
			PostDiscovery: func(resources *checks.DiscoveredResources) {
				csv := builder.NewCSV("test-operator.v1.0.0", resources.Namespaces[0]).
					WithOLMInstalled().
					Build()
				resources.CSVs = append(resources.CSVs, csv)
			},
		},
	)
}
