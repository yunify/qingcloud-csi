# 快照
快照功能包括创建和删除快照，从快照恢复存储卷功能。示例 YAML 文件在 https://github.com/yunify/qingcloud-csi/tree/master/deploy/disk/example/snapshot 内。

## 准备工作
- Kubernetes 1.14+ 集群
- 在 apiserver, controller-manager 的 `feature-gate` 增加 `VolumeSnapshotDataSource=true`
- 安装 QingCloud CSI 存储插件
- 配置了 QingCloud CSI storageclass
- 创建一个带数据的存储卷

### 创建带数据的存储卷 `pvc-snap-1`
- 创建存储卷 
```
$ kubectl create -f pvc.yaml 
persistentvolumeclaim/pvc-snap-1 created
```
- 检查存储卷
```
$ kubectl get pvc
NAME         STATUS   VOLUME                                     CAPACITY   ACCESS MODES   STORAGECLASS    AGE
pvc-snap-1   Bound    pvc-28090960-9eeb-11e9-9a6a-5254ef68c8c1   20Gi       RWO            csi-qingcloud   37s
```
- 向存储卷写数据
```
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

## 创建快照
每个 Kubernetes 快照对应于一个 QingCloud 全量备份，请确保有足够全量备份链配额。

### 创建快照类型
```
$ kubectl create -f snapshot-class.yaml 
volumesnapshotclass.snapshot.storage.k8s.io/csi-qingcloud created

$ kubectl get volumesnapshotclass
NAME            AGE
csi-qingcloud   16s
```

### 创建快照
```
$ kubectl create -f volume-snapshot.yaml 
volumesnapshot.snapshot.storage.k8s.io/snap-1 created

$ kubectl get volumesnapshot
NAME     AGE
snap-1   91s
```

## 从快照恢复存储卷
### 恢复存储卷 `pvc-snap-2`
```
$ kubectl create -f restore-pvc.yaml 
persistentvolumeclaim/pvc-snap-2 created
```

```
$ kubectl get pvc pvc-snap-2
NAME         STATUS   VOLUME                                     CAPACITY   ACCESS MODES   STORAGECLASS    AGE
pvc-snap-2   Bound    pvc-b8a05427-9eef-11e9-9a6a-5254ef68c8c1   20Gi       RWO            csi-qingcloud   52s
```

### 检查存储卷数据
从快照恢复的存储卷 `pvc-snap-2` 与在创建快照时的 `pvc-snap-1` 内容应一致。

```
$ kubectl create -f deploy-viewer.yaml 
deployment.apps/nginx created

$ kubectl get po |grep nginx
nginx-7b98f8c4d4-fmjzf   1/1     Running   0          3m6s

$ kubectl exec -ti nginx-7b98f8c4d4-fmjzf /bin/bash
root@nginx-7b98f8c4d4-fmjzf:/# ls /mnt -lh
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

## 删除快照

```
$ kubectl delete volumesnapshot snap-1
volumesnapshot.snapshot.storage.k8s.io "snap-1" deleted
```