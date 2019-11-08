# QingCloud-CSI

[![Build Status](https://travis-ci.org/yunify/qingcloud-csi.svg?branch=master)](https://travis-ci.org/yunify/qingcloud-csi)
[![Go Report Card](https://goreportcard.com/badge/github.com/yunify/qingcloud-csi)](https://goreportcard.com/report/github.com/yunify/qingcloud-csi)

> [English](README.md) | 中文

<!-- START doctoc generated TOC please keep comment here to allow auto update -->
<!-- DON'T EDIT THIS SECTION, INSTEAD RE-RUN doctoc TO UPDATE -->
- [描述](#%E6%8F%8F%E8%BF%B0)
- [块存储插件](#%E5%9D%97%E5%AD%98%E5%82%A8%E6%8F%92%E4%BB%B6)
  - [Kubernetes 适配性](#kubernetes-%E9%80%82%E9%85%8D%E6%80%A7)
  - [功能](#%E5%8A%9F%E8%83%BD)
  - [安装](#%E5%AE%89%E8%A3%85)
  - [卸载](#%E5%8D%B8%E8%BD%BD)
  - [使用文档](#%E4%BD%BF%E7%94%A8%E6%96%87%E6%A1%A3)
- [支持](#%E6%94%AF%E6%8C%81)

<!-- END doctoc generated TOC please keep comment here to allow auto update -->

---

## 描述
QingCloud CSI 插件实现了 [CSI](https://github.com/container-storage-interface/) 接口，并使容器编排平台能够使用 QingCloud 云平台的存储资源。目前 QingCloud CSI 实现了块存储插件，可以对接云平台块存储资源。

## 块存储插件

插件的设计和安装使用 Kubernetes 社区推荐的 CSI 插件[架构](https://github.com/kubernetes/community/blob/master/contributors/design-proposals/storage/container-storage-interface.md#recommended-mechanism-for-deploying-csi-drivers-on-kubernetes)，插件架构包含 Controller 和 Node 两部分，在 Controller 部分，由 Deployment 在 Kubernetes 集群内创建一个 Pod 副本。在 Node 部分，每个可调度的节点由 DaemonSet 创建一个 Pod 副本。插件已经在 Kubernetes v1.15 环境中通过了 [CSI 测试](https://github.com/kubernetes-csi/csi-test)。

块存储插件部署后, 用户可创建访问模式（Access Mode）为单节点读写（ReadWriteOnce）的基于 QingCloud 的基础型、SSD 企业型、性能型、超高性能型、超高性能容量型（NeonSAN）、NeonSAN HDD 型、容量型硬盘的存储卷并挂载至工作负载。

### Kubernetes 适配性

| |Kubernetes v1.10-v1.13|Kubernetes v1.14-1.15|
|:---:|:---:|:---:|
|QingCloud CSI v0.2.x|✓|-|
|QingCloud CSI v1.1.0|-|✓|

### 功能

| | 存储卷管理* | 存储卷扩容 | 存储卷监控 |存储卷克隆| 快照管理**|拓扑|
|:---:|:---:|:---:|:---:|:---:|:---:|:---:|
|QingCloud CSI v0.2.x |✓|-|-|-|-|-|
|QingCloud CSI v1.1.0 |✓|✓|✓|✓|✓|✓|

注：
- `*`: 存储卷管理包括存储卷创建/删除和存储卷挂载/卸载至容器组
- `**`: 快照管理包括快照创建/删除和从快照恢复存储卷

### 安装
此安装指南将 CSI 插件安装在 Kubernetes v1.14+ 的 *kube-system* namespace 内。用户也可以将插件部署在其他 namespace 内。

- 设置 Kubernetes 参数
  - kube-apiserver, kube-controller-manager, kube-scheduler, kubelet 设置 `--allow-privileged=true`。
  - 启用（默认开启）[Mount Propagation](https://kubernetes.io/docs/concepts/storage/volumes/#mount-propagation) 特性。
  - kube-apiserver, kube-controller-manager, kube-scheduler, kubelet 设置 `--feature-gates=CSINodeInfo=true,CSIDriverRegistry=true,KubeletPluginsWatcher=true,VolumeSnapshotDataSource=true,ExpandCSIVolumes=true,VolumePVCDataSource=true（仅限 Kubernetes v1.15）` 
  - kubelet 设置 `--read-only-port=10255`

- 下载安装文件并解压
```
$ wget https://raw.githubusercontent.com/yunify/qingcloud-csi/master/deploy/disk/kubernetes/releases/qingcloud-csi-disk-v1.1.0.yaml
```

- 修改 QingCloud 云平台配置参数

    QingCloud CSI 插件通过 QingCloud 云平台 API 调用云平台资源，用户应首先通过 [QingCloud CLI](https://docs.qingcloud.com/product/cli/) 测试 QingCloud 云平台 API 连通性和 QingCloud 云平台配置参数。
  * 修改安装文件内 `csi-qingcloud` ConfigMap 配置项
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
    - `qy_access_key_id`, `qy_secret_access_key`: 在 QingCloud 控制台创建 Access key 密钥. 此密钥需要有操作 QingCloud 云平台资源的权限。

    - `zone`: `zone` 应与 Kubernetes 集群所在区或可用区相同。CSI 插件将会操作此区或可用区的存储卷资源。例如：`zone` 可以设置为 `ap2a` 或 `sh1`。
    
    - `host`, `port`. `protocol`, `uri`: 共同构成 QingCloud IaaS 平台服务的 url。

- 部署 CSI 插件
> 注:  如果 Kubernetes 集群的 [kubelet](https://kubernetes.io/docs/reference/command-line-tools-reference/kubelet/) 设置了 `--root-dir` 选项（默认值为 *"/var/lib/kubelet"*），请将 DaemonSet 的 `spec.template.spec.containers[name=csi-qingcloud].volumeMounts[name=mount-dir].mountPath` 和 `spec.template.spec.volumes[name=mount-dir].hostPath.path` 的值 *"/var/lib/kubelet"* 替换为 `--root-dir` 选项的值。例如：在通过 QingCloud AppCenter 创建的 Kubernetes 集群内, 需要将 DaemonSet 的 *"/var/lib/kubelet"* 字段替换为 *"/data/var/lib/kubelet"*。

```
$ kubectl apply -f qingcloud-csi-disk-v1.1.0.yaml
```

- 检查 CSI 插件状态
```
$ kubectl get pods -n kube-system --selector=app=csi-qingcloud
  NAME                                       READY   STATUS    RESTARTS   AGE
  csi-qingcloud-controller-5bd48bb49-dw9rs   5/5     Running   0          3h16m
  csi-qingcloud-node-d2kdt                   2/2     Running   0          3h16m
  csi-qingcloud-node-hvtq7                   2/2     Running   0          3h16m
  csi-qingcloud-node-njghb                   2/2     Running   0          3h16m
  csi-qingcloud-node-wssdt                   2/2     Running   0          3h16m
```

### 卸载
```
$ kubectl delete -f qingcloud-csi-disk-v1.1.0.yaml
```

### 文档
参数配置和功能用法请参考[使用文档](docs/user-guide-zh.md)。
开发者请参考[开发文档](docs/developer-guide-zh.md)。

## 支持
如果有任何问题或建议, 请在 [qingcloud-csi](https://github.com/yunify/qingcloud-csi/issues) 项目提 issue。
