module github.com/yunify/qingcloud-csi

go 1.16

require (
	github.com/container-storage-interface/spec v1.5.0
	github.com/golang/protobuf v1.4.2
	github.com/kubernetes-csi/csi-lib-utils v0.6.1
	github.com/prometheus/client_golang v1.1.0 // indirect
	github.com/yunify/qingcloud-sdk-go v0.0.0-20230406022709-e32d107bcab7
	golang.org/x/net v0.0.0-20200707034311-ab3426394381
	google.golang.org/grpc v1.27.0
	k8s.io/apiextensions-apiserver v0.0.0-20190823014223-07b4561f8b0e // indirect
	k8s.io/apimachinery v0.20.0-alpha.2
	k8s.io/apiserver v0.0.0-20190823053033-1316076af51c // indirect
	k8s.io/client-go v0.20.0-alpha.2
	k8s.io/cloud-provider v0.0.0-20191212015549-86a326830157 // indirect
	k8s.io/klog v0.4.0
	k8s.io/kubernetes v1.14.1
)
