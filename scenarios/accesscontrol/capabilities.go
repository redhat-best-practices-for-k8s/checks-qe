package accesscontrol

import (
	checks "github.com/redhat-best-practices-for-k8s/checks"
	"github.com/redhat-best-practices-for-k8s/checks-qe/pkg/builder"
	"github.com/redhat-best-practices-for-k8s/checks-qe/pkg/cluster"
	"github.com/redhat-best-practices-for-k8s/checks-qe/pkg/scenario"
	corev1 "k8s.io/api/core/v1"
)

func registerCapabilities() {
	caps := []struct {
		capability corev1.Capability
		checkName  string
	}{
		{"SYS_ADMIN", "access-control-sys-admin-capability-check"},
		{"NET_ADMIN", "access-control-net-admin-capability-check"},
		{"NET_RAW", "access-control-net-raw-capability-check"},
		{"IPC_LOCK", "access-control-ipc-lock-capability-check"},
		{"BPF", "access-control-bpf-capability-check"},
	}

	for _, c := range caps {
		cap := c // capture loop variable
		scenario.Register(
			scenario.Scenario{
				Name:           "accesscontrol/" + string(cap.capability) + "/compliant",
				CheckName:      cap.checkName,
				Category:       checks.CategoryAccessControl,
				Description:    "Deployment without " + string(cap.capability) + " capability should be compliant",
				ExpectedStatus: checks.StatusCompliant,
				Setup: func(ctx *scenario.RunContext) error {
					dep := builder.NewDeployment("test-dep", ctx.Namespace).Build()
					return cluster.CreateAndWaitForDeployment(ctx.Ctx, ctx.Client, dep, cluster.DefaultTimeout)
				},
			},
			scenario.Scenario{
				Name:           "accesscontrol/" + string(cap.capability) + "/non-compliant",
				CheckName:      cap.checkName,
				Category:       checks.CategoryAccessControl,
				Description:    "Deployment with " + string(cap.capability) + " capability should be non-compliant",
				ExpectedStatus: checks.StatusNonCompliant,
				Privileged:     true,
				Setup: func(ctx *scenario.RunContext) error {
					dep := builder.NewDeployment("test-dep", ctx.Namespace).
						WithCapability(cap.capability).
						Build()
					return cluster.CreateAndWaitForDeployment(ctx.Ctx, ctx.Client, dep, cluster.DefaultTimeout)
				},
			},
			scenario.Scenario{
				Name:           "accesscontrol/" + string(cap.capability) + "/mixed-two-deployments",
				CheckName:      cap.checkName,
				Category:       checks.CategoryAccessControl,
				Description:    "Two deployments, one with " + string(cap.capability) + " should be non-compliant",
				ExpectedStatus: checks.StatusNonCompliant,
				Privileged:     true,
				Setup: scenario.TwoDeploymentSetup(func(b *builder.DeploymentBuilder) *builder.DeploymentBuilder {
					return b.WithCapability(cap.capability)
				}),
			},
		)
	}
}
