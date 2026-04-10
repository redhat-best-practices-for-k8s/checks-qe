package builder

import (
	policyv1 "k8s.io/api/policy/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/util/intstr"
)

type PDBBuilder struct {
	pdb *policyv1.PodDisruptionBudget
}

func NewPDB(name, namespace string) *PDBBuilder {
	return &PDBBuilder{pdb: &policyv1.PodDisruptionBudget{
		ObjectMeta: metav1.ObjectMeta{
			Name:      name,
			Namespace: namespace,
		},
		Spec: policyv1.PodDisruptionBudgetSpec{},
	}}
}

func (b *PDBBuilder) WithMinAvailable(val int32) *PDBBuilder {
	v := intstr.FromInt32(val)
	b.pdb.Spec.MinAvailable = &v
	return b
}

func (b *PDBBuilder) WithMaxUnavailable(val int32) *PDBBuilder {
	v := intstr.FromInt32(val)
	b.pdb.Spec.MaxUnavailable = &v
	return b
}

func (b *PDBBuilder) WithSelector(labels map[string]string) *PDBBuilder {
	b.pdb.Spec.Selector = &metav1.LabelSelector{MatchLabels: labels}
	return b
}

func (b *PDBBuilder) Build() *policyv1.PodDisruptionBudget {
	return b.pdb.DeepCopy()
}
