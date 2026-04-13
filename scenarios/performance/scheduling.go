package performance

import (
	checks "github.com/redhat-best-practices-for-k8s/checks"
	"github.com/redhat-best-practices-for-k8s/checks-qe/pkg/scenario"
)

func registerSchedulingPolicy() {
	scenario.Register(
		scenario.Scenario{
			Name:           "performance/exclusive-cpu-pool-rt-scheduling/error-no-probe",
			CheckName:      "performance-exclusive-cpu-pool-rt-scheduling-policy",
			Category:       checks.CategoryPerformance,
			Description:    "No probe executor returns error for exclusive CPU pool check",
			ExpectedStatus: checks.StatusError,
			Setup:          scenario.VanillaDeploymentSetup,
		},
		scenario.Scenario{
			Name:           "performance/isolated-cpu-pool-rt-scheduling/error-no-probe",
			CheckName:      "performance-isolated-cpu-pool-rt-scheduling-policy",
			Category:       checks.CategoryPerformance,
			Description:    "No probe executor returns error for isolated CPU pool check",
			ExpectedStatus: checks.StatusError,
			Setup:          scenario.VanillaDeploymentSetup,
		},
		scenario.Scenario{
			Name:           "performance/shared-cpu-pool-scheduling/error-no-probe",
			CheckName:      "performance-shared-cpu-pool-non-rt-scheduling-policy",
			Category:       checks.CategoryPerformance,
			Description:    "No probe executor returns error for shared CPU pool check",
			ExpectedStatus: checks.StatusError,
			Setup:          scenario.VanillaDeploymentSetup,
		},
	)
}
