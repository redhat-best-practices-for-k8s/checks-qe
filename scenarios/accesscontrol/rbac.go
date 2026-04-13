package accesscontrol

import (
	checks "github.com/redhat-best-practices-for-k8s/checks"
	"github.com/redhat-best-practices-for-k8s/checks-qe/pkg/builder"
	"github.com/redhat-best-practices-for-k8s/checks-qe/pkg/cluster"
	"github.com/redhat-best-practices-for-k8s/checks-qe/pkg/scenario"
	rbacv1 "k8s.io/api/rbac/v1"
)

func registerRBAC() {
	// cluster-role-bindings
	scenario.Register(
		scenario.Scenario{
			Name:           "accesscontrol/cluster-role-bindings/compliant",
			CheckName:      "access-control-cluster-role-bindings",
			Category:       checks.CategoryAccessControl,
			Description:    "Deployment without cluster role binding should be compliant",
			ExpectedStatus: checks.StatusCompliant,
			Setup: func(ctx *scenario.RunContext) error {
				sa := builder.NewServiceAccount("test-sa", ctx.Namespace)
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
			Name:           "accesscontrol/cluster-role-bindings/non-compliant",
			CheckName:      "access-control-cluster-role-bindings",
			Category:       checks.CategoryAccessControl,
			Description:    "Deployment with cluster role binding should be non-compliant",
			ExpectedStatus: checks.StatusNonCompliant,
			Setup: func(ctx *scenario.RunContext) error {
				sa := builder.NewServiceAccount("test-sa", ctx.Namespace)
				if err := cluster.CreateServiceAccount(ctx.Ctx, ctx.Client, sa); err != nil {
					return err
				}
				crb := builder.NewClusterRoleBinding("cqe-test-crb-"+ctx.Namespace, "test-sa", ctx.Namespace, "view")
				if err := cluster.CreateClusterRoleBinding(ctx.Ctx, ctx.Client, crb); err != nil {
					return err
				}
				dep := builder.NewDeployment("test-dep", ctx.Namespace).
					WithServiceAccountName("test-sa").
					Build()
				return cluster.CreateAndWaitForDeployment(ctx.Ctx, ctx.Client, dep, cluster.DefaultTimeout)
			},
		},
	)

	// pod-role-bindings
	scenario.Register(
		scenario.Scenario{
			Name:           "accesscontrol/pod-role-bindings/compliant",
			CheckName:      "access-control-pod-role-bindings",
			Category:       checks.CategoryAccessControl,
			Description:    "Pod with valid role binding should be compliant",
			ExpectedStatus: checks.StatusCompliant,
			Setup: func(ctx *scenario.RunContext) error {
				sa := builder.NewServiceAccount("test-sa", ctx.Namespace)
				if err := cluster.CreateServiceAccount(ctx.Ctx, ctx.Client, sa); err != nil {
					return err
				}
				role := builder.NewRole("test-role", ctx.Namespace, []rbacv1.PolicyRule{{
					APIGroups: []string{""},
					Resources: []string{"pods"},
					Verbs:     []string{"get", "list"},
				}})
				if err := cluster.CreateRole(ctx.Ctx, ctx.Client, role); err != nil {
					return err
				}
				rb := builder.NewRoleBinding("test-rb", ctx.Namespace, "test-sa", "test-role")
				if err := cluster.CreateRoleBinding(ctx.Ctx, ctx.Client, rb); err != nil {
					return err
				}
				pod := builder.NewPod("test-pod", ctx.Namespace).Build()
				pod.Spec.ServiceAccountName = "test-sa"
				return cluster.CreateAndWaitForPod(ctx.Ctx, ctx.Client, pod, cluster.DefaultTimeout)
			},
		},
		scenario.Scenario{
			Name:           "accesscontrol/pod-role-bindings/non-compliant",
			CheckName:      "access-control-pod-role-bindings",
			Category:       checks.CategoryAccessControl,
			Description:    "Pod using default service account should be non-compliant",
			ExpectedStatus: checks.StatusNonCompliant,
			Setup: func(ctx *scenario.RunContext) error {
				pod := builder.NewPod("test-pod", ctx.Namespace).Build()
				return cluster.CreateAndWaitForPod(ctx.Ctx, ctx.Client, pod, cluster.DefaultTimeout)
			},
		},
	)

	// pod-service-account
	scenario.Register(
		scenario.Scenario{
			Name:           "accesscontrol/pod-service-account/compliant",
			CheckName:      "access-control-pod-service-account",
			Category:       checks.CategoryAccessControl,
			Description:    "Pod with custom service account should be compliant",
			ExpectedStatus: checks.StatusCompliant,
			Privileged:     true,
			Setup: func(ctx *scenario.RunContext) error {
				sa := builder.NewServiceAccount("test-sa", ctx.Namespace)
				if err := cluster.CreateServiceAccount(ctx.Ctx, ctx.Client, sa); err != nil {
					return err
				}
				pod := builder.NewPod("test-pod", ctx.Namespace).Build()
				pod.Spec.ServiceAccountName = "test-sa"
				return cluster.CreateAndWaitForPod(ctx.Ctx, ctx.Client, pod, cluster.DefaultTimeout)
			},
		},
		scenario.Scenario{
			Name:           "accesscontrol/pod-service-account/non-compliant",
			CheckName:      "access-control-pod-service-account",
			Category:       checks.CategoryAccessControl,
			Description:    "Pod with default service account should be non-compliant",
			ExpectedStatus: checks.StatusNonCompliant,
			Privileged:     true,
			Setup: func(ctx *scenario.RunContext) error {
				pod := builder.NewPod("test-pod", ctx.Namespace).Build()
				return cluster.CreateAndWaitForPod(ctx.Ctx, ctx.Client, pod, cluster.DefaultTimeout)
			},
		},
	)
}
