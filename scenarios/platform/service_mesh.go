package platform

import (
	checks "github.com/redhat-best-practices-for-k8s/checks"
	"github.com/redhat-best-practices-for-k8s/checks-qe/pkg/builder"
	"github.com/redhat-best-practices-for-k8s/checks-qe/pkg/cluster"
	"github.com/redhat-best-practices-for-k8s/checks-qe/pkg/scenario"
)

func init() {
	scenario.Register(
		scenario.Scenario{
			Name:           "platform/service-mesh-usage/compliant",
			CheckName:      "platform-alteration-service-mesh-usage",
			Category:       checks.CategoryPlatformAlteration,
			Description:    "Deployment without Istio sidecar should be compliant",
			ExpectedStatus: checks.StatusCompliant,
			Setup:          scenario.VanillaDeploymentSetup,
		},
		scenario.Scenario{
			Name:           "platform/service-mesh-usage/non-compliant",
			CheckName:      "platform-alteration-service-mesh-usage",
			Category:       checks.CategoryPlatformAlteration,
			Description:    "Deployment with Istio sidecar injection annotation should be non-compliant",
			ExpectedStatus: checks.StatusNonCompliant,
			Setup: func(ctx *scenario.RunContext) error {
				dep := builder.NewDeployment("test-dep", ctx.Namespace).
					WithPodAnnotation("sidecar.istio.io/inject", "true").
					Build()
				return cluster.CreateAndWaitForDeployment(ctx.Ctx, ctx.Client, dep, cluster.DefaultTimeout)
			},
		},
	)
}
