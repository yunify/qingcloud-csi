# QingCloud-CSI

[![Build Status](https://travis-ci.org/yunify/qingcloud-csi.svg?branch=master)](https://travis-ci.org/yunify/qingcloud-csi)
[![Go Report Card](https://goreportcard.com/badge/github.com/yunify/qingcloud-csi)](https://goreportcard.com/report/github.com/yunify/qingcloud-csi)

> English | [中文](README_zh.md)

<!-- START doctoc generated TOC please keep comment here to allow auto update -->
<!-- DON'T EDIT THIS SECTION, INSTEAD RE-RUN doctoc TO UPDATE -->
- [Description](#description)
- [Disk Plugin](#disk-plugin)
  - [Kubernetes Compatibility Matrix](#kubernetes-compatibility-matrix)
  - [Feature Matrix](#feature-matrix)
  - [Installation](#installation)
  - [Uninstall](#uninstall)
  - [StorageClass Parameters](#storageclass-parameters)
- [Support](#support)

<!-- END doctoc generated TOC please keep comment here to allow auto update -->

---

## Description
QingCloud CSI plugin implements an interface between Container Storage Interface ([CSI](https://github.com/container-storage-interface/)) enabled Container Orchestrator (CO) and the storage of QingCloud. Currently, QingCloud CSI disk plugin has been developed and manages disk volume in QingCloud platform.

## Disk Plugin

Disk plugin's design and installation use Kubernetes community recommended CSI plugin [architecture](https://github.com/kubernetes/community/blob/master/contributors/design-proposals/storage/container-storage-interface.md#recommended-mechanism-for-deploying-csi-drivers-on-kubernetes). Plugin architecture contains Controller part and Node part. In the part of Controller, one Pod is created by Deployment in Kubernetes cluster. In the part of Node, one Pod is created by DaemonSet on every node. Now, it has been passed the [CSI test](https://github.com/kubernetes-csi/csi-test) in Kubernetes v1.15 environment.

After plugin installation completes, user can create volumes based on several types of disk, such as Standard disk, SSD Enterprise disk, High Performance disk, Super High Performance disk, NeonSAN disk, NeonSAN HDD disk and High Capacity disk, with ReadWriteOnce access mode and mount volumes on workloads.

### Kubernetes Compatibility Matrix

| |Kubernetes v1.10-v1.13|Kubernetes v1.14-1.15|
|:---:|:---:|:---:|
|QingCloud CSI v0.2.x|✓|-|
|QingCloud CSI v1.1.0|-|✓|

### Feature Matrix

| | Volume Management* | Volume Expansion | Volume Monitor | Volume Cloning| Snapshot Management**| Topology Awareness|
|:---:|:---:|:---:|:---:|:---:|:---:|:---:|
|QingCloud CSI v0.2.x |✓|-|-|-|-|-|
|QingCloud CSI v1.1.0 |✓|✓|✓|✓|✓|✓|

Notes:
- `*`: Volume Management including creating/deleting volume and mounting/unmount volume on Pod.
- `**`: Snapshot management including creating/deleting snapshot and restoring volume from snapshot.

### Installation
This guide will install CSI plugin in the *kube-system* namespace of Kubernetes v1.14+. You can also deploy the plugin in other namespace. 

- Set Kubernetes Parameters
  - Enable `--allow-privileged=true` on kube-apiserver, kube-controller-manager, kube-scheduler, kubelet
  - Enable (Default enabled) [Mount Propagation](https://kubernetes.io/docs/concepts/storage/volumes/#mount-propagation) feature gate。
  - Enable `--feature-gates=CSINodeInfo=true,CSIDriverRegistry=true,KubeletPluginsWatcher=true,VolumeSnapshotDataSource=true,ExpandCSIVolumes=true,VolumePVCDataSource=true（Only for Kubernetes v1.15）` option on kube-apiserver, kube-controller-manager, kube-scheduler, kubelet
  - Enable `--read-only-port=10255` on kubelet
- Download installation file 
```
$ wget https://raw.githubusercontent.com/yunify/qingcloud-csi/master/deploy/disk/kubernetes/releases/qingcloud-csi-disk-v1.1.0.yaml
```
- Add QingCloud platform parameter on ConfigMap
QingCloud CSI plugin manipulates cloud resource by QingCloud platform API. User must test the connection between QingCloud platform API and user's own instance by and check QingCloud platform configuration by [QingCloud CLI](https://docs.qingcloud.com/product/cli/).
  - Modify `csi-qingcloud` ConfigMap parameters in installation file
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

    - `zone`: `zone` should be the same as Kubernetes cluster. CSI plugin will manipulate resources in this region or zone. For example, `zone` can be set to `sh1` or `ap2a`.

    - `host`, `port`. `protocol`, `uri`: QingCloud IaaS platform service url.

- Deploy CSI plugin
> IMPORTANT: If kubelet, a component of Kubernetes, set the `--root-dir` option (default: *"/var/lib/kubelet"*), please replace *"/var/lib/kubelet"* with the value of `--root-dir` at the CSI [DaemonSet](deploy/disk/kubernetes/csi-node-ds.yaml) YAML file's `spec.template.spec.containers[name=csi-qingcloud].volumeMounts[name=mount-dir].mountPath` and `spec.template.spec.volumes[name=mount-dir].hostPath.path` fields. For instance, in Kubernetes cluster based on QingCloud AppCenter, you should replace *"/var/lib/kubelet"* with *"/data/var/lib/kubelet"* in the CSI [DaemonSet](deploy/disk/kubernetes/csi-node-ds.yaml) YAML file.

```
$ kubectl apply -f qingcloud-csi-disk-v1.1.0.yaml
```

- Check CSI plugin
```
$ kubectl get pods -n kube-system --selector=app=csi-qingcloud
  NAME                                       READY   STATUS    RESTARTS   AGE
  csi-qingcloud-controller-5bd48bb49-dw9rs   5/5     Running   0          3h16m
  csi-qingcloud-node-d2kdt                   2/2     Running   0          3h16m
  csi-qingcloud-node-hvtq7                   2/2     Running   0          3h16m
  csi-qingcloud-node-njghb                   2/2     Running   0          3h16m
  csi-qingcloud-node-wssdt                   2/2     Running   0          3h16m
```

### Uninstall
```
$ kubectl delete -f qingcloud-csi-disk-v1.1.0.yaml
```

### StorageClass Parameters

StorageClass definition [file](deploy/disk/example/sc.yaml) shown below is used to create StorageClass object.
```
apiVersion: storage.k8s.io/v1
kind: StorageClass
metadata:
  name: csi-qingcloud
provisioner: disk.csi.qingcloud.com
parameters:
  type: "0"
  maxSize: "500"
  minSize: "10"
  stepSize: "10"
  fsType: "ext4"
  replica: "2"
  tags: "tag-y7uu1q2a"
reclaimPolicy: Delete 
```

- `type`: The type of volume in QingCloud IaaS platform. In QingCloud public cloud platform, `0` represents high performance volume. `3` respresents super high performance volume. `2` represents high capacity volume depending on cluster‘s zone. `5` represents enterprise distributed SAN (NeonSAN) volume. `6` represents NeonSAN HDD volume. `100` represents standard volume. `200` represents SSD enterprise volume. See [QingCloud docs](https://docs.qingcloud.com/product/api/action/volume/create_volumes.html) for details.

- `maxSize`, `minSize`: Limit the range of volume size in GiB.

- `stepSize`: Set the increment of volumes size in GiB.

- `fsType`: `ext3`, `ext4`, `xfs`. Default `ext4`.

- `replica`: `1` means single replica, `2` means multiple replicas. Default `2`.

- `tags`: Tags ID in QingCloud IaaS platform split with a comma.

## Support
If you have any qustions or suggestions, please submit an issue at [qingcloud-csi](https://github.com/yunify/qingcloud-csi/issues)
