package spoke

import (
	"github.com/spf13/cobra"
	"k8s.io/klog/v2"

	"github.com/openshift/library-go/pkg/controller/controllercmd"

	"open-cluster-management.io/scalar-agent/pkg/spoke"
	"open-cluster-management.io/scalar-agent/pkg/version"
)

func NewAgent() *cobra.Command {
	klog.Info("Start new agent")
	agentOptions := spoke.NewSpokeAgentOptions()
	klog.Info("finish new an agent option")
	cmd := controllercmd.
		NewControllerCommandConfig("scalar-agent", version.Get(), agentOptions.RunSpokeAgent).
		NewCommand()
	cmd.Use = "agent"
	cmd.Short = "Start the Scalar Agent"

	klog.Info("finish new cmd")
	agentOptions.AddFlags(cmd.Flags())
	return cmd
}
