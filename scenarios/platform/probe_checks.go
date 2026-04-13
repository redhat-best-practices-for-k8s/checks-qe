package platform

import (
	checks "github.com/redhat-best-practices-for-k8s/checks"
	"github.com/redhat-best-practices-for-k8s/checks-qe/pkg/scenario"
)

func registerProbeChecks() {
	scenario.Register(
		scenario.Scenario{
			Name:           "platform/boot-params/compliant-no-probe",
			CheckName:      "platform-alteration-boot-params",
			Category:       checks.CategoryPlatformAlteration,
			Description:    "No probe executor means boot params check returns compliant",
			ExpectedStatus: checks.StatusCompliant,
			Setup:          scenario.VanillaDeploymentSetup,
		},
		scenario.Scenario{
			Name:           "platform/sysctl-config/compliant-no-probe",
			CheckName:      "platform-alteration-sysctl-config",
			Category:       checks.CategoryPlatformAlteration,
			Description:    "No probe executor means sysctl check returns compliant",
			ExpectedStatus: checks.StatusCompliant,
			Setup:          scenario.VanillaDeploymentSetup,
		},
		scenario.Scenario{
			Name:           "platform/tainted-node-kernel/compliant-no-probe",
			CheckName:      "platform-alteration-tainted-node-kernel",
			Category:       checks.CategoryPlatformAlteration,
			Description:    "No probe executor means tainted kernel check returns compliant",
			ExpectedStatus: checks.StatusCompliant,
			Setup:          scenario.VanillaDeploymentSetup,
		},
		scenario.Scenario{
			Name:           "platform/hugepages-config/compliant-no-probe",
			CheckName:      "platform-alteration-hugepages-config",
			Category:       checks.CategoryPlatformAlteration,
			Description:    "No probe executor means hugepages config check returns compliant",
			ExpectedStatus: checks.StatusCompliant,
			Setup:          scenario.VanillaDeploymentSetup,
		},
		scenario.Scenario{
			Name:           "platform/selinux-enforcing/compliant-no-probe",
			CheckName:      "platform-alteration-is-selinux-enforcing",
			Category:       checks.CategoryPlatformAlteration,
			Description:    "No probe executor means SELinux check returns compliant",
			ExpectedStatus: checks.StatusCompliant,
			Setup:          scenario.VanillaDeploymentSetup,
		},
		scenario.Scenario{
			Name:           "platform/isredhat-release/error-no-probe",
			CheckName:      "platform-alteration-isredhat-release",
			Category:       checks.CategoryPlatformAlteration,
			Description:    "No probe executor returns error for Red Hat release check",
			ExpectedStatus: checks.StatusError,
			Setup:          scenario.VanillaDeploymentSetup,
		},
		scenario.Scenario{
			Name:           "platform/base-image/compliant-not-ocp",
			CheckName:      "platform-alteration-base-image",
			Category:       checks.CategoryPlatformAlteration,
			Description:    "Non-OCP cluster should be compliant for base image check",
			ExpectedStatus: checks.StatusCompliant,
			Setup:          scenario.VanillaDeploymentSetup,
		},
		scenario.Scenario{
			Name:           "platform/hyperthread-enable/error-no-probe",
			CheckName:      "platform-alteration-hyperthread-enable",
			Category:       checks.CategoryPlatformAlteration,
			Description:    "No probe executor returns error for hyperthread check",
			ExpectedStatus: checks.StatusError,
			Setup:          scenario.VanillaDeploymentSetup,
		},
	)
}
