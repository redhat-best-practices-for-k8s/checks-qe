package builder

import (
	apiextv1 "k8s.io/apiextensions-apiserver/pkg/apis/apiextensions/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

type CRDBuilder struct {
	crd *apiextv1.CustomResourceDefinition
}

func NewCRD(name, group string) *CRDBuilder {
	return &CRDBuilder{crd: &apiextv1.CustomResourceDefinition{
		ObjectMeta: metav1.ObjectMeta{Name: name},
		Spec: apiextv1.CustomResourceDefinitionSpec{
			Group: group,
			Names: apiextv1.CustomResourceDefinitionNames{
				Plural:   name,
				Singular: name,
				Kind:     name,
			},
			Scope: apiextv1.NamespaceScoped,
			Versions: []apiextv1.CustomResourceDefinitionVersion{{
				Name:    "v1",
				Served:  true,
				Storage: true,
				Schema: &apiextv1.CustomResourceValidation{
					OpenAPIV3Schema: &apiextv1.JSONSchemaProps{
						Type: "object",
					},
				},
			}},
		},
	}}
}

func (b *CRDBuilder) WithVersion(name string, withSchema bool) *CRDBuilder {
	v := apiextv1.CustomResourceDefinitionVersion{
		Name:    name,
		Served:  true,
		Storage: false,
	}
	if withSchema {
		v.Schema = &apiextv1.CustomResourceValidation{
			OpenAPIV3Schema: &apiextv1.JSONSchemaProps{Type: "object"},
		}
	}
	b.crd.Spec.Versions = append(b.crd.Spec.Versions, v)
	return b
}

func (b *CRDBuilder) WithoutSchema() *CRDBuilder {
	for i := range b.crd.Spec.Versions {
		b.crd.Spec.Versions[i].Schema = nil
	}
	return b
}

func (b *CRDBuilder) Build() apiextv1.CustomResourceDefinition {
	return *b.crd.DeepCopy()
}
