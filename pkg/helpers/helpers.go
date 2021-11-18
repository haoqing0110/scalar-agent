package helpers

import (
	"context"

	"k8s.io/apimachinery/pkg/api/equality"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/util/retry"
	clusterclientset "open-cluster-management.io/api/client/cluster/clientset/versioned"
	clusterv1alpha1 "open-cluster-management.io/api/cluster/v1alpha1"
)

type UpdateManagedClusterScalarStatusFunc func(status *clusterv1alpha1.ManagedClusterScalarStatus) error

func UpdateManagedClusterScalarStatus(
	ctx context.Context,
	client clusterclientset.Interface,
	spokeClusterName string,
	crName string,
	updateFuncs ...UpdateManagedClusterScalarStatusFunc) (*clusterv1alpha1.ManagedClusterScalarStatus, bool, error) {
	updated := false
	var updatedManagedClusterScalarStatus *clusterv1alpha1.ManagedClusterScalarStatus

	err := retry.RetryOnConflict(retry.DefaultBackoff, func() error {
		managedClusterScalar, err := client.ClusterV1alpha1().ManagedClusterScalars(spokeClusterName).Get(ctx, crName, metav1.GetOptions{})
		if err != nil {
			return err
		}
		oldStatus := &managedClusterScalar.Status

		newStatus := oldStatus.DeepCopy()
		for _, update := range updateFuncs {
			if err := update(newStatus); err != nil {
				return err
			}
		}
		if equality.Semantic.DeepEqual(oldStatus, newStatus) {
			// We return the newStatus which is a deep copy of oldStatus but with all update funcs applied.
			updatedManagedClusterScalarStatus = newStatus
			return nil
		}

		managedClusterScalar.Status = *newStatus
		updatedManagedClusterScalar, err := client.ClusterV1alpha1().ManagedClusterScalars(spokeClusterName).UpdateStatus(ctx, managedClusterScalar, metav1.UpdateOptions{})
		if err != nil {
			return err
		}
		updatedManagedClusterScalarStatus = &updatedManagedClusterScalar.Status
		updated = err == nil
		return err
	})

	return updatedManagedClusterScalarStatus, updated, err
}
