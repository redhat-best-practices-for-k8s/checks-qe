package cluster

import (
	"context"
	"fmt"

	appsv1 "k8s.io/api/apps/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/types"
	"k8s.io/client-go/rest"
	"sigs.k8s.io/controller-runtime/pkg/client"

	bpsv1alpha1 "github.com/sebrandon1/bps-operator/api/v1alpha1"
)

func NewOperatorClient(config *rest.Config) (client.Client, error) {
	s := runtime.NewScheme()
	if err := appsv1.AddToScheme(s); err != nil {
		return nil, err
	}
	if err := bpsv1alpha1.AddToScheme(s); err != nil {
		return nil, err
	}
	return client.New(config, client.Options{Scheme: s})
}

func VerifyOperatorRunning(ctx context.Context, c client.Client, namespace string) error {
	var deploy appsv1.Deployment
	key := types.NamespacedName{
		Name:      "bps-operator-controller-manager",
		Namespace: namespace,
	}
	if err := c.Get(ctx, key, &deploy); err != nil {
		return fmt.Errorf("operator deployment not found in %s: %w", namespace, err)
	}
	if deploy.Status.AvailableReplicas == 0 {
		return fmt.Errorf("operator deployment %s/%s has no available replicas", namespace, deploy.Name)
	}
	return nil
}
