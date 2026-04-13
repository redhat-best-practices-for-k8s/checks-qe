package performance

import (
	checks "github.com/redhat-best-practices-for-k8s/checks"
	"github.com/redhat-best-practices-for-k8s/checks-qe/pkg/builder"
	"github.com/redhat-best-practices-for-k8s/checks-qe/pkg/cluster"
	"github.com/redhat-best-practices-for-k8s/checks-qe/pkg/scenario"
	corev1 "k8s.io/api/core/v1"
)

func execProbeWithPeriod(period int32) *corev1.Probe {
	return &corev1.Probe{
		ProbeHandler:        corev1.ProbeHandler{Exec: &corev1.ExecAction{Command: []string{"true"}}},
		PeriodSeconds:       period,
		InitialDelaySeconds: 1,
	}
}

func registerExecProbes() {
	scenario.Register(
		scenario.Scenario{
			Name:           "performance/cpu-pinning-no-exec-probes/compliant",
			CheckName:      "performance-cpu-pinning-no-exec-probes",
			Category:       checks.CategoryPerformance,
			Description:    "CPU-pinned deployment without exec probes should be compliant",
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
			Name:           "performance/cpu-pinning-no-exec-probes/non-compliant",
			CheckName:      "performance-cpu-pinning-no-exec-probes",
			Category:       checks.CategoryPerformance,
			Description:    "CPU-pinned deployment with exec liveness probe should be non-compliant",
			ExpectedStatus: checks.StatusNonCompliant,
			Setup: func(ctx *scenario.RunContext) error {
				dep := builder.NewDeployment("test-dep", ctx.Namespace).
					WithResourceRequests("1", "128Mi").
					WithResourceLimits("1", "128Mi").
					WithLivenessProbe(execProbeWithPeriod(10)).
					Build()
				return cluster.CreateAndWaitForDeployment(ctx.Ctx, ctx.Client, dep, cluster.DefaultTimeout)
			},
		},
	)

	scenario.Register(
		scenario.Scenario{
			Name:           "performance/limited-use-of-exec-probes/compliant",
			CheckName:      "performance-limited-use-of-exec-probes",
			Category:       checks.CategoryPerformance,
			Description:    "Deployment without exec probes should be compliant",
			ExpectedStatus: checks.StatusCompliant,
			Setup:          scenario.VanillaDeploymentSetup,
		},
	)

	scenario.Register(
		scenario.Scenario{
			Name:           "performance/max-resources-exec-probes/compliant",
			CheckName:      "performance-max-resources-exec-probes",
			Category:       checks.CategoryPerformance,
			Description:    "Exec probe with periodSeconds >= 10 should be compliant",
			ExpectedStatus: checks.StatusCompliant,
			Setup: func(ctx *scenario.RunContext) error {
				dep := builder.NewDeployment("test-dep", ctx.Namespace).
					WithLivenessProbe(execProbeWithPeriod(15)).
					Build()
				return cluster.CreateAndWaitForDeployment(ctx.Ctx, ctx.Client, dep, cluster.DefaultTimeout)
			},
		},
		scenario.Scenario{
			Name:           "performance/max-resources-exec-probes/non-compliant",
			CheckName:      "performance-max-resources-exec-probes",
			Category:       checks.CategoryPerformance,
			Description:    "Exec probe with periodSeconds < 10 should be non-compliant",
			ExpectedStatus: checks.StatusNonCompliant,
			Setup: func(ctx *scenario.RunContext) error {
				dep := builder.NewDeployment("test-dep", ctx.Namespace).
					WithLivenessProbe(execProbeWithPeriod(5)).
					Build()
				return cluster.CreateAndWaitForDeployment(ctx.Ctx, ctx.Client, dep, cluster.DefaultTimeout)
			},
		},
	)

	scenario.Register(
		scenario.Scenario{
			Name:           "performance/rt-apps-no-exec-probes/compliant",
			CheckName:      "performance-rt-apps-no-exec-probes",
			Category:       checks.CategoryPerformance,
			Description:    "RT-annotated pod without exec probes should be compliant",
			ExpectedStatus: checks.StatusCompliant,
			Setup: func(ctx *scenario.RunContext) error {
				dep := builder.NewDeployment("test-dep", ctx.Namespace).
					WithPodAnnotation("rt-app", "true").
					Build()
				return cluster.CreateAndWaitForDeployment(ctx.Ctx, ctx.Client, dep, cluster.DefaultTimeout)
			},
		},
		scenario.Scenario{
			Name:           "performance/rt-apps-no-exec-probes/non-compliant",
			CheckName:      "performance-rt-apps-no-exec-probes",
			Category:       checks.CategoryPerformance,
			Description:    "RT-annotated pod with exec probe should be non-compliant",
			ExpectedStatus: checks.StatusNonCompliant,
			Setup: func(ctx *scenario.RunContext) error {
				dep := builder.NewDeployment("test-dep", ctx.Namespace).
					WithPodAnnotation("rt-app", "true").
					WithLivenessProbe(execProbeWithPeriod(10)).
					Build()
				return cluster.CreateAndWaitForDeployment(ctx.Ctx, ctx.Client, dep, cluster.DefaultTimeout)
			},
		},
	)
}
