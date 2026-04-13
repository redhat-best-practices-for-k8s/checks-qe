package platform

import (
	checks "github.com/redhat-best-practices-for-k8s/checks"
	"github.com/redhat-best-practices-for-k8s/checks-qe/pkg/builder"
	"github.com/redhat-best-practices-for-k8s/checks-qe/pkg/cluster"
	"github.com/redhat-best-practices-for-k8s/checks-qe/pkg/scenario"
)

func registerHugepages() {
	scenario.Register(
		scenario.Scenario{
			Name:           "platform/hugepages-1g-only/compliant",
			CheckName:      "platform-alteration-hugepages-1g-only",
			Category:       checks.CategoryPlatformAlteration,
			Description:    "Deployment without 2Mi hugepages should be compliant",
			ExpectedStatus: checks.StatusCompliant,
			Setup:          scenario.VanillaDeploymentSetup,
		},
		scenario.Scenario{
			Name:           "platform/hugepages-1g-only/non-compliant",
			CheckName:      "platform-alteration-hugepages-1g-only",
			Category:       checks.CategoryPlatformAlteration,
			Description:    "Deployment requesting 2Mi hugepages should be non-compliant",
			ExpectedStatus: checks.StatusNonCompliant,
			Privileged:     true,
			Tags:        []string{"hugepages"},
			RequiresProbe: true,
			Setup: func(ctx *scenario.RunContext) error {
				dep := builder.NewDeployment("test-dep", ctx.Namespace).
					WithResourceRequests("100m", "128Mi").
					WithHugepagesRequest("2Mi", "128Mi").
					Build()
				return cluster.CreateAndWaitForDeployment(ctx.Ctx, ctx.Client, dep, cluster.DefaultTimeout)
			},
		},
	)

	scenario.Register(
		scenario.Scenario{
			Name:           "platform/hugepages-2m-only/compliant",
			CheckName:      "platform-alteration-hugepages-2m-only",
			Category:       checks.CategoryPlatformAlteration,
			Description:    "Deployment without 1Gi hugepages should be compliant",
			ExpectedStatus: checks.StatusCompliant,
			Setup:          scenario.VanillaDeploymentSetup,
		},
		scenario.Scenario{
			Name:           "platform/hugepages-2m-only/non-compliant",
			CheckName:      "platform-alteration-hugepages-2m-only",
			Category:       checks.CategoryPlatformAlteration,
			Description:    "Deployment requesting 1Gi hugepages should be non-compliant",
			ExpectedStatus: checks.StatusNonCompliant,
			Privileged:     true,
			Tags:           []string{"hugepages"},
			RequiresProbe:  true,
			Setup: func(ctx *scenario.RunContext) error {
				dep := builder.NewDeployment("test-dep", ctx.Namespace).
					WithResourceRequests("100m", "128Mi").
					WithHugepagesRequest("1Gi", "1Gi").
					Build()
				return cluster.CreateAndWaitForDeployment(ctx.Ctx, ctx.Client, dep, cluster.DefaultTimeout)
			},
		},
	)
}
