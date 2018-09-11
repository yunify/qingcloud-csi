# 静态创建存储卷

静态创建存储卷也称为预分配存储卷，整体流程为：先在青云云平台手动创建块存储，创建PV管理块存储，创建PVC关联PV。删除 PVC 时可以关联删除 PV 和底层块存储。

## 步骤

### 准备资源

- 基于青云云平台在某区创建的性能型或超高性能型 Kubernetes 1.10+ 集群
- 安装 QingCloud CSI 插件
- 在青云云平台某区已有的性能型或超高性能型块存储卷，

### 创建 StorageClass

```
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

### 创建 PV
- 本次实验是在 AP2A 区创建的性能型 Kubernetes 1.10 集群，使用在同一区创建的性能型块存储卷，块存储卷名为static-volume，ID为vol-jjtedp2i，容量为 20 GiB。

- 编辑 PV 资源定义文件
```
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
```
$ kubectl create -f pv.yaml
```

- 查看 PV 状态
```
$ kubectl get pv pv-static
NAME        CAPACITY   ACCESS MODES   RECLAIM POLICY   STATUS      CLAIM     STORAGECLASS    REASON    AGE
pv-static   20Gi       RWO            Delete           Available             csi-qingcloud             8m
```

### 创建 PVC

- 编辑 PVC 资源定义文件

```
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

```
$ kubectl create -f pvc.yaml
```

- 查看 PVC 和 PV 状态

```
$  kubectl get pvc pvc-static
NAME         STATUS    VOLUME      CAPACITY   ACCESS MODES   STORAGECLASS    AGE
pvc-static   Bound     pv-static   20Gi       RWO            csi-qingcloud   11s

$ kubectl get pv pv-static 
NAME        CAPACITY   ACCESS MODES   RECLAIM POLICY   STATUS    CLAIM                STORAGECLASS    REASON    AGE
pv-static   20Gi       RWO            Delete           Bound     default/pvc-static   csi-qingcloud             12m
```

## 使用场景

### 存储插件升级，迁移块存储

- 现有块存储
    - flex-volume 插件创建青云块存储名为 pvc-55754c8c-b577-11e8-a480-525445c0b555，ID 为 vol-djwgkjil，容量为10 GiB。

    - PVC
    ```
    $ kubectl get pvc old-pvc
    NAME      STATUS    VOLUME                                     CAPACITY   ACCESS MODES   STORAGECLASS             AGE
    old-pvc   Bound     pvc-55754c8c-b577-11e8-a480-525445c0b555   10Gi       RWO            qingcloud-storageclass   25s
    ```
    - PV
    ```
    $ kubectl get pv pvc-55754c8c-b577-11e8-a480-525445c0b555
    NAME                                       CAPACITY   ACCESS MODES   RECLAIM POLICY   STATUS    CLAIM             STORAGECLASS             REASON    AGE
    pvc-55754c8c-b577-11e8-a480-525445c0b555   10Gi       RWO            Delete           Bound     default/old-pvc   qingcloud-storageclass             4m
    ```

    ```
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
  ```
  # ls
  lost+found  tmp
  # cat tmp
  Tue Sep 11 04:11:11 UTC 2018
  ```

- 块存储与原 PVC 解绑
    - 将 PV 的 spec.persistentVolumeReclaimPolicy 的值从 Delete 改为 Retain
    ```
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
    ```
    $ kubectl delete pvc old-pvc
    persistentvolumeclaim "old-pvc" deleted
    $ kubectl delete pv pvc-55754c8c-b577-11e8-a480-525445c0b555
    persistentvolume "pvc-55754c8c-b577-11e8-a480-525445c0b555" deleted
    ```

- 静态创建存储卷绑定块存储

    - 编辑 PV 资源定义文件
    ```
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
    ```
    $ kubectl create -f pv.yaml

    ```

    - 编辑 PVC 资源定义文件
    ```
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
    ```
    i-lcsolq8c# kubectl get pv new-pv
    NAME      CAPACITY   ACCESS MODES   RECLAIM POLICY   STATUS    CLAIM             STORAGECLASS    REASON    AGE
    new-pv    10Gi       RWO            Delete           Bound     default/new-pvc   csi-qingcloud             44s
    i-lcsolq8c# kubectl get pvc new-pvc
    NAME      STATUS    VOLUME    CAPACITY   ACCESS MODES   STORAGECLASS    AGE
    new-pvc   Bound     new-pv    10Gi       RWO            csi-qingcloud   32s
    ```

- 创建 Deployment 查看PVC里内容

```
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
