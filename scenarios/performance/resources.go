package performance

import (
	checks "github.com/redhat-best-practices-for-k8s/checks"
	"github.com/redhat-best-practices-for-k8s/checks-qe/pkg/builder"
	"github.com/redhat-best-practices-for-k8s/checks-qe/pkg/cluster"
	"github.com/redhat-best-practices-for-k8s/checks-qe/pkg/scenario"
	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/resource"
)

func registerResources() {
	scenario.Register(
		scenario.Scenario{
			Name:           "performance/exclusive-cpu-pool/compliant",
			CheckName:      "performance-exclusive-cpu-pool",
			Category:       checks.CategoryPerformance,
			Description:    "Pod with all containers in the same CPU pool should be compliant",
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
			Description:    "Pod mixing exclusive and shared CPU pool containers should be non-compliant",
			ExpectedStatus: checks.StatusNonCompliant,
			Setup: func(ctx *scenario.RunContext) error {
				dep := builder.NewDeployment("test-dep", ctx.Namespace).
					WithResourceRequests("1", "128Mi").
					WithResourceLimits("1", "128Mi").
					Build()
				dep.Spec.Template.Spec.Containers = append(dep.Spec.Template.Spec.Containers,
					corev1.Container{
						Name:    "shared",
						Image:   builder.DefaultImage,
						Command: []string{"/bin/sh", "-c", "sleep infinity"},
						Resources: corev1.ResourceRequirements{
							Requests: corev1.ResourceList{
								corev1.ResourceCPU:    resource.MustParse("100m"),
								corev1.ResourceMemory: resource.MustParse("64Mi"),
							},
							Limits: corev1.ResourceList{
								corev1.ResourceCPU:    resource.MustParse("500m"),
								corev1.ResourceMemory: resource.MustParse("128Mi"),
							},
						},
					},
				)
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
