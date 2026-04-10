package cluster

import (
	"context"
	"fmt"
	"time"

	corev1 "k8s.io/api/core/v1"
	networkingv1 "k8s.io/api/networking/v1"
	policyv1 "k8s.io/api/policy/v1"
	rbacv1 "k8s.io/api/rbac/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/util/wait"
	"k8s.io/client-go/kubernetes"
)

func CreateServiceAccount(ctx context.Context, client kubernetes.Interface, sa *corev1.ServiceAccount) error {
	_, err := client.CoreV1().ServiceAccounts(sa.Namespace).Create(ctx, sa, metav1.CreateOptions{})
	if err != nil {
		return fmt.Errorf("creating service account %s/%s: %w", sa.Namespace, sa.Name, err)
	}
	return nil
}

func CreateRole(ctx context.Context, client kubernetes.Interface, role *rbacv1.Role) error {
	_, err := client.RbacV1().Roles(role.Namespace).Create(ctx, role, metav1.CreateOptions{})
	if err != nil {
		return fmt.Errorf("creating role %s/%s: %w", role.Namespace, role.Name, err)
	}
	return nil
}

func CreateRoleBinding(ctx context.Context, client kubernetes.Interface, rb *rbacv1.RoleBinding) error {
	_, err := client.RbacV1().RoleBindings(rb.Namespace).Create(ctx, rb, metav1.CreateOptions{})
	if err != nil {
		return fmt.Errorf("creating role binding %s/%s: %w", rb.Namespace, rb.Name, err)
	}
	return nil
}

func CreateClusterRoleBinding(ctx context.Context, client kubernetes.Interface, crb *rbacv1.ClusterRoleBinding) error {
	_, err := client.RbacV1().ClusterRoleBindings().Create(ctx, crb, metav1.CreateOptions{})
	if err != nil {
		return fmt.Errorf("creating cluster role binding %s: %w", crb.Name, err)
	}
	return nil
}

func DeleteClusterRoleBinding(ctx context.Context, client kubernetes.Interface, name string) error {
	return client.RbacV1().ClusterRoleBindings().Delete(ctx, name, metav1.DeleteOptions{})
}

func CreateService(ctx context.Context, client kubernetes.Interface, svc *corev1.Service) error {
	_, err := client.CoreV1().Services(svc.Namespace).Create(ctx, svc, metav1.CreateOptions{})
	if err != nil {
		return fmt.Errorf("creating service %s/%s: %w", svc.Namespace, svc.Name, err)
	}
	return nil
}

func CreateResourceQuota(ctx context.Context, client kubernetes.Interface, rq *corev1.ResourceQuota) error {
	_, err := client.CoreV1().ResourceQuotas(rq.Namespace).Create(ctx, rq, metav1.CreateOptions{})
	if err != nil {
		return fmt.Errorf("creating resource quota %s/%s: %w", rq.Namespace, rq.Name, err)
	}
	return nil
}

func CreateNetworkPolicy(ctx context.Context, client kubernetes.Interface, np *networkingv1.NetworkPolicy) error {
	_, err := client.NetworkingV1().NetworkPolicies(np.Namespace).Create(ctx, np, metav1.CreateOptions{})
	if err != nil {
		return fmt.Errorf("creating network policy %s/%s: %w", np.Namespace, np.Name, err)
	}
	return nil
}

func CreatePDB(ctx context.Context, client kubernetes.Interface, pdb *policyv1.PodDisruptionBudget) error {
	_, err := client.PolicyV1().PodDisruptionBudgets(pdb.Namespace).Create(ctx, pdb, metav1.CreateOptions{})
	if err != nil {
		return fmt.Errorf("creating PDB %s/%s: %w", pdb.Namespace, pdb.Name, err)
	}
	return nil
}

func CreateAndWaitForPod(ctx context.Context, client kubernetes.Interface, pod *corev1.Pod, timeout time.Duration) error {
	_, err := client.CoreV1().Pods(pod.Namespace).Create(ctx, pod, metav1.CreateOptions{})
	if err != nil {
		return fmt.Errorf("creating pod %s/%s: %w", pod.Namespace, pod.Name, err)
	}
	return wait.PollUntilContextTimeout(ctx, 2*time.Second, timeout, true,
		func(ctx context.Context) (bool, error) {
			p, err := client.CoreV1().Pods(pod.Namespace).Get(ctx, pod.Name, metav1.GetOptions{})
			if err != nil {
				return false, nil
			}
			for _, c := range p.Status.Conditions {
				if c.Type == corev1.PodReady && c.Status == corev1.ConditionTrue {
					return true, nil
				}
			}
			return false, nil
		},
	)
}
