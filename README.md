# QingCloud-CSI

[![Build Status](https://travis-ci.org/yunify/qingcloud-csi.svg?branch=master)](https://travis-ci.org/yunify/qingcloud-csi)
[![Go Report Card](https://goreportcard.com/badge/github.com/yunify/qingcloud-csi)](https://goreportcard.com/report/github.com/yunify/qingcloud-csi)

> English | [中文](README_zh.md)

## Description
QingCloud CSI plugin implements an interface between Container Storage Interface ([CSI](https://github.com/container-storage-interface/)) enabled Container Orchestrator (CO) and the storage of QingCloud. Currently, QingCloud CSI plugin has been passed the [CSI test](https://github.com/kubernetes-csi/csi-test) in Kubernetes v1.10 environment.

## Block Storage Plugin

Block storage plugin's design and installation use Kubernetes community recommended CSI plugin [architecture](https://github.com/kubernetes/community/blob/master/contributors/design-proposals/storage/container-storage-interface.md#recommended-mechanism-for-deploying-csi-drivers-on-kubernetes). Plugin architecture contains Controller part and Node part. In the part of Controller, one Pod is created by StatefulSet in Kubernetes cluster. In the part of Node, one Pod is created by DaemonSet on every node. 

After plugin installation completes, user can create volumes based on several types of disk, such as super high performance disk, high performance disk and high capacity disk, with ReadWriteOnce access mode and mount volumes on workloads.

### Installation
This guide will install CSI plugin in *kube-system* namespace. You can also deploy the plugin in other namespace. To use this CSI plugin, please ensure `--allow-privileged` flag set to `true` and enable [Mount Propagation](https://kubernetes.io/docs/concepts/storage/volumes/#mount-propagation) (Default enalbed) feature gate in Kubernetes control plane.

- Download and decompress installation package 
```
$ wget $(curl --silent "https://api.github.com/repos/yunify/qingcloud-csi/releases/latest" | \
  grep browser_download_url | grep install|cut -d '"' -f 4)
$ tar -xvf csi-qingcloud-install.tar.gz
$ cd csi-qingcloud-install
```

- Create ConfigMap
  * In Kubernetes cluster based on QingCloud IaaS platform
    1. Modify config file (./config.yaml) in installation package
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

    - `zone`: `Zone` should be the same as Kubernetes cluster. CSI plugin will operate block storage volumes in this zone. For example, `zone` can be set to `sh1a` and `ap2a`.

    - `host`, `port`. `protocol`, `uri`: QingCloud IaaS platform service url.

    2. Create ConfigMap
    ```
    $ kubectl create configmap csi-qingcloud --from-file=config.yaml=./config.yaml --namespace=kube-system
    ```
  * In Kubernetes cluster based on QingCloud AppCenter

    1. Create ConfigMap
    ```
    $ kubectl create configmap csi-qingcloud --from-file=config.yaml=/etc/qingcloud/client.yaml --namespace=kube-system
    ```

- Create Docker image registry secret
```
$ kubectl apply -f ./csi-secret.yaml
```

- Create access control objects
```
$ kubectl apply -f ./csi-controller-rbac.yaml
$ kubectl apply -f ./csi-node-rbac.yaml
```

- Deploy CSI plugin
> IMPORTANT: If kubelet, a component of Kubernetes, set the `--root-dir` option (default: *"/var/lib/kubelet"*), please replace *"/var/lib/kubelet"* with the value of `--root-dir` at the CSI [DaemonSet](deploy/disk/kubernetes/csi-node-ds.yaml) YAML file's `spec.template.spec.containers[name=csi-qingcloud].volumeMounts[name=mount-dir].mountPath` and `spec.template.spec.volumes[name=mount-dir].hostPath.path` fields. For instance, in Kubernetes cluster based on QingCloud AppCenter, you should replace *"/var/lib/kubelet"* with *"/data/var/lib/kubelet"* in the CSI [DaemonSet](deploy/disk/kubernetes/csi-node-ds.yaml) YAML file.

```
$ kubectl apply -f ./csi-controller-sts.yaml
$ kubectl apply -f ./csi-node-ds.yaml
```

- Check CSI plugin
```
$ kubectl get pods -n kube-system --selector=app=csi-qingcloud
NAME                            READY     STATUS        RESTARTS   AGE
csi-qingcloud-controller-0      3/3       Running       0          5m
csi-qingcloud-node-kks3q        2/2       Running       0          2m
csi-qingcloud-node-pgsbn        2/2       Running       0          2m
```

### Verification
- Create a StorageClass by Kubernetes cluster administrator
> NOTE: This guide will create a StorageClass which sets `type` to `0`. User could set StorageClass parameters according to following instruction.
```
$ kubectl apply -f https://raw.githubusercontent.com/yunify/qingcloud-csi/master/deploy/block/example/sc.yaml
```

- Create a PVC
```
$ kubectl apply -f https://raw.githubusercontent.com/yunify/qingcloud-csi/master/deploy/block/example/pvc.yaml
```

- Create a Deployment mounting the PVC
```
$ kubectl apply -f https://raw.githubusercontent.com/yunify/qingcloud-csi/master/deploy/block/example/deploy.yaml
```

- Check Pod status
```
$ kubectl get po | grep nginx
nginx-84474cf674-zfhbs   1/1       Running   0          1m
```

- Access container's directory which mounting volume
```
$ kubectl exec -ti deploy-nginx-qingcloud-84474cf674-zfhbs /bin/bash
# cd /mnt
# ls
lost+found
```

### StorageClass Parameters

StorageClass definition [file](deploy/disk/example/sc.yaml) shown below is used to create StorageClass object.
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
  stepSize: "10"
  fsType: "ext4"
  replica: "2"
reclaimPolicy: Delete 
```

- `type`: The type of volume in QingCloud IaaS platform. In QingCloud public cloud platform, `0` represents high performance volume. `3` respresents super high performance volume. `1` or `2` represents high capacity volume depending on cluster‘s zone. `5` represents enterprise distributed SAN (NeonSAN) volume. `100` represents standard volume. `200` represents SSD enterprise volume. See [QingCloud docs](https://docs.qingcloud.com/product/api/action/volume/create_volumes.html) for details.

- `maxSize`, `minSize`: Limit the range of volume size in GiB.

- `stepSize`: Set the increment of volumes size in GiB.

- `fsType`: `ext3`, `ext4`, `xfs`. Default `ext4`.

- `replica`: `1` means single replica, `2` means multiple replicas. Default `2`.

## Support
If you have any qustions or suggestions, please submit an issue at [qingcloud-csi](https://github.com/yunify/qingcloud-csi/issues)
