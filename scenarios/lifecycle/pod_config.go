package lifecycle

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
			Name:           "lifecycle/image-pull-policy/compliant",
			CheckName:      "lifecycle-image-pull-policy",
			Category:       checks.CategoryLifecycle,
			Description:    "Deployment with IfNotPresent pull policy should be compliant",
			ExpectedStatus: checks.StatusCompliant,
			Setup: func(ctx *scenario.RunContext) error {
				dep := builder.NewDeployment("test-dep", ctx.Namespace).
					WithImagePullPolicy(corev1.PullIfNotPresent).
					Build()
				return cluster.CreateAndWaitForDeployment(ctx.Ctx, ctx.Client, dep, cluster.DefaultTimeout)
			},
		},
		scenario.Scenario{
			Name:           "lifecycle/image-pull-policy/non-compliant",
			CheckName:      "lifecycle-image-pull-policy",
			Category:       checks.CategoryLifecycle,
			Description:    "Deployment with Always pull policy should be non-compliant",
			ExpectedStatus: checks.StatusNonCompliant,
			Setup: func(ctx *scenario.RunContext) error {
				dep := builder.NewDeployment("test-dep", ctx.Namespace).
					WithImagePullPolicy(corev1.PullAlways).
					Build()
				return cluster.CreateAndWaitForDeployment(ctx.Ctx, ctx.Client, dep, cluster.DefaultTimeout)
			},
		},
	)

	scenario.Register(
		scenario.Scenario{
			Name:           "lifecycle/pod-owner-type/compliant",
			CheckName:      "lifecycle-pod-owner-type",
			Category:       checks.CategoryLifecycle,
			Description:    "Pod owned by deployment should be compliant",
			ExpectedStatus: checks.StatusCompliant,
			Setup:          scenario.VanillaDeploymentSetup,
		},
		scenario.Scenario{
			Name:           "lifecycle/pod-owner-type/non-compliant",
			CheckName:      "lifecycle-pod-owner-type",
			Category:       checks.CategoryLifecycle,
			Description:    "Bare pod without owner should be non-compliant",
			ExpectedStatus: checks.StatusNonCompliant,
			Setup: func(ctx *scenario.RunContext) error {
				pod := builder.NewPod("bare-pod", ctx.Namespace).Build()
				return cluster.CreateAndWaitForPod(ctx.Ctx, ctx.Client, pod, cluster.DefaultTimeout)
			},
		},
	)

	scenario.Register(
		scenario.Scenario{
			Name:           "lifecycle/cpu-isolation/compliant",
			CheckName:      "lifecycle-cpu-isolation",
			Category:       checks.CategoryLifecycle,
			Description:    "Deployment with CPU requests == limits should be compliant",
			ExpectedStatus: checks.StatusCompliant,
			Setup: func(ctx *scenario.RunContext) error {
				dep := builder.NewDeployment("test-dep", ctx.Namespace).
					WithResourceRequests("100m", "128Mi").
					WithResourceLimits("100m", "128Mi").
					Build()
				return cluster.CreateAndWaitForDeployment(ctx.Ctx, ctx.Client, dep, cluster.DefaultTimeout)
			},
		},
		scenario.Scenario{
			Name:           "lifecycle/cpu-isolation/non-compliant",
			CheckName:      "lifecycle-cpu-isolation",
			Category:       checks.CategoryLifecycle,
			Description:    "Deployment with CPU requests != limits should be non-compliant",
			ExpectedStatus: checks.StatusNonCompliant,
			Setup: func(ctx *scenario.RunContext) error {
				dep := builder.NewDeployment("test-dep", ctx.Namespace).
					WithResourceRequests("100m", "128Mi").
					WithResourceLimits("500m", "256Mi").
					Build()
				return cluster.CreateAndWaitForDeployment(ctx.Ctx, ctx.Client, dep, cluster.DefaultTimeout)
			},
		},
	)

	scenario.Register(
		scenario.Scenario{
			Name:           "lifecycle/pod-high-availability/compliant",
			CheckName:      "lifecycle-pod-high-availability",
			Category:       checks.CategoryLifecycle,
			Description:    "Deployment with 2+ replicas should be compliant",
			ExpectedStatus: checks.StatusCompliant,
			Setup: func(ctx *scenario.RunContext) error {
				dep := builder.NewDeployment("test-dep", ctx.Namespace).
					WithReplicas(2).
					Build()
				return cluster.CreateAndWaitForDeployment(ctx.Ctx, ctx.Client, dep, cluster.DefaultTimeout)
			},
		},
		scenario.Scenario{
			Name:           "lifecycle/pod-high-availability/non-compliant",
			CheckName:      "lifecycle-pod-high-availability",
			Category:       checks.CategoryLifecycle,
			Description:    "Deployment with 1 replica should be non-compliant",
			ExpectedStatus: checks.StatusNonCompliant,
			Setup: func(ctx *scenario.RunContext) error {
				dep := builder.NewDeployment("test-dep", ctx.Namespace).
					WithReplicas(1).
					Build()
				return cluster.CreateAndWaitForDeployment(ctx.Ctx, ctx.Client, dep, cluster.DefaultTimeout)
			},
		},
	)
}
