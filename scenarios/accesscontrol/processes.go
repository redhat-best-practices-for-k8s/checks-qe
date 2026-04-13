package accesscontrol

import (
	checks "github.com/redhat-best-practices-for-k8s/checks"
	"github.com/redhat-best-practices-for-k8s/checks-qe/pkg/scenario"
)

func registerProcesses() {
	scenario.Register(
		scenario.Scenario{
			Name:           "accesscontrol/one-process-per-container/compliant-no-probe",
			CheckName:      "access-control-one-process-per-container",
			Category:       checks.CategoryAccessControl,
			Description:    "No probe executor means process check returns compliant",
			ExpectedStatus: checks.StatusCompliant,
			Setup:          scenario.VanillaDeploymentSetup,
		},
		scenario.Scenario{
			Name:           "accesscontrol/ssh-daemons/compliant-no-probe",
			CheckName:      "access-control-ssh-daemons",
			Category:       checks.CategoryAccessControl,
			Description:    "No probe executor means SSH daemon check returns compliant",
			ExpectedStatus: checks.StatusCompliant,
			Setup:          scenario.VanillaDeploymentSetup,
		},
	)
}
