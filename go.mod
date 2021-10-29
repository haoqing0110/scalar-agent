module open-cluster-management.io/score-agent

go 1.16

replace open-cluster-management.io/api v0.0.0-20210916013819-2e58cdb938f9 => github.com/haoqing0110/api v0.0.0-20211029061752-bcd4aa6a346a

require (
	github.com/hashicorp/golang-lru v0.5.4 // indirect
	github.com/imdario/mergo v0.3.10 // indirect
	github.com/onsi/ginkgo v1.14.0 // indirect
	github.com/openshift/build-machinery-go v0.0.0-20210209125900-0da259a2c359
	github.com/openshift/library-go v0.0.0-20210407140145-f831e911c638
	github.com/spf13/cobra v1.2.1
	github.com/spf13/pflag v1.0.5
	go.uber.org/tools v0.0.0-20190618225709-2cfd321de3ee // indirect
	k8s.io/apimachinery v0.22.2
	k8s.io/apiserver v0.22.2 // indirect
	k8s.io/client-go v0.22.2
	k8s.io/component-base v0.22.2
	k8s.io/klog/v2 v2.9.0
	open-cluster-management.io/api v0.0.0-20210916013819-2e58cdb938f9
)
