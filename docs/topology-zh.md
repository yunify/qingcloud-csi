# 拓扑
在跨可用区 Kubernetes 集群和拥有不同类型节点的 Kubernetes 集群中创建和挂载存储卷需要拓扑功能。示例 YAML 文件在 https://github.com/yunify/qingcloud-csi/tree/master/deploy/disk/example/topology 内。

## 准备工作
- Kubernetes 1.14+ 集群
- 在 Kubernetes 控制平面和 Kubelet 的 `feature-gate` 增加 `CSINodeInfo=true`，默认为 `true`
- 安装 QingCloud CSI 存储插件，`external-provisioner` 边车容器的 `feature-gate` 增加 `Topology=true`，默认为 `true`
- 配置存储类型

### Kubernetes 集群
本例使用跨 Pek3 可用区的 Kubernetes v1.15 集群。集群中 node1 和 node2 在 Pek3c, node3 和 node4 在 Pek3b，node5 和 node6 在 Pek3d。node 类型均为基础型。

### 配置存储类型
- 拓扑功能的存储类型中 `volumeBindingMode` 字段的值默认设置为 `WaitForFirstConsumer`，这样可以按照 Kubernetes 调度容器组情况，在相应的可用区创建 存储卷。如果设置为 `Immediate` 将会在容器组调度之前创建存储卷，会限制容器组调度。

```
$ kubectl create -f sc.yaml
```

```
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

## 创建存储卷
- 创建存储卷
```
$ kubectl create -f pvc.yaml
persistentvolumeclaim/pvc-topology created
```

- 存储卷创建好后 Pending 状态是正常现象，等待容器组调度后存储卷就会自动创建
```
$ kubectl get pvc pvc-topology
NAME           STATUS    VOLUME   CAPACITY   ACCESS MODES   STORAGECLASS    AGE
pvc-topology   Pending                                      csi-qingcloud   31s
```

## 创建工作负载
- 创建实例工作负载
```
$ kubectl create -f deploy.yaml
deployment.apps/nginx-topology created
```

- 查看容器组，调度到 node3 上
```
kubectl get po -o wide
NAME                      READY   STATUS    RESTARTS   AGE   IP               NODE    NOMINATED NODE   READINESS GATES
nginx-topology-79d8d5d86d-4lvcl    1/1     Running   0          52s   10.233.92.27     node3   <none>           <none>
```

- 查看存储卷状态，此时会自动创建基于 Pek3b 的硬盘的存储卷
```
$ kubectl get pvc pvc-topology
NAME           STATUS   VOLUME                                     CAPACITY   ACCESS MODES   STORAGECLASS    AGE
pvc-topology   Bound    pvc-5b34120c-6119-4c86-b9de-e152304683e6   20Gi       RWO            csi-qingcloud   2m48s
```

- pvc-topology 这个存储卷包含了硬盘的拓扑信息，之后挂载这个存储卷的容器组将会自动调度到可挂载此存储卷的节点上，在此示例中是 node3 或 node4.