package lifecycle

import (
	checks "github.com/redhat-best-practices-for-k8s/checks"
	"github.com/redhat-best-practices-for-k8s/checks-qe/pkg/builder"
	"github.com/redhat-best-practices-for-k8s/checks-qe/pkg/cluster"
	"github.com/redhat-best-practices-for-k8s/checks-qe/pkg/scenario"
	corev1 "k8s.io/api/core/v1"
)

var execProbe = &corev1.Probe{
	ProbeHandler: corev1.ProbeHandler{
		Exec: &corev1.ExecAction{Command: []string{"true"}},
	},
	InitialDelaySeconds: 1,
	PeriodSeconds:       5,
}

func init() {
	probeChecks := []struct {
		name      string
		checkName string
		withProbe func(b *builder.DeploymentBuilder) *builder.DeploymentBuilder
	}{
		{"liveness-probe", "lifecycle-liveness-probe",
			func(b *builder.DeploymentBuilder) *builder.DeploymentBuilder { return b.WithLivenessProbe(execProbe) }},
		{"readiness-probe", "lifecycle-readiness-probe",
			func(b *builder.DeploymentBuilder) *builder.DeploymentBuilder { return b.WithReadinessProbe(execProbe) }},
		{"startup-probe", "lifecycle-startup-probe",
			func(b *builder.DeploymentBuilder) *builder.DeploymentBuilder { return b.WithStartupProbe(execProbe) }},
	}

	for _, pc := range probeChecks {
		p := pc
		scenario.Register(
			scenario.Scenario{
				Name:           "lifecycle/" + p.name + "/compliant",
				CheckName:      p.checkName,
				Category:       checks.CategoryLifecycle,
				Description:    "Deployment with " + p.name + " should be compliant",
				ExpectedStatus: checks.StatusCompliant,
				Setup: func(ctx *scenario.RunContext) error {
					dep := p.withProbe(builder.NewDeployment("test-dep", ctx.Namespace)).Build()
					return cluster.CreateAndWaitForDeployment(ctx.Ctx, ctx.Client, dep, cluster.DefaultTimeout)
				},
			},
			scenario.Scenario{
				Name:           "lifecycle/" + p.name + "/non-compliant",
				CheckName:      p.checkName,
				Category:       checks.CategoryLifecycle,
				Description:    "Deployment without " + p.name + " should be non-compliant",
				ExpectedStatus: checks.StatusNonCompliant,
				Setup: scenario.VanillaDeploymentSetup,
			},
		)
	}
}
