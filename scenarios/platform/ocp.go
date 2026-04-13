package platform

import (
	"fmt"

	checks "github.com/redhat-best-practices-for-k8s/checks"
	"github.com/redhat-best-practices-for-k8s/checks-qe/pkg/scenario"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

func registerOCP() {
	scenario.Register(
		scenario.Scenario{
			Name:           "platform/ocp-lifecycle/compliant-not-ocp",
			CheckName:      "platform-alteration-ocp-lifecycle",
			Category:       checks.CategoryPlatformAlteration,
			Description:    "Non-OCP cluster should be compliant",
			ExpectedStatus: checks.StatusCompliant,
			Setup:          scenario.VanillaDeploymentSetup,
		},
		scenario.Scenario{
			Name:           "platform/ocp-lifecycle/compliant-ga",
			CheckName:      "platform-alteration-ocp-lifecycle",
			Category:       checks.CategoryPlatformAlteration,
			Description:    "OCP version in GA status should be compliant",
			ExpectedStatus: checks.StatusCompliant,
			Setup:          scenario.VanillaDeploymentSetup,
			PostDiscovery: func(resources *checks.DiscoveredResources) {
				resources.OpenshiftVersion = "4.14.5"
				resources.OCPStatus = "GA"
			},
		},
		scenario.Scenario{
			Name:           "platform/ocp-lifecycle/non-compliant-eol",
			CheckName:      "platform-alteration-ocp-lifecycle",
			Category:       checks.CategoryPlatformAlteration,
			Description:    "OCP version in EOL status should be non-compliant",
			ExpectedStatus: checks.StatusNonCompliant,
			Setup:          scenario.VanillaDeploymentSetup,
			PostDiscovery: func(resources *checks.DiscoveredResources) {
				resources.OpenshiftVersion = "4.10.0"
				resources.OCPStatus = "EOL"
			},
		},
	)

	scenario.Register(
		scenario.Scenario{
			Name:           "platform/ocp-node-count/compliant",
			CheckName:      "platform-alteration-ocp-node-count",
			Category:       checks.CategoryPlatformAlteration,
			Description:    "Cluster with 3+ worker nodes should be compliant",
			ExpectedStatus: checks.StatusCompliant,
			Setup:          scenario.VanillaDeploymentSetup,
			PostDiscovery: func(resources *checks.DiscoveredResources) {
				resources.Nodes = makeWorkerNodes(3)
			},
		},
		scenario.Scenario{
			Name:           "platform/ocp-node-count/non-compliant",
			CheckName:      "platform-alteration-ocp-node-count",
			Category:       checks.CategoryPlatformAlteration,
			Description:    "Cluster with fewer than 3 worker nodes should be non-compliant",
			ExpectedStatus: checks.StatusNonCompliant,
			Setup:          scenario.VanillaDeploymentSetup,
			PostDiscovery: func(resources *checks.DiscoveredResources) {
				resources.Nodes = makeWorkerNodes(1)
			},
		},
	)

	scenario.Register(
		scenario.Scenario{
			Name:           "platform/ocp-node-os/compliant-not-ocp",
			CheckName:      "platform-alteration-ocp-node-os-lifecycle",
			Category:       checks.CategoryPlatformAlteration,
			Description:    "Non-OCP cluster should be compliant",
			ExpectedStatus: checks.StatusCompliant,
			Setup:          scenario.VanillaDeploymentSetup,
		},
		scenario.Scenario{
			Name:           "platform/ocp-node-os/compliant-rhcos",
			CheckName:      "platform-alteration-ocp-node-os-lifecycle",
			Category:       checks.CategoryPlatformAlteration,
			Description:    "Control plane nodes with RHCOS should be compliant",
			ExpectedStatus: checks.StatusCompliant,
			Setup:          scenario.VanillaDeploymentSetup,
			PostDiscovery: func(resources *checks.DiscoveredResources) {
				resources.OpenshiftVersion = "4.14.5"
				resources.Nodes = []corev1.Node{
					makeControlPlaneNode("cp-1", "Red Hat Enterprise Linux CoreOS 414.92.202401041559-0"),
				}
			},
		},
		scenario.Scenario{
			Name:           "platform/ocp-node-os/non-compliant",
			CheckName:      "platform-alteration-ocp-node-os-lifecycle",
			Category:       checks.CategoryPlatformAlteration,
			Description:    "Control plane nodes with incompatible OS should be non-compliant",
			ExpectedStatus: checks.StatusNonCompliant,
			Setup:          scenario.VanillaDeploymentSetup,
			PostDiscovery: func(resources *checks.DiscoveredResources) {
				resources.OpenshiftVersion = "4.14.5"
				resources.Nodes = []corev1.Node{
					makeControlPlaneNode("cp-1", "Ubuntu 22.04.3 LTS"),
				}
			},
		},
	)

	scenario.Register(
		scenario.Scenario{
			Name:           "platform/cluster-operator-health/compliant-no-operators",
			CheckName:      "platform-alteration-cluster-operator-health",
			Category:       checks.CategoryPlatformAlteration,
			Description:    "No cluster operators should be compliant",
			ExpectedStatus: checks.StatusCompliant,
			Setup:          scenario.VanillaDeploymentSetup,
		},
	)
}

func makeWorkerNodes(count int) []corev1.Node {
	nodes := make([]corev1.Node, count)
	for i := range count {
		nodes[i] = corev1.Node{
			ObjectMeta: metav1.ObjectMeta{
				Name: fmt.Sprintf("worker-%d", i),
			},
		}
	}
	return nodes
}

func makeControlPlaneNode(name, osImage string) corev1.Node {
	return corev1.Node{
		ObjectMeta: metav1.ObjectMeta{
			Name: name,
			Labels: map[string]string{
				"node-role.kubernetes.io/control-plane": "",
			},
		},
		Status: corev1.NodeStatus{
			NodeInfo: corev1.NodeSystemInfo{OSImage: osImage},
		},
	}
}
