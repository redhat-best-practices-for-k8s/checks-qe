package accesscontrol

import (
	checks "github.com/redhat-best-practices-for-k8s/checks"
	"github.com/redhat-best-practices-for-k8s/checks-qe/pkg/scenario"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

func syntheticPod(name, namespace string) corev1.Pod {
	return corev1.Pod{
		ObjectMeta: metav1.ObjectMeta{Name: name, Namespace: namespace},
		Spec: corev1.PodSpec{
			Containers: []corev1.Container{{Name: "test", Image: "test:latest"}},
		},
	}
}

func registerNamespace() {
	scenario.Register(
		scenario.Scenario{
			Name:           "accesscontrol/namespace/compliant",
			CheckName:      "access-control-namespace",
			Category:       checks.CategoryAccessControl,
			Description:    "Pod in valid namespace should be compliant",
			ExpectedStatus: checks.StatusCompliant,
			Setup:          scenario.VanillaDeploymentSetup,
		},
		scenario.Scenario{
			Name:           "accesscontrol/namespace/non-compliant-openshift-prefix",
			CheckName:      "access-control-namespace",
			Category:       checks.CategoryAccessControl,
			Description:    "Pod in openshift- namespace should be non-compliant",
			ExpectedStatus: checks.StatusNonCompliant,
			Setup:          scenario.VanillaDeploymentSetup,
			PostDiscovery: func(resources *checks.DiscoveredResources) {
				resources.Pods = append(resources.Pods, syntheticPod("bad-pod", "openshift-test"))
			},
		},
		scenario.Scenario{
			Name:           "accesscontrol/namespace/non-compliant-default",
			CheckName:      "access-control-namespace",
			Category:       checks.CategoryAccessControl,
			Description:    "Pod in default namespace should be non-compliant",
			ExpectedStatus: checks.StatusNonCompliant,
			Setup:          scenario.VanillaDeploymentSetup,
			PostDiscovery: func(resources *checks.DiscoveredResources) {
				resources.Pods = append(resources.Pods, syntheticPod("bad-pod", "default"))
			},
		},
		scenario.Scenario{
			Name:           "accesscontrol/namespace/mixed-two-pods",
			CheckName:      "access-control-namespace",
			Category:       checks.CategoryAccessControl,
			Description:    "Pods in valid and invalid namespaces should be non-compliant",
			ExpectedStatus: checks.StatusNonCompliant,
			Setup:          scenario.VanillaDeploymentSetup,
			PostDiscovery: func(resources *checks.DiscoveredResources) {
				resources.Pods = append(resources.Pods, syntheticPod("bad-pod", "openshift-bad"))
			},
		},
		scenario.Scenario{
			Name:           "accesscontrol/namespace/compliant-cr-in-valid-ns",
			CheckName:      "access-control-namespace",
			Category:       checks.CategoryAccessControl,
			Description:    "CR in configured namespace should be compliant",
			ExpectedStatus: checks.StatusCompliant,
			Setup:          scenario.VanillaDeploymentSetup,
			PostDiscovery: func(resources *checks.DiscoveredResources) {
				ns := resources.Namespaces[0]
				resources.CRInstances = map[string]map[string][]string{
					"widgets.example.com": {ns: {"my-widget"}},
				}
			},
		},
		scenario.Scenario{
			Name:           "accesscontrol/namespace/non-compliant-cr-in-invalid-ns",
			CheckName:      "access-control-namespace",
			Category:       checks.CategoryAccessControl,
			Description:    "CR in unconfigured namespace should be non-compliant",
			ExpectedStatus: checks.StatusNonCompliant,
			Setup:          scenario.VanillaDeploymentSetup,
			PostDiscovery: func(resources *checks.DiscoveredResources) {
				resources.CRInstances = map[string]map[string][]string{
					"widgets.example.com": {"other-namespace": {"my-widget"}},
				}
			},
		},
	)
}
