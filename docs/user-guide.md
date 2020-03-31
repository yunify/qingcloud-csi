<!-- START doctoc generated TOC please keep comment here to allow auto update -->
<!-- DON'T EDIT THIS SECTION, INSTEAD RE-RUN doctoc TO UPDATE -->
- [User Guide](#user-guide)
  - [Set Storage Class](#set-storage-class)
    - [An Example of Storage Class](#an-example-of-storage-class)
    - [Parameters in Storage Class](#parameters-in-storage-class)
      - [type, maxSize, minSize, stepSize](#type-maxsize-minsize-stepsize)
      - [fsType](#fstype)
      - [replica](#replica)
      - [tags](#tags)
    - [Other Parameters](#other-parameters)
      - [Set Default Storage Class](#set-default-storage-class)
      - [Expand Volume](#expand-volume)
      - [Topology Awareness](#topology-awareness)
    - [Disk Type Matrix](#disk-type-matrix)
    - [Instance Type Matrix](#instance-type-matrix)
    - [Disk Compatiblity Matrix](#disk-compatiblity-matrix)
  - [Volume Management](#volume-management)
    - [Prerequisite](#prerequisite)
      - [Create Storage Class](#create-storage-class)
    - [Create Volume](#create-volume)
    - [Mount Volume](#mount-volume)
    - [Unmount Volume](#unmount-volume)
    - [Delete Volume](#delete-volume)
  - [Expand Volume](#expand-volume-1)
    - [Prerequisite](#prerequisite-1)
    - [Unmount Volume](#unmount-volume-1)
    - [Expand Volume](#expand-volume-2)
    - [Check](#check)
  - [Clone Volume](#clone-volume)
    - [Prerequisite](#prerequisite-2)
    - [Cloning](#cloning)
  - [Snapshot Management](#snapshot-management)
    - [Prerequisite](#prerequisite-3)
      - [Create a volume pre-populated data](#create-a-volume-pre-populated-data)
    - [Create Snapshot](#create-snapshot)
      - [Create Volume Snapshot Class](#create-volume-snapshot-class)
      - [Craete Volume Snapshot](#craete-volume-snapshot)
    - [Restore Volume from Snapshot](#restore-volume-from-snapshot)
      - [Restore](#restore)
      - [Check Data](#check-data)
    - [Delete Snapshot](#delete-snapshot)
  - [Topology Awareness](#topology-awareness-1)
    - [Prerequisite](#prerequisite-4)
      - [Exmaple Kubernetes Cluster](#exmaple-kubernetes-cluster)
      - [Create Storage Class](#create-storage-class-1)
    - [Create Volume](#create-volume-1)
    - [Create Workload](#create-workload)
  - [Static Volume Provisioning](#static-volume-provisioning)
    - [Process](#process)
      - [Prerequisite](#prerequisite-5)
      - [Create Storage Class](#create-storage-class-2)
      - [Create PV](#create-pv)
      - [Create PVC](#create-pvc)

<!-- END doctoc generated TOC please keep comment here to allow auto update -->

# User Guide

## Set Storage Class
### An Example of Storage Class

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

### Parameters in Storage Class

#### type, maxSize, minSize, stepSize
See details in [QingCloud docs](https://docs.qingcloud.com/product/api/action/volume/create_volumes.html)。

|Disk|type|maxSize|minSize|stepSize|
|:---:|:---:|:---:|:---:|:---:|
|High Performance|0|2000|10|10|
|High Capacity|2|5000|100|50|
|Super High Performance|3|2000|10|10|
| NeonSAN|5|50000|100|100|
|NeonSAN HDD|5|50000|100|100|
| Standard|100|2000|10|10|
| SSD Enterprise|200| 2000|10|10|

#### fsType
Support `ext3`, `ext4`, `xfs`. Default is `ext4`.

#### replica
`1` represents single duplication disk，`2` represents multiple duplication disk. Default is `2`.

#### tags
The ID of QingCloud Tag resource, split by a single comma. Disks and snapshots created by this plugin will be attached with the specified tags.

### Other Parameters

#### Set Default Storage Class
In annotation, please set the value of `.metadata.annotations.storageclass.beta.kubernetes.io/is-default-class` as `true`. See details in [Kubernetes docs](https://kubernetes.io/docs/tasks/administer-cluster/change-default-storage-class/)

#### Expand Volume
Set the value of `.allowVolumeExpansion` as `true`. See details in [Kubernetes docs](https://kubernetes.io/docs/concepts/storage/storage-classes/#allow-volume-expansion)

#### Topology Awareness
We can set `Immediate` or `WaitForFirstConsumer` as the value of `.volumeBindingMode`. See details in [Kubernetes docs](https://kubernetes.io/docs/concepts/storage/storage-classes/#volume-binding-mode)

### Disk Type Matrix

|Disk|type|
|:----:|:----:|
| High Performance|0|
| High Capacity|2|
|Super High Performance|3|
| NeonSAN|5|
|NeonSAN HDD|6|
| Standard|100|
| SSD Enterprise|200|

### Instance Type Matrix
|Instance|type|
|:----:|:----:|
|High Performance|0|
|Super High Performance|1|
|Super High Performance SAN|6|
|High Performance SAN|7|
|Standard|101|
|Enterprise1|201|
|Enterprise2|202|
|Premium|301|

### Disk Compatiblity Matrix

 |          | High Performance Disk    | High Capacity Disk  | Super High Performance Disk | NeonSAN Disk | NeonSAN HDD Disk | Standard Disk| SSD Enterprise Disk|
|-----------|------------------|------------------|-----------------|---------|----------|-------|-------|
|High Performance Instance| ✓        | ✓                | -               | ✓      |  ✓      | ✓     | -     |
|Super High Performance Instance| -       | ✓                | ✓               |✓      |✓     |-  |✓  |
|Super High Performance SAN Instance| -       | -                | -              |✓      |-     |-  |-  |
|High Performance SAN Instance| -       | -                | -               |-     |✓     |-  |-  |
|Standard Instance| ✓        | ✓                | -               |✓  |✓ |✓ |-  |
|Enterprise1 Instance| -       | ✓                | ✓               |✓  |✓ |-  |✓  |
|Enterprise2 Instance| -       | ✓                | ✓               |✓  |✓ |-  |✓  |
|Premium Instance| -       | ✓                | ✓               |✓  |✓ |-  |✓  |

## Volume Management
Volume management including dynamical provisioning/deleting volume, attaching/detaching volume. Please reference [Example YAML Files](https://github.com/yunify/qingcloud-csi/tree/master/deploy/disk/example/volume)。

### Prerequisite
- Kubernetes 1.14+ Cluster
- Installed QingCloud CSI plugin
- Created QingCloud CSI storage class

#### Create Storage Class
- Create
```console
$ kubectl create -f sc.yaml
```
- Check
```console
$ kubectl get sc
NAME            PROVISIONER              AGE
csi-qingcloud   disk.csi.qingcloud.com   14m
```

### Create Volume
- Create
```console
$ kubectl create -f pvc.yaml 
persistentvolumeclaim/pvc-example created
```
- Check
```console
$ kubectl get pvc pvc-example
NAME          STATUS   VOLUME                                     CAPACITY   ACCESS MODES   STORAGECLASS    AGE
pvc-example   Bound    pvc-76429525-a930-11e9-9a6a-5254ef68c8c1   20Gi       RWO            csi-qingcloud   25m
```

### Mount Volume
- Create Deployment
```console
$ kubectl create -f deploy-nginx.yaml 
deployment.apps/deploy-nginx created
```
- Check
```console
$ kubectl exec -ti deploy-nginx-84474cf674-zfhbs /bin/bash
# cd /mnt
# ls
lost+found
```

### Unmount Volume
- Delete Deployment
```console
$ kubectl delete deploy deploy-nginx
deployment.extensions "deploy-nginx" deleted
```

### Delete Volume
- Delete
```console
$ kubectl delete pvc pvc-example
persistentvolumeclaim "pvc-example" deleted
```
- Check
```console
$ kubectl get pvc pvc-example
Error from server (NotFound): persistentvolumeclaims "pvc-example" not found
```

## Expand Volume
This feature could expand the capacity of volume. This plugin only supports offline volume expansion. The procedure of offline volume expansion is shown as follows. 
1. Ensure volume in unmounted status
2. Edit the capacity of PVC
3. Mount volume on workload
Please reference [Example YAML files](https://github.com/yunify/qingcloud-csi/tree/master/deploy/disk/volume)。

### Prerequisite
- Kubernetes 1.14+ cluster
- Add `ExpandCSIVolumes=true` in `feature-gate` 
- Set `allowVolumeExpansion` as `true` in storage class
- Create a Pod mounting a volume

### Unmount Volume
```console
$ kubectl scale deploy deploy-nginx --replicas=0
```

### Expand Volume
- Change volume capacity
```console
$ kubectl patch pvc pvc-example -p '{"spec":{"resources":{"requests":{"storage": "40Gi"}}}}'
persistentvolumeclaim/pvc-example patched
```
- Mount volume
```console
$ kubectl scale deploy deploy-nginx --replicas=1
```
- Check volume capacity
```console
$ kubectl get pvc pvc-example
NAME          STATUS   VOLUME                                     CAPACITY   ACCESS MODES   STORAGECLASS    AGE
pvc-example   Bound    pvc-906f5760-a935-11e9-9a6a-5254ef68c8c1   40Gi       RWO            csi-qingcloud   6m7s
$ kubectl get po
NAME                            READY   STATUS    RESTARTS   AGE
deploy-nginx-6c444c9b7f-d6n29   1/1     Running   0          3m38s
```

### Check
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

## Clone Volume
A Clone is defined as a duplicate of an existing Kubernetes Volume. Please reference [Example YAML files](https://github.com/yunify/qingcloud-csi/tree/master/deploy/disk/example/volume/pvc-clone.yaml).

### Prerequisite
- Kubernetes 1.15+ cluster
- Enable `VolumePVCDataSource=true` feature gate
- Install QingCloud CSI plugin
- Create QingCloud CSI storage class
- Create a volume

### Cloning
- Find volume
```console
$ kubectl get pvc
NAME          STATUS   VOLUME                                     CAPACITY   ACCESS MODES   STORAGECLASS    AGE
pvc-example   Bound    pvc-d1fb263e-b368-4339-8f8b-448446f4b840   20Gi       RWO            csi-qingcloud   32s
```

- Clone volume
```console
$ kubectl create -f pvc-clone.yaml
persistentvolumeclaim/pvc-clone created
```

- Check
```console
$ kubectl get pvc pvc-clone
NAME        STATUS   VOLUME                                     CAPACITY   ACCESS MODES   STORAGECLASS    AGE
pvc-clone   Bound    pvc-529d2502-02bd-442b-a69f-d3eff28316a8   20Gi       RWO            csi-qingcloud   31s
```

## Snapshot Management
Snapshot management contains creating/deleting snapshot and restoring volume from snapshpot. Please reference [Example YAML files](https://github.com/yunify/qingcloud-csi/tree/master/deploy/disk/example/snapshot).

### Prerequisite
- Kubernetes 1.14+ cluster
- Enable `VolumeSnapshotDataSource=true` feature gate at kube-apiserver and kube-controller-manager
- Install QingCloud CSI plugin
- Create QingCloud CSI storage class
- Create a volume

#### Create a volume pre-populated data
- Create volume
```console
$ kubectl create -f original-pvc.yaml
persistentvolumeclaim/pvc-snap-1 created
```
- Check
```console
$ kubectl get pvc
NAME         STATUS   VOLUME                                     CAPACITY   ACCESS MODES   STORAGECLASS    AGE
pvc-snap-1   Bound    pvc-28090960-9eeb-11e9-9a6a-5254ef68c8c1   20Gi       RWO            csi-qingcloud   37s
```
- Write data
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

### Create Snapshot

#### Create Volume Snapshot Class
```console
$ kubectl create -f snapshot-class.yaml 
volumesnapshotclass.snapshot.storage.k8s.io/csi-qingcloud created

$ kubectl get volumesnapshotclass
NAME            AGE
csi-qingcloud   16s
```

#### Craete Volume Snapshot
```console
$ kubectl create -f volume-snapshot.yaml 
volumesnapshot.snapshot.storage.k8s.io/snap-1 created

$ kubectl get volumesnapshot
NAME     AGE
snap-1   91s
```

### Restore Volume from Snapshot
#### Restore
```console
$ kubectl create -f restore-pvc.yaml 
persistentvolumeclaim/pvc-snap-2 created
```

```console
$ kubectl get pvc pvc-snap-2
NAME         STATUS   VOLUME                                     CAPACITY   ACCESS MODES   STORAGECLASS    AGE
pvc-snap-2   Bound    pvc-b8a05427-9eef-11e9-9a6a-5254ef68c8c1   20Gi       RWO            csi-qingcloud   52s
```

#### Check Data
Compare the difference between restored volume with original volume.

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

### Delete Snapshot

```console
$ kubectl delete volumesnapshot snap-1
volumesnapshot.snapshot.storage.k8s.io "snap-1" deleted
```

## Topology Awareness
Topology awareness is used at Kubernetes clusters whose nodes across different available zones or having different types of instance.  Please reference [Example YAML files](https://github.com/yunify/qingcloud-csi/tree/master/deploy/disk/example/topology).

### Prerequisite
- Kubernetes 1.14+ cluster
- Enable `CSINodeInfo=true` feature gate at Kubernetes control plane and Kubelet
- Install QingCloud CSI plugin and enable `Topology=true` feature gate at `external-provisioner` sidecar container
- Set QingCloud CSI storage class

#### Exmaple Kubernetes Cluster
In QingCloud Pek3 zone, a Kubernetes v1.15 cluster with same types of instance is created and node1 and node2 running in Pek3c, node3 and node4 running in Pek3b, node5 and node6 running in Pek3d.

#### Create Storage Class
- Volume binding mode can be set as `WaitForFirstConsumer` or `Immediate`. Please reference [Kubernetes docs](https://kubernetes.io/docs/concepts/storage/storage-classes/#volume-binding-mode).

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

### Create Volume
- Create
```console
$ kubectl create -f pvc.yaml
persistentvolumeclaim/pvc-topology created
```

- If `VolumeBindingMode` set as `WaitForFirstConsumer`, the status of PVC is shown as Pending. After Pod mounted PVC is sheduled, the PVC status will change to Bound.
```console
$ kubectl get pvc pvc-topology
NAME           STATUS    VOLUME   CAPACITY   ACCESS MODES   STORAGECLASS    AGE
pvc-topology   Pending                                      csi-qingcloud   31s
```

### Create Workload
- Create Deployment
```console
$ kubectl create -f deploy.yaml
deployment.apps/nginx-topology created
```

- Check if Pods scheduled on node 3
```console
$ kubectl get po -o wide
NAME                      READY   STATUS    RESTARTS   AGE   IP               NODE    NOMINATED NODE   READINESS GATES
nginx-topology-79d8d5d86d-4lvcl    1/1     Running   0          52s   10.233.92.27     node3   <none>           <none>
```

- Check volume bound
```console
$ kubectl get pvc pvc-topology
NAME           STATUS   VOLUME                                     CAPACITY   ACCESS MODES   STORAGECLASS    AGE
pvc-topology   Bound    pvc-5b34120c-6119-4c86-b9de-e152304683e6   20Gi       RWO            csi-qingcloud   2m48s
```

- The volume named pvc-topology contains topology information and can be mounted on special nodes. In this example, the volume can be mounted on node3 or node4.

## Static Volume Provisioning

Static volume provisioning is also called pre-provisioning volume. The process is shown below.
1. Create QingCloud disk
2. Create PV
3. Create PVC

### Process

#### Prerequisite

- Create Kubernetes cluster on QingCloud IaaS platform
- Install QingCloud CSI plugin
- Create QingCloud disk

#### Create Storage Class

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

#### Create PV

- Edit
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

- Create
```console
$ kubectl create -f pv.yaml
```

- Check
```console
$ kubectl get pv pv-static
NAME        CAPACITY   ACCESS MODES   RECLAIM POLICY   STATUS      CLAIM     STORAGECLASS    REASON    AGE
pv-static   20Gi       RWO            Delete           Available             csi-qingcloud             8m
```

#### Create PVC

- Edit

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

- Create

```console
$ kubectl create -f pvc.yaml
```

- Check

```console
$  kubectl get pvc pvc-static
NAME         STATUS    VOLUME      CAPACITY   ACCESS MODES   STORAGECLASS    AGE
pvc-static   Bound     pv-static   20Gi       RWO            csi-qingcloud   11s

$ kubectl get pv pv-static 
NAME        CAPACITY   ACCESS MODES   RECLAIM POLICY   STATUS    CLAIM                STORAGECLASS    REASON    AGE
pv-static   20Gi       RWO            Delete           Bound     default/pvc-static   csi-qingcloud             12m
```
