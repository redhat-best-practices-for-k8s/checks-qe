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
			Name:           "accesscontrol/sys-ptrace/compliant",
			CheckName:      "access-control-sys-ptrace-capability",
			Category:       checks.CategoryAccessControl,
			Description:    "Deployment with shareProcessNamespace + SYS_PTRACE should be compliant",
			ExpectedStatus: checks.StatusCompliant,
			Privileged:     true,
			Setup: func(ctx *scenario.RunContext) error {
				dep := builder.NewDeployment("test-dep", ctx.Namespace).
					WithShareProcessNamespace(true).
					WithCapability("SYS_PTRACE").
					Build()
				return cluster.CreateAndWaitForDeployment(ctx.Ctx, ctx.Client, dep, cluster.DefaultTimeout)
			},
		},
		scenario.Scenario{
			Name:           "accesscontrol/sys-ptrace/non-compliant",
			CheckName:      "access-control-sys-ptrace-capability",
			Category:       checks.CategoryAccessControl,
			Description:    "Deployment with shareProcessNamespace but without SYS_PTRACE should be non-compliant",
			ExpectedStatus: checks.StatusNonCompliant,
			Privileged:     true,
			Setup: func(ctx *scenario.RunContext) error {
				dep := builder.NewDeployment("test-dep", ctx.Namespace).
					WithShareProcessNamespace(true).
					Build()
				return cluster.CreateAndWaitForDeployment(ctx.Ctx, ctx.Client, dep, cluster.DefaultTimeout)
			},
		},
		scenario.Scenario{
			Name:           "accesscontrol/sys-ptrace/mixed-two-deployments",
			CheckName:      "access-control-sys-ptrace-capability",
			Category:       checks.CategoryAccessControl,
			Description:    "Two deployments with shareProcessNamespace, one without SYS_PTRACE",
			ExpectedStatus: checks.StatusNonCompliant,
			Privileged:     true,
			Setup: func(ctx *scenario.RunContext) error {
				dep1 := builder.NewDeployment("test-dep-1", ctx.Namespace).
					WithShareProcessNamespace(true).
					Build()
				if err := cluster.CreateAndWaitForDeployment(ctx.Ctx, ctx.Client, dep1, cluster.DefaultTimeout); err != nil {
					return err
				}
				dep2 := builder.NewDeployment("test-dep-2", ctx.Namespace).
					WithShareProcessNamespace(true).
					WithCapability("SYS_PTRACE").
					Build()
				return cluster.CreateAndWaitForDeployment(ctx.Ctx, ctx.Client, dep2, cluster.DefaultTimeout)
			},
		},
	)
}
