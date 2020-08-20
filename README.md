# QingCloud-CSI

[![Build Status](https://travis-ci.org/yunify/qingcloud-csi.svg?branch=master)](https://travis-ci.org/yunify/qingcloud-csi)
[![Go Report Card](https://goreportcard.com/badge/github.com/yunify/qingcloud-csi)](https://goreportcard.com/report/github.com/yunify/qingcloud-csi)

<!-- START doctoc generated TOC please keep comment here to allow auto update -->
<!-- DON'T EDIT THIS SECTION, INSTEAD RE-RUN doctoc TO UPDATE -->

- [Description](#description)
- [Disk Plugin](#disk-plugin)
  - [Kubernetes Compatibility Matrix](#kubernetes-compatibility-matrix)
  - [Feature Matrix](#feature-matrix)
  - [Installation](#installation)
  - [Document](#document)
- [Support](#support)

<!-- END doctoc generated TOC please keep comment here to allow auto update -->

---

## Description
QingCloud CSI plugin implements an interface between Container Storage Interface ([CSI](https://github.com/container-storage-interface/)) enabled Container Orchestrator (CO) and the storage of QingCloud. Currently, QingCloud CSI disk plugin has been developed and manages disk volume in QingCloud platform.

## Disk Plugin

Disk plugin's design and installation use Kubernetes community recommended CSI plugin [architecture](https://github.com/kubernetes/community/blob/master/contributors/design-proposals/storage/container-storage-interface.md#recommended-mechanism-for-deploying-csi-drivers-on-kubernetes). Plugin architecture contains Controller part and Node part. In the part of Controller, one Pod is created by Deployment in Kubernetes cluster. In the part of Node, one Pod is created by DaemonSet on every node. Now, it has been passed the [CSI test](https://github.com/kubernetes-csi/csi-test) in Kubernetes v1.15 environment.

After plugin installation completes, user can create volumes based on several types of disk, such as Standard disk, SSD Enterprise disk, High Performance disk, Super High Performance disk, NeonSAN disk, NeonSAN HDD disk and High Capacity disk, with ReadWriteOnce access mode and mount volumes on workloads.

### Kubernetes Compatibility Matrix

|QingCloud CSI|Kubernetes v1.10-v1.13|Kubernetes v1.14-1.15|Kubernetes v1.16|Kubernetes v1.17|
|:---:|:---:|:---:|:---:|:---:|
|v0.2.x|✓|-|-|-|
|v1.1.0|-|✓|-|-|
|v1.2.0|-|-|✓|✓|

### Feature Matrix

|QingCloud CSI | Volume Management* | Volume Expansion | Volume Monitor | Volume Cloning| Snapshot Management**| Topology Awareness|
|:---:|:---:|:---:|:---:|:---:|:---:|:---:|
|v0.2.x |✓|-|-|-|-|-|
|v1.1.0 |✓|✓|✓|✓|✓|✓|
|v1.2.0 |✓|✓|✓|✓|✓***|✓|

Notes:
- `*`: Volume Management including creating/deleting volume and mounting/unmount volume on Pod.
- `**`: Snapshot management including creating/deleting snapshot and restoring volume from snapshot.
- `***`: Only supports Snapshot Management on Kubernetes v1.17+ because snapshot features goes into Beta on this version.

### Installation 
From v1.2.0, QingCloud-CSI will be installed by helm. See [Helm Charts](https://github.com/kubesphere/helm-charts/tree/master/src/test/csi-qingcloud) for details.

### Document
- [User Guide](docs/user-guide.md)
- [Developer Guide](docs/developer-guide.md)

## Support
If you have any questions or suggestions, please submit an issue at [qingcloud-csi](https://github.com/yunify/qingcloud-csi/issues)
