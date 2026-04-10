package lifecycle

import (
	checks "github.com/redhat-best-practices-for-k8s/checks"
	"github.com/redhat-best-practices-for-k8s/checks-qe/pkg/builder"
	"github.com/redhat-best-practices-for-k8s/checks-qe/pkg/cluster"
	"github.com/redhat-best-practices-for-k8s/checks-qe/pkg/scenario"
	corev1 "k8s.io/api/core/v1"
)

var execHandler = &corev1.LifecycleHandler{
	Exec: &corev1.ExecAction{Command: []string{"true"}},
}

func init() {
	hookChecks := []struct {
		name      string
		checkName string
		withHook  func(b *builder.DeploymentBuilder) *builder.DeploymentBuilder
	}{
		{"poststart", "lifecycle-container-poststart",
			func(b *builder.DeploymentBuilder) *builder.DeploymentBuilder { return b.WithPostStartHook(execHandler) }},
		{"prestop", "lifecycle-container-prestop",
			func(b *builder.DeploymentBuilder) *builder.DeploymentBuilder { return b.WithPreStopHook(execHandler) }},
	}

	for _, hc := range hookChecks {
		h := hc
		scenario.Register(
			scenario.Scenario{
				Name:           "lifecycle/" + h.name + "/compliant",
				CheckName:      h.checkName,
				Category:       checks.CategoryLifecycle,
				Description:    "Deployment with " + h.name + " hook should be compliant",
				ExpectedStatus: checks.StatusCompliant,
				Setup: func(ctx *scenario.RunContext) error {
					dep := h.withHook(builder.NewDeployment("test-dep", ctx.Namespace)).Build()
					return cluster.CreateAndWaitForDeployment(ctx.Ctx, ctx.Client, dep, cluster.DefaultTimeout)
				},
			},
			scenario.Scenario{
				Name:           "lifecycle/" + h.name + "/non-compliant",
				CheckName:      h.checkName,
				Category:       checks.CategoryLifecycle,
				Description:    "Deployment without " + h.name + " hook should be non-compliant",
				ExpectedStatus: checks.StatusNonCompliant,
				Setup: scenario.VanillaDeploymentSetup,
			},
		)
	}
}
