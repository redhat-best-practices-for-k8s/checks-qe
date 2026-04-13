package certification

import (
	checks "github.com/redhat-best-practices-for-k8s/checks"
	"github.com/redhat-best-practices-for-k8s/checks-qe/pkg/builder"
	"github.com/redhat-best-practices-for-k8s/checks-qe/pkg/scenario"
)

func Register() {
	scenario.Register(
		scenario.Scenario{
			Name:           "certification/container-certified/compliant",
			CheckName:      "affiliated-certification-container-is-certified-digest",
			Category:       checks.CategoryAffiliatedCertification,
			Description:    "Certified container should be compliant",
			ExpectedStatus: checks.StatusCompliant,
			Setup:          scenario.VanillaDeploymentSetup,
			PostDiscovery: func(resources *checks.DiscoveredResources) {
				resources.CertValidator = &builder.MockCertValidator{ContainerCertified: true}
			},
		},
		scenario.Scenario{
			Name:           "certification/container-certified/non-compliant",
			CheckName:      "affiliated-certification-container-is-certified-digest",
			Category:       checks.CategoryAffiliatedCertification,
			Description:    "Non-certified container should be non-compliant",
			ExpectedStatus: checks.StatusNonCompliant,
			Setup:          scenario.VanillaDeploymentSetup,
			PostDiscovery: func(resources *checks.DiscoveredResources) {
				resources.CertValidator = &builder.MockCertValidator{ContainerCertified: false}
			},
		},
	)

	scenario.Register(
		scenario.Scenario{
			Name:           "certification/operator-certified/compliant",
			CheckName:      "affiliated-certification-operator-is-certified",
			Category:       checks.CategoryAffiliatedCertification,
			Description:    "Certified operator should be compliant",
			ExpectedStatus: checks.StatusCompliant,
			Setup:          scenario.VanillaDeploymentSetup,
			PostDiscovery: func(resources *checks.DiscoveredResources) {
				resources.CertValidator = &builder.MockCertValidator{OperatorCertified: true}
				csv := builder.NewCSV("test-operator.v1.0.0", resources.Namespaces[0]).Build()
				resources.CSVs = append(resources.CSVs, csv)
				resources.OpenshiftVersion = "4.14.5"
			},
		},
		scenario.Scenario{
			Name:           "certification/operator-certified/non-compliant",
			CheckName:      "affiliated-certification-operator-is-certified",
			Category:       checks.CategoryAffiliatedCertification,
			Description:    "Non-certified operator should be non-compliant",
			ExpectedStatus: checks.StatusNonCompliant,
			Setup:          scenario.VanillaDeploymentSetup,
			PostDiscovery: func(resources *checks.DiscoveredResources) {
				resources.CertValidator = &builder.MockCertValidator{OperatorCertified: false}
				csv := builder.NewCSV("test-operator.v1.0.0", resources.Namespaces[0]).Build()
				resources.CSVs = append(resources.CSVs, csv)
				resources.OpenshiftVersion = "4.14.5"
			},
		},
	)

	scenario.Register(
		scenario.Scenario{
			Name:           "certification/helmchart-certified/compliant-no-helm",
			CheckName:      "affiliated-certification-helmchart-is-certified",
			Category:       checks.CategoryAffiliatedCertification,
			Description:    "No Helm chart releases should be compliant",
			ExpectedStatus: checks.StatusCompliant,
			Setup:          scenario.VanillaDeploymentSetup,
			PostDiscovery: func(resources *checks.DiscoveredResources) {
				resources.CertValidator = &builder.MockCertValidator{}
			},
		},
		scenario.Scenario{
			Name:           "certification/helmchart-certified/compliant",
			CheckName:      "affiliated-certification-helmchart-is-certified",
			Category:       checks.CategoryAffiliatedCertification,
			Description:    "Certified Helm chart should be compliant",
			ExpectedStatus: checks.StatusCompliant,
			Setup:          scenario.VanillaDeploymentSetup,
			PostDiscovery: func(resources *checks.DiscoveredResources) {
				resources.CertValidator = &builder.MockCertValidator{HelmChartCertified: true}
				resources.HelmChartReleases = []checks.HelmChartRelease{
					{Name: "my-chart", Version: "1.0.0"},
				}
			},
		},
		scenario.Scenario{
			Name:           "certification/helmchart-certified/non-compliant",
			CheckName:      "affiliated-certification-helmchart-is-certified",
			Category:       checks.CategoryAffiliatedCertification,
			Description:    "Non-certified Helm chart should be non-compliant",
			ExpectedStatus: checks.StatusNonCompliant,
			Setup:          scenario.VanillaDeploymentSetup,
			PostDiscovery: func(resources *checks.DiscoveredResources) {
				resources.CertValidator = &builder.MockCertValidator{HelmChartCertified: false}
				resources.HelmChartReleases = []checks.HelmChartRelease{
					{Name: "my-chart", Version: "1.0.0"},
				}
			},
		},
	)

	scenario.Register(
		scenario.Scenario{
			Name:           "certification/helm-version/compliant-no-helm",
			CheckName:      "affiliated-certification-helm-version",
			Category:       checks.CategoryAffiliatedCertification,
			Description:    "No Helm releases should be compliant",
			ExpectedStatus: checks.StatusCompliant,
			Setup:          scenario.VanillaDeploymentSetup,
		},
		scenario.Scenario{
			Name:           "certification/helm-version/compliant-v3",
			CheckName:      "affiliated-certification-helm-version",
			Category:       checks.CategoryAffiliatedCertification,
			Description:    "Helm v3 (no Tiller pods) should be compliant",
			ExpectedStatus: checks.StatusCompliant,
			Setup:          scenario.VanillaDeploymentSetup,
			PostDiscovery: func(resources *checks.DiscoveredResources) {
				resources.HelmChartReleases = []checks.HelmChartRelease{
					{Name: "my-chart", Version: "1.0.0"},
				}
			},
		},
	)
}
