package helpers

import (
	"context"

	"k8s.io/apimachinery/pkg/api/equality"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/util/retry"
	clusterclientset "open-cluster-management.io/api/client/cluster/clientset/versioned"
	clusterv1alpha1 "open-cluster-management.io/api/cluster/v1alpha1"
)

type UpdateManagedClusterScoreStatusFunc func(status *clusterv1alpha1.ManagedClusterScoreStatus) error

func UpdateManagedClusterScoreStatus(
	ctx context.Context,
	client clusterclientset.Interface,
	spokeClusterName string,
	updateFuncs ...UpdateManagedClusterScoreStatusFunc) (*clusterv1alpha1.ManagedClusterScoreStatus, bool, error) {
	updated := false
	var updatedManagedClusterScoreStatus *clusterv1alpha1.ManagedClusterScoreStatus

	err := retry.RetryOnConflict(retry.DefaultBackoff, func() error {
		managedClusterScore, err := client.ClusterV1alpha1().ManagedClusterScores(spokeClusterName).Get(ctx, spokeClusterName, metav1.GetOptions{})
		if err != nil {
			return err
		}
		oldStatus := &managedClusterScore.Status

		newStatus := oldStatus.DeepCopy()
		for _, update := range updateFuncs {
			if err := update(newStatus); err != nil {
				return err
			}
		}
		if equality.Semantic.DeepEqual(oldStatus, newStatus) {
			// We return the newStatus which is a deep copy of oldStatus but with all update funcs applied.
			updatedManagedClusterScoreStatus = newStatus
			return nil
		}

		managedClusterScore.Status = *newStatus
		updatedManagedClusterScore, err := client.ClusterV1alpha1().ManagedClusterScores(spokeClusterName).UpdateStatus(ctx, managedClusterScore, metav1.UpdateOptions{})
		if err != nil {
			return err
		}
		updatedManagedClusterScoreStatus = &updatedManagedClusterScore.Status
		updated = err == nil
		return err
	})

	return updatedManagedClusterScoreStatus, updated, err
}
