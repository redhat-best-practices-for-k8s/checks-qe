package lifecycle

import (
	"fmt"

	checks "github.com/redhat-best-practices-for-k8s/checks"
	"github.com/redhat-best-practices-for-k8s/checks-qe/pkg/scenario"
	corev1 "k8s.io/api/core/v1"
	storagev1 "k8s.io/api/storage/v1"
	"k8s.io/apimachinery/pkg/api/resource"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

func registerStorage() {
	scenario.Register(
		scenario.Scenario{
			Name:           "lifecycle/pv-reclaim-policy/compliant-no-pvs",
			CheckName:      "lifecycle-persistent-volume-reclaim-policy",
			Category:       checks.CategoryLifecycle,
			Description:    "No PersistentVolumes should be compliant",
			ExpectedStatus: checks.StatusCompliant,
			Setup:          scenario.VanillaDeploymentSetup,
			PostDiscovery: func(resources *checks.DiscoveredResources) {
				resources.PersistentVolumes = nil
			},
		},
		scenario.Scenario{
			Name:           "lifecycle/pv-reclaim-policy/compliant",
			CheckName:      "lifecycle-persistent-volume-reclaim-policy",
			Category:       checks.CategoryLifecycle,
			Description:    "PV with Delete reclaim policy should be compliant",
			ExpectedStatus: checks.StatusCompliant,
			Setup:          scenario.VanillaDeploymentSetup,
			PostDiscovery: func(resources *checks.DiscoveredResources) {
				resources.PersistentVolumes = []corev1.PersistentVolume{
					makePV("test-pv", corev1.PersistentVolumeReclaimDelete),
				}
			},
		},
		scenario.Scenario{
			Name:           "lifecycle/pv-reclaim-policy/non-compliant",
			CheckName:      "lifecycle-persistent-volume-reclaim-policy",
			Category:       checks.CategoryLifecycle,
			Description:    "PV with Retain reclaim policy should be non-compliant",
			ExpectedStatus: checks.StatusNonCompliant,
			Setup:          scenario.VanillaDeploymentSetup,
			PostDiscovery: func(resources *checks.DiscoveredResources) {
				resources.PersistentVolumes = append(resources.PersistentVolumes,
					makePV("test-pv", corev1.PersistentVolumeReclaimRetain))
			},
		},
	)

	provisioner := "ebs.csi.aws.com"
	localProvisioner := "kubernetes.io/no-provisioner"

	scenario.Register(
		scenario.Scenario{
			Name:           "lifecycle/storage-provisioner/compliant-no-pvcs",
			CheckName:      "lifecycle-storage-provisioner",
			Category:       checks.CategoryLifecycle,
			Description:    "Pods without PVC volumes should be compliant",
			ExpectedStatus: checks.StatusCompliant,
			Setup:          scenario.VanillaDeploymentSetup,
		},
		scenario.Scenario{
			Name:           "lifecycle/storage-provisioner/compliant-non-local-multi-node",
			CheckName:      "lifecycle-storage-provisioner",
			Category:       checks.CategoryLifecycle,
			Description:    "Non-local provisioner in multi-node cluster should be compliant",
			ExpectedStatus: checks.StatusCompliant,
			Setup:          scenario.VanillaDeploymentSetup,
			PostDiscovery: func(resources *checks.DiscoveredResources) {
				injectStorageScenario(resources, provisioner, 3)
			},
		},
		scenario.Scenario{
			Name:           "lifecycle/storage-provisioner/non-compliant-local-multi-node",
			CheckName:      "lifecycle-storage-provisioner",
			Category:       checks.CategoryLifecycle,
			Description:    "Local provisioner in multi-node cluster should be non-compliant",
			ExpectedStatus: checks.StatusNonCompliant,
			Setup:          scenario.VanillaDeploymentSetup,
			PostDiscovery: func(resources *checks.DiscoveredResources) {
				injectStorageScenario(resources, localProvisioner, 3)
			},
		},
	)
}

func makePV(name string, policy corev1.PersistentVolumeReclaimPolicy) corev1.PersistentVolume {
	return corev1.PersistentVolume{
		ObjectMeta: metav1.ObjectMeta{Name: name},
		Spec: corev1.PersistentVolumeSpec{
			PersistentVolumeReclaimPolicy: policy,
			Capacity: corev1.ResourceList{
				corev1.ResourceStorage: resource.MustParse("1Gi"),
			},
			AccessModes: []corev1.PersistentVolumeAccessMode{corev1.ReadWriteOnce},
			PersistentVolumeSource: corev1.PersistentVolumeSource{
				HostPath: &corev1.HostPathVolumeSource{Path: "/tmp/test"},
			},
		},
	}
}

func injectStorageScenario(resources *checks.DiscoveredResources, prov string, nodeCount int) {
	scName := "test-sc"

	resources.StorageClasses = append(resources.StorageClasses, storagev1.StorageClass{
		ObjectMeta:  metav1.ObjectMeta{Name: scName},
		Provisioner: prov,
	})

	resources.PersistentVolumeClaims = append(resources.PersistentVolumeClaims,
		corev1.PersistentVolumeClaim{
			ObjectMeta: metav1.ObjectMeta{
				Name:      "test-pvc",
				Namespace: resources.Namespaces[0],
			},
			Spec: corev1.PersistentVolumeClaimSpec{
				StorageClassName: &scName,
				AccessModes:      []corev1.PersistentVolumeAccessMode{corev1.ReadWriteOnce},
				Resources: corev1.VolumeResourceRequirements{
					Requests: corev1.ResourceList{
						corev1.ResourceStorage: resource.MustParse("1Gi"),
					},
				},
			},
		},
	)

	for i := range resources.Pods {
		resources.Pods[i].Spec.Volumes = append(resources.Pods[i].Spec.Volumes,
			corev1.Volume{
				Name: "test-vol",
				VolumeSource: corev1.VolumeSource{
					PersistentVolumeClaim: &corev1.PersistentVolumeClaimVolumeSource{
						ClaimName: "test-pvc",
					},
				},
			},
		)
	}

	resources.Nodes = make([]corev1.Node, nodeCount)
	for i := range nodeCount {
		resources.Nodes[i] = corev1.Node{
			ObjectMeta: metav1.ObjectMeta{Name: fmt.Sprintf("worker-%d", i)},
		}
	}
}
