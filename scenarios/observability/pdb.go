package observability

import (
	checks "github.com/redhat-best-practices-for-k8s/checks"
	"github.com/redhat-best-practices-for-k8s/checks-qe/pkg/builder"
	"github.com/redhat-best-practices-for-k8s/checks-qe/pkg/cluster"
	"github.com/redhat-best-practices-for-k8s/checks-qe/pkg/scenario"
)

func init() {
	scenario.Register(
		scenario.Scenario{
			Name:           "observability/pod-disruption-budget/compliant",
			CheckName:      "observability-pod-disruption-budget",
			Category:       checks.CategoryObservability,
			Description:    "Deployment with matching PDB should be compliant",
			ExpectedStatus: checks.StatusCompliant,
			Setup: func(ctx *scenario.RunContext) error {
				dep := builder.NewDeployment("test-dep", ctx.Namespace).
					WithReplicas(2).
					Build()
				if err := cluster.CreateAndWaitForDeployment(ctx.Ctx, ctx.Client, dep, cluster.DefaultTimeout); err != nil {
					return err
				}
				pdb := builder.NewPDB("test-pdb", ctx.Namespace).
					WithMinAvailable(1).
					WithSelector(map[string]string{"app": "test-dep"}).
					Build()
				return cluster.CreatePDB(ctx.Ctx, ctx.Client, pdb)
			},
		},
		scenario.Scenario{
			Name:           "observability/pod-disruption-budget/non-compliant",
			CheckName:      "observability-pod-disruption-budget",
			Category:       checks.CategoryObservability,
			Description:    "Deployment without PDB should be non-compliant",
			ExpectedStatus: checks.StatusNonCompliant,
			Setup: func(ctx *scenario.RunContext) error {
				dep := builder.NewDeployment("test-dep", ctx.Namespace).
					WithReplicas(2).
					Build()
				return cluster.CreateAndWaitForDeployment(ctx.Ctx, ctx.Client, dep, cluster.DefaultTimeout)
			},
		},
	)
}
