package accesscontrol

import (
	checks "github.com/redhat-best-practices-for-k8s/checks"
	"github.com/redhat-best-practices-for-k8s/checks-qe/pkg/builder"
	"github.com/redhat-best-practices-for-k8s/checks-qe/pkg/scenario"
	rbacv1 "k8s.io/api/rbac/v1"
)

func registerCRDRoles() {
	scenario.Register(
		scenario.Scenario{
			Name:           "accesscontrol/crd-roles/compliant",
			CheckName:      "access-control-crd-roles",
			Category:       checks.CategoryAccessControl,
			Description:    "Role granting access only to CRD resources should be compliant",
			ExpectedStatus: checks.StatusCompliant,
			Setup:          scenario.VanillaDeploymentSetup,
			PostDiscovery: func(resources *checks.DiscoveredResources) {
				crd := builder.NewCRD("widgets", "example.com").Build()
				resources.CRDs = append(resources.CRDs, crd)
				resources.Roles = append(resources.Roles, rbacv1.Role{})
				resources.Roles[len(resources.Roles)-1].Name = "widget-role"
				resources.Roles[len(resources.Roles)-1].Namespace = resources.Namespaces[0]
				resources.Roles[len(resources.Roles)-1].Rules = []rbacv1.PolicyRule{{
					APIGroups: []string{"example.com"},
					Resources: []string{"widgets"},
					Verbs:     []string{"get", "list", "watch"},
				}}
			},
		},
		scenario.Scenario{
			Name:           "accesscontrol/crd-roles/non-compliant-multiple-api-groups",
			CheckName:      "access-control-crd-roles",
			Category:       checks.CategoryAccessControl,
			Description:    "Role with CRD access plus extra API groups should be non-compliant",
			ExpectedStatus: checks.StatusNonCompliant,
			Setup:          scenario.VanillaDeploymentSetup,
			PostDiscovery: func(resources *checks.DiscoveredResources) {
				crd := builder.NewCRD("widgets", "example.com").Build()
				resources.CRDs = append(resources.CRDs, crd)
				resources.Roles = append(resources.Roles, rbacv1.Role{})
				role := &resources.Roles[len(resources.Roles)-1]
				role.Name = "overpermissioned-role"
				role.Namespace = resources.Namespaces[0]
				role.Rules = []rbacv1.PolicyRule{
					{APIGroups: []string{"example.com"}, Resources: []string{"widgets"}, Verbs: []string{"get"}},
					{APIGroups: []string{""}, Resources: []string{"pods"}, Verbs: []string{"get"}},
				}
			},
		},
		scenario.Scenario{
			Name:           "accesscontrol/crd-roles/non-compliant-multiple-resources",
			CheckName:      "access-control-crd-roles",
			Category:       checks.CategoryAccessControl,
			Description:    "Role with CRD access plus non-CRD resources should be non-compliant",
			ExpectedStatus: checks.StatusNonCompliant,
			Setup:          scenario.VanillaDeploymentSetup,
			PostDiscovery: func(resources *checks.DiscoveredResources) {
				crd := builder.NewCRD("widgets", "example.com").Build()
				resources.CRDs = append(resources.CRDs, crd)
				resources.Roles = append(resources.Roles, rbacv1.Role{})
				role := &resources.Roles[len(resources.Roles)-1]
				role.Name = "mixed-role"
				role.Namespace = resources.Namespaces[0]
				role.Rules = []rbacv1.PolicyRule{{
					APIGroups: []string{"example.com"},
					Resources: []string{"widgets", "somethingelse"},
					Verbs:     []string{"get"},
				}}
			},
		},
	)
}
