package manageability

import (
	checks "github.com/redhat-best-practices-for-k8s/checks"
	"github.com/redhat-best-practices-for-k8s/checks-qe/pkg/builder"
	"github.com/redhat-best-practices-for-k8s/checks-qe/pkg/cluster"
	"github.com/redhat-best-practices-for-k8s/checks-qe/pkg/scenario"
)

func registerContainers() {
	scenario.Register(
		scenario.Scenario{
			Name:           "manageability/container-port-name-format/compliant",
			CheckName:      "manageability-container-port-name-format",
			Category:       checks.CategoryManageability,
			Description:    "Container port with valid protocol prefix should be compliant",
			ExpectedStatus: checks.StatusCompliant,
			Setup: func(ctx *scenario.RunContext) error {
				dep := builder.NewDeployment("test-dep", ctx.Namespace).
					WithNamedContainerPort("http-web", 8080).
					Build()
				return cluster.CreateAndWaitForDeployment(ctx.Ctx, ctx.Client, dep, cluster.DefaultTimeout)
			},
		},
		scenario.Scenario{
			Name:           "manageability/container-port-name-format/non-compliant",
			CheckName:      "manageability-container-port-name-format",
			Category:       checks.CategoryManageability,
			Description:    "Container port with invalid name should be non-compliant",
			ExpectedStatus: checks.StatusNonCompliant,
			Setup: func(ctx *scenario.RunContext) error {
				dep := builder.NewDeployment("test-dep", ctx.Namespace).
					WithNamedContainerPort("myport", 8080).
					Build()
				return cluster.CreateAndWaitForDeployment(ctx.Ctx, ctx.Client, dep, cluster.DefaultTimeout)
			},
		},
	)

	scenario.Register(
		scenario.Scenario{
			Name:           "manageability/containers-image-tag/compliant",
			CheckName:      "manageability-containers-image-tag",
			Category:       checks.CategoryManageability,
			Description:    "Container with tagged image should be compliant",
			ExpectedStatus: checks.StatusCompliant,
			Setup:          scenario.VanillaDeploymentSetup,
		},
	)
}
