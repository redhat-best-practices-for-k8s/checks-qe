package lifecycle

import (
	checks "github.com/redhat-best-practices-for-k8s/checks"
	"github.com/redhat-best-practices-for-k8s/checks-qe/pkg/builder"
	"github.com/redhat-best-practices-for-k8s/checks-qe/pkg/cluster"
	"github.com/redhat-best-practices-for-k8s/checks-qe/pkg/scenario"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

func registerAffinity() {
	scenario.Register(
		scenario.Scenario{
			Name:           "lifecycle/affinity-required/compliant-no-label",
			CheckName:      "lifecycle-affinity-required-pods",
			Category:       checks.CategoryLifecycle,
			Description:    "Pods without AffinityRequired label should be compliant",
			ExpectedStatus: checks.StatusCompliant,
			Setup:          scenario.VanillaDeploymentSetup,
		},
		scenario.Scenario{
			Name:           "lifecycle/affinity-required/compliant",
			CheckName:      "lifecycle-affinity-required-pods",
			Category:       checks.CategoryLifecycle,
			Description:    "Pod with AffinityRequired label and pod affinity should be compliant",
			ExpectedStatus: checks.StatusCompliant,
			Setup: func(ctx *scenario.RunContext) error {
				dep := builder.NewDeployment("test-dep", ctx.Namespace).
					WithLabel("AffinityRequired", "true").
					Build()
				return cluster.CreateAndWaitForDeployment(ctx.Ctx, ctx.Client, dep, cluster.DefaultTimeout)
			},
			PostDiscovery: func(resources *checks.DiscoveredResources) {
				for i := range resources.Pods {
					resources.Pods[i].Labels["AffinityRequired"] = "true"
					resources.Pods[i].Spec.Affinity = &corev1.Affinity{
						PodAffinity: &corev1.PodAffinity{
							RequiredDuringSchedulingIgnoredDuringExecution: []corev1.PodAffinityTerm{{
								LabelSelector: &metav1.LabelSelector{
									MatchLabels: map[string]string{"app": "test-dep"},
								},
								TopologyKey: "kubernetes.io/hostname",
							}},
						},
					}
				}
			},
		},
		scenario.Scenario{
			Name:           "lifecycle/affinity-required/non-compliant",
			CheckName:      "lifecycle-affinity-required-pods",
			Category:       checks.CategoryLifecycle,
			Description:    "Pod with AffinityRequired label but no affinity should be non-compliant",
			ExpectedStatus: checks.StatusNonCompliant,
			Setup:          scenario.VanillaDeploymentSetup,
			PostDiscovery: func(resources *checks.DiscoveredResources) {
				for i := range resources.Pods {
					resources.Pods[i].Labels["AffinityRequired"] = "true"
				}
			},
		},
	)
}
