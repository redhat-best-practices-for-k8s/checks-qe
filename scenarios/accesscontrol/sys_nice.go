package accesscontrol

import (
	checks "github.com/redhat-best-practices-for-k8s/checks"
	"github.com/redhat-best-practices-for-k8s/checks-qe/pkg/builder"
	"github.com/redhat-best-practices-for-k8s/checks-qe/pkg/cluster"
	"github.com/redhat-best-practices-for-k8s/checks-qe/pkg/scenario"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

func registerSysNice() {
	scenario.Register(
		scenario.Scenario{
			Name:           "accesscontrol/sys-nice-rt/compliant-no-rt-nodes",
			CheckName:      "access-control-sys-nice-realtime-capability",
			Category:       checks.CategoryAccessControl,
			Description:    "No RT nodes means SYS_NICE check is compliant",
			ExpectedStatus: checks.StatusCompliant,
			Setup:          scenario.VanillaDeploymentSetup,
		},
		scenario.Scenario{
			Name:           "accesscontrol/sys-nice-rt/compliant",
			CheckName:      "access-control-sys-nice-realtime-capability",
			Category:       checks.CategoryAccessControl,
			Description:    "Container on RT node with SYS_NICE should be compliant",
			ExpectedStatus: checks.StatusCompliant,
			Privileged:     true,
			Setup: func(ctx *scenario.RunContext) error {
				dep := builder.NewDeployment("test-dep", ctx.Namespace).
					WithCapability("SYS_NICE").
					Build()
				return cluster.CreateAndWaitForDeployment(ctx.Ctx, ctx.Client, dep, cluster.DefaultTimeout)
			},
			PostDiscovery: injectRTNode,
		},
		scenario.Scenario{
			Name:           "accesscontrol/sys-nice-rt/non-compliant",
			CheckName:      "access-control-sys-nice-realtime-capability",
			Category:       checks.CategoryAccessControl,
			Description:    "Container on RT node without SYS_NICE should be non-compliant",
			ExpectedStatus: checks.StatusNonCompliant,
			Setup:          scenario.VanillaDeploymentSetup,
			PostDiscovery:  injectRTNode,
		},
	)
}

func injectRTNode(resources *checks.DiscoveredResources) {
	resources.Nodes = append(resources.Nodes, corev1.Node{
		ObjectMeta: metav1.ObjectMeta{Name: "rt-node-1"},
		Status: corev1.NodeStatus{
			NodeInfo: corev1.NodeSystemInfo{KernelVersion: "5.14.0-362.24.1.el9_3.x86_64+rt"},
		},
	})
	for i := range resources.Pods {
		resources.Pods[i].Spec.NodeName = "rt-node-1"
	}
}
