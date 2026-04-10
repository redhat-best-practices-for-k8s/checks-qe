package builder

import (
	"k8s.io/utils/ptr"

	appsv1 "k8s.io/api/apps/v1"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

type StatefulSetBuilder struct {
	sts *appsv1.StatefulSet
}

func NewStatefulSet(name, namespace string) *StatefulSetBuilder {
	labels := map[string]string{"app": name}
	return &StatefulSetBuilder{sts: &appsv1.StatefulSet{
		ObjectMeta: metav1.ObjectMeta{
			Name:      name,
			Namespace: namespace,
			Labels:    labels,
		},
		Spec: appsv1.StatefulSetSpec{
			Replicas: ptr.To[int32](1),
			Selector: &metav1.LabelSelector{MatchLabels: labels},
			Template: corev1.PodTemplateSpec{
				ObjectMeta: metav1.ObjectMeta{Labels: labels},
				Spec: corev1.PodSpec{
					TerminationGracePeriodSeconds: ptr.To[int64](0),
					Containers: []corev1.Container{{
						Name:    "test",
						Image:   DefaultImage,
						Command: []string{"/bin/sh", "-c", "sleep infinity"},
					}},
				},
			},
		},
	}}
}

func (b *StatefulSetBuilder) WithReplicas(n int32) *StatefulSetBuilder {
	b.sts.Spec.Replicas = ptr.To(n)
	return b
}

func (b *StatefulSetBuilder) WithLabel(key, value string) *StatefulSetBuilder {
	b.sts.Labels[key] = value
	b.sts.Spec.Template.Labels[key] = value
	return b
}

func (b *StatefulSetBuilder) Build() *appsv1.StatefulSet {
	return b.sts.DeepCopy()
}
