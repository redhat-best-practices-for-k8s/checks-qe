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
			Name:           "accesscontrol/host-pid/compliant",
			CheckName:      "access-control-pod-host-pid",
			Category:       checks.CategoryAccessControl,
			Description:    "Deployment without hostPID should be compliant",
			ExpectedStatus: checks.StatusCompliant,
			Setup: func(ctx *scenario.RunContext) error {
				dep := builder.NewDeployment("test-dep", ctx.Namespace).Build()
				return cluster.CreateAndWaitForDeployment(ctx.Ctx, ctx.Client, dep, cluster.DefaultTimeout)
			},
		},
		scenario.Scenario{
			Name:           "accesscontrol/host-pid/non-compliant",
			CheckName:      "access-control-pod-host-pid",
			Category:       checks.CategoryAccessControl,
			Description:    "Deployment with hostPID=true should be non-compliant",
			ExpectedStatus: checks.StatusNonCompliant,
			Privileged:     true,
			Setup: func(ctx *scenario.RunContext) error {
				dep := builder.NewDeployment("test-dep", ctx.Namespace).
					WithHostPID(true).
					Build()
				return cluster.CreateAndWaitForDeployment(ctx.Ctx, ctx.Client, dep, cluster.DefaultTimeout)
			},
		},
		scenario.Scenario{
			Name:           "accesscontrol/host-pid/mixed-two-deployments",
			CheckName:      "access-control-pod-host-pid",
			Category:       checks.CategoryAccessControl,
			Description:    "Two deployments, one with hostPID=true should be non-compliant",
			ExpectedStatus: checks.StatusNonCompliant,
			Privileged:     true,
			Setup: scenario.TwoDeploymentSetup(func(b *builder.DeploymentBuilder) *builder.DeploymentBuilder {
				return b.WithHostPID(true)
			}),
		},
	)
}
