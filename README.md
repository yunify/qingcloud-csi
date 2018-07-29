# QingCloud-CSI

[![Build Status](https://travis-ci.org/yunify/qingcloud-csi.svg?branch=master)](https://travis-ci.org/yunify/qingcloud-csi)
[![Go Report Card](https://goreportcard.com/badge/github.com/yunify/qingcloud-csi)](https://goreportcard.com/report/github.com/yunify/qingcloud-csi)

Kubernetes volume plugin based on CSI specification which support block storage of qingcloud

## Description
QingCloud CSI plugin implements an interface between Container Storage Interface([CSI](https://github.com/container-storage-interface/)) enabled Container Orchestrator(CO) and the storage of QingCloud. Currently, QingCloud CSI plugin has been passed the [CSI test](https://github.com/kubernetes-csi/csi-test) in Kubernetes v1.10 environment.

## Block Plugin

### Compiling
QingCloud CSI plugin can be complied as a binary file or a container.  We can get a binary file in _output folder. When compiled as a container, the image is stored in a local Docker's image store.

To compile a binary file:
```
$ make blockplugin
```

To compile a Docker image:
```
$ make blockplugin-container
```

You can find image in your local image store
```
$ docker images | grep csi-qingcloud
dockerhub.qingcloud.com/wiley/csi-qingcloud	v0.2.0.1	640a9519e59b		55 minutes ago		40MB
```

### Configuration
- [ConfigMap](deploy/block/kubernetes/config.yaml): Set parameters about accessing storage server.
- [StorageClass](deploy/block/kubernetes/sc.yaml): Set creating volume parameters.
- [Mount Propagation](https://kubernetes.io/docs/concepts/storage/volumes/#mount-propagation): DO NOT disable this feature gate.

> Notes: When deploy this plugin in QingCloud AppCenter, you must follow guides below.
> 1. Modify creating ConfigMap [script](deploy/block/kubernetes/create-cm.sh) and create a ConfigMap which references the YAML file(*/etc/qingcloud/client.yaml*) on the host machine.
> 1. In the [DaemonSet](deploy/block/kubernetes/csi-node-ds.yaml) YAML file, please replace *"/var/lib/kubelet"* with *"/data/var/lib/kubelet"* .

### Deployment
This guide will create a namespace named csi-qingcloud and deploy CSI plugin in this namespace. You can modify YAML files mentioned below and deploy the plugin in other namespace.

- Create csi-qingcloud namespace
```
$ kubectl create -f deploy/block/kubernetes/csi-ns.yaml
```

- Create ConfigMap from file
```
$ chmod +x deploy/block/kubernetes/create-cm.sh
$ ./create-cm.sh
```


- Create Docker image registry secret
```
kubectl create secret docker-registry csi-registry --docker-server=dockerhub.qingcloud.com --docker-username=<YOUR_USERNAME> --docker-password=<YOUR_PASSWORD> --docker-email=<YOUR_EMAIL> --namespace=csi-qingcloud
```

- Create access control objects
```
$ kubectl create -f deploy/block/kubernetes/csi-controller-rbac.yaml
$ kubectl create -f deploy/block/kubernetes/csi-node-rbac.yaml
```

- Deploy CSI plugin
```
$ kubectl create -f deploy/block/kubernetes/csi-controller-sts.yaml
$ kubectl create -f deploy/block/kubernetes/csi-node-ds.yaml
```

- Check CSI plugin
```
$ kubectl get pods -n csi-qingcloud | grep csi
csi-qingcloud-controller-0      3/3       Running       0          5m
csi-qingcloud-node-kks3q        2/2       Running       0          2m
csi-qingcloud-node-pgsbn        2/2       Running       0          2m
```

### Verification
- Create a StorageClass by Kubernetes cluster administrator
```
$ kubectl create -f deploy/block/kubernetes/sc.yaml
```

- Create a PVC
```
$ kubectl create -f deploy/block/kubernetes/pvc.yaml
```

- Create a Deployment mounting the PVC
```
$ kubectl create -f deploy/block/kubernetes/deploy.yaml
```

- Check Pod status
```
$ kubectl get po | grep deploy
nginx-84474cf674-zfhbs   1/1       Running   0          1m
```

```
$ kubectl exec -ti deploy-nginx-qingcloud-84474cf674-zfhbs /bin/bash
// We can access the directoriy mounting persistent volume in container
# cd /mnt
# ls
lost+found
```

## Support
If you have any qustions or suggestions, please submit an issue at [qingcloud-csi](https://github.com/yunify/qingcloud-csi/issues)
