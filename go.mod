module github.com/alexeldeib/trace-operator

go 1.13

require (
	github.com/go-logr/logr v0.1.0
	github.com/hashicorp/golang-lru v0.5.3 // indirect
	github.com/iovisor/kubectl-trace v0.1.0-rc.1
	github.com/onsi/ginkgo v1.10.1
	github.com/onsi/gomega v1.7.0
	k8s.io/api v0.17.0
	k8s.io/apimachinery v0.17.0
	k8s.io/client-go v11.0.0+incompatible
	sigs.k8s.io/controller-runtime v0.4.0
)

replace (
	github.com/iovisor/kubectl-trace => github.com/alexeldeib/kubectl-trace v0.1.0-rc.1.0.20200114233816-21e275138fdc
	k8s.io/api => k8s.io/api v0.17.0
	k8s.io/apimachinery => k8s.io/apimachinery v0.17.0
	k8s.io/client-go => k8s.io/client-go v0.17.0
)
