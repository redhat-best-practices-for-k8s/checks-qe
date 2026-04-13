package networking

import (
	checks "github.com/redhat-best-practices-for-k8s/checks"
	"github.com/redhat-best-practices-for-k8s/checks-qe/pkg/builder"
	"github.com/redhat-best-practices-for-k8s/checks-qe/pkg/cluster"
	"github.com/redhat-best-practices-for-k8s/checks-qe/pkg/scenario"
	corev1 "k8s.io/api/core/v1"
)

func registerNetworkServices() {
	scenario.Register(
		scenario.Scenario{
			Name:              "networking/dual-stack-service/compliant",
			CheckName:         "networking-dual-stack-service",
			Category:          checks.CategoryNetworking,
			Description:       "Service with dual-stack IP family policy should be compliant",
			ExpectedStatus:    checks.StatusCompliant,
			RequiresDualStack: true,
			Tags:              []string{"dual-stack"},
			Setup: func(ctx *scenario.RunContext) error {
				dep := builder.NewDeployment("test-dep", ctx.Namespace).
					WithContainerPort(8080).
					Build()
				if err := cluster.CreateAndWaitForDeployment(ctx.Ctx, ctx.Client, dep, cluster.DefaultTimeout); err != nil {
					return err
				}
				svc := builder.NewService("test-svc", ctx.Namespace).
					WithSelector("app", "test-dep").
					WithIPFamilyPolicy(corev1.IPFamilyPolicyPreferDualStack).
					Build()
				return cluster.CreateService(ctx.Ctx, ctx.Client, svc)
			},
		},
		scenario.Scenario{
			Name:           "networking/dual-stack-service/non-compliant",
			CheckName:      "networking-dual-stack-service",
			Category:       checks.CategoryNetworking,
			Description:    "Service with single-stack IP family policy should be non-compliant",
			ExpectedStatus: checks.StatusNonCompliant,
			Setup: func(ctx *scenario.RunContext) error {
				dep := builder.NewDeployment("test-dep", ctx.Namespace).
					WithContainerPort(8080).
					Build()
				if err := cluster.CreateAndWaitForDeployment(ctx.Ctx, ctx.Client, dep, cluster.DefaultTimeout); err != nil {
					return err
				}
				svc := builder.NewService("test-svc", ctx.Namespace).
					WithSelector("app", "test-dep").
					WithIPFamilyPolicy(corev1.IPFamilyPolicySingleStack).
					Build()
				return cluster.CreateService(ctx.Ctx, ctx.Client, svc)
			},
		},
	)
}
