package discovery

import (
	"context"
	"fmt"

	checks "github.com/redhat-best-practices-for-k8s/checks"
	corev1 "k8s.io/api/core/v1"
	storagev1 "k8s.io/api/storage/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
)

func Targeted(ctx context.Context, client kubernetes.Interface, namespace string) (*checks.DiscoveredResources, error) {
	resources := &checks.DiscoveredResources{
		Namespaces: []string{namespace},
	}

	pods, err := client.CoreV1().Pods(namespace).List(ctx, metav1.ListOptions{})
	if err != nil {
		return nil, fmt.Errorf("listing pods: %w", err)
	}
	resources.Pods = pods.Items

	deps, err := client.AppsV1().Deployments(namespace).List(ctx, metav1.ListOptions{})
	if err != nil {
		return nil, fmt.Errorf("listing deployments: %w", err)
	}
	resources.Deployments = deps.Items

	sts, err := client.AppsV1().StatefulSets(namespace).List(ctx, metav1.ListOptions{})
	if err != nil {
		return nil, fmt.Errorf("listing statefulsets: %w", err)
	}
	resources.StatefulSets = sts.Items

	ds, err := client.AppsV1().DaemonSets(namespace).List(ctx, metav1.ListOptions{})
	if err != nil {
		return nil, fmt.Errorf("listing daemonsets: %w", err)
	}
	resources.DaemonSets = ds.Items

	svcs, err := client.CoreV1().Services(namespace).List(ctx, metav1.ListOptions{})
	if err != nil {
		return nil, fmt.Errorf("listing services: %w", err)
	}
	resources.Services = svcs.Items

	sas, err := client.CoreV1().ServiceAccounts(namespace).List(ctx, metav1.ListOptions{})
	if err != nil {
		return nil, fmt.Errorf("listing service accounts: %w", err)
	}
	resources.ServiceAccounts = sas.Items

	roles, err := client.RbacV1().Roles(namespace).List(ctx, metav1.ListOptions{})
	if err != nil {
		return nil, fmt.Errorf("listing roles: %w", err)
	}
	resources.Roles = roles.Items

	rbs, err := client.RbacV1().RoleBindings(namespace).List(ctx, metav1.ListOptions{})
	if err != nil {
		return nil, fmt.Errorf("listing role bindings: %w", err)
	}
	resources.RoleBindings = rbs.Items

	nps, err := client.NetworkingV1().NetworkPolicies(namespace).List(ctx, metav1.ListOptions{})
	if err != nil {
		return nil, fmt.Errorf("listing network policies: %w", err)
	}
	resources.NetworkPolicies = nps.Items

	rqs, err := client.CoreV1().ResourceQuotas(namespace).List(ctx, metav1.ListOptions{})
	if err != nil {
		return nil, fmt.Errorf("listing resource quotas: %w", err)
	}
	resources.ResourceQuotas = rqs.Items

	pdbs, err := client.PolicyV1().PodDisruptionBudgets(namespace).List(ctx, metav1.ListOptions{})
	if err != nil {
		return nil, fmt.Errorf("listing pod disruption budgets: %w", err)
	}
	resources.PodDisruptionBudgets = pdbs.Items

	return resources, nil
}

type ClusterSnapshot struct {
	Nodes             []corev1.Node
	PersistentVolumes []corev1.PersistentVolume
	StorageClasses    []storagev1.StorageClass
	IsDualStack       bool
	HasMultus         bool
}

func FetchClusterSnapshot(ctx context.Context, client kubernetes.Interface) (*ClusterSnapshot, error) {
	nodes, err := client.CoreV1().Nodes().List(ctx, metav1.ListOptions{})
	if err != nil {
		return nil, fmt.Errorf("listing nodes: %w", err)
	}
	pvs, err := client.CoreV1().PersistentVolumes().List(ctx, metav1.ListOptions{})
	if err != nil {
		return nil, fmt.Errorf("listing persistent volumes: %w", err)
	}
	scs, err := client.StorageV1().StorageClasses().List(ctx, metav1.ListOptions{})
	if err != nil {
		return nil, fmt.Errorf("listing storage classes: %w", err)
	}

	dualStack := false
	kubeSvc, err := client.CoreV1().Services("default").Get(ctx, "kubernetes", metav1.GetOptions{})
	if err == nil && len(kubeSvc.Spec.IPFamilies) > 1 {
		dualStack = true
	}

	hasMultus := false
	_, err = client.Discovery().ServerResourcesForGroupVersion("k8s.cni.cncf.io/v1")
	if err == nil {
		hasMultus = true
	}

	return &ClusterSnapshot{
		Nodes:             nodes.Items,
		PersistentVolumes: pvs.Items,
		StorageClasses:    scs.Items,
		IsDualStack:       dualStack,
		HasMultus:         hasMultus,
	}, nil
}

func ApplyClusterSnapshot(snap *ClusterSnapshot, resources *checks.DiscoveredResources) {
	resources.Nodes = snap.Nodes
	resources.PersistentVolumes = snap.PersistentVolumes
	resources.StorageClasses = snap.StorageClasses
}

func WithClusterRoleBindings(ctx context.Context, client kubernetes.Interface, resources *checks.DiscoveredResources) error {
	crbs, err := client.RbacV1().ClusterRoleBindings().List(ctx, metav1.ListOptions{})
	if err != nil {
		return fmt.Errorf("listing cluster role bindings: %w", err)
	}
	resources.ClusterRoleBindings = crbs.Items
	return nil
}
