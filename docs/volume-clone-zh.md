# 存储卷克隆
存储卷克隆可以创建现有存储卷的副本，示例 YAML 文件在 https://github.com/yunify/qingcloud-csi/tree/master/deploy/disk/example/volume/pvc-clone.yaml 内。

## 准备工作
- Kubernetes 1.15+ 集群
- Kubernetes 组件的 `feature-gate` 增加 `ExpandCSIVolumes=true`
- 安装 QingCloud CSI v1.1.0
- 配置 QingCloud CSI 存储类型
- 创建一个存储卷，参考[存储卷管理](volume-zh.md)

## 克隆存储卷
- 查询已存在存储卷
```
$ kubectl get pvc
NAME          STATUS   VOLUME                                     CAPACITY   ACCESS MODES   STORAGECLASS    AGE
pvc-example   Bound    pvc-d1fb263e-b368-4339-8f8b-448446f4b840   20Gi       RWO            csi-qingcloud   32s
```

- 克隆存储卷
```
$ kubectl create -f pvc-clone.yaml
persistentvolumeclaim/pvc-clone created
```

- 查询克隆存储卷
```
$ kubectl get pvc pvc-clone
NAME        STATUS   VOLUME                                     CAPACITY   ACCESS MODES   STORAGECLASS    AGE
pvc-clone   Bound    pvc-529d2502-02bd-442b-a69f-d3eff28316a8   20Gi       RWO            csi-qingcloud   31s
```