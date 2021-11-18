package managedclusterScalar

import (
	"context"
	"fmt"
	"math/rand"
	"time"

	"k8s.io/apimachinery/pkg/api/meta"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
	"k8s.io/klog/v2"
	clientset "open-cluster-management.io/api/client/cluster/clientset/versioned"
	clusterv1informer "open-cluster-management.io/api/client/cluster/informers/externalversions/cluster/v1"
	clusterv1alpha1informer "open-cluster-management.io/api/client/cluster/informers/externalversions/cluster/v1alpha1"
	clusterv1listers "open-cluster-management.io/api/client/cluster/listers/cluster/v1"
	clusterv1alpha1listers "open-cluster-management.io/api/client/cluster/listers/cluster/v1alpha1"
	clusterv1 "open-cluster-management.io/api/cluster/v1"
	clusterv1alpha1 "open-cluster-management.io/api/cluster/v1alpha1"
	"open-cluster-management.io/scalar-agent/pkg/helpers"

	"github.com/openshift/library-go/pkg/controller/factory"
	"github.com/openshift/library-go/pkg/operator/events"
)

const (
	ResyncDuration = 5 * time.Minute
	ValidDuration  = 10 * time.Minute
)

// managedClusterClaimController exposes cluster claims created on managed cluster on hub after it joins the hub.
type managedClusterScalarController struct {
	clusterName      string
	hubKubeClient    *kubernetes.Clientset
	hubClusterClient clientset.Interface
	hubClusterLister clusterv1listers.ManagedClusterLister
	scalarLister     clusterv1alpha1listers.ManagedClusterScalarLister
}

// NewManagedClusterClaimController creates a new managed cluster claim controller on the managed cluster.
func NewManagedClusterScalarController(
	clusterName string,
	hubKubeClient *kubernetes.Clientset,
	hubClusterClient clientset.Interface,
	hubManagedClusterInformer clusterv1informer.ManagedClusterInformer,
	hubManagedClusterScalarInformer clusterv1alpha1informer.ManagedClusterScalarInformer,
	recorder events.Recorder) factory.Controller {
	c := &managedClusterScalarController{
		clusterName:      clusterName,
		hubKubeClient:    hubKubeClient,
		hubClusterClient: hubClusterClient,
		hubClusterLister: hubManagedClusterInformer.Lister(),
		scalarLister:     hubManagedClusterScalarInformer.Lister(),
	}

	return factory.New().
		ResyncEvery(ResyncDuration).
		WithInformers(hubManagedClusterScalarInformer.Informer()).
		WithSync(c.sync).
		ToController("NewManagedClusterScalarController", recorder)
}

// sync maintains the cluster claims in status of the managed cluster on hub once it joins the hub.
func (c *managedClusterScalarController) sync(ctx context.Context, syncCtx factory.SyncContext) error {
	managedCluster, err := c.hubClusterLister.Get(c.clusterName)
	if err != nil {
		return fmt.Errorf("unable to get managed cluster with name %q from hub: %w", c.clusterName, err)
	}

	// current managed cluster has not joined the hub yet, do nothing.
	if !meta.IsStatusConditionTrue(managedCluster.Status.Conditions, clusterv1.ManagedClusterConditionJoined) {
		return fmt.Errorf("managed cluster %q does not join the hub yet", c.clusterName)
	}

	// each collector update Scalar to corresponding CR.
	names := []string{"customizeresourceallocatablememory", "customizeresourceallocatablecpu"}
	for _, name := range names {
		fakescalar := int64(rand.Intn(100))
		//c.resourceAllocatableCPU()
		if _, err := c.scalarLister.ManagedClusterScalars(c.clusterName).Get(name); err == nil {
			if err := c.UpdateScalar(ctx, syncCtx, name, fakescalar); err != nil {
				return err
			}
		}
	}
	return nil
}

// TODO
/* func (c *managedClusterScalarController) resourceAllocatableCPU() {
	nd := describe.NodeDescriber{c.hubKubeClient}
	if result, err := nd.Describe("cluster1-control-plane", "kube-system", describe.DescriberSettings{ShowEvents: true}); err != nil {
		klog.Infof("describe node : %s", result)
	} else {
		klog.Warningf("describe node failed : %s", err)
	}
}*/

// TODO
func (c *managedClusterScalarController) UpdateScalar(ctx context.Context, syncCtx factory.SyncContext, crName string, Scalar int64) error {
	// update the status of the managed cluster Scalar
	updateStatusFuncs := []helpers.UpdateManagedClusterScalarStatusFunc{
		updateClusterScalarsFn(clusterv1alpha1.ManagedClusterScalarStatus{
			Scalar:     Scalar,
			ValidUntil: &metav1.Time{Time: time.Now().Add(ValidDuration)},
		}),
		updateManagedClusterScalarConditionFn(metav1.Condition{
			Type:    "ManagedClusterScalarUpdated",
			Status:  "True",
			Reason:  "ManagedClusterScalarUpdated",
			Message: "ManagedClusterScalar updated successfully",
		}),
	}

	_, updated, err := helpers.UpdateManagedClusterScalarStatus(ctx, c.hubClusterClient, c.clusterName, crName, updateStatusFuncs...)
	if err != nil {
		return fmt.Errorf("unable to update status of managed cluster Scalar %q: %w", c.clusterName, err)
	}
	if updated {
		klog.V(4).Infof("The managed cluster Scalar status %q has been updated", c.clusterName)
	}

	return nil
}

func updateClusterScalarsFn(status clusterv1alpha1.ManagedClusterScalarStatus) helpers.UpdateManagedClusterScalarStatusFunc {
	return func(oldStatus *clusterv1alpha1.ManagedClusterScalarStatus) error {
		oldStatus.Scalar = status.Scalar
		oldStatus.ValidUntil = status.ValidUntil
		return nil
	}
}

func updateManagedClusterScalarConditionFn(cond metav1.Condition) helpers.UpdateManagedClusterScalarStatusFunc {
	return func(oldStatus *clusterv1alpha1.ManagedClusterScalarStatus) error {
		setStatusCondition(&oldStatus.Conditions, cond)
		return nil
	}
}

func setStatusCondition(conditions *[]metav1.Condition, newCondition metav1.Condition) {
	if conditions == nil {
		return
	}
	existingCondition := findStatusCondition(*conditions, newCondition.Type)
	if existingCondition == nil {
		if newCondition.LastTransitionTime.IsZero() {
			newCondition.LastTransitionTime = metav1.NewTime(time.Now())
		}
		*conditions = append(*conditions, newCondition)
		return
	}

	if existingCondition.Status != newCondition.Status {
		existingCondition.Status = newCondition.Status
		if !newCondition.LastTransitionTime.IsZero() {
			existingCondition.LastTransitionTime = newCondition.LastTransitionTime
		} else {
			existingCondition.LastTransitionTime = metav1.NewTime(time.Now())
		}
	}

	existingCondition.Reason = newCondition.Reason
	existingCondition.Message = newCondition.Message
	existingCondition.ObservedGeneration = newCondition.ObservedGeneration
}

// FindStatusCondition finds the conditionType in conditions.
func findStatusCondition(conditions []metav1.Condition, conditionType string) *metav1.Condition {
	for i := range conditions {
		if conditions[i].Type == conditionType {
			return &conditions[i]
		}
	}

	return nil
}
