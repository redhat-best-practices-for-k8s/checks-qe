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
			Name:           "networking/network-policy-deny-all/compliant",
			CheckName:      "networking-network-policy-deny-all",
			Category:       checks.CategoryNetworking,
			Description:    "Pod with deny-all network policy should be compliant",
			ExpectedStatus: checks.StatusCompliant,
			Setup: func(ctx *scenario.RunContext) error {
				dep := builder.NewDeployment("test-dep", ctx.Namespace).Build()
				if err := cluster.CreateAndWaitForDeployment(ctx.Ctx, ctx.Client, dep, cluster.DefaultTimeout); err != nil {
					return err
				}
				np := builder.NewDenyAllNetworkPolicy("deny-all", ctx.Namespace, map[string]string{"app": "test-dep"})
				return cluster.CreateNetworkPolicy(ctx.Ctx, ctx.Client, np)
			},
		},
		scenario.Scenario{
			Name:           "networking/network-policy-deny-all/non-compliant",
			CheckName:      "networking-network-policy-deny-all",
			Category:       checks.CategoryNetworking,
			Description:    "Pod without network policy should be non-compliant",
			ExpectedStatus: checks.StatusNonCompliant,
			Setup:          scenario.VanillaDeploymentSetup,
		},
	)
}
