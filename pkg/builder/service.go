package builder

import (
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/util/intstr"
)

type ServiceBuilder struct {
	svc *corev1.Service
}

func NewService(name, namespace string) *ServiceBuilder {
	return &ServiceBuilder{svc: &corev1.Service{
		ObjectMeta: metav1.ObjectMeta{
			Name:      name,
			Namespace: namespace,
		},
		Spec: corev1.ServiceSpec{
			Selector: map[string]string{"app": name},
			Ports: []corev1.ServicePort{{
				Port:       80,
				TargetPort: intstr.FromInt32(8080),
			}},
		},
	}}
}

func (b *ServiceBuilder) WithSelector(key, value string) *ServiceBuilder {
	b.svc.Spec.Selector[key] = value
	return b
}

func (b *ServiceBuilder) WithType(t corev1.ServiceType) *ServiceBuilder {
	b.svc.Spec.Type = t
	return b
}

func (b *ServiceBuilder) WithIPFamilyPolicy(p corev1.IPFamilyPolicy) *ServiceBuilder {
	b.svc.Spec.IPFamilyPolicy = &p
	return b
}

func (b *ServiceBuilder) WithNodePort(port int32) *ServiceBuilder {
	b.svc.Spec.Type = corev1.ServiceTypeNodePort
	b.svc.Spec.Ports[0].NodePort = port
	return b
}

func (b *ServiceBuilder) Build() *corev1.Service {
	return b.svc.DeepCopy()
}
