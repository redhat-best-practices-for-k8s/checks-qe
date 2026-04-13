package accesscontrol

import (
	checks "github.com/redhat-best-practices-for-k8s/checks"
	"github.com/redhat-best-practices-for-k8s/checks-qe/pkg/builder"
	"github.com/redhat-best-practices-for-k8s/checks-qe/pkg/cluster"
	"github.com/redhat-best-practices-for-k8s/checks-qe/pkg/scenario"
)

func registerResourceQuota() {
	scenario.Register(
		scenario.Scenario{
			Name:           "accesscontrol/resource-quota/compliant",
			CheckName:      "access-control-namespace-resource-quota",
			Category:       checks.CategoryAccessControl,
			Description:    "Namespace with resource quota should be compliant",
			ExpectedStatus: checks.StatusCompliant,
			Setup: func(ctx *scenario.RunContext) error {
				rq := builder.NewResourceQuota("test-quota", ctx.Namespace, "2", "2Gi")
				if err := cluster.CreateResourceQuota(ctx.Ctx, ctx.Client, rq); err != nil {
					return err
				}
				dep := builder.NewDeployment("test-dep", ctx.Namespace).
					WithResourceRequests("100m", "128Mi").
					WithResourceLimits("500m", "256Mi").
					Build()
				return cluster.CreateAndWaitForDeployment(ctx.Ctx, ctx.Client, dep, cluster.DefaultTimeout)
			},
		},
		scenario.Scenario{
			Name:           "accesscontrol/resource-quota/non-compliant",
			CheckName:      "access-control-namespace-resource-quota",
			Category:       checks.CategoryAccessControl,
			Description:    "Namespace without resource quota should be non-compliant",
			ExpectedStatus: checks.StatusNonCompliant,
			Setup: func(ctx *scenario.RunContext) error {
				dep := builder.NewDeployment("test-dep", ctx.Namespace).Build()
				return cluster.CreateAndWaitForDeployment(ctx.Ctx, ctx.Client, dep, cluster.DefaultTimeout)
			},
		},
		scenario.Scenario{
			Name:           "accesscontrol/resource-quota/compliant-two-deployments",
			CheckName:      "access-control-namespace-resource-quota",
			Category:       checks.CategoryAccessControl,
			Description:    "Two deployments both in namespace with quota should be compliant",
			ExpectedStatus: checks.StatusCompliant,
			Setup: func(ctx *scenario.RunContext) error {
				rq := builder.NewResourceQuota("test-quota", ctx.Namespace, "4", "4Gi")
				if err := cluster.CreateResourceQuota(ctx.Ctx, ctx.Client, rq); err != nil {
					return err
				}
				dep1 := builder.NewDeployment("test-dep-1", ctx.Namespace).
					WithResourceRequests("100m", "128Mi").WithResourceLimits("500m", "256Mi").Build()
				if err := cluster.CreateAndWaitForDeployment(ctx.Ctx, ctx.Client, dep1, cluster.DefaultTimeout); err != nil {
					return err
				}
				dep2 := builder.NewDeployment("test-dep-2", ctx.Namespace).
					WithResourceRequests("100m", "128Mi").WithResourceLimits("500m", "256Mi").Build()
				return cluster.CreateAndWaitForDeployment(ctx.Ctx, ctx.Client, dep2, cluster.DefaultTimeout)
			},
		},
	)
}
