package builder

import (
	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/resource"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

func NewResourceQuota(name, namespace, cpu, memory string) *corev1.ResourceQuota {
	return &corev1.ResourceQuota{
		ObjectMeta: metav1.ObjectMeta{Name: name, Namespace: namespace},
		Spec: corev1.ResourceQuotaSpec{
			Hard: corev1.ResourceList{
				corev1.ResourceRequestsCPU:    resource.MustParse(cpu),
				corev1.ResourceRequestsMemory: resource.MustParse(memory),
				corev1.ResourceLimitsCPU:      resource.MustParse(cpu),
				corev1.ResourceLimitsMemory:   resource.MustParse(memory),
			},
		},
	}
}
