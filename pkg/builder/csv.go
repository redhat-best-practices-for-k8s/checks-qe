package builder

import (
	"github.com/blang/semver/v4"
	olmversion "github.com/operator-framework/api/pkg/lib/version"
	olmv1alpha1 "github.com/operator-framework/api/pkg/operators/v1alpha1"
	rbacv1 "k8s.io/api/rbac/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

type CSVBuilder struct {
	csv *olmv1alpha1.ClusterServiceVersion
}

func NewCSV(name, namespace string) *CSVBuilder {
	return &CSVBuilder{csv: &olmv1alpha1.ClusterServiceVersion{
		ObjectMeta: metav1.ObjectMeta{
			Name:        name,
			Namespace:   namespace,
			Annotations: map[string]string{},
		},
		Spec: olmv1alpha1.ClusterServiceVersionSpec{
			Version: olmversion.OperatorVersion{Version: semver.MustParse("1.0.0")},
			InstallModes: []olmv1alpha1.InstallMode{
				{Type: olmv1alpha1.InstallModeTypeOwnNamespace, Supported: true},
				{Type: olmv1alpha1.InstallModeTypeSingleNamespace, Supported: true},
				{Type: olmv1alpha1.InstallModeTypeMultiNamespace, Supported: false},
				{Type: olmv1alpha1.InstallModeTypeAllNamespaces, Supported: false},
			},
		},
		Status: olmv1alpha1.ClusterServiceVersionStatus{
			Phase: olmv1alpha1.CSVPhaseSucceeded,
		},
	}}
}

func (b *CSVBuilder) WithPhase(phase olmv1alpha1.ClusterServiceVersionPhase) *CSVBuilder {
	b.csv.Status.Phase = phase
	return b
}

func (b *CSVBuilder) WithVersion(v string) *CSVBuilder {
	b.csv.Spec.Version = olmversion.OperatorVersion{Version: semver.MustParse(v)}
	return b
}

func (b *CSVBuilder) WithAnnotation(key, value string) *CSVBuilder {
	b.csv.Annotations[key] = value
	return b
}

func (b *CSVBuilder) WithOLMInstalled() *CSVBuilder {
	b.csv.Annotations["olm.operatorNamespace"] = b.csv.Namespace
	return b
}

func (b *CSVBuilder) WithSkipRange(r string) *CSVBuilder {
	b.csv.Annotations["olm.skipRange"] = r
	return b
}

func (b *CSVBuilder) WithOwnedCRD(crdName string) *CSVBuilder {
	b.csv.Spec.CustomResourceDefinitions.Owned = append(
		b.csv.Spec.CustomResourceDefinitions.Owned,
		olmv1alpha1.CRDDescription{Name: crdName},
	)
	return b
}

func (b *CSVBuilder) WithClusterPermissionSCC() *CSVBuilder {
	b.csv.Spec.InstallStrategy.StrategySpec.ClusterPermissions = append(
		b.csv.Spec.InstallStrategy.StrategySpec.ClusterPermissions,
		olmv1alpha1.StrategyDeploymentPermissions{
			Rules: []rbacv1.PolicyRule{{
				APIGroups: []string{"security.openshift.io"},
				Resources: []string{"securitycontextconstraints"},
				Verbs:     []string{"use"},
			}},
		},
	)
	return b
}

func (b *CSVBuilder) WithInstallMode(mode olmv1alpha1.InstallModeType, supported bool) *CSVBuilder {
	for i := range b.csv.Spec.InstallModes {
		if b.csv.Spec.InstallModes[i].Type == mode {
			b.csv.Spec.InstallModes[i].Supported = supported
			return b
		}
	}
	b.csv.Spec.InstallModes = append(b.csv.Spec.InstallModes,
		olmv1alpha1.InstallMode{Type: mode, Supported: supported})
	return b
}

func (b *CSVBuilder) Build() olmv1alpha1.ClusterServiceVersion {
	return *b.csv.DeepCopy()
}
