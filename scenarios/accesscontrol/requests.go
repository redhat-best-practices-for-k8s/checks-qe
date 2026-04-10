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
			Name:           "accesscontrol/requests/compliant",
			CheckName:      "access-control-requests",
			Category:       checks.CategoryAccessControl,
			Description:    "Deployment with CPU and memory requests should be compliant",
			ExpectedStatus: checks.StatusCompliant,
			Setup: func(ctx *scenario.RunContext) error {
				dep := builder.NewDeployment("test-dep", ctx.Namespace).
					WithResourceRequests("100m", "128Mi").
					Build()
				return cluster.CreateAndWaitForDeployment(ctx.Ctx, ctx.Client, dep, cluster.DefaultTimeout)
			},
		},
		scenario.Scenario{
			Name:           "accesscontrol/requests/non-compliant-no-requests",
			CheckName:      "access-control-requests",
			Category:       checks.CategoryAccessControl,
			Description:    "Deployment with no resource requests should be non-compliant",
			ExpectedStatus: checks.StatusNonCompliant,
			Setup: func(ctx *scenario.RunContext) error {
				dep := builder.NewDeployment("test-dep", ctx.Namespace).Build()
				return cluster.CreateAndWaitForDeployment(ctx.Ctx, ctx.Client, dep, cluster.DefaultTimeout)
			},
		},
		scenario.Scenario{
			Name:           "accesscontrol/requests/non-compliant-no-memory",
			CheckName:      "access-control-requests",
			Category:       checks.CategoryAccessControl,
			Description:    "Deployment with CPU request only should be non-compliant",
			ExpectedStatus: checks.StatusNonCompliant,
			Setup: func(ctx *scenario.RunContext) error {
				dep := builder.NewDeployment("test-dep", ctx.Namespace).
					WithResourceRequests("100m", "").
					Build()
				return cluster.CreateAndWaitForDeployment(ctx.Ctx, ctx.Client, dep, cluster.DefaultTimeout)
			},
		},
		scenario.Scenario{
			Name:           "accesscontrol/requests/non-compliant-no-cpu",
			CheckName:      "access-control-requests",
			Category:       checks.CategoryAccessControl,
			Description:    "Deployment with memory request only should be non-compliant",
			ExpectedStatus: checks.StatusNonCompliant,
			Setup: func(ctx *scenario.RunContext) error {
				dep := builder.NewDeployment("test-dep", ctx.Namespace).
					WithResourceRequests("", "128Mi").
					Build()
				return cluster.CreateAndWaitForDeployment(ctx.Ctx, ctx.Client, dep, cluster.DefaultTimeout)
			},
		},
		scenario.Scenario{
			Name:           "accesscontrol/requests/compliant-two-deployments",
			CheckName:      "access-control-requests",
			Category:       checks.CategoryAccessControl,
			Description:    "Two deployments both with requests should be compliant",
			ExpectedStatus: checks.StatusCompliant,
			Setup: func(ctx *scenario.RunContext) error {
				dep1 := builder.NewDeployment("test-dep-1", ctx.Namespace).WithResourceRequests("100m", "128Mi").Build()
				if err := cluster.CreateAndWaitForDeployment(ctx.Ctx, ctx.Client, dep1, cluster.DefaultTimeout); err != nil {
					return err
				}
				dep2 := builder.NewDeployment("test-dep-2", ctx.Namespace).WithResourceRequests("200m", "256Mi").Build()
				return cluster.CreateAndWaitForDeployment(ctx.Ctx, ctx.Client, dep2, cluster.DefaultTimeout)
			},
		},
		scenario.Scenario{
			Name:           "accesscontrol/requests/mixed-two-deployments",
			CheckName:      "access-control-requests",
			Category:       checks.CategoryAccessControl,
			Description:    "Two deployments, one without memory request should be non-compliant",
			ExpectedStatus: checks.StatusNonCompliant,
			Setup: func(ctx *scenario.RunContext) error {
				dep1 := builder.NewDeployment("test-dep-1", ctx.Namespace).WithResourceRequests("100m", "128Mi").Build()
				if err := cluster.CreateAndWaitForDeployment(ctx.Ctx, ctx.Client, dep1, cluster.DefaultTimeout); err != nil {
					return err
				}
				dep2 := builder.NewDeployment("test-dep-2", ctx.Namespace).WithResourceRequests("100m", "").Build()
				return cluster.CreateAndWaitForDeployment(ctx.Ctx, ctx.Client, dep2, cluster.DefaultTimeout)
			},
		},
	)
}
