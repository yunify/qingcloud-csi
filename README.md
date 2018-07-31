# QingCloud-CSI

[![Build Status](https://travis-ci.org/yunify/qingcloud-csi.svg?branch=master)](https://travis-ci.org/yunify/qingcloud-csi)
[![Go Report Card](https://goreportcard.com/badge/github.com/yunify/qingcloud-csi)](https://goreportcard.com/report/github.com/yunify/qingcloud-csi)

## Description
QingCloud CSI plugin implements an interface between Container Storage Interface([CSI](https://github.com/container-storage-interface/)) enabled Container Orchestrator(CO) and the storage of QingCloud. Currently, QingCloud CSI plugin has been passed the [CSI test](https://github.com/kubernetes-csi/csi-test) in Kubernetes v1.10 environment.

## Block Plugin

### Compiling
QingCloud CSI plugin can be complied as a binary file or a container.  We can get a binary file in _output folder. When compiled as a container, the image is stored in a local Docker's image store. 

To compile a binary file:
```
$ make blockplugin
```

To build a Docker image:
```
$ make blockplugin-container
```

You can find image in your local image store
```
$ docker images | grep csi-qingcloud
dockerhub.qingcloud.com/csiplugin/csi-qingcloud		v0.2.0.1	640a9519e59b		55 minutes ago		40MB
```

### Configuration
#### Config File

Config [file](deploy/block/kubernetes/config.yaml) shown below would be referenced by a ConfigMap.
> IMPORTANT: In QingCloud AppCenter, please modify [script](deploy/block/kubernetes/create-cm.sh) and create a ConfigMap which references another config file(*/etc/qingcloud/client.yaml*) on the host machine.

```
qy_access_key_id: 'ACCESS_KEY_ID'
qy_secret_access_key: 'ACCESS_KEY_SECRET'
zone: 'ZONE'
host: 'api.qingcloud.com'
port: 443
protocol: 'https'
uri: '/iaas'
connection_retries: 3
connection_timeout: 30
```

- `qy_access_key_id`, `qy_secret_access_key`: Access key pair can be created in QingCloud console. The access key pair must have the power to manipulate QingCloud IaaS platform resource.

- `zone`: Zone should be the same as Kubernetes cluster. CSI plugin will operate block storage volumes in this zone.

- `host`, `prot`. `protocol`, `uri`: QingCloud IaaS platform service url.

### StorageClass

SotrageClass definition [file](deploy/block/example/sc.yaml) shown below is used to create StorageClass object.
```
apiVersion: storage.k8s.io/v1
kind: StorageClass
metadata:
  name: csi-qingcloud
provisioner: csi-qingcloud
parameters:
  type: "0"
  maxSize: "500"
  minSize: "10"
  fsType: "ext4"
reclaimPolicy: Delete 
```

- `type`: The type of volume in QingCloud IaaS platform. See [QingCloud docs](https://docs.qingcloud.com/product/api/action/volume/create_volumes.html) for details.

- `maxSize`, `minSize`: The maximum and minimum volume size with specific volume type.

- `fsType`: `ext3`, `ext4`, `xfs`. Default `ext4`.



### Installation
This guide will deploy CSI plugin in *kube-system* namespace. You can deploy the plugin in other namespace. DO NOT disable [Mount Propagation](https://kubernetes.io/docs/concepts/storage/volumes/#mount-propagation) feature gate in Kubernetes control plane.

- Create ConfigMap
```
$ chmod +x deploy/block/kubernetes/create-cm.sh
$ ./create-cm.sh
```

- Create Docker image registry secret
```
kubectl create secret docker-registry csi-registry --docker-server=dockerhub.qingcloud.com --docker-username=<YOUR_USERNAME> --docker-password=<YOUR_PASSWORD> --docker-email=<YOUR_EMAIL> --namespace=kube-system
```

- Create access control objects
```
$ kubectl create -f deploy/block/kubernetes/csi-controller-rbac.yaml
$ kubectl create -f deploy/block/kubernetes/csi-node-rbac.yaml
```

- Deploy CSI plugin
> IMPORTANT: In QingCloud AppCenter, please replace *"/var/lib/kubelet"* with *"/data/var/lib/kubelet"* in [DaemonSet](deploy/block/kubernetes/csi-node-ds.yaml) YAML file,.

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
$ kubectl create -f deploy/block/example/sc.yaml
```

- Create a PVC
```
$ kubectl create -f deploy/block/example/pvc.yaml
```

- Create a Deployment mounting the PVC
```
$ kubectl create -f deploy/block/example/deploy.yaml
```

- Check Pod status
```
$ kubectl get po | grep deploy
nginx-84474cf674-zfhbs   1/1       Running   0          1m
```

- Access container's directory which mounting volume
```
$ kubectl exec -ti deploy-nginx-qingcloud-84474cf674-zfhbs /bin/bash
# cd /mnt
# ls
lost+found
```

## Support
If you have any qustions or suggestions, please submit an issue at [qingcloud-csi](https://github.com/yunify/qingcloud-csi/issues)
