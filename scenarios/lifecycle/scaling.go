package lifecycle

import (
	checks "github.com/redhat-best-practices-for-k8s/checks"
	"github.com/redhat-best-practices-for-k8s/checks-qe/pkg/builder"
	"github.com/redhat-best-practices-for-k8s/checks-qe/pkg/cluster"
	"github.com/redhat-best-practices-for-k8s/checks-qe/pkg/scenario"
)

func registerScaling() {
	scenario.Register(
		scenario.Scenario{
			Name:           "lifecycle/deployment-scaling/compliant",
			CheckName:      "lifecycle-deployment-scaling",
			Category:       checks.CategoryLifecycle,
			Description:    "Deployment that can scale up and down should be compliant",
			ExpectedStatus: checks.StatusCompliant,
			Setup:          scenario.VanillaDeploymentSetup,
		},
	)

	scenario.Register(
		scenario.Scenario{
			Name:           "lifecycle/statefulset-scaling/compliant",
			CheckName:      "lifecycle-statefulset-scaling",
			Category:       checks.CategoryLifecycle,
			Description:    "StatefulSet that can scale up and down should be compliant",
			ExpectedStatus: checks.StatusCompliant,
			Setup: func(ctx *scenario.RunContext) error {
				sts := builder.NewStatefulSet("test-sts", ctx.Namespace).Build()
				return cluster.CreateAndWaitForStatefulSet(ctx.Ctx, ctx.Client, sts, cluster.DefaultTimeout)
			},
		},
	)

	scenario.Register(
		scenario.Scenario{
			Name:           "lifecycle/crd-scaling/error-no-scale-client",
			CheckName:      "lifecycle-crd-scaling",
			Category:       checks.CategoryLifecycle,
			Description:    "No ScaleClient returns error for CRD scaling check",
			ExpectedStatus: checks.StatusError,
			Setup:          scenario.VanillaDeploymentSetup,
		},
	)
}
