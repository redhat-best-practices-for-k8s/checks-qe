package operator

import (
	checks "github.com/redhat-best-practices-for-k8s/checks"
	"github.com/redhat-best-practices-for-k8s/checks-qe/pkg/builder"
	"github.com/redhat-best-practices-for-k8s/checks-qe/pkg/scenario"
	olmv1alpha1 "github.com/operator-framework/api/pkg/operators/v1alpha1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

func registerCatalog() {
	scenario.Register(
		scenario.Scenario{
			Name:           "operator/catalogsource-bundle-count/compliant-no-catalogs",
			CheckName:      "operator-catalogsource-bundle-count",
			Category:       checks.CategoryOperator,
			Description:    "No catalog sources should be compliant",
			ExpectedStatus: checks.StatusCompliant,
			Setup:          scenario.VanillaDeploymentSetup,
		},
		scenario.Scenario{
			Name:           "operator/catalogsource-bundle-count/compliant",
			CheckName:      "operator-catalogsource-bundle-count",
			Category:       checks.CategoryOperator,
			Description:    "Catalog source with operator should be compliant",
			ExpectedStatus: checks.StatusCompliant,
			Setup:          scenario.VanillaDeploymentSetup,
			PostDiscovery: func(resources *checks.DiscoveredResources) {
				csv := builder.NewCSV("test-operator.v1.0.0", resources.Namespaces[0]).
					WithAnnotation("operatorframework.io/properties",
						`[{"type":"olm.package","value":{"packageName":"test-operator","version":"1.0.0"}}]`).
					Build()
				resources.CSVs = append(resources.CSVs, csv)
				resources.CatalogSources = append(resources.CatalogSources,
					olmv1alpha1.CatalogSource{
						ObjectMeta: metav1.ObjectMeta{
							Name:      "test-catalog",
							Namespace: resources.Namespaces[0],
						},
						Spec: olmv1alpha1.CatalogSourceSpec{
							SourceType: olmv1alpha1.SourceTypeGrpc,
						},
					},
				)
			},
		},
	)
}
