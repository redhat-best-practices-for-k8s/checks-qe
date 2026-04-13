package networking

import (
	checks "github.com/redhat-best-practices-for-k8s/checks"
	"github.com/redhat-best-practices-for-k8s/checks-qe/pkg/scenario"
	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/resource"
)

func registerSRIOV() {
	scenario.Register(
		scenario.Scenario{
			Name:           "networking/sriov-restart-label/compliant-no-sriov",
			CheckName:      "networking-restart-on-reboot-sriov-pod",
			Category:       checks.CategoryNetworking,
			Description:    "Pods without SR-IOV should be compliant",
			ExpectedStatus: checks.StatusCompliant,
			Setup:          scenario.VanillaDeploymentSetup,
		},
		scenario.Scenario{
			Name:           "networking/sriov-restart-label/compliant",
			CheckName:      "networking-restart-on-reboot-sriov-pod",
			Category:       checks.CategoryNetworking,
			Description:    "SR-IOV pod with restart-on-reboot label should be compliant",
			ExpectedStatus: checks.StatusCompliant,
			Setup:          scenario.VanillaDeploymentSetup,
			PostDiscovery: func(resources *checks.DiscoveredResources) {
				for i := range resources.Pods {
					makeSRIOVPod(&resources.Pods[i])
					resources.Pods[i].Labels["restart-on-reboot"] = "true"
				}
			},
		},
		scenario.Scenario{
			Name:           "networking/sriov-restart-label/non-compliant",
			CheckName:      "networking-restart-on-reboot-sriov-pod",
			Category:       checks.CategoryNetworking,
			Description:    "SR-IOV pod without restart-on-reboot label should be non-compliant",
			ExpectedStatus: checks.StatusNonCompliant,
			Setup:          scenario.VanillaDeploymentSetup,
			PostDiscovery: func(resources *checks.DiscoveredResources) {
				for i := range resources.Pods {
					makeSRIOVPod(&resources.Pods[i])
				}
			},
		},
	)

	scenario.Register(
		scenario.Scenario{
			Name:           "networking/sriov-nad-mtu/compliant-no-nads",
			CheckName:      "networking-network-attachment-definition-sriov-mtu",
			Category:       checks.CategoryNetworking,
			Description:    "No network attachment definitions should be compliant",
			ExpectedStatus: checks.StatusCompliant,
			Setup:          scenario.VanillaDeploymentSetup,
		},
	)
}

func makeSRIOVPod(pod *corev1.Pod) {
	if pod.Annotations == nil {
		pod.Annotations = map[string]string{}
	}
	pod.Annotations["k8s.v1.cni.cncf.io/networks"] = "sriov-net"
	for i := range pod.Spec.Containers {
		if pod.Spec.Containers[i].Resources.Requests == nil {
			pod.Spec.Containers[i].Resources.Requests = corev1.ResourceList{}
		}
		pod.Spec.Containers[i].Resources.Requests["openshift.io/sriov-nic"] = resource.MustParse("1")
	}
}
