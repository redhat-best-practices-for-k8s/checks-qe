package accesscontrol

import (
	checks "github.com/redhat-best-practices-for-k8s/checks"
	"github.com/redhat-best-practices-for-k8s/checks-qe/pkg/builder"
	"github.com/redhat-best-practices-for-k8s/checks-qe/pkg/cluster"
	"github.com/redhat-best-practices-for-k8s/checks-qe/pkg/scenario"
)

func registerHostIPC() {
	scenario.Register(
		scenario.Scenario{
			Name:           "accesscontrol/host-ipc/compliant",
			CheckName:      "access-control-pod-host-ipc",
			Category:       checks.CategoryAccessControl,
			Description:    "Deployment with hostIPC=false should be compliant",
			ExpectedStatus: checks.StatusCompliant,
			Setup: func(ctx *scenario.RunContext) error {
				dep := builder.NewDeployment("test-dep", ctx.Namespace).Build()
				return cluster.CreateAndWaitForDeployment(ctx.Ctx, ctx.Client, dep, cluster.DefaultTimeout)
			},
		},
		scenario.Scenario{
			Name:           "accesscontrol/host-ipc/non-compliant",
			CheckName:      "access-control-pod-host-ipc",
			Category:       checks.CategoryAccessControl,
			Description:    "Deployment with hostIPC=true should be non-compliant",
			ExpectedStatus: checks.StatusNonCompliant,
			Privileged:     true,
			Setup: func(ctx *scenario.RunContext) error {
				dep := builder.NewDeployment("test-dep", ctx.Namespace).
					WithHostIPC(true).
					Build()
				return cluster.CreateAndWaitForDeployment(ctx.Ctx, ctx.Client, dep, cluster.DefaultTimeout)
			},
		},
		scenario.Scenario{
			Name:           "accesscontrol/host-ipc/mixed-two-deployments",
			CheckName:      "access-control-pod-host-ipc",
			Category:       checks.CategoryAccessControl,
			Description:    "Two deployments, one with hostIPC=true should be non-compliant",
			ExpectedStatus: checks.StatusNonCompliant,
			Privileged:     true,
			Setup: scenario.TwoDeploymentSetup(func(b *builder.DeploymentBuilder) *builder.DeploymentBuilder {
				return b.WithHostIPC(true)
			}),
		},
	)
}
