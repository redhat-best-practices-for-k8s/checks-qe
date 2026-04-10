package builder

import (
	"k8s.io/utils/ptr"

	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

type PodBuilder struct {
	pod *corev1.Pod
}

func NewPod(name, namespace string) *PodBuilder {
	labels := map[string]string{"app": name}
	return &PodBuilder{pod: &corev1.Pod{
		ObjectMeta: metav1.ObjectMeta{
			Name:      name,
			Namespace: namespace,
			Labels:    labels,
		},
		Spec: corev1.PodSpec{
			TerminationGracePeriodSeconds: ptr.To[int64](0),
			Containers: []corev1.Container{{
				Name:    "test",
				Image:   DefaultImage,
				Command: []string{"/bin/sh", "-c", "sleep infinity"},
			}},
		},
	}}
}

func (b *PodBuilder) WithHostNetwork(v bool) *PodBuilder {
	b.pod.Spec.HostNetwork = v
	return b
}

func (b *PodBuilder) WithHostIPC(v bool) *PodBuilder {
	b.pod.Spec.HostIPC = v
	return b
}

func (b *PodBuilder) WithHostPID(v bool) *PodBuilder {
	b.pod.Spec.HostPID = v
	return b
}

func (b *PodBuilder) WithCapability(cap corev1.Capability) *PodBuilder {
	c := &b.pod.Spec.Containers[0]
	if c.SecurityContext == nil {
		c.SecurityContext = &corev1.SecurityContext{}
	}
	if c.SecurityContext.Capabilities == nil {
		c.SecurityContext.Capabilities = &corev1.Capabilities{}
	}
	c.SecurityContext.Capabilities.Add = append(c.SecurityContext.Capabilities.Add, cap)
	return b
}

func (b *PodBuilder) WithLabel(key, value string) *PodBuilder {
	b.pod.Labels[key] = value
	return b
}

func (b *PodBuilder) Build() *corev1.Pod {
	return b.pod.DeepCopy()
}
