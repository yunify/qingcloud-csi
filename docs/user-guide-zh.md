<!-- START doctoc generated TOC please keep comment here to allow auto update -->
<!-- DON'T EDIT THIS SECTION, INSTEAD RE-RUN doctoc TO UPDATE -->
- [使用指南](#%E4%BD%BF%E7%94%A8%E6%8C%87%E5%8D%97)
  - [如何设置存储类型](#%E5%A6%82%E4%BD%95%E8%AE%BE%E7%BD%AE%E5%AD%98%E5%82%A8%E7%B1%BB%E5%9E%8B)
    - [存储类型模版](#%E5%AD%98%E5%82%A8%E7%B1%BB%E5%9E%8B%E6%A8%A1%E7%89%88)
    - [存储卷参数](#%E5%AD%98%E5%82%A8%E5%8D%B7%E5%8F%82%E6%95%B0)
      - [type, maxSize, minSize, stepSize](#type-maxsize-minsize-stepsize)
      - [fsType](#fstype)
      - [replica](#replica)
      - [tags](#tags)
    - [其他参数](#%E5%85%B6%E4%BB%96%E5%8F%82%E6%95%B0)
      - [设置默认存储类型](#%E8%AE%BE%E7%BD%AE%E9%BB%98%E8%AE%A4%E5%AD%98%E5%82%A8%E7%B1%BB%E5%9E%8B)
      - [扩容参数](#%E6%89%A9%E5%AE%B9%E5%8F%82%E6%95%B0)
      - [拓扑参数](#%E6%8B%93%E6%89%91%E5%8F%82%E6%95%B0)
    - [硬盘类型与 type 参数对应关系](#%E7%A1%AC%E7%9B%98%E7%B1%BB%E5%9E%8B%E4%B8%8E-type-%E5%8F%82%E6%95%B0%E5%AF%B9%E5%BA%94%E5%85%B3%E7%B3%BB)
    - [主机类型与 type 参数对应关系](#%E4%B8%BB%E6%9C%BA%E7%B1%BB%E5%9E%8B%E4%B8%8E-type-%E5%8F%82%E6%95%B0%E5%AF%B9%E5%BA%94%E5%85%B3%E7%B3%BB)
    - [硬盘类型与主机适配性](#%E7%A1%AC%E7%9B%98%E7%B1%BB%E5%9E%8B%E4%B8%8E%E4%B8%BB%E6%9C%BA%E9%80%82%E9%85%8D%E6%80%A7)
  - [存储卷管理](#%E5%AD%98%E5%82%A8%E5%8D%B7%E7%AE%A1%E7%90%86)
    - [存储卷管理环境](#%E5%AD%98%E5%82%A8%E5%8D%B7%E7%AE%A1%E7%90%86%E7%8E%AF%E5%A2%83)
      - [存储卷管理安装 QingCloud CSI 存储类型](#%E5%AD%98%E5%82%A8%E5%8D%B7%E7%AE%A1%E7%90%86%E5%AE%89%E8%A3%85-qingcloud-csi-%E5%AD%98%E5%82%A8%E7%B1%BB%E5%9E%8B)
    - [存储卷管理创建存储卷](#%E5%AD%98%E5%82%A8%E5%8D%B7%E7%AE%A1%E7%90%86%E5%88%9B%E5%BB%BA%E5%AD%98%E5%82%A8%E5%8D%B7)
    - [存储卷管理挂载存储卷](#%E5%AD%98%E5%82%A8%E5%8D%B7%E7%AE%A1%E7%90%86%E6%8C%82%E8%BD%BD%E5%AD%98%E5%82%A8%E5%8D%B7)
    - [存储卷管理卸载存储卷](#%E5%AD%98%E5%82%A8%E5%8D%B7%E7%AE%A1%E7%90%86%E5%8D%B8%E8%BD%BD%E5%AD%98%E5%82%A8%E5%8D%B7)
    - [存储卷管理删除存储卷](#%E5%AD%98%E5%82%A8%E5%8D%B7%E7%AE%A1%E7%90%86%E5%88%A0%E9%99%A4%E5%AD%98%E5%82%A8%E5%8D%B7)
  - [存储卷扩容](#%E5%AD%98%E5%82%A8%E5%8D%B7%E6%89%A9%E5%AE%B9)
    - [存储卷扩容准备环境](#%E5%AD%98%E5%82%A8%E5%8D%B7%E6%89%A9%E5%AE%B9%E5%87%86%E5%A4%87%E7%8E%AF%E5%A2%83)
    - [存储卷扩容卸载存储卷](#%E5%AD%98%E5%82%A8%E5%8D%B7%E6%89%A9%E5%AE%B9%E5%8D%B8%E8%BD%BD%E5%AD%98%E5%82%A8%E5%8D%B7)
    - [扩容存储卷](#%E6%89%A9%E5%AE%B9%E5%AD%98%E5%82%A8%E5%8D%B7)
    - [存储卷扩容检查](#%E5%AD%98%E5%82%A8%E5%8D%B7%E6%89%A9%E5%AE%B9%E6%A3%80%E6%9F%A5)
  - [存储卷克隆](#%E5%AD%98%E5%82%A8%E5%8D%B7%E5%85%8B%E9%9A%86)
    - [存储卷克隆准备环境](#%E5%AD%98%E5%82%A8%E5%8D%B7%E5%85%8B%E9%9A%86%E5%87%86%E5%A4%87%E7%8E%AF%E5%A2%83)
    - [克隆存储卷](#%E5%85%8B%E9%9A%86%E5%AD%98%E5%82%A8%E5%8D%B7)
  - [快照管理](#%E5%BF%AB%E7%85%A7%E7%AE%A1%E7%90%86)
    - [快照管理准备环境](#%E5%BF%AB%E7%85%A7%E7%AE%A1%E7%90%86%E5%87%86%E5%A4%87%E7%8E%AF%E5%A2%83)
      - [创建带数据的存储卷 pvc-snap-1](#%E5%88%9B%E5%BB%BA%E5%B8%A6%E6%95%B0%E6%8D%AE%E7%9A%84%E5%AD%98%E5%82%A8%E5%8D%B7-pvc-snap-1)
    - [创建快照](#%E5%88%9B%E5%BB%BA%E5%BF%AB%E7%85%A7)
      - [创建快照类型](#%E5%88%9B%E5%BB%BA%E5%BF%AB%E7%85%A7%E7%B1%BB%E5%9E%8B)
      - [创建快照资源对象](#%E5%88%9B%E5%BB%BA%E5%BF%AB%E7%85%A7%E8%B5%84%E6%BA%90%E5%AF%B9%E8%B1%A1)
    - [从快照恢复存储卷](#%E4%BB%8E%E5%BF%AB%E7%85%A7%E6%81%A2%E5%A4%8D%E5%AD%98%E5%82%A8%E5%8D%B7)
      - [恢复存储卷 pvc-snap-2](#%E6%81%A2%E5%A4%8D%E5%AD%98%E5%82%A8%E5%8D%B7-pvc-snap-2)
      - [从快照恢复存储卷检查存储卷数据](#%E4%BB%8E%E5%BF%AB%E7%85%A7%E6%81%A2%E5%A4%8D%E5%AD%98%E5%82%A8%E5%8D%B7%E6%A3%80%E6%9F%A5%E5%AD%98%E5%82%A8%E5%8D%B7%E6%95%B0%E6%8D%AE)
    - [删除快照](#%E5%88%A0%E9%99%A4%E5%BF%AB%E7%85%A7)
  - [拓扑感知](#%E6%8B%93%E6%89%91%E6%84%9F%E7%9F%A5)
    - [拓扑感知准备环境](#%E6%8B%93%E6%89%91%E6%84%9F%E7%9F%A5%E5%87%86%E5%A4%87%E7%8E%AF%E5%A2%83)
      - [拓扑感知 Kubernetes 集群](#%E6%8B%93%E6%89%91%E6%84%9F%E7%9F%A5-kubernetes-%E9%9B%86%E7%BE%A4)
      - [拓扑感知配置存储类型](#%E6%8B%93%E6%89%91%E6%84%9F%E7%9F%A5%E9%85%8D%E7%BD%AE%E5%AD%98%E5%82%A8%E7%B1%BB%E5%9E%8B)
    - [拓扑感知创建存储卷](#%E6%8B%93%E6%89%91%E6%84%9F%E7%9F%A5%E5%88%9B%E5%BB%BA%E5%AD%98%E5%82%A8%E5%8D%B7)
    - [拓扑感知创建工作负载](#%E6%8B%93%E6%89%91%E6%84%9F%E7%9F%A5%E5%88%9B%E5%BB%BA%E5%B7%A5%E4%BD%9C%E8%B4%9F%E8%BD%BD)
  - [静态创建存储卷](#%E9%9D%99%E6%80%81%E5%88%9B%E5%BB%BA%E5%AD%98%E5%82%A8%E5%8D%B7)
    - [静态创建存储卷步骤](#%E9%9D%99%E6%80%81%E5%88%9B%E5%BB%BA%E5%AD%98%E5%82%A8%E5%8D%B7%E6%AD%A5%E9%AA%A4)
      - [静态创建存储卷准备环境](#%E9%9D%99%E6%80%81%E5%88%9B%E5%BB%BA%E5%AD%98%E5%82%A8%E5%8D%B7%E5%87%86%E5%A4%87%E7%8E%AF%E5%A2%83)
      - [静态创建存储卷创建存储类型](#%E9%9D%99%E6%80%81%E5%88%9B%E5%BB%BA%E5%AD%98%E5%82%A8%E5%8D%B7%E5%88%9B%E5%BB%BA%E5%AD%98%E5%82%A8%E7%B1%BB%E5%9E%8B)
      - [静态创建存储卷创建 PV](#%E9%9D%99%E6%80%81%E5%88%9B%E5%BB%BA%E5%AD%98%E5%82%A8%E5%8D%B7%E5%88%9B%E5%BB%BA-pv)
      - [静态创建存储卷创建 PVC](#%E9%9D%99%E6%80%81%E5%88%9B%E5%BB%BA%E5%AD%98%E5%82%A8%E5%8D%B7%E5%88%9B%E5%BB%BA-pvc)
    - [静态创建存储卷使用场景](#%E9%9D%99%E6%80%81%E5%88%9B%E5%BB%BA%E5%AD%98%E5%82%A8%E5%8D%B7%E4%BD%BF%E7%94%A8%E5%9C%BA%E6%99%AF)
      - [存储插件升级，迁移块存储](#%E5%AD%98%E5%82%A8%E6%8F%92%E4%BB%B6%E5%8D%87%E7%BA%A7%E8%BF%81%E7%A7%BB%E5%9D%97%E5%AD%98%E5%82%A8)

<!-- END doctoc generated TOC please keep comment here to allow auto update -->

# 使用指南

## 如何设置存储类型
### 存储类型模版

如下所示的 StorageClass 资源定义可用来创建 StorageClass 对象。
```
apiVersion: storage.k8s.io/v1
kind: StorageClass
metadata:
  annotations:
    storageclass.beta.kubernetes.io/is-default-class: "true"
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
allowVolumeExpansion: true
volumeBindingMode: Immediate
```

### 存储卷参数
存储卷类型模板中 `.parameters` 设置存储卷参数

#### type, maxSize, minSize, stepSize
详情见 [QingCloud 文档](https://docs.qingcloud.com/product/api/action/volume/create_volumes.html)。

|硬盘类型|type|maxSize|minSize|stepSize|
|:---:|:---:|:---:|:---:|:---:|
|性能型|0|2000|10|10|
|容量型|2|5000|100|50|
|超高性能型|3|2000|10|10|
| NeonSAN|5|50000|100|100|
|NeonSAN HDD|5|50000|100|100|
| 基础型|100|2000|10|10|
| SSD 企业型|200| 2000|10|10|

#### fsType
支持 `ext3`, `ext4`, `xfs`. 默认为 `ext4`。

#### replica
`1` 代表单副本硬盘，`2` 代表多副本硬盘。 默认为 `2`。

#### tags
青云云平台 tag ID，多个 tag 用逗号分割，可以将插件创建的硬盘或快照与 tag 绑定。

### 其他参数

#### 设置默认存储类型
存储类型模版中 `.metadata.annotations.storageclass.beta.kubernetes.io/is-default-class` 的值设置为 `true` 表明此存储类型设置为默认存储类型。详见 [Kubernetes 官方文档](https://kubernetes.io/docs/tasks/administer-cluster/change-default-storage-class/)

#### 扩容参数
存储类型模版中 `.allowVolumeExpansion` 的值可填写 `true` 或 `false`, 设置是否支持扩容存储卷。详见 [Kubernetes 官方文档](https://kubernetes.io/docs/concepts/storage/storage-classes/#allow-volume-expansion)

#### 拓扑参数
存储类型模版中 `.volumeBindingMode` 的值可填写 `Immediate` 或 `WaitForFirstConsumer`，通常设置为立即绑定存储卷 `Immediate`，如果 Kubernetes 节点是不同类型主机或跨可用区主机，应设置为延迟绑定存储卷 `WaitForFirstConsumer`。详见 [Kubernetes 官方文档](https://kubernetes.io/docs/concepts/storage/storage-classes/#volume-binding-mode)

### 硬盘类型与 type 参数对应关系

 |硬盘|Volume|type 值|
|:---:|:----:|:----:|
|性能型| High Performance|0|
|容量型| High Capacity|2|
|超高性能型|Super High Performance|3|
|NeonSAN| NeonSAN|5|
|NeonSAN HDD|NeonSAN HDD|6|
|基础型| Standard|100|
|SSD 企业型| SSD Enterprise|200|

### 主机类型与 type 参数对应关系
|主机|英文名|type 值|
|:---:|:----:|:----:|
|性能型|High Performance|0|
|超高性能型|Super High Performance|1|
|SAN存储超高性能型|Super High Performance SAN|6|
|SAN存储性能型|High Performance SAN|7|
|基础型|Standard|101|
|企业型e1|Enterprise1|201|
|企业型e2|Enterprise2|202|
|专业增强型|Premium|301|

### 硬盘类型与主机适配性

 |          | 性能型硬盘    | 容量型硬盘  | 超高性能型硬盘 | NeonSAN 硬盘 | NeonSAN HDD 硬盘 | 基础型硬盘| SSD 企业型硬盘|
|-----------|------------------|------------------|-----------------|---------|----------|-------|-------|
|性能型主机| ✓        | ✓                | -               | ✓      |  ✓      | ✓     | -     |
|超高性能型主机| -       | ✓                | ✓               |✓      |✓     |-  |✓  |
|SAN存储超高性能型| -       | -                | -              |✓      |-     |-  |-  |
|SAN存储性能型| -       | -                | -               |-     |✓     |-  |-  |
|基础型主机| ✓        | ✓                | -               |✓  |✓ |✓ |-  |
|企业型e1主机| -       | ✓                | ✓               |✓  |✓ |-  |✓  |
|企业型e1主机| -       | ✓                | ✓               |✓  |✓ |-  |✓  |
|专业增强型| -       | ✓                | ✓               |✓  |✓ |-  |✓  |

## 存储卷管理
存储卷（PVC，PersistentVolumeClaim）管理功能包括动态分配存储卷，删除存储卷，挂载存储卷到 Pod，从 Pod 卸载存储卷。用户可参考[示例 YAML 文件](https://github.com/yunify/qingcloud-csi/tree/master/deploy/disk/example/volume)。

### 存储卷管理环境
- Kubernetes 1.14+ 集群
- 安装了 QingCloud CSI 存储插件
- 安装了 QingCloud CSI 存储类型

#### 存储卷管理安装 QingCloud CSI 存储类型
- 安装
```console
$ kubectl create -f sc.yaml
```
- 检查
```console
$ kubectl get sc
NAME            PROVISIONER              AGE
csi-qingcloud   disk.csi.qingcloud.com   14m
```

### 存储卷管理创建存储卷
- 创建存储卷
```console
$ kubectl create -f pvc.yaml 
persistentvolumeclaim/pvc-example created
```
- 检查存储卷
```console
$ kubectl get pvc pvc-example
NAME          STATUS   VOLUME                                     CAPACITY   ACCESS MODES   STORAGECLASS    AGE
pvc-example   Bound    pvc-76429525-a930-11e9-9a6a-5254ef68c8c1   20Gi       RWO            csi-qingcloud   25m
```

### 存储卷管理挂载存储卷
- 创建 Deployment 挂载存储卷
```console
$ kubectl create -f deploy-nginx.yaml 
deployment.apps/deploy-nginx created
```
- 访问容器内挂载存储卷的目录
```console
$ kubectl exec -ti deploy-nginx-84474cf674-zfhbs /bin/bash
# cd /mnt
# ls
lost+found
```

### 存储卷管理卸载存储卷
- 删除 deployment 卸载存储卷
```console
$ kubectl delete deploy deploy-nginx
deployment.extensions "deploy-nginx" deleted
```

### 存储卷管理删除存储卷
- 删除存储卷
```console
$ kubectl delete pvc pvc-example
persistentvolumeclaim "pvc-example" deleted
```
- 检查存储卷
```console
$ kubectl get pvc pvc-example
Error from server (NotFound): persistentvolumeclaims "pvc-example" not found
```

## 存储卷扩容
扩容功能将扩大存储卷可用容量。由于云平台限制，本存储插件仅支持离线扩容硬盘。离线扩容硬盘流程是 1. 存储卷处于未挂载状态，2. 扩容存储卷，3. 挂载一次存储卷。示例 YAML 文件在 https://github.com/yunify/qingcloud-csi/tree/master/deploy/disk/volume 内。

### 存储卷扩容准备环境
- Kubernetes 1.14+ 集群
- Kubernetes 组件的 `feature-gate` 增加 `ExpandCSIVolumes=true`
- 配置 QingCloud CSI 存储类型，并将其 `allowVolumeExpansion` 字段值设置为 `true`
- 创建一个存储卷并挂载至 Pod，参考存储卷管理

### 存储卷扩容卸载存储卷
```console
$ kubectl scale deploy deploy-nginx --replicas=0
```

### 扩容存储卷
- 修改存储卷容量
```console
$ kubectl patch pvc pvc-example -p '{"spec":{"resources":{"requests":{"storage": "40Gi"}}}}'
persistentvolumeclaim/pvc-example patched
```
- 挂载存储卷
```console
$ kubectl scale deploy deploy-nginx --replicas=1
```
- 完成扩容
```console
$ kubectl get pvc pvc-example
NAME          STATUS   VOLUME                                     CAPACITY   ACCESS MODES   STORAGECLASS    AGE
pvc-example   Bound    pvc-906f5760-a935-11e9-9a6a-5254ef68c8c1   40Gi       RWO            csi-qingcloud   6m7s
$ kubectl get po
NAME                            READY   STATUS    RESTARTS   AGE
deploy-nginx-6c444c9b7f-d6n29   1/1     Running   0          3m38s
```

### 存储卷扩容检查
- 进入 Pod 查看
```console
$ kubectl exec -ti deploy-nginx-6c444c9b7f-d6n29 /bin/bash
root@deploy-nginx-6c444c9b7f-d6n29:/# s
bash: s: command not found
root@deploy-nginx-6c444c9b7f-d6n29:/# df -ah
Filesystem      Size  Used Avail Use% Mounted on
...
/dev/vdc         40G   49M   40G   1% /mnt
...
```

## 存储卷克隆
存储卷克隆可以创建现有存储卷的副本，示例 YAML 文件在 https://github.com/yunify/qingcloud-csi/tree/master/deploy/disk/example/volume/pvc-clone.yaml 内。

### 存储卷克隆准备环境
- Kubernetes 1.15+ 集群
- Kubernetes 组件的 `feature-gate` 增加 `VolumePVCDataSource=true`
- 安装 QingCloud CSI v1.1.0
- 配置 QingCloud CSI 存储类型
- 创建一个存储卷，参考存储卷管理

### 克隆存储卷
- 查询已存在存储卷
```console
$ kubectl get pvc
NAME          STATUS   VOLUME                                     CAPACITY   ACCESS MODES   STORAGECLASS    AGE
pvc-example   Bound    pvc-d1fb263e-b368-4339-8f8b-448446f4b840   20Gi       RWO            csi-qingcloud   32s
```

- 克隆存储卷
```console
$ kubectl create -f pvc-clone.yaml
persistentvolumeclaim/pvc-clone created
```

- 查询克隆存储卷
```console
$ kubectl get pvc pvc-clone
NAME        STATUS   VOLUME                                     CAPACITY   ACCESS MODES   STORAGECLASS    AGE
pvc-clone   Bound    pvc-529d2502-02bd-442b-a69f-d3eff28316a8   20Gi       RWO            csi-qingcloud   31s
```

## 快照管理
快照功能包括创建和删除快照，从快照恢复存储卷功能。示例 YAML 文件在 https://github.com/yunify/qingcloud-csi/tree/master/deploy/disk/example/snapshot 内。

### 快照管理准备环境
- Kubernetes 1.14+ 集群
- 在 kube-apiserver, kube-controller-manager 的 `feature-gate` 增加 `VolumeSnapshotDataSource=true`
- 安装 QingCloud CSI v1.1.0
- 配置 QingCloud CSI 存储类型
- 创建一个带数据的存储卷

#### 创建带数据的存储卷 pvc-snap-1
- 创建存储卷 
```console
$ kubectl create -f original-pvc.yaml
persistentvolumeclaim/pvc-snap-1 created
```
- 检查存储卷
```console
$ kubectl get pvc
NAME         STATUS   VOLUME                                     CAPACITY   ACCESS MODES   STORAGECLASS    AGE
pvc-snap-1   Bound    pvc-28090960-9eeb-11e9-9a6a-5254ef68c8c1   20Gi       RWO            csi-qingcloud   37s
```
- 向存储卷写数据
```console
$ kubectl create -f deploy-writer.yaml 
deployment.apps/fio created

$ kubectl get po
NAME                   READY   STATUS    RESTARTS   AGE
fio-645b5d6499-8tc7p   1/1     Running   0          23s

$ kubectl exec -ti fio-645b5d6499-8tc7p /bin/bash
root@fio-645b5d6499-8tc7p:/# cd root
root@fio-645b5d6499-8tc7p:/# ./start-test.sh
crtl+c (5 秒后执行此命令，停止写数据)
root@fio-645b5d6499-8tc7p:/# ls -lh /mnt
total 20G
drwx------ 2 root root  16K Jul  5 06:09 lost+found
-rw-r--r-- 1 root root    0 Jul  5 06:10 rand-write.0.0
-rw-r--r-- 1 root root    0 Jul  5 06:10 rand-write.1.0
-rw-r--r-- 1 root root    0 Jul  5 06:10 rand-write.2.0
-rw-r--r-- 1 root root    0 Jul  5 06:10 rand-write.3.0
-rw-r--r-- 1 root root  10G Jul  5 06:10 seq-write.0.0
-rw-r--r-- 1 root root    0 Jul  5 06:10 seq-write.1.0
-rw-r--r-- 1 root root    0 Jul  5 06:10 seq-write.2.0
-rw-r--r-- 1 root root 9.6G Jul  5 06:10 seq-write.3.0
```

### 创建快照
注意：每个 Kubernetes 快照对应于一个 QingCloud 全量备份，请确保有足够全量备份链配额。

#### 创建快照类型
```console
$ kubectl create -f snapshot-class.yaml 
volumesnapshotclass.snapshot.storage.k8s.io/csi-qingcloud created

$ kubectl get volumesnapshotclass
NAME            AGE
csi-qingcloud   16s
```

#### 创建快照资源对象
```console
$ kubectl create -f volume-snapshot.yaml 
volumesnapshot.snapshot.storage.k8s.io/snap-1 created

$ kubectl get volumesnapshot
NAME     AGE
snap-1   91s
```

### 从快照恢复存储卷
#### 恢复存储卷 pvc-snap-2
```console
$ kubectl create -f restore-pvc.yaml 
persistentvolumeclaim/pvc-snap-2 created
```

```console
$ kubectl get pvc pvc-snap-2
NAME         STATUS   VOLUME                                     CAPACITY   ACCESS MODES   STORAGECLASS    AGE
pvc-snap-2   Bound    pvc-b8a05427-9eef-11e9-9a6a-5254ef68c8c1   20Gi       RWO            csi-qingcloud   52s
```

#### 从快照恢复存储卷检查存储卷数据
从快照恢复的存储卷 `pvc-snap-2` 与在创建快照时的 `pvc-snap-1` 内容应一致。

```console
$ kubectl create -f deploy-viewer.yaml 
deployment.apps/nginx created

$ kubectl get po |grep snap-example 
snap-example-85dd9b646c-56g85   1/1     Running   0          3m6s

$ kubectl exec -ti snap-example-85dd9b646c-56g85 /bin/bash
root@snap-example-85dd9b646c-56g85:/# ls /mnt -lh
total 20G
drwx------ 2 root root  16K Jul  5 06:09 lost+found
-rw-r--r-- 1 root root    0 Jul  5 06:10 rand-write.0.0
-rw-r--r-- 1 root root    0 Jul  5 06:10 rand-write.1.0
-rw-r--r-- 1 root root    0 Jul  5 06:10 rand-write.2.0
-rw-r--r-- 1 root root    0 Jul  5 06:10 rand-write.3.0
-rw-r--r-- 1 root root  10G Jul  5 06:10 seq-write.0.0
-rw-r--r-- 1 root root    0 Jul  5 06:10 seq-write.1.0
-rw-r--r-- 1 root root    0 Jul  5 06:10 seq-write.2.0
-rw-r--r-- 1 root root 9.6G Jul  5 06:10 seq-write.3.0
```

### 删除快照

```console
$ kubectl delete volumesnapshot snap-1
volumesnapshot.snapshot.storage.k8s.io "snap-1" deleted
```

## 拓扑感知
在跨可用区 Kubernetes 集群和拥有不同类型节点的 Kubernetes 集群中创建和挂载存储卷需要拓扑功能。示例 YAML 文件在 https://github.com/yunify/qingcloud-csi/tree/master/deploy/disk/example/topology 内。

### 拓扑感知准备环境
- Kubernetes 1.14+ 集群
- 在 Kubernetes 控制平面和 Kubelet 的 `feature-gate` 增加 `CSINodeInfo=true`，默认为 `true`
- 安装 QingCloud CSI v1.1.0 存储插件，`external-provisioner` 边车容器的 `feature-gate` 增加 `Topology=true`，默认为 `true`
- 配置 QingCloud CSI 存储类型

#### 拓扑感知 Kubernetes 集群
本例使用跨 Pek3 可用区的 Kubernetes v1.15 集群。集群中 node1 和 node2 在 Pek3c, node3 和 node4 在 Pek3b，node5 和 node6 在 Pek3d。node 类型均为基础型。

#### 拓扑感知配置存储类型
- 拓扑功能的存储类型中 `volumeBindingMode` 字段的值默认设置为 `WaitForFirstConsumer`，这样可以按照 Kubernetes 调度容器组情况，在相应的可用区创建 存储卷。如果设置为 `Immediate` 将会在容器组调度之前创建存储卷，会限制容器组调度。

```console
$ kubectl create -f sc.yaml
```

```console
$ kubectl get sc csi-qingcloud -oyaml
allowVolumeExpansion: true
apiVersion: storage.k8s.io/v1
kind: StorageClass
metadata:
  name: csi-qingcloud
parameters:
  fsType: ext4
  maxSize: "5000"
  minSize: "10"
  replica: "2"
  stepSize: "10"
  type: "100"
provisioner: disk.csi.qingcloud.com
reclaimPolicy: Delete
volumeBindingMode: WaitForFirstConsumer
```

### 拓扑感知创建存储卷
- 创建存储卷
```console
$ kubectl create -f pvc.yaml
persistentvolumeclaim/pvc-topology created
```

- 存储卷创建好后 Pending 状态是正常现象，等待容器组调度后存储卷就会自动创建
```console
$ kubectl get pvc pvc-topology
NAME           STATUS    VOLUME   CAPACITY   ACCESS MODES   STORAGECLASS    AGE
pvc-topology   Pending                                      csi-qingcloud   31s
```

### 拓扑感知创建工作负载
- 创建实例工作负载
```console
$ kubectl create -f deploy.yaml
deployment.apps/nginx-topology created
```

- 查看容器组，调度到 node3 上
```console
$ kubectl get po -o wide
NAME                      READY   STATUS    RESTARTS   AGE   IP               NODE    NOMINATED NODE   READINESS GATES
nginx-topology-79d8d5d86d-4lvcl    1/1     Running   0          52s   10.233.92.27     node3   <none>           <none>
```

- 查看存储卷状态，此时会自动创建基于 Pek3b 的硬盘的存储卷
```console
$ kubectl get pvc pvc-topology
NAME           STATUS   VOLUME                                     CAPACITY   ACCESS MODES   STORAGECLASS    AGE
pvc-topology   Bound    pvc-5b34120c-6119-4c86-b9de-e152304683e6   20Gi       RWO            csi-qingcloud   2m48s
```

- pvc-topology 这个存储卷包含了硬盘的拓扑信息，之后挂载这个存储卷的容器组将会自动调度到可挂载此存储卷的节点上，在此示例中是 node3 或 node4.

## 静态创建存储卷

静态创建存储卷也称为预分配存储卷，整体流程为：先在青云云平台手动创建块存储，创建PV管理块存储，创建PVC关联PV。删除 PVC 时可以关联删除 PV 和底层块存储。

### 静态创建存储卷步骤

#### 静态创建存储卷准备环境

- 基于青云云平台 Kubernetes 集群
- 安装 QingCloud CSI 插件
- 在青云云平台某区已有块存储卷，

#### 静态创建存储卷创建存储类型

```console
$ kubectl get sc csi-qingcloud -o yaml
apiVersion: storage.k8s.io/v1
kind: StorageClass
metadata:
  annotations:
    kubectl.kubernetes.io/last-applied-configuration: |
      {"apiVersion":"storage.k8s.io/v1","kind":"StorageClass","metadata":{"annotations":{},"name":"csi-qingcloud","namespace":""},"parameters":{"fsType":"ext4","maxSize":"500","minSize":"10","stepSize":"10","type":"0"},"provisioner":"csi-qingcloud","reclaimPolicy":"Delete"}
  creationTimestamp: 2018-08-06T02:20:19Z
  name: csi-qingcloud
  resourceVersion: "1355065"
  selfLink: /apis/storage.k8s.io/v1/storageclasses/csi-qingcloud
  uid: 43f25337-991f-11e8-b5aa-525445c0b555
parameters:
  fsType: ext4
  maxSize: "500"
  minSize: "10"
  stepSize: "10"
  type: "0"
provisioner: csi-qingcloud
reclaimPolicy: Delete
volumeBindingMode: Immediate
```

#### 静态创建存储卷创建 PV
- 本次实验是在 AP2A 区创建的性能型 Kubernetes 1.10 集群，使用在同一区创建的性能型块存储卷，块存储卷名为static-volume，ID为vol-jjtedp2i，容量为 20 GiB。

- 编辑 PV 资源定义文件
```console
$ vi pv.yaml
apiVersion: v1
kind: PersistentVolume
metadata:
  annotations:
    pv.kubernetes.io/provisioned-by: csi-qingcloud
  name: pv-static
spec:
  capacity:
    storage: 20Gi
  volumeMode: Filesystem
  accessModes:
    - ReadWriteOnce
  persistentVolumeReclaimPolicy: Delete
  storageClassName: csi-qingcloud
  csi:
    driver: csi-qingcloud
    fsType: ext4
    volumeAttributes:
      fsType: ext4
      maxSize: "500"
      minSize: "10"
      stepSize: "10"
      type: "0"
    volumeHandle: vol-jjtedp2i
```

- 创建 PV
```console
$ kubectl create -f pv.yaml
```

- 查看 PV 状态
```console
$ kubectl get pv pv-static
NAME        CAPACITY   ACCESS MODES   RECLAIM POLICY   STATUS      CLAIM     STORAGECLASS    REASON    AGE
pv-static   20Gi       RWO            Delete           Available             csi-qingcloud             8m
```

#### 静态创建存储卷创建 PVC

- 编辑 PVC 资源定义文件

```console
$ vi pvc.yaml
apiVersion: v1
kind: PersistentVolumeClaim
metadata:
  annotations:
    volume.beta.kubernetes.io/storage-provisioner: csi-qingcloud
  name: pvc-static
spec:
  accessModes:
  - ReadWriteOnce
  resources:
    requests:
      storage: 20Gi
  storageClassName: csi-qingcloud
```

- 创建 PVC

```console
$ kubectl create -f pvc.yaml
```

- 查看 PVC 和 PV 状态

```console
$  kubectl get pvc pvc-static
NAME         STATUS    VOLUME      CAPACITY   ACCESS MODES   STORAGECLASS    AGE
pvc-static   Bound     pv-static   20Gi       RWO            csi-qingcloud   11s

$ kubectl get pv pv-static 
NAME        CAPACITY   ACCESS MODES   RECLAIM POLICY   STATUS    CLAIM                STORAGECLASS    REASON    AGE
pv-static   20Gi       RWO            Delete           Bound     default/pvc-static   csi-qingcloud             12m
```

### 静态创建存储卷使用场景

#### 存储插件升级，迁移块存储

- 现有块存储
    - flex-volume 插件创建青云块存储名为 pvc-55754c8c-b577-11e8-a480-525445c0b555，ID 为 vol-djwgkjil，容量为10 GiB。

    - PVC
    ```console
    $ kubectl get pvc old-pvc
    NAME      STATUS    VOLUME                                     CAPACITY   ACCESS MODES   STORAGECLASS             AGE
    old-pvc   Bound     pvc-55754c8c-b577-11e8-a480-525445c0b555   10Gi       RWO            qingcloud-storageclass   25s
    ```
    - PV
    ```console
    $ kubectl get pv pvc-55754c8c-b577-11e8-a480-525445c0b555
    NAME                                       CAPACITY   ACCESS MODES   RECLAIM POLICY   STATUS    CLAIM             STORAGECLASS             REASON    AGE
    pvc-55754c8c-b577-11e8-a480-525445c0b555   10Gi       RWO            Delete           Bound     default/old-pvc   qingcloud-storageclass             4m
    ```

    ```console
    $ kubectl get pv pvc-55754c8c-b577-11e8-a480-525445c0b555 -oyaml
    apiVersion: v1
    kind: PersistentVolume
    metadata:
      annotations:
        Provisioner_Id: qingcloud/volume-provisioner
        kubernetes.io/createdby: qingcloud-volume-provisioner
        pv.kubernetes.io/provisioned-by: qingcloud/volume-provisioner
      creationTimestamp: 2018-09-11T04:01:34Z
      finalizers:
      - kubernetes.io/pv-protection
      name: pvc-55754c8c-b577-11e8-a480-525445c0b555
      resourceVersion: "14041782"
      selfLink: /api/v1/persistentvolumes/pvc-55754c8c-b577-11e8-a480-525445c0b555
      uid: 5fa6cb8e-b577-11e8-a480-525445c0b555
    spec:
      accessModes:
      - ReadWriteOnce
      capacity:
        storage: 10Gi
      claimRef:
        apiVersion: v1
        kind: PersistentVolumeClaim
        name: old-pvc
        namespace: default
        resourceVersion: "14041596"
        uid: 55754c8c-b577-11e8-a480-525445c0b555
      flexVolume:
        driver: qingcloud/flex-volume
        fsType: ext4
        options:
          volumeID: vol-djwgkjil
      persistentVolumeReclaimPolicy: Delete
      storageClassName: qingcloud-storageclass
      volumeMode: Filesystem
    status:
      phase: Bound
    ```

  - 块存储内已有tmp文件
  ```console
  # ls
  lost+found  tmp
  # cat tmp
  Tue Sep 11 04:11:11 UTC 2018
  ```

- 块存储与原 PVC 解绑
    - 将 PV 的 spec.persistentVolumeReclaimPolicy 的值从 Delete 改为 Retain
    ```console
    $ kubectl edit pv pvc-55754c8c-b577-11e8-a480-525445c0b555
    apiVersion: v1
    kind: PersistentVolume
    metadata:
      ...
      name: pvc-55754c8c-b577-11e8-a480-525445c0b555
    spec
      ...
      persistentVolumeReclaimPolicy: Retain
      ...
    ```

    - 删除 PVC 和 PV
    ```console
    $ kubectl delete pvc old-pvc
    persistentvolumeclaim "old-pvc" deleted
    $ kubectl delete pv pvc-55754c8c-b577-11e8-a480-525445c0b555
    persistentvolume "pvc-55754c8c-b577-11e8-a480-525445c0b555" deleted
    ```

- 静态创建存储卷绑定块存储

    - 编辑 PV 资源定义文件
    ```console
    $ vi pv.yaml
    apiVersion: v1
    kind: PersistentVolume
    metadata:
      annotations:
        pv.kubernetes.io/provisioned-by: csi-qingcloud
      name: new-pv
    spec:
      capacity:
        storage: 10Gi
      volumeMode: Filesystem
      accessModes:
      - ReadWriteOnce
      persistentVolumeReclaimPolicy: Delete
      storageClassName: csi-qingcloud
      csi:
        driver: csi-qingcloud
        fsType: ext4
        volumeAttributes:
          fsType: ext4
          maxSize: "500"
          minSize: "10"
          stepSize: "10"
          type: "0"
        volumeHandle: vol-djwgkjil
    ```

    - 创建 PV
    ```console
    $ kubectl create -f pv.yaml
    ```

    - 编辑 PVC 资源定义文件
    ```console
    $ vi pvc.yaml
    apiVersion: v1
    kind: PersistentVolumeClaim
    metadata:
      annotations:
        volume.beta.kubernetes.io/storage-provisioner: csi-qingcloud
      name: pvc-static
    spec:
      accessModes:
      - ReadWriteOnce
      resources:
        requests:
          storage: 10Gi
      storageClassName: csi-qingcloud
      volumeMode: Filesystem
      volumeName: new-pv
    ```
   
    - 查看 PVC 和 PV
    ```console
    $ kubectl get pv new-pv
    NAME      CAPACITY   ACCESS MODES   RECLAIM POLICY   STATUS    CLAIM             STORAGECLASS    REASON    AGE
    new-pv    10Gi       RWO            Delete           Bound     default/new-pvc   csi-qingcloud             44s

    $ kubectl get pvc new-pvc
    NAME      STATUS    VOLUME    CAPACITY   ACCESS MODES   STORAGECLASS    AGE
    new-pvc   Bound     new-pv    10Gi       RWO            csi-qingcloud   32s
    ```

- 创建 Deployment 查看PVC里内容

```console
$ kubectl exec -ti new-pvc-55f77cfb9-c7tdd  /bin/bash
root@new-pvc-55f77cfb9-c7tdd:/# ls
bin  boot  dev	etc  home  lib	lib64  media  mnt  opt	proc  root  run  sbin  srv  sys  tmp  usr  var
root@new-pvc-55f77cfb9-c7tdd:/# cd mnt/
root@new-pvc-55f77cfb9-c7tdd:/mnt# ls
lost+found  tmp
root@new-pvc-55f77cfb9-c7tdd:/mnt# cat tmp 
Tue Sep 11 04:11:11 UTC 2018
```

- 删除 PVC
    - 删除 PVC 后 PV 和块存储均能够自动删除
