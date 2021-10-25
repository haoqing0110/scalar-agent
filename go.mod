module open-cluster-management.io/score-agent

go 1.16

replace open-cluster-management.io/api v0.0.0-20210927063308-2c6896161c48 => github.com/haoqing0110/api v0.0.0-20211021030927-af39ebdbc930
require (
	github.com/go-openapi/spec v0.19.5 // indirect
	github.com/openshift/generic-admission-server v1.14.1-0.20200903115324-4ddcdd976480 // indirect
	github.com/openshift/library-go v0.0.0-20211018074344-7fcf688c505e
	github.com/spf13/cobra v1.2.1
	github.com/spf13/pflag v1.0.5
	go.etcd.io/etcd v0.5.0-alpha.5.0.20200910180754-dd1b699fc489 // indirect
	k8s.io/apimachinery v0.22.2
	k8s.io/client-go v0.22.2
	k8s.io/component-base v0.22.2
	k8s.io/klog/v2 v2.20.0
	open-cluster-management.io/api v0.0.0-20210927063308-2c6896161c48 
)
