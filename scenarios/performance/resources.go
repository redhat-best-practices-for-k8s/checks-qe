package performance

import (
	checks "github.com/redhat-best-practices-for-k8s/checks"
	"github.com/redhat-best-practices-for-k8s/checks-qe/pkg/builder"
	"github.com/redhat-best-practices-for-k8s/checks-qe/pkg/cluster"
	"github.com/redhat-best-practices-for-k8s/checks-qe/pkg/scenario"
)

func init() {
	scenario.Register(
		scenario.Scenario{
			Name:           "performance/exclusive-cpu-pool/compliant",
			CheckName:      "performance-exclusive-cpu-pool",
			Category:       checks.CategoryPerformance,
			Description:    "Deployment with whole CPU requests matching limits should be compliant",
			ExpectedStatus: checks.StatusCompliant,
			Setup: func(ctx *scenario.RunContext) error {
				dep := builder.NewDeployment("test-dep", ctx.Namespace).
					WithResourceRequests("1", "128Mi").
					WithResourceLimits("1", "128Mi").
					Build()
				return cluster.CreateAndWaitForDeployment(ctx.Ctx, ctx.Client, dep, cluster.DefaultTimeout)
			},
		},
		scenario.Scenario{
			Name:           "performance/exclusive-cpu-pool/non-compliant",
			CheckName:      "performance-exclusive-cpu-pool",
			Category:       checks.CategoryPerformance,
			Description:    "Deployment with whole CPU requests not matching limits should be non-compliant",
			ExpectedStatus: checks.StatusNonCompliant,
			Setup: func(ctx *scenario.RunContext) error {
				dep := builder.NewDeployment("test-dep", ctx.Namespace).
					WithResourceRequests("1", "128Mi").
					WithResourceLimits("2", "128Mi").
					Build()
				return cluster.CreateAndWaitForDeployment(ctx.Ctx, ctx.Client, dep, cluster.DefaultTimeout)
			},
		},
	)

	scenario.Register(
		scenario.Scenario{
			Name:           "performance/limit-memory-allocation/compliant",
			CheckName:      "performance-limit-memory-allocation",
			Category:       checks.CategoryPerformance,
			Description:    "Deployment with memory limits should be compliant",
			ExpectedStatus: checks.StatusCompliant,
			Setup: func(ctx *scenario.RunContext) error {
				dep := builder.NewDeployment("test-dep", ctx.Namespace).
					WithResourceLimits("100m", "128Mi").
					Build()
				return cluster.CreateAndWaitForDeployment(ctx.Ctx, ctx.Client, dep, cluster.DefaultTimeout)
			},
		},
		scenario.Scenario{
			Name:           "performance/limit-memory-allocation/non-compliant",
			CheckName:      "performance-limit-memory-allocation",
			Category:       checks.CategoryPerformance,
			Description:    "Deployment without memory limits should be non-compliant",
			ExpectedStatus: checks.StatusNonCompliant,
			Setup:          scenario.VanillaDeploymentSetup,
		},
	)
}
