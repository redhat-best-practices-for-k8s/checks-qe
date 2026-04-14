package observability

import (
	checks "github.com/redhat-best-practices-for-k8s/checks"
	"github.com/redhat-best-practices-for-k8s/checks-qe/pkg/builder"
	"github.com/redhat-best-practices-for-k8s/checks-qe/pkg/cluster"
	"github.com/redhat-best-practices-for-k8s/checks-qe/pkg/scenario"
)

func registerContainerLogging() {
	scenario.Register(
		scenario.Scenario{
			Name:           "observability/container-logging/compliant",
			CheckName:      "observability-container-logging",
			Category:       checks.CategoryObservability,
			Description:    "Container producing stdout output should be compliant",
			ExpectedStatus: checks.StatusCompliant,
			Setup: func(ctx *scenario.RunContext) error {
				dep := builder.NewDeployment("test-dep", ctx.Namespace).
					WithCommand("/bin/sh", "-c", "echo started; while true; do echo alive; sleep 1; done").
					Build()
				return cluster.CreateAndWaitForDeployment(ctx.Ctx, ctx.Client, dep, cluster.DefaultTimeout)
			},
		},
		scenario.Scenario{
			Name:           "observability/container-logging/non-compliant",
			CheckName:      "observability-container-logging",
			Category:       checks.CategoryObservability,
			Description:    "Container producing no stdout output should be non-compliant",
			ExpectedStatus: checks.StatusNonCompliant,
			Setup: func(ctx *scenario.RunContext) error {
				dep := builder.NewDeployment("test-dep", ctx.Namespace).
					WithCommand("/bin/sh", "-c", "sleep infinity").
					Build()
				return cluster.CreateAndWaitForDeployment(ctx.Ctx, ctx.Client, dep, cluster.DefaultTimeout)
			},
		},
	)
}
