package managedclusterscore

import (
	"context"
	"fmt"

	"k8s.io/apimachinery/pkg/api/meta"
	"k8s.io/klog/v2"
	clientset "open-cluster-management.io/api/client/cluster/clientset/versioned"
	clusterv1informer "open-cluster-management.io/api/client/cluster/informers/externalversions/cluster/v1"
	clusterv1alpha1informer "open-cluster-management.io/api/client/cluster/informers/externalversions/cluster/v1alpha1"
	clusterv1listers "open-cluster-management.io/api/client/cluster/listers/cluster/v1"
	clusterv1alpha1listers "open-cluster-management.io/api/client/cluster/listers/cluster/v1alpha1"
	clusterv1 "open-cluster-management.io/api/cluster/v1"
	clusterv1alpha1 "open-cluster-management.io/api/cluster/v1alpha1"
	"open-cluster-management.io/score-agent/pkg/helpers"

	"github.com/openshift/library-go/pkg/controller/factory"
	"github.com/openshift/library-go/pkg/operator/events"
)

//const labelCustomizedOnly = "open-cluster-management.io/spoke-only"

// managedClusterClaimController exposes cluster claims created on managed cluster on hub after it joins the hub.
type managedClusterScoreController struct {
	clusterName      string
	hubClusterClient clientset.Interface
	hubClusterLister clusterv1listers.ManagedClusterLister
	scoreLister      clusterv1alpha1listers.ManagedClusterScoreLister
}

// NewManagedClusterClaimController creates a new managed cluster claim controller on the managed cluster.
func NewManagedClusterScoreController(
	clusterName string,
	hubClusterClient clientset.Interface,
	hubManagedClusterInformer clusterv1informer.ManagedClusterInformer,
	hubManagedClusterScoreInformer clusterv1alpha1informer.ManagedClusterScoreInformer,
	recorder events.Recorder) factory.Controller {
	c := &managedClusterScoreController{
		clusterName:      clusterName,
		hubClusterClient: hubClusterClient,
		hubClusterLister: hubManagedClusterInformer.Lister(),
		scoreLister:      hubManagedClusterScoreInformer.Lister(),
	}

	return factory.New().
		WithInformers(hubManagedClusterScoreInformer.Informer()).
		/*		WithInformersQueueKeyFunc(func(obj runtime.Object) string {
					accessor, _ := meta.Accessor(obj)
					return accessor.GetName()
				}, hubManagedClusterInformer.Informer()).
		*/
		WithSync(c.sync).
		ToController("NewManagedClusterScoreController", recorder)
}

// sync maintains the cluster claims in status of the managed cluster on hub once it joins the hub.
func (c *managedClusterScoreController) sync(ctx context.Context, syncCtx factory.SyncContext) error {
	klog.Info("HQ : Start syncing spoke score")
	managedCluster, err := c.hubClusterLister.Get(c.clusterName)
	if err != nil {
		return fmt.Errorf("unable to get managed cluster with name %q from hub: %w", c.clusterName, err)
	}

	// current managed cluster has not joined the hub yet, do nothing.
	if !meta.IsStatusConditionTrue(managedCluster.Status.Conditions, clusterv1.ManagedClusterConditionJoined) {
		syncCtx.Recorder().Eventf("ManagedClusterIsNotAccepted", "Managed cluster %q does not join the hub yet", c.clusterName)
		return nil
	}

	return c.updateScore(ctx, syncCtx)
}

// TODO
func (c *managedClusterScoreController) updateScore(ctx context.Context, syncCtx factory.SyncContext) error {
	klog.Info("HQ : Start updating spoke score")
	// update the status of the managed cluster score
	score := (int64)(100)
	updateStatusFuncs := []helpers.UpdateManagedClusterScoreStatusFunc{updateClusterScoresFn(clusterv1alpha1.ManagedClusterScoreStatus{
		Score: score,
	})}

	_, updated, err := helpers.UpdateManagedClusterScoreStatus(ctx, c.hubClusterClient, c.clusterName, updateStatusFuncs...)
	if err != nil {
		return fmt.Errorf("unable to update status of managed cluster score %q: %w", c.clusterName, err)
	}
	if updated {
		klog.V(4).Infof("The managed cluster score status %q has been updated", c.clusterName)
	}
	return nil
}

func updateClusterScoresFn(status clusterv1alpha1.ManagedClusterScoreStatus) helpers.UpdateManagedClusterScoreStatusFunc {
	return func(oldStatus *clusterv1alpha1.ManagedClusterScoreStatus) error {
		oldStatus.Score = status.Score
		return nil
	}
}
