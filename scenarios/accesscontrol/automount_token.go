package accesscontrol

import (
	checks "github.com/redhat-best-practices-for-k8s/checks"
	"github.com/redhat-best-practices-for-k8s/checks-qe/pkg/builder"
	"github.com/redhat-best-practices-for-k8s/checks-qe/pkg/cluster"
	"github.com/redhat-best-practices-for-k8s/checks-qe/pkg/scenario"
)

func registerAutomountToken() {
	scenario.Register(
		scenario.Scenario{
			Name:           "accesscontrol/automount-token/compliant",
			CheckName:      "access-control-pod-automount-service-account-token",
			Category:       checks.CategoryAccessControl,
			Description:    "Deployment with automountServiceAccountToken=false and custom SA should be compliant",
			ExpectedStatus: checks.StatusCompliant,
			Setup: func(ctx *scenario.RunContext) error {
				sa := builder.NewServiceAccount("test-sa", ctx.Namespace)
				if err := cluster.CreateServiceAccount(ctx.Ctx, ctx.Client, sa); err != nil {
					return err
				}
				dep := builder.NewDeployment("test-dep", ctx.Namespace).
					WithServiceAccountName("test-sa").
					WithAutomountServiceAccountToken(false).
					Build()
				return cluster.CreateAndWaitForDeployment(ctx.Ctx, ctx.Client, dep, cluster.DefaultTimeout)
			},
		},
		scenario.Scenario{
			Name:           "accesscontrol/automount-token/non-compliant",
			CheckName:      "access-control-pod-automount-service-account-token",
			Category:       checks.CategoryAccessControl,
			Description:    "Deployment with automountServiceAccountToken=true should be non-compliant",
			ExpectedStatus: checks.StatusNonCompliant,
			Setup: func(ctx *scenario.RunContext) error {
				dep := builder.NewDeployment("test-dep", ctx.Namespace).
					WithAutomountServiceAccountToken(true).
					Build()
				return cluster.CreateAndWaitForDeployment(ctx.Ctx, ctx.Client, dep, cluster.DefaultTimeout)
			},
		},
		scenario.Scenario{
			Name:           "accesscontrol/automount-token/compliant-sa-false",
			CheckName:      "access-control-pod-automount-service-account-token",
			Category:       checks.CategoryAccessControl,
			Description:    "SA with automount=false, pod unset should be compliant",
			ExpectedStatus: checks.StatusCompliant,
			Setup: func(ctx *scenario.RunContext) error {
				sa := builder.NewServiceAccountWithToken("test-sa", ctx.Namespace, false)
				if err := cluster.CreateServiceAccount(ctx.Ctx, ctx.Client, sa); err != nil {
					return err
				}
				dep := builder.NewDeployment("test-dep", ctx.Namespace).
					WithServiceAccountName("test-sa").
					Build()
				return cluster.CreateAndWaitForDeployment(ctx.Ctx, ctx.Client, dep, cluster.DefaultTimeout)
			},
		},
		scenario.Scenario{
			Name:           "accesscontrol/automount-token/non-compliant-sa-true",
			CheckName:      "access-control-pod-automount-service-account-token",
			Category:       checks.CategoryAccessControl,
			Description:    "SA with automount=true, pod unset should be non-compliant",
			ExpectedStatus: checks.StatusNonCompliant,
			Setup: func(ctx *scenario.RunContext) error {
				sa := builder.NewServiceAccountWithToken("test-sa", ctx.Namespace, true)
				if err := cluster.CreateServiceAccount(ctx.Ctx, ctx.Client, sa); err != nil {
					return err
				}
				dep := builder.NewDeployment("test-dep", ctx.Namespace).
					WithServiceAccountName("test-sa").
					Build()
				return cluster.CreateAndWaitForDeployment(ctx.Ctx, ctx.Client, dep, cluster.DefaultTimeout)
			},
		},
		scenario.Scenario{
			Name:           "accesscontrol/automount-token/compliant-pod-overrides-sa",
			CheckName:      "access-control-pod-automount-service-account-token",
			Category:       checks.CategoryAccessControl,
			Description:    "Pod automount=false overrides SA automount=true should be compliant",
			ExpectedStatus: checks.StatusCompliant,
			Setup: func(ctx *scenario.RunContext) error {
				sa := builder.NewServiceAccountWithToken("test-sa", ctx.Namespace, true)
				if err := cluster.CreateServiceAccount(ctx.Ctx, ctx.Client, sa); err != nil {
					return err
				}
				dep := builder.NewDeployment("test-dep", ctx.Namespace).
					WithServiceAccountName("test-sa").
					WithAutomountServiceAccountToken(false).
					Build()
				return cluster.CreateAndWaitForDeployment(ctx.Ctx, ctx.Client, dep, cluster.DefaultTimeout)
			},
		},
		scenario.Scenario{
			Name:           "accesscontrol/automount-token/compliant-two-deployments",
			CheckName:      "access-control-pod-automount-service-account-token",
			Category:       checks.CategoryAccessControl,
			Description:    "Two deployments both with automount=false should be compliant",
			ExpectedStatus: checks.StatusCompliant,
			Setup: func(ctx *scenario.RunContext) error {
				sa := builder.NewServiceAccount("test-sa", ctx.Namespace)
				if err := cluster.CreateServiceAccount(ctx.Ctx, ctx.Client, sa); err != nil {
					return err
				}
				dep1 := builder.NewDeployment("test-dep-1", ctx.Namespace).
					WithServiceAccountName("test-sa").WithAutomountServiceAccountToken(false).Build()
				if err := cluster.CreateAndWaitForDeployment(ctx.Ctx, ctx.Client, dep1, cluster.DefaultTimeout); err != nil {
					return err
				}
				dep2 := builder.NewDeployment("test-dep-2", ctx.Namespace).
					WithServiceAccountName("test-sa").WithAutomountServiceAccountToken(false).Build()
				return cluster.CreateAndWaitForDeployment(ctx.Ctx, ctx.Client, dep2, cluster.DefaultTimeout)
			},
		},
		scenario.Scenario{
			Name:           "accesscontrol/automount-token/mixed-two-deployments",
			CheckName:      "access-control-pod-automount-service-account-token",
			Category:       checks.CategoryAccessControl,
			Description:    "Two deployments, one with automount=true should be non-compliant",
			ExpectedStatus: checks.StatusNonCompliant,
			Setup: func(ctx *scenario.RunContext) error {
				sa := builder.NewServiceAccount("test-sa", ctx.Namespace)
				if err := cluster.CreateServiceAccount(ctx.Ctx, ctx.Client, sa); err != nil {
					return err
				}
				dep1 := builder.NewDeployment("test-dep-1", ctx.Namespace).
					WithServiceAccountName("test-sa").WithAutomountServiceAccountToken(false).Build()
				if err := cluster.CreateAndWaitForDeployment(ctx.Ctx, ctx.Client, dep1, cluster.DefaultTimeout); err != nil {
					return err
				}
				dep2 := builder.NewDeployment("test-dep-2", ctx.Namespace).
					WithAutomountServiceAccountToken(true).Build()
				return cluster.CreateAndWaitForDeployment(ctx.Ctx, ctx.Client, dep2, cluster.DefaultTimeout)
			},
		},
	)
}
