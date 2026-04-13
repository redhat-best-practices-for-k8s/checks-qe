package networking

import (
	checks "github.com/redhat-best-practices-for-k8s/checks"
	"github.com/redhat-best-practices-for-k8s/checks-qe/pkg/scenario"
)

func registerICMP() {
	scenario.Register(
		scenario.Scenario{
			Name:           "networking/icmpv4-connectivity/error-no-probe",
			CheckName:      "networking-icmpv4-connectivity",
			Category:       checks.CategoryNetworking,
			Description:    "No probe executor returns error for ICMP connectivity check",
			ExpectedStatus: checks.StatusError,
			Setup:          scenario.VanillaDeploymentSetup,
		},
		scenario.Scenario{
			Name:           "networking/icmpv4-connectivity-multus/error-no-probe",
			CheckName:      "networking-icmpv4-connectivity-multus",
			Category:       checks.CategoryNetworking,
			Description:    "No probe executor returns error for Multus ICMP check",
			ExpectedStatus: checks.StatusError,
			RequiresMultus: true,
			Setup:          scenario.VanillaDeploymentSetup,
		},
		scenario.Scenario{
			Name:              "networking/icmpv6-connectivity/error-no-probe",
			CheckName:         "networking-icmpv6-connectivity",
			Category:          checks.CategoryNetworking,
			Description:       "No probe executor returns error for IPv6 ICMP check",
			ExpectedStatus:    checks.StatusError,
			RequiresDualStack: true,
			Setup:             scenario.VanillaDeploymentSetup,
		},
		scenario.Scenario{
			Name:              "networking/icmpv6-connectivity-multus/error-no-probe",
			CheckName:         "networking-icmpv6-connectivity-multus",
			Category:          checks.CategoryNetworking,
			Description:       "No probe executor returns error for IPv6 Multus ICMP check",
			ExpectedStatus:    checks.StatusError,
			RequiresDualStack: true,
			RequiresMultus:    true,
			Setup:             scenario.VanillaDeploymentSetup,
		},
	)
}
