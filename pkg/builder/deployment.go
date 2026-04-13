package builder

import (
	"k8s.io/apimachinery/pkg/api/resource"
	"k8s.io/utils/ptr"

	appsv1 "k8s.io/api/apps/v1"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

const DefaultImage = "registry.access.redhat.com/ubi9/ubi-micro:latest"

type DeploymentBuilder struct {
	dep *appsv1.Deployment
}

func NewDeployment(name, namespace string) *DeploymentBuilder {
	labels := map[string]string{"app": name}
	return &DeploymentBuilder{dep: &appsv1.Deployment{
		ObjectMeta: metav1.ObjectMeta{
			Name:      name,
			Namespace: namespace,
			Labels:    labels,
		},
		Spec: appsv1.DeploymentSpec{
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

func (b *DeploymentBuilder) WithReplicas(n int32) *DeploymentBuilder {
	b.dep.Spec.Replicas = ptr.To(n)
	return b
}

func (b *DeploymentBuilder) WithImage(image string) *DeploymentBuilder {
	b.dep.Spec.Template.Spec.Containers[0].Image = image
	return b
}

func (b *DeploymentBuilder) WithHostNetwork(v bool) *DeploymentBuilder {
	b.dep.Spec.Template.Spec.HostNetwork = v
	return b
}

func (b *DeploymentBuilder) WithHostIPC(v bool) *DeploymentBuilder {
	b.dep.Spec.Template.Spec.HostIPC = v
	return b
}

func (b *DeploymentBuilder) WithHostPID(v bool) *DeploymentBuilder {
	b.dep.Spec.Template.Spec.HostPID = v
	return b
}

func (b *DeploymentBuilder) WithCapability(cap corev1.Capability) *DeploymentBuilder {
	sc := b.secCtx()
	if sc.Capabilities == nil {
		sc.Capabilities = &corev1.Capabilities{}
	}
	sc.Capabilities.Add = append(sc.Capabilities.Add, cap)
	return b
}

func (b *DeploymentBuilder) WithRunAsNonRoot() *DeploymentBuilder {
	b.secCtx().RunAsNonRoot = ptr.To(true)
	return b
}

func (b *DeploymentBuilder) WithRunAsUser(uid int64) *DeploymentBuilder {
	b.secCtx().RunAsUser = ptr.To(uid)
	return b
}

func (b *DeploymentBuilder) WithPrivileged(v bool) *DeploymentBuilder {
	b.secCtx().Privileged = ptr.To(v)
	return b
}

func (b *DeploymentBuilder) WithReadOnlyRootFS() *DeploymentBuilder {
	b.secCtx().ReadOnlyRootFilesystem = ptr.To(true)
	return b
}

func (b *DeploymentBuilder) WithAllowPrivilegeEscalation(v bool) *DeploymentBuilder {
	b.secCtx().AllowPrivilegeEscalation = ptr.To(v)
	return b
}

func (b *DeploymentBuilder) secCtx() *corev1.SecurityContext {
	c := &b.dep.Spec.Template.Spec.Containers[0]
	if c.SecurityContext == nil {
		c.SecurityContext = &corev1.SecurityContext{}
	}
	return c.SecurityContext
}

func (b *DeploymentBuilder) WithLivenessProbe(probe *corev1.Probe) *DeploymentBuilder {
	b.dep.Spec.Template.Spec.Containers[0].LivenessProbe = probe
	return b
}

func (b *DeploymentBuilder) WithReadinessProbe(probe *corev1.Probe) *DeploymentBuilder {
	b.dep.Spec.Template.Spec.Containers[0].ReadinessProbe = probe
	return b
}

func (b *DeploymentBuilder) WithStartupProbe(probe *corev1.Probe) *DeploymentBuilder {
	b.dep.Spec.Template.Spec.Containers[0].StartupProbe = probe
	return b
}

func (b *DeploymentBuilder) WithPreStopHook(handler *corev1.LifecycleHandler) *DeploymentBuilder {
	b.lifecycle().PreStop = handler
	return b
}

func (b *DeploymentBuilder) WithPostStartHook(handler *corev1.LifecycleHandler) *DeploymentBuilder {
	b.lifecycle().PostStart = handler
	return b
}

func (b *DeploymentBuilder) lifecycle() *corev1.Lifecycle {
	c := &b.dep.Spec.Template.Spec.Containers[0]
	if c.Lifecycle == nil {
		c.Lifecycle = &corev1.Lifecycle{}
	}
	return c.Lifecycle
}

func (b *DeploymentBuilder) WithContainerPort(port int32) *DeploymentBuilder {
	b.dep.Spec.Template.Spec.Containers[0].Ports = append(
		b.dep.Spec.Template.Spec.Containers[0].Ports,
		corev1.ContainerPort{ContainerPort: port},
	)
	return b
}

func (b *DeploymentBuilder) WithNamedContainerPort(name string, port int32) *DeploymentBuilder {
	b.dep.Spec.Template.Spec.Containers[0].Ports = append(
		b.dep.Spec.Template.Spec.Containers[0].Ports,
		corev1.ContainerPort{Name: name, ContainerPort: port},
	)
	return b
}

func (b *DeploymentBuilder) WithHostPort(containerPort, hostPort int32) *DeploymentBuilder {
	b.dep.Spec.Template.Spec.Containers[0].Ports = append(
		b.dep.Spec.Template.Spec.Containers[0].Ports,
		corev1.ContainerPort{ContainerPort: containerPort, HostPort: hostPort},
	)
	return b
}

func (b *DeploymentBuilder) WithServiceAccountName(name string) *DeploymentBuilder {
	b.dep.Spec.Template.Spec.ServiceAccountName = name
	return b
}

func (b *DeploymentBuilder) WithAutomountServiceAccountToken(v bool) *DeploymentBuilder {
	b.dep.Spec.Template.Spec.AutomountServiceAccountToken = ptr.To(v)
	return b
}

func (b *DeploymentBuilder) WithPodAnnotation(key, value string) *DeploymentBuilder {
	if b.dep.Spec.Template.Annotations == nil {
		b.dep.Spec.Template.Annotations = map[string]string{}
	}
	b.dep.Spec.Template.Annotations[key] = value
	return b
}

func (b *DeploymentBuilder) WithLabel(key, value string) *DeploymentBuilder {
	b.dep.Labels[key] = value
	b.dep.Spec.Template.Labels[key] = value
	return b
}

func (b *DeploymentBuilder) WithImagePullPolicy(policy corev1.PullPolicy) *DeploymentBuilder {
	b.dep.Spec.Template.Spec.Containers[0].ImagePullPolicy = policy
	return b
}

func (b *DeploymentBuilder) WithHostPathVolume(name, path string) *DeploymentBuilder {
	b.dep.Spec.Template.Spec.Volumes = append(b.dep.Spec.Template.Spec.Volumes,
		corev1.Volume{
			Name:         name,
			VolumeSource: corev1.VolumeSource{HostPath: &corev1.HostPathVolumeSource{Path: path}},
		},
	)
	b.dep.Spec.Template.Spec.Containers[0].VolumeMounts = append(
		b.dep.Spec.Template.Spec.Containers[0].VolumeMounts,
		corev1.VolumeMount{Name: name, MountPath: path},
	)
	return b
}

func (b *DeploymentBuilder) WithResourceRequests(cpu, memory string) *DeploymentBuilder {
	c := &b.dep.Spec.Template.Spec.Containers[0]
	if c.Resources.Requests == nil {
		c.Resources.Requests = corev1.ResourceList{}
	}
	if cpu != "" {
		c.Resources.Requests[corev1.ResourceCPU] = resource.MustParse(cpu)
	}
	if memory != "" {
		c.Resources.Requests[corev1.ResourceMemory] = resource.MustParse(memory)
	}
	return b
}

func (b *DeploymentBuilder) WithResourceLimits(cpu, memory string) *DeploymentBuilder {
	c := &b.dep.Spec.Template.Spec.Containers[0]
	if c.Resources.Limits == nil {
		c.Resources.Limits = corev1.ResourceList{}
	}
	if cpu != "" {
		c.Resources.Limits[corev1.ResourceCPU] = resource.MustParse(cpu)
	}
	if memory != "" {
		c.Resources.Limits[corev1.ResourceMemory] = resource.MustParse(memory)
	}
	return b
}

func (b *DeploymentBuilder) WithHugepagesRequest(pageSize, amount string) *DeploymentBuilder {
	c := &b.dep.Spec.Template.Spec.Containers[0]
	if c.Resources.Requests == nil {
		c.Resources.Requests = corev1.ResourceList{}
	}
	if c.Resources.Limits == nil {
		c.Resources.Limits = corev1.ResourceList{}
	}
	key := corev1.ResourceName("hugepages-" + pageSize)
	c.Resources.Requests[key] = resource.MustParse(amount)
	c.Resources.Limits[key] = resource.MustParse(amount)
	return b
}

func (b *DeploymentBuilder) WithShareProcessNamespace(v bool) *DeploymentBuilder {
	b.dep.Spec.Template.Spec.ShareProcessNamespace = ptr.To(v)
	return b
}

func (b *DeploymentBuilder) WithPodRunAsUser(uid int64) *DeploymentBuilder {
	if b.dep.Spec.Template.Spec.SecurityContext == nil {
		b.dep.Spec.Template.Spec.SecurityContext = &corev1.PodSecurityContext{}
	}
	b.dep.Spec.Template.Spec.SecurityContext.RunAsUser = ptr.To(uid)
	return b
}

func (b *DeploymentBuilder) WithNodeSelector(key, value string) *DeploymentBuilder {
	if b.dep.Spec.Template.Spec.NodeSelector == nil {
		b.dep.Spec.Template.Spec.NodeSelector = map[string]string{}
	}
	b.dep.Spec.Template.Spec.NodeSelector[key] = value
	return b
}

func (b *DeploymentBuilder) WithToleration(key, value string, effect corev1.TaintEffect) *DeploymentBuilder {
	b.dep.Spec.Template.Spec.Tolerations = append(b.dep.Spec.Template.Spec.Tolerations,
		corev1.Toleration{Key: key, Value: value, Effect: effect, Operator: corev1.TolerationOpEqual},
	)
	return b
}

func (b *DeploymentBuilder) WithTopologySpreadConstraint(key string, maxSkew int32, whenUnsatisfiable corev1.UnsatisfiableConstraintAction) *DeploymentBuilder {
	b.dep.Spec.Template.Spec.TopologySpreadConstraints = append(
		b.dep.Spec.Template.Spec.TopologySpreadConstraints,
		corev1.TopologySpreadConstraint{
			MaxSkew:           maxSkew,
			TopologyKey:       key,
			WhenUnsatisfiable: whenUnsatisfiable,
			LabelSelector:     b.dep.Spec.Selector,
		},
	)
	return b
}

func (b *DeploymentBuilder) WithSecondContainer(name, image string) *DeploymentBuilder {
	b.dep.Spec.Template.Spec.Containers = append(b.dep.Spec.Template.Spec.Containers,
		corev1.Container{
			Name:    name,
			Image:   image,
			Command: []string{"/bin/sh", "-c", "sleep infinity"},
		},
	)
	return b
}

func (b *DeploymentBuilder) WithCommand(command ...string) *DeploymentBuilder {
	b.dep.Spec.Template.Spec.Containers[0].Command = command
	return b
}

func (b *DeploymentBuilder) WithTerminationMessagePolicy(p corev1.TerminationMessagePolicy) *DeploymentBuilder {
	b.dep.Spec.Template.Spec.Containers[0].TerminationMessagePolicy = p
	return b
}

func (b *DeploymentBuilder) Build() *appsv1.Deployment {
	return b.dep.DeepCopy()
}
