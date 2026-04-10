package observability

import (
	checks "github.com/redhat-best-practices-for-k8s/checks"
	"github.com/redhat-best-practices-for-k8s/checks-qe/pkg/builder"
	"github.com/redhat-best-practices-for-k8s/checks-qe/pkg/cluster"
	"github.com/redhat-best-practices-for-k8s/checks-qe/pkg/scenario"
	corev1 "k8s.io/api/core/v1"
)

func init() {
	scenario.Register(
		scenario.Scenario{
			Name:           "observability/termination-policy/compliant",
			CheckName:      "observability-termination-policy",
			Category:       checks.CategoryObservability,
			Description:    "Deployment with FallbackToLogsOnError should be compliant",
			ExpectedStatus: checks.StatusCompliant,
			Setup: func(ctx *scenario.RunContext) error {
				dep := builder.NewDeployment("test-dep", ctx.Namespace).
					WithTerminationMessagePolicy(corev1.TerminationMessageFallbackToLogsOnError).
					Build()
				return cluster.CreateAndWaitForDeployment(ctx.Ctx, ctx.Client, dep, cluster.DefaultTimeout)
			},
		},
		scenario.Scenario{
			Name:           "observability/termination-policy/non-compliant",
			CheckName:      "observability-termination-policy",
			Category:       checks.CategoryObservability,
			Description:    "Deployment without FallbackToLogsOnError should be non-compliant",
			ExpectedStatus: checks.StatusNonCompliant,
			Setup: scenario.VanillaDeploymentSetup,
		},
	)
}
