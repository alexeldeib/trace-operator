module github.com/alexeldeib/trace-operator

go 1.13

require (
	github.com/go-logr/logr v0.1.0
	github.com/hashicorp/golang-lru v0.5.3 // indirect
	github.com/iovisor/kubectl-trace v0.1.0-rc.1
	github.com/onsi/ginkgo v1.10.1
	github.com/onsi/gomega v1.7.0
	golang.org/x/crypto v0.0.0-20190820162420-60c769a6c586 // indirect
	golang.org/x/time v0.0.0-20190308202827-9d24e82272b4 // indirect
	k8s.io/api v0.0.0-20190918195907-bd6ac527cfd2
	k8s.io/apimachinery v0.0.0-20190817020851-f2f3a405f61d
	k8s.io/client-go v11.0.0+incompatible
	sigs.k8s.io/controller-runtime v0.3.0
)

replace (
	github.com/iovisor/kubectl-trace => github.com/alexeldeib/kubectl-trace v0.1.0-rc.1.0.20200114233816-21e275138fdc
	k8s.io/api => k8s.io/api v0.0.0-20190409021203-6e4e0e4f393b
	k8s.io/apiextensions-apiserver => k8s.io/apiextensions-apiserver v0.0.0-20190409022649-727a075fdec8
	k8s.io/apimachinery => k8s.io/apimachinery v0.0.0-20190404173353-6a84e37a896d
	k8s.io/client-go => k8s.io/client-go v11.0.1-0.20190409021438-1a26190bd76a+incompatible
)
