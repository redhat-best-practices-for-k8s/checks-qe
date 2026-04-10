package networking

import (
	checks "github.com/redhat-best-practices-for-k8s/checks"
	"github.com/redhat-best-practices-for-k8s/checks-qe/pkg/builder"
	"github.com/redhat-best-practices-for-k8s/checks-qe/pkg/cluster"
	"github.com/redhat-best-practices-for-k8s/checks-qe/pkg/scenario"
)

func init() {
	scenario.Register(
		scenario.Scenario{
			Name:           "networking/ocp-reserved-ports/compliant",
			CheckName:      "networking-ocp-reserved-ports-usage",
			Category:       checks.CategoryNetworking,
			Description:    "Deployment without OCP reserved ports should be compliant",
			ExpectedStatus: checks.StatusCompliant,
			Setup: func(ctx *scenario.RunContext) error {
				dep := builder.NewDeployment("test-dep", ctx.Namespace).
					WithContainerPort(8080).
					Build()
				return cluster.CreateAndWaitForDeployment(ctx.Ctx, ctx.Client, dep, cluster.DefaultTimeout)
			},
		},
		scenario.Scenario{
			Name:           "networking/ocp-reserved-ports/non-compliant",
			CheckName:      "networking-ocp-reserved-ports-usage",
			Category:       checks.CategoryNetworking,
			Description:    "Deployment using OCP reserved port 22623 should be non-compliant",
			ExpectedStatus: checks.StatusNonCompliant,
			Setup: func(ctx *scenario.RunContext) error {
				dep := builder.NewDeployment("test-dep", ctx.Namespace).
					WithContainerPort(22623).
					Build()
				return cluster.CreateAndWaitForDeployment(ctx.Ctx, ctx.Client, dep, cluster.DefaultTimeout)
			},
		},
	)

	scenario.Register(
		scenario.Scenario{
			Name:           "networking/reserved-partner-ports/compliant",
			CheckName:      "networking-reserved-partner-ports",
			Category:       checks.CategoryNetworking,
			Description:    "Deployment without Istio reserved ports should be compliant",
			ExpectedStatus: checks.StatusCompliant,
			Setup: func(ctx *scenario.RunContext) error {
				dep := builder.NewDeployment("test-dep", ctx.Namespace).
					WithContainerPort(8080).
					Build()
				return cluster.CreateAndWaitForDeployment(ctx.Ctx, ctx.Client, dep, cluster.DefaultTimeout)
			},
		},
		scenario.Scenario{
			Name:           "networking/reserved-partner-ports/non-compliant",
			CheckName:      "networking-reserved-partner-ports",
			Category:       checks.CategoryNetworking,
			Description:    "Deployment using Istio reserved port 15001 should be non-compliant",
			ExpectedStatus: checks.StatusNonCompliant,
			Setup: func(ctx *scenario.RunContext) error {
				dep := builder.NewDeployment("test-dep", ctx.Namespace).
					WithContainerPort(15001).
					Build()
				return cluster.CreateAndWaitForDeployment(ctx.Ctx, ctx.Client, dep, cluster.DefaultTimeout)
			},
		},
	)
}
