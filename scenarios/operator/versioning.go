package operator

import (
	checks "github.com/redhat-best-practices-for-k8s/checks"
	"github.com/redhat-best-practices-for-k8s/checks-qe/pkg/builder"
	"github.com/redhat-best-practices-for-k8s/checks-qe/pkg/scenario"
)

func registerVersioning() {
	scenario.Register(
		scenario.Scenario{
			Name:           "operator/semantic-versioning/compliant",
			CheckName:      "operator-semantic-versioning",
			Category:       checks.CategoryOperator,
			Description:    "CSV with valid semver should be compliant",
			ExpectedStatus: checks.StatusCompliant,
			Setup:          scenario.VanillaDeploymentSetup,
			PostDiscovery: func(resources *checks.DiscoveredResources) {
				csv := builder.NewCSV("test-operator.v1.2.3", resources.Namespaces[0]).
					WithVersion("1.2.3").
					Build()
				resources.CSVs = append(resources.CSVs, csv)
			},
		},
	)

	scenario.Register(
		scenario.Scenario{
			Name:           "operator/olm-skip-range/compliant",
			CheckName:      "operator-olm-skip-range",
			Category:       checks.CategoryOperator,
			Description:    "CSV with olm.skipRange annotation should be compliant",
			ExpectedStatus: checks.StatusCompliant,
			Setup:          scenario.VanillaDeploymentSetup,
			PostDiscovery: func(resources *checks.DiscoveredResources) {
				csv := builder.NewCSV("test-operator.v1.2.3", resources.Namespaces[0]).
					WithSkipRange(">=1.0.0 <2.0.0").
					Build()
				resources.CSVs = append(resources.CSVs, csv)
			},
		},
		scenario.Scenario{
			Name:           "operator/olm-skip-range/non-compliant",
			CheckName:      "operator-olm-skip-range",
			Category:       checks.CategoryOperator,
			Description:    "CSV without olm.skipRange annotation should be non-compliant",
			ExpectedStatus: checks.StatusNonCompliant,
			Setup:          scenario.VanillaDeploymentSetup,
			PostDiscovery: func(resources *checks.DiscoveredResources) {
				csv := builder.NewCSV("test-operator.v1.2.3", resources.Namespaces[0]).
					Build()
				resources.CSVs = append(resources.CSVs, csv)
			},
		},
	)

	scenario.Register(
		scenario.Scenario{
			Name:           "operator/crd-versioning/compliant",
			CheckName:      "operator-crd-versioning",
			Category:       checks.CategoryOperator,
			Description:    "CRD with valid K8s version naming should be compliant",
			ExpectedStatus: checks.StatusCompliant,
			Setup:          scenario.VanillaDeploymentSetup,
			PostDiscovery: func(resources *checks.DiscoveredResources) {
				crd := builder.NewCRD("testresources", "example.com").Build()
				resources.CRDs = append(resources.CRDs, crd)
			},
		},
		scenario.Scenario{
			Name:           "operator/crd-versioning/non-compliant",
			CheckName:      "operator-crd-versioning",
			Category:       checks.CategoryOperator,
			Description:    "CRD with invalid version naming should be non-compliant",
			ExpectedStatus: checks.StatusNonCompliant,
			Setup:          scenario.VanillaDeploymentSetup,
			PostDiscovery: func(resources *checks.DiscoveredResources) {
				crd := builder.NewCRD("testresources", "example.com").
					WithVersion("invalid", true).
					Build()
				resources.CRDs = append(resources.CRDs, crd)
			},
		},
	)

	scenario.Register(
		scenario.Scenario{
			Name:           "operator/crd-openapi-schema/compliant",
			CheckName:      "operator-crd-openapi-schema",
			Category:       checks.CategoryOperator,
			Description:    "CRD with OpenAPI schema should be compliant",
			ExpectedStatus: checks.StatusCompliant,
			Setup:          scenario.VanillaDeploymentSetup,
			PostDiscovery: func(resources *checks.DiscoveredResources) {
				crd := builder.NewCRD("testresources", "example.com").Build()
				resources.CRDs = append(resources.CRDs, crd)
			},
		},
		scenario.Scenario{
			Name:           "operator/crd-openapi-schema/non-compliant",
			CheckName:      "operator-crd-openapi-schema",
			Category:       checks.CategoryOperator,
			Description:    "CRD without OpenAPI schema should be non-compliant",
			ExpectedStatus: checks.StatusNonCompliant,
			Setup:          scenario.VanillaDeploymentSetup,
			PostDiscovery: func(resources *checks.DiscoveredResources) {
				crd := builder.NewCRD("testresources", "example.com").
					WithoutSchema().
					Build()
				resources.CRDs = append(resources.CRDs, crd)
			},
		},
	)
}
