// Code generated by client-gen. DO NOT EDIT.

package versioned

import (
	"fmt"

	discovery "k8s.io/client-go/discovery"
	rest "k8s.io/client-go/rest"
	flowcontrol "k8s.io/client-go/util/flowcontrol"
	clusterv1 "open-cluster-management.io/api/client/cluster/clientset/versioned/typed/cluster/v1"
	clusterv1alpha1 "open-cluster-management.io/api/client/cluster/clientset/versioned/typed/cluster/v1alpha1"
	clusterv1beta1 "open-cluster-management.io/api/client/cluster/clientset/versioned/typed/cluster/v1beta1"
)

type Interface interface {
	Discovery() discovery.DiscoveryInterface
	ClusterV1() clusterv1.ClusterV1Interface
	ClusterV1alpha1() clusterv1alpha1.ClusterV1alpha1Interface
	ClusterV1beta1() clusterv1beta1.ClusterV1beta1Interface
}

// Clientset contains the clients for groups. Each group has exactly one
// version included in a Clientset.
type Clientset struct {
	*discovery.DiscoveryClient
	clusterV1       *clusterv1.ClusterV1Client
	clusterV1alpha1 *clusterv1alpha1.ClusterV1alpha1Client
	clusterV1beta1  *clusterv1beta1.ClusterV1beta1Client
}

// ClusterV1 retrieves the ClusterV1Client
func (c *Clientset) ClusterV1() clusterv1.ClusterV1Interface {
	return c.clusterV1
}

// ClusterV1alpha1 retrieves the ClusterV1alpha1Client
func (c *Clientset) ClusterV1alpha1() clusterv1alpha1.ClusterV1alpha1Interface {
	return c.clusterV1alpha1
}

// ClusterV1beta1 retrieves the ClusterV1beta1Client
func (c *Clientset) ClusterV1beta1() clusterv1beta1.ClusterV1beta1Interface {
	return c.clusterV1beta1
}

// Discovery retrieves the DiscoveryClient
func (c *Clientset) Discovery() discovery.DiscoveryInterface {
	if c == nil {
		return nil
	}
	return c.DiscoveryClient
}

// NewForConfig creates a new Clientset for the given config.
// If config's RateLimiter is not set and QPS and Burst are acceptable,
// NewForConfig will generate a rate-limiter in configShallowCopy.
func NewForConfig(c *rest.Config) (*Clientset, error) {
	configShallowCopy := *c
	if configShallowCopy.RateLimiter == nil && configShallowCopy.QPS > 0 {
		if configShallowCopy.Burst <= 0 {
			return nil, fmt.Errorf("burst is required to be greater than 0 when RateLimiter is not set and QPS is set to greater than 0")
		}
		configShallowCopy.RateLimiter = flowcontrol.NewTokenBucketRateLimiter(configShallowCopy.QPS, configShallowCopy.Burst)
	}
	var cs Clientset
	var err error
	cs.clusterV1, err = clusterv1.NewForConfig(&configShallowCopy)
	if err != nil {
		return nil, err
	}
	cs.clusterV1alpha1, err = clusterv1alpha1.NewForConfig(&configShallowCopy)
	if err != nil {
		return nil, err
	}
	cs.clusterV1beta1, err = clusterv1beta1.NewForConfig(&configShallowCopy)
	if err != nil {
		return nil, err
	}

	cs.DiscoveryClient, err = discovery.NewDiscoveryClientForConfig(&configShallowCopy)
	if err != nil {
		return nil, err
	}
	return &cs, nil
}

// NewForConfigOrDie creates a new Clientset for the given config and
// panics if there is an error in the config.
func NewForConfigOrDie(c *rest.Config) *Clientset {
	var cs Clientset
	cs.clusterV1 = clusterv1.NewForConfigOrDie(c)
	cs.clusterV1alpha1 = clusterv1alpha1.NewForConfigOrDie(c)
	cs.clusterV1beta1 = clusterv1beta1.NewForConfigOrDie(c)

	cs.DiscoveryClient = discovery.NewDiscoveryClientForConfigOrDie(c)
	return &cs
}

// New creates a new Clientset for the given RESTClient.
func New(c rest.Interface) *Clientset {
	var cs Clientset
	cs.clusterV1 = clusterv1.New(c)
	cs.clusterV1alpha1 = clusterv1alpha1.New(c)
	cs.clusterV1beta1 = clusterv1beta1.New(c)

	cs.DiscoveryClient = discovery.NewDiscoveryClient(c)
	return &cs
}
