package accesscontrol

import (
	checks "github.com/redhat-best-practices-for-k8s/checks"
	"github.com/redhat-best-practices-for-k8s/checks-qe/pkg/builder"
	"github.com/redhat-best-practices-for-k8s/checks-qe/pkg/cluster"
	"github.com/redhat-best-practices-for-k8s/checks-qe/pkg/scenario"
	corev1 "k8s.io/api/core/v1"
)

func init() {
	scenario.Register(
		scenario.Scenario{
			Name:           "accesscontrol/service-type/compliant",
			CheckName:      "access-control-service-type",
			Category:       checks.CategoryAccessControl,
			Description:    "ClusterIP service should be compliant",
			ExpectedStatus: checks.StatusCompliant,
			Setup: func(ctx *scenario.RunContext) error {
				dep := builder.NewDeployment("test-dep", ctx.Namespace).
					WithReplicas(1).
					WithContainerPort(8080).
					Build()
				if err := cluster.CreateAndWaitForDeployment(ctx.Ctx, ctx.Client, dep, cluster.DefaultTimeout); err != nil {
					return err
				}
				svc := builder.NewService("test-svc", ctx.Namespace).
					WithSelector("app", "test-dep").
					WithType(corev1.ServiceTypeClusterIP).
					Build()
				return cluster.CreateService(ctx.Ctx, ctx.Client, svc)
			},
		},
		scenario.Scenario{
			Name:           "accesscontrol/service-type/non-compliant",
			CheckName:      "access-control-service-type",
			Category:       checks.CategoryAccessControl,
			Description:    "NodePort service should be non-compliant",
			ExpectedStatus: checks.StatusNonCompliant,
			Setup: func(ctx *scenario.RunContext) error {
				dep := builder.NewDeployment("test-dep", ctx.Namespace).
					WithReplicas(1).
					WithContainerPort(8080).
					Build()
				if err := cluster.CreateAndWaitForDeployment(ctx.Ctx, ctx.Client, dep, cluster.DefaultTimeout); err != nil {
					return err
				}
				svc := builder.NewService("test-svc", ctx.Namespace).
					WithSelector("app", "test-dep").
					WithType(corev1.ServiceTypeNodePort).
					Build()
				return cluster.CreateService(ctx.Ctx, ctx.Client, svc)
			},
		},
		scenario.Scenario{
			Name:           "accesscontrol/service-type/compliant-multiple-services",
			CheckName:      "access-control-service-type",
			Category:       checks.CategoryAccessControl,
			Description:    "Multiple ClusterIP services should be compliant",
			ExpectedStatus: checks.StatusCompliant,
			Setup: func(ctx *scenario.RunContext) error {
				dep := builder.NewDeployment("test-dep", ctx.Namespace).WithContainerPort(8080).Build()
				if err := cluster.CreateAndWaitForDeployment(ctx.Ctx, ctx.Client, dep, cluster.DefaultTimeout); err != nil {
					return err
				}
				svc1 := builder.NewService("test-svc-1", ctx.Namespace).WithSelector("app", "test-dep").Build()
				if err := cluster.CreateService(ctx.Ctx, ctx.Client, svc1); err != nil {
					return err
				}
				svc2 := builder.NewService("test-svc-2", ctx.Namespace).WithSelector("app", "test-dep").Build()
				return cluster.CreateService(ctx.Ctx, ctx.Client, svc2)
			},
		},
		scenario.Scenario{
			Name:           "accesscontrol/service-type/mixed-services",
			CheckName:      "access-control-service-type",
			Category:       checks.CategoryAccessControl,
			Description:    "Multiple services, one NodePort should be non-compliant",
			ExpectedStatus: checks.StatusNonCompliant,
			Setup: func(ctx *scenario.RunContext) error {
				dep := builder.NewDeployment("test-dep", ctx.Namespace).WithContainerPort(8080).Build()
				if err := cluster.CreateAndWaitForDeployment(ctx.Ctx, ctx.Client, dep, cluster.DefaultTimeout); err != nil {
					return err
				}
				svc1 := builder.NewService("test-svc-1", ctx.Namespace).WithSelector("app", "test-dep").Build()
				if err := cluster.CreateService(ctx.Ctx, ctx.Client, svc1); err != nil {
					return err
				}
				svc2 := builder.NewService("test-svc-2", ctx.Namespace).
					WithSelector("app", "test-dep").
					WithType(corev1.ServiceTypeNodePort).
					Build()
				return cluster.CreateService(ctx.Ctx, ctx.Client, svc2)
			},
		},
	)
}
