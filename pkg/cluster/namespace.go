package cluster

import (
	"context"
	"crypto/rand"
	"encoding/hex"
	"fmt"
	"time"

	corev1 "k8s.io/api/core/v1"
	rbacv1 "k8s.io/api/rbac/v1"
	apierrors "k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/util/wait"
	"k8s.io/client-go/kubernetes"
)

func CreateNamespace(ctx context.Context, client kubernetes.Interface, prefix string, privileged bool) (string, error) {
	name := uniqueName(prefix)
	ns := &corev1.Namespace{
		ObjectMeta: metav1.ObjectMeta{
			Name:   name,
			Labels: map[string]string{},
		},
	}
	if privileged {
		ns.Labels["pod-security.kubernetes.io/enforce"] = "privileged"
		ns.Labels["pod-security.kubernetes.io/warn"] = "privileged"
	}
	_, err := client.CoreV1().Namespaces().Create(ctx, ns, metav1.CreateOptions{})
	if err != nil {
		return "", fmt.Errorf("creating namespace %s: %w", name, err)
	}

	if privileged {
		if err := grantPrivilegedSCC(ctx, client, name); err != nil {
			return name, fmt.Errorf("granting privileged SCC in %s: %w", name, err)
		}
	}

	return name, nil
}

func grantPrivilegedSCC(ctx context.Context, client kubernetes.Interface, namespace string) error {
	rb := &rbacv1.RoleBinding{
		ObjectMeta: metav1.ObjectMeta{
			Name:      "privileged-scc",
			Namespace: namespace,
		},
		RoleRef: rbacv1.RoleRef{
			APIGroup: "rbac.authorization.k8s.io",
			Kind:     "ClusterRole",
			Name:     "system:openshift:scc:privileged",
		},
		Subjects: []rbacv1.Subject{
			{
				Kind:      "ServiceAccount",
				Name:      "default",
				Namespace: namespace,
			},
		},
	}
	_, err := client.RbacV1().RoleBindings(namespace).Create(ctx, rb, metav1.CreateOptions{})
	if apierrors.IsNotFound(err) {
		return nil
	}
	return err
}

func DeleteNamespace(ctx context.Context, client kubernetes.Interface, name string) error {
	err := client.CoreV1().Namespaces().Delete(ctx, name, metav1.DeleteOptions{})
	if err != nil {
		return fmt.Errorf("deleting namespace %s: %w", name, err)
	}

	return wait.PollUntilContextTimeout(ctx, 2*time.Second, 2*time.Minute, true,
		func(ctx context.Context) (bool, error) {
			_, err := client.CoreV1().Namespaces().Get(ctx, name, metav1.GetOptions{})
			if apierrors.IsNotFound(err) {
				return true, nil
			}
			return false, err
		},
	)
}

func uniqueName(prefix string) string {
	b := make([]byte, 4)
	_, _ = rand.Read(b)
	return fmt.Sprintf("cqe-%s-%s", prefix, hex.EncodeToString(b))
}
