package lifecycle

import (
	checks "github.com/redhat-best-practices-for-k8s/checks"
	"github.com/redhat-best-practices-for-k8s/checks-qe/pkg/builder"
	"github.com/redhat-best-practices-for-k8s/checks-qe/pkg/cluster"
	"github.com/redhat-best-practices-for-k8s/checks-qe/pkg/scenario"
	corev1 "k8s.io/api/core/v1"
)

func registerScheduling() {
	scenario.Register(
		scenario.Scenario{
			Name:           "lifecycle/pod-scheduling/compliant",
			CheckName:      "lifecycle-pod-scheduling",
			Category:       checks.CategoryLifecycle,
			Description:    "Deployment without nodeSelector should be compliant",
			ExpectedStatus: checks.StatusCompliant,
			Setup:          scenario.VanillaDeploymentSetup,
		},
		scenario.Scenario{
			Name:           "lifecycle/pod-scheduling/non-compliant",
			CheckName:      "lifecycle-pod-scheduling",
			Category:       checks.CategoryLifecycle,
			Description:    "Deployment with nodeSelector should be non-compliant",
			ExpectedStatus: checks.StatusNonCompliant,
			Setup: func(ctx *scenario.RunContext) error {
				dep := builder.NewDeployment("test-dep", ctx.Namespace).
					WithNodeSelector("kubernetes.io/os", "linux").
					Build()
				return cluster.CreateAndWaitForDeployment(ctx.Ctx, ctx.Client, dep, cluster.DefaultTimeout)
			},
		},
	)

	scenario.Register(
		scenario.Scenario{
			Name:           "lifecycle/toleration-bypass/compliant",
			CheckName:      "lifecycle-pod-toleration-bypass",
			Category:       checks.CategoryLifecycle,
			Description:    "Deployment with no custom tolerations should be compliant",
			ExpectedStatus: checks.StatusCompliant,
			Setup:          scenario.VanillaDeploymentSetup,
		},
		scenario.Scenario{
			Name:           "lifecycle/toleration-bypass/non-compliant",
			CheckName:      "lifecycle-pod-toleration-bypass",
			Category:       checks.CategoryLifecycle,
			Description:    "Deployment with custom toleration should be non-compliant",
			ExpectedStatus: checks.StatusNonCompliant,
			Setup: func(ctx *scenario.RunContext) error {
				dep := builder.NewDeployment("test-dep", ctx.Namespace).
					WithToleration("node-role.kubernetes.io/master", "", corev1.TaintEffectNoSchedule).
					Build()
				return cluster.CreateAndWaitForDeployment(ctx.Ctx, ctx.Client, dep, cluster.DefaultTimeout)
			},
		},
	)

	scenario.Register(
		scenario.Scenario{
			Name:           "lifecycle/topology-spread/compliant",
			CheckName:      "lifecycle-topology-spread-constraint",
			Category:       checks.CategoryLifecycle,
			Description:    "Deployment with both hostname and zone topology keys should be compliant",
			ExpectedStatus: checks.StatusCompliant,
			Setup: func(ctx *scenario.RunContext) error {
				dep := builder.NewDeployment("test-dep", ctx.Namespace).
					WithReplicas(1).
					WithTopologySpreadConstraint("kubernetes.io/hostname", 1, corev1.ScheduleAnyway).
					WithTopologySpreadConstraint("topology.kubernetes.io/zone", 1, corev1.ScheduleAnyway).
					Build()
				return cluster.CreateAndWaitForDeployment(ctx.Ctx, ctx.Client, dep, cluster.DefaultTimeout)
			},
		},
		scenario.Scenario{
			Name:           "lifecycle/topology-spread/non-compliant",
			CheckName:      "lifecycle-topology-spread-constraint",
			Category:       checks.CategoryLifecycle,
			Description:    "Deployment with only hostname topology key should be non-compliant",
			ExpectedStatus: checks.StatusNonCompliant,
			Setup: func(ctx *scenario.RunContext) error {
				dep := builder.NewDeployment("test-dep", ctx.Namespace).
					WithReplicas(1).
					WithTopologySpreadConstraint("kubernetes.io/hostname", 1, corev1.ScheduleAnyway).
					Build()
				return cluster.CreateAndWaitForDeployment(ctx.Ctx, ctx.Client, dep, cluster.DefaultTimeout)
			},
		},
	)
}
