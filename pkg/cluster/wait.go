package cluster

import (
	"context"
	"fmt"
	"time"

	appsv1 "k8s.io/api/apps/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/util/wait"
	"k8s.io/client-go/kubernetes"
)

const DefaultTimeout = 2 * time.Minute

func CreateAndWaitForDeployment(ctx context.Context, client kubernetes.Interface, dep *appsv1.Deployment, timeout time.Duration) error {
	created, err := client.AppsV1().Deployments(dep.Namespace).Create(ctx, dep, metav1.CreateOptions{})
	if err != nil {
		return fmt.Errorf("creating deployment %s/%s: %w", dep.Namespace, dep.Name, err)
	}
	return waitForDeploymentReady(ctx, client, created.Namespace, created.Name, timeout)
}

func waitForDeploymentReady(ctx context.Context, client kubernetes.Interface, namespace, name string, timeout time.Duration) error {
	return wait.PollUntilContextTimeout(ctx, 2*time.Second, timeout, true,
		func(ctx context.Context) (bool, error) {
			dep, err := client.AppsV1().Deployments(namespace).Get(ctx, name, metav1.GetOptions{})
			if err != nil {
				return false, nil
			}
			desired := int32(1)
			if dep.Spec.Replicas != nil {
				desired = *dep.Spec.Replicas
			}
			return dep.Status.ReadyReplicas >= desired, nil
		},
	)
}

func CreateAndWaitForStatefulSet(ctx context.Context, client kubernetes.Interface, sts *appsv1.StatefulSet, timeout time.Duration) error {
	created, err := client.AppsV1().StatefulSets(sts.Namespace).Create(ctx, sts, metav1.CreateOptions{})
	if err != nil {
		return fmt.Errorf("creating statefulset %s/%s: %w", sts.Namespace, sts.Name, err)
	}
	return wait.PollUntilContextTimeout(ctx, 2*time.Second, timeout, true,
		func(ctx context.Context) (bool, error) {
			s, err := client.AppsV1().StatefulSets(created.Namespace).Get(ctx, created.Name, metav1.GetOptions{})
			if err != nil {
				return false, nil
			}
			desired := int32(1)
			if s.Spec.Replicas != nil {
				desired = *s.Spec.Replicas
			}
			return s.Status.ReadyReplicas >= desired, nil
		},
	)
}

func CreateAndWaitForDaemonSet(ctx context.Context, client kubernetes.Interface, ds *appsv1.DaemonSet, timeout time.Duration) error {
	created, err := client.AppsV1().DaemonSets(ds.Namespace).Create(ctx, ds, metav1.CreateOptions{})
	if err != nil {
		return fmt.Errorf("creating daemonset %s/%s: %w", ds.Namespace, ds.Name, err)
	}
	return wait.PollUntilContextTimeout(ctx, 2*time.Second, timeout, true,
		func(ctx context.Context) (bool, error) {
			d, err := client.AppsV1().DaemonSets(created.Namespace).Get(ctx, created.Name, metav1.GetOptions{})
			if err != nil {
				return false, nil
			}
			return d.Status.NumberReady > 0 && d.Status.NumberReady == d.Status.DesiredNumberScheduled, nil
		},
	)
}
