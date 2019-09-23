# 如何设置存储类型

## 存储类型模版

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

## 存储卷参数
存储卷类型模板中 `.parameters` 设置存储卷参数

### `type`, `maxSize`, `minSize`, `stepSize`
详情见 [QingCloud 文档](https://docs.qingcloud.com/product/api/action/volume/create_volumes.html)。

|硬盘类型|type|maxSize|minSize|stepSize|
|:---:|:---:|:---:|:---:|:---:|
|性能型|0|1000|10|10|
|容量型|2|5000|100|50|
|超高性能型|3|1000|10|10|
| NeonSAN|5|5000|100|100|
| 基础型|100|2000|10|10|
| SSD 企业型|200| 2000|10|10|

### `fsType`
支持 `ext3`, `ext4`, `xfs`. 默认为 `ext4`。

### `replica`
`1` 代表单副本硬盘，`2` 代表多副本硬盘。 默认为 `2`。

### `tags`
青云云平台 tag ID，多个 tag 用逗号分割，可以将插件创建的硬盘或快照与 tag 绑定。

## 其他参数

### 设置默认存储类型
存储类型模版中 `.metadata.annotations.storageclass.beta.kubernetes.io/is-default-class` 的值设置为 `true` 表明此存储类型设置为默认存储类型。详见 [Kubernetes 官方文档](https://kubernetes.io/docs/tasks/administer-cluster/change-default-storage-class/)

### 扩容
存储类型模版中 `.allowVolumeExpansion` 的值可填写 `true` 或 `false`, 设置是否支持扩容存储卷。详见 [Kubernetes 官方文档](https://kubernetes.io/docs/concepts/storage/storage-classes/#allow-volume-expansion)

### 拓扑
存储类型模版中 `.volumeBindingMode` 的值可填写 `Immediate` 或 `WaitForFirstConsumer`，通常设置为立即绑定存储卷 `Immediate`，如果 Kubernetes 节点是不同类型主机或跨可用区主机，应设置为延迟绑定存储卷 `WaitForFirstConsumer`。详见 [Kubernetes 官方文档](https://kubernetes.io/docs/concepts/storage/storage-classes/#volume-binding-mode)

## 硬盘类型与 type 参数对应关系

 |硬盘|Volume|type 值|
|:---:|:----:|:----:|
|性能型| High Performance|0|
|容量型| High Capacity|2|
|超高性能型|Super High Performance|3|
|NeonSAN| NeonSAN|5|
|基础型| Standard|100|
|SSD 企业型| SSD Enterprise|200|

## 主机类型与 type 参数对应关系
|主机|英文名|type 值|
|:---:|:----:|:----:|
|性能型|High Performance|0|
|超高性能型|Super High Performance|1|
|基础型|Standard|101|
|企业型|Enterprise|201|
|专业增强型|Premium|301|

 ## 硬盘类型与主机适配性

 |          | 性能型硬盘    | 容量型硬盘  | 超高性能型硬盘 | NeonSAN 硬盘 |基础型硬盘| SSD 企业型硬盘|
|-----------|------------------|------------------|-----------------|---------|----------|-------|
|性能型主机| ✓        | ✓                | -               | ✓      | -     | -     |
|超高性能型主机| -       | ✓                | ✓               |✓  |-  |-  |
|基础型主机| -       | ✓                | -               |✓  |✓  |-  |
|企业型主机| -       | ✓                | -               |✓  |-  |✓  |
|专业增强型| -       | ✓                | -               |✓  |-  |✓  |