package accesscontrol

import (
	checks "github.com/redhat-best-practices-for-k8s/checks"
	"github.com/redhat-best-practices-for-k8s/checks-qe/pkg/builder"
	"github.com/redhat-best-practices-for-k8s/checks-qe/pkg/cluster"
	"github.com/redhat-best-practices-for-k8s/checks-qe/pkg/scenario"
)

func registerSecurityContext() {
	scenario.Register(
		scenario.Scenario{
			Name:           "accesscontrol/no-1337-uid/compliant",
			CheckName:      "access-control-no-1337-uid",
			Category:       checks.CategoryAccessControl,
			Description:    "Deployment with non-1337 UID should be compliant",
			ExpectedStatus: checks.StatusCompliant,
			Privileged:     true,
			Setup: func(ctx *scenario.RunContext) error {
				dep := builder.NewDeployment("test-dep", ctx.Namespace).
					WithPodRunAsUser(1000).
					Build()
				return cluster.CreateAndWaitForDeployment(ctx.Ctx, ctx.Client, dep, cluster.DefaultTimeout)
			},
		},
		scenario.Scenario{
			Name:           "accesscontrol/no-1337-uid/non-compliant",
			CheckName:      "access-control-no-1337-uid",
			Category:       checks.CategoryAccessControl,
			Description:    "Deployment with UID 1337 should be non-compliant",
			ExpectedStatus: checks.StatusNonCompliant,
			Privileged:     true,
			Setup: func(ctx *scenario.RunContext) error {
				dep := builder.NewDeployment("test-dep", ctx.Namespace).
					WithPodRunAsUser(1337).
					Build()
				return cluster.CreateAndWaitForDeployment(ctx.Ctx, ctx.Client, dep, cluster.DefaultTimeout)
			},
		},
		scenario.Scenario{
			Name:           "accesscontrol/no-1337-uid/mixed-two-deployments",
			CheckName:      "access-control-no-1337-uid",
			Category:       checks.CategoryAccessControl,
			Description:    "Two deployments, one with UID 1337 should be non-compliant",
			ExpectedStatus: checks.StatusNonCompliant,
			Privileged:     true,
			Setup: scenario.TwoDeploymentSetup(func(b *builder.DeploymentBuilder) *builder.DeploymentBuilder {
				return b.WithPodRunAsUser(1337)
			}),
		},
	)

	scenario.Register(
		scenario.Scenario{
			Name:           "accesscontrol/privilege-escalation/compliant",
			CheckName:      "access-control-security-context-privilege-escalation",
			Category:       checks.CategoryAccessControl,
			Description:    "Deployment with allowPrivilegeEscalation=false should be compliant",
			ExpectedStatus: checks.StatusCompliant,
			Privileged:     true,
			Setup: func(ctx *scenario.RunContext) error {
				dep := builder.NewDeployment("test-dep", ctx.Namespace).
					WithAllowPrivilegeEscalation(false).
					Build()
				return cluster.CreateAndWaitForDeployment(ctx.Ctx, ctx.Client, dep, cluster.DefaultTimeout)
			},
		},
		scenario.Scenario{
			Name:           "accesscontrol/privilege-escalation/non-compliant",
			CheckName:      "access-control-security-context-privilege-escalation",
			Category:       checks.CategoryAccessControl,
			Description:    "Deployment with allowPrivilegeEscalation=true should be non-compliant",
			ExpectedStatus: checks.StatusNonCompliant,
			Privileged:     true,
			Setup: func(ctx *scenario.RunContext) error {
				dep := builder.NewDeployment("test-dep", ctx.Namespace).
					WithAllowPrivilegeEscalation(true).
					Build()
				return cluster.CreateAndWaitForDeployment(ctx.Ctx, ctx.Client, dep, cluster.DefaultTimeout)
			},
		},
		scenario.Scenario{
			Name:           "accesscontrol/privilege-escalation/mixed-two-deployments",
			CheckName:      "access-control-security-context-privilege-escalation",
			Category:       checks.CategoryAccessControl,
			Description:    "Two deployments, one with privilege escalation should be non-compliant",
			ExpectedStatus: checks.StatusNonCompliant,
			Privileged:     true,
			Setup: scenario.TwoDeploymentSetup(func(b *builder.DeploymentBuilder) *builder.DeploymentBuilder {
				return b.WithAllowPrivilegeEscalation(true)
			}),
		},
	)

	scenario.Register(
		scenario.Scenario{
			Name:           "accesscontrol/read-only-fs/compliant",
			CheckName:      "access-control-security-context-read-only-file-system",
			Category:       checks.CategoryAccessControl,
			Description:    "Deployment with readOnlyRootFilesystem=true should be compliant",
			ExpectedStatus: checks.StatusCompliant,
			Setup: func(ctx *scenario.RunContext) error {
				dep := builder.NewDeployment("test-dep", ctx.Namespace).
					WithReadOnlyRootFS().
					Build()
				return cluster.CreateAndWaitForDeployment(ctx.Ctx, ctx.Client, dep, cluster.DefaultTimeout)
			},
		},
		scenario.Scenario{
			Name:           "accesscontrol/read-only-fs/non-compliant",
			CheckName:      "access-control-security-context-read-only-file-system",
			Category:       checks.CategoryAccessControl,
			Description:    "Deployment without readOnlyRootFilesystem should be non-compliant",
			ExpectedStatus: checks.StatusNonCompliant,
			Setup:          scenario.VanillaDeploymentSetup,
		},
		scenario.Scenario{
			Name:           "accesscontrol/read-only-fs/mixed-two-deployments",
			CheckName:      "access-control-security-context-read-only-file-system",
			Category:       checks.CategoryAccessControl,
			Description:    "Two deployments, one without readOnlyRootFS should be non-compliant",
			ExpectedStatus: checks.StatusNonCompliant,
			Setup: func(ctx *scenario.RunContext) error {
				dep1 := builder.NewDeployment("test-dep-1", ctx.Namespace).Build()
				if err := cluster.CreateAndWaitForDeployment(ctx.Ctx, ctx.Client, dep1, cluster.DefaultTimeout); err != nil {
					return err
				}
				dep2 := builder.NewDeployment("test-dep-2", ctx.Namespace).WithReadOnlyRootFS().Build()
				return cluster.CreateAndWaitForDeployment(ctx.Ctx, ctx.Client, dep2, cluster.DefaultTimeout)
			},
		},
	)

	scenario.Register(
		scenario.Scenario{
			Name:           "accesscontrol/non-root-user-id/compliant",
			CheckName:      "access-control-security-context-non-root-user-id-check",
			Category:       checks.CategoryAccessControl,
			Description:    "Deployment with runAsNonRoot=true should be compliant",
			ExpectedStatus: checks.StatusCompliant,
			Privileged:     true,
			Setup: func(ctx *scenario.RunContext) error {
				dep := builder.NewDeployment("test-dep", ctx.Namespace).
					WithRunAsNonRoot().
					WithRunAsUser(1000).
					Build()
				return cluster.CreateAndWaitForDeployment(ctx.Ctx, ctx.Client, dep, cluster.DefaultTimeout)
			},
		},
		scenario.Scenario{
			Name:           "accesscontrol/non-root-user-id/non-compliant",
			CheckName:      "access-control-security-context-non-root-user-id-check",
			Category:       checks.CategoryAccessControl,
			Description:    "Deployment with runAsUser=0 should be non-compliant",
			ExpectedStatus: checks.StatusNonCompliant,
			Privileged:     true,
			Setup: func(ctx *scenario.RunContext) error {
				dep := builder.NewDeployment("test-dep", ctx.Namespace).
					WithRunAsUser(0).
					Build()
				return cluster.CreateAndWaitForDeployment(ctx.Ctx, ctx.Client, dep, cluster.DefaultTimeout)
			},
		},
		scenario.Scenario{
			Name:           "accesscontrol/non-root-user-id/mixed-two-deployments",
			CheckName:      "access-control-security-context-non-root-user-id-check",
			Category:       checks.CategoryAccessControl,
			Description:    "Two deployments, one with runAsUser=0 should be non-compliant",
			ExpectedStatus: checks.StatusNonCompliant,
			Privileged:     true,
			Setup: scenario.TwoDeploymentSetup(func(b *builder.DeploymentBuilder) *builder.DeploymentBuilder {
				return b.WithRunAsUser(0)
			}),
		},
	)

	scenario.Register(
		scenario.Scenario{
			Name:           "accesscontrol/security-context/non-compliant",
			CheckName:      "access-control-security-context",
			Category:       checks.CategoryAccessControl,
			Description:    "Deployment with IPC_LOCK capability should require elevated SCC",
			ExpectedStatus: checks.StatusNonCompliant,
			Privileged:     true,
			Setup: func(ctx *scenario.RunContext) error {
				dep := builder.NewDeployment("test-dep", ctx.Namespace).
					WithCapability("IPC_LOCK").
					Build()
				return cluster.CreateAndWaitForDeployment(ctx.Ctx, ctx.Client, dep, cluster.DefaultTimeout)
			},
		},
		scenario.Scenario{
			Name:           "accesscontrol/security-context/mixed-two-deployments",
			CheckName:      "access-control-security-context",
			Category:       checks.CategoryAccessControl,
			Description:    "Two deployments, one with IPC_LOCK should be non-compliant",
			ExpectedStatus: checks.StatusNonCompliant,
			Privileged:     true,
			Setup: scenario.TwoDeploymentSetup(func(b *builder.DeploymentBuilder) *builder.DeploymentBuilder {
				return b.WithCapability("IPC_LOCK")
			}),
		},
	)
}
