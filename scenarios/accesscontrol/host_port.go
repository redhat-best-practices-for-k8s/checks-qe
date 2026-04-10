package accesscontrol

import (
	checks "github.com/redhat-best-practices-for-k8s/checks"
	"github.com/redhat-best-practices-for-k8s/checks-qe/pkg/builder"
	"github.com/redhat-best-practices-for-k8s/checks-qe/pkg/cluster"
	"github.com/redhat-best-practices-for-k8s/checks-qe/pkg/scenario"
)

func init() {
	scenario.Register(
		scenario.Scenario{
			Name:           "accesscontrol/host-port/compliant",
			CheckName:      "access-control-container-host-port",
			Category:       checks.CategoryAccessControl,
			Description:    "Deployment without host port should be compliant",
			ExpectedStatus: checks.StatusCompliant,
			Setup: func(ctx *scenario.RunContext) error {
				dep := builder.NewDeployment("test-dep", ctx.Namespace).
					WithContainerPort(8080).
					Build()
				return cluster.CreateAndWaitForDeployment(ctx.Ctx, ctx.Client, dep, cluster.DefaultTimeout)
			},
		},
		scenario.Scenario{
			Name:           "accesscontrol/host-port/non-compliant",
			CheckName:      "access-control-container-host-port",
			Category:       checks.CategoryAccessControl,
			Description:    "Deployment with host port should be non-compliant",
			ExpectedStatus: checks.StatusNonCompliant,
			Privileged:     true,
			Setup: func(ctx *scenario.RunContext) error {
				dep := builder.NewDeployment("test-dep", ctx.Namespace).
					WithHostPort(8080, 8080).
					Build()
				return cluster.CreateAndWaitForDeployment(ctx.Ctx, ctx.Client, dep, cluster.DefaultTimeout)
			},
		},
		scenario.Scenario{
			Name:           "accesscontrol/host-port/compliant-two-containers",
			CheckName:      "access-control-container-host-port",
			Category:       checks.CategoryAccessControl,
			Description:    "Two containers without host ports should be compliant",
			ExpectedStatus: checks.StatusCompliant,
			Setup: func(ctx *scenario.RunContext) error {
				dep := builder.NewDeployment("test-dep", ctx.Namespace).
					WithContainerPort(8080).
					WithSecondContainer("sidecar", builder.DefaultImage).
					Build()
				return cluster.CreateAndWaitForDeployment(ctx.Ctx, ctx.Client, dep, cluster.DefaultTimeout)
			},
		},
		scenario.Scenario{
			Name:           "accesscontrol/host-port/mixed-two-deployments",
			CheckName:      "access-control-container-host-port",
			Category:       checks.CategoryAccessControl,
			Description:    "Two deployments, one with host port should be non-compliant",
			ExpectedStatus: checks.StatusNonCompliant,
			Privileged:     true,
			Setup: scenario.TwoDeploymentSetup(func(b *builder.DeploymentBuilder) *builder.DeploymentBuilder {
				return b.WithHostPort(8080, 8080)
			}),
		},
	)
}
