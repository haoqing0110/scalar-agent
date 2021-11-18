package spoke

import (
	"context"
	"path"
	"time"

	"github.com/openshift/library-go/pkg/controller/controllercmd"
	"github.com/openshift/library-go/pkg/controller/factory"
	"github.com/spf13/pflag"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/fields"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/tools/clientcmd"
	"k8s.io/klog/v2"
	clusterv1client "open-cluster-management.io/api/client/cluster/clientset/versioned"
	clusterv1informers "open-cluster-management.io/api/client/cluster/informers/externalversions"
	managedclusterScalar "open-cluster-management.io/scalar-agent/pkg/spoke/managedclusterscalar"
)

const (
	// spokeAgentNameLength is the length of the spoke agent name which is generated automatically
	spokeAgentNameLength = 5
	// defaultSpokeComponentNamespace is the default namespace in which the spoke agent is deployed
	defaultSpokeComponentNamespace = "open-cluster-management"
)

// AddOnLeaseControllerSyncInterval is exposed so that integration tests can crank up the constroller sync speed.
// TODO if we register the lease informer to the lease controller, we need to increase this time
var AddOnLeaseControllerSyncInterval = 30 * time.Second

// SpokeAgentOptions holds configuration for spoke cluster agent
type SpokeAgentOptions struct {
	ComponentNamespace       string
	ClusterName              string
	AgentName                string
	BootstrapKubeconfig      string
	HubKubeconfigSecret      string
	HubKubeconfigDir         string
	SpokeExternalServerURLs  []string
	ClusterHealthCheckPeriod time.Duration
	MaxCustomClusterClaims   int
}

// NewSpokeAgentOptions returns a SpokeAgentOptions
func NewSpokeAgentOptions() *SpokeAgentOptions {
	return &SpokeAgentOptions{
		HubKubeconfigSecret:      "hub-kubeconfig-secret",
		HubKubeconfigDir:         "/spoke/hub-kubeconfig",
		ClusterHealthCheckPeriod: 1 * time.Minute,
		MaxCustomClusterClaims:   20,
	}
}

// RunSpokeAgent starts the controllers on spoke agent to register to the hub.
//
// The spoke agent uses three kubeconfigs for different concerns:
// - The 'spoke' kubeconfig: used to communicate with the spoke cluster where
//   the agent is running.
// - The 'bootstrap' kubeconfig: used to communicate with the hub in order to
//   submit a CertificateSigningRequest, begin the join flow with the hub, and
//   to write the 'hub' kubeconfig.
// - The 'hub' kubeconfig: used to communicate with the hub using a signed
//   certificate from the hub.
//
// RunSpokeAgent handles the following scenarios:
//   #1. Bootstrap kubeconfig is valid and there is no valid hub kubeconfig in secret
//   #2. Both bootstrap kubeconfig and hub kubeconfig are valid
//   #3. Bootstrap kubeconfig is invalid (e.g. certificate expired) and hub kubeconfig is valid
//   #4. Neither bootstrap kubeconfig nor hub kubeconfig is valid
//
// A temporary ClientCertForHubController with bootstrap kubeconfig is created
// and started if the hub kubeconfig does not exist or is invalid and used to
// create a valid hub kubeconfig. Once the hub kubeconfig is valid, the
// temporary controller is stopped and the main controllers are started.
func (o *SpokeAgentOptions) RunSpokeAgent(ctx context.Context, controllerContext *controllercmd.ControllerContext) error {
	klog.Info("Start running spoke agent")
	spokeKubeClient, err := kubernetes.NewForConfig(controllerContext.KubeConfig)
	if err != nil {
		return err
	}

	// create hub clients and shared informer factories from hub kube config
	KubeconfigFile := "kubeconfig"
	hubClientConfig, err := clientcmd.BuildConfigFromFlags("", path.Join(o.HubKubeconfigDir, KubeconfigFile))
	if err != nil {
		return err
	}

	hubClusterClient, err := clusterv1client.NewForConfig(hubClientConfig)
	if err != nil {
		return err
	}

	// create a cluster informer factory with name field selector because we just need to handle the current spoke cluster
	hubClusterInformerFactory := clusterv1informers.NewSharedInformerFactoryWithOptions(
		hubClusterClient,
		10*time.Minute,
		clusterv1informers.WithTweakListOptions(func(listOptions *metav1.ListOptions) {
			listOptions.FieldSelector = fields.OneTermEqualSelector("metadata.name", o.ClusterName).String()
		}),
	)
	hubClusterNamespaceInformerFactory := clusterv1informers.NewSharedInformerFactoryWithOptions(
		hubClusterClient,
		10*time.Minute,
		clusterv1informers.WithNamespace(o.ClusterName),
	)

	controllerContext.EventRecorder.Event("HubClientConfigReady", "Client config for hub is ready.")

	var managedClusterScalarController factory.Controller
	// create managedClusterClaimController to sync cluster claims
	managedClusterScalarController = managedclusterScalar.NewManagedClusterScalarController(
		o.ClusterName,
		spokeKubeClient,
		hubClusterClient,
		hubClusterInformerFactory.Cluster().V1().ManagedClusters(),
		hubClusterNamespaceInformerFactory.Cluster().V1alpha1().ManagedClusterScalars(),
		controllerContext.EventRecorder,
	)

	go hubClusterInformerFactory.Start(ctx.Done())
	go hubClusterNamespaceInformerFactory.Start(ctx.Done())

	go managedClusterScalarController.Run(ctx, 1)

	<-ctx.Done()
	return nil
}

// AddFlags registers flags for Agent
func (o *SpokeAgentOptions) AddFlags(fs *pflag.FlagSet) {
	//	features.DefaultMutableFeatureGate.AddFlag(fs)
	fs.StringVar(&o.ClusterName, "cluster-name", o.ClusterName,
		"If non-empty, will use as cluster name instead of generated random name.")
	fs.StringVar(&o.BootstrapKubeconfig, "bootstrap-kubeconfig", o.BootstrapKubeconfig,
		"The path of the kubeconfig file for agent bootstrap.")
	fs.StringVar(&o.HubKubeconfigSecret, "hub-kubeconfig-secret", o.HubKubeconfigSecret,
		"The name of secret in component namespace storing kubeconfig for hub.")
	fs.StringVar(&o.HubKubeconfigDir, "hub-kubeconfig-dir", o.HubKubeconfigDir,
		"The mount path of hub-kubeconfig-secret in the container.")
	fs.StringArrayVar(&o.SpokeExternalServerURLs, "spoke-external-server-urls", o.SpokeExternalServerURLs,
		"A list of reachable spoke cluster api server URLs for hub cluster.")
	fs.DurationVar(&o.ClusterHealthCheckPeriod, "cluster-healthcheck-period", o.ClusterHealthCheckPeriod,
		"The period to check managed cluster kube-apiserver health")
}
