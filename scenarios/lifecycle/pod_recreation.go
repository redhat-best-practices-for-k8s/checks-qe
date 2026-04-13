package lifecycle

import (
	checks "github.com/redhat-best-practices-for-k8s/checks"
	"github.com/redhat-best-practices-for-k8s/checks-qe/pkg/scenario"
)

func registerPodRecreation() {
	scenario.Register(
		scenario.Scenario{
			Name:           "lifecycle/pod-recreation/compliant",
			CheckName:      "lifecycle-pod-recreation",
			Category:       checks.CategoryLifecycle,
			Description:    "Pods managed by a Deployment should be recreated after deletion",
			ExpectedStatus: checks.StatusCompliant,
			Setup:          scenario.VanillaDeploymentSetup,
		},
	)
}
