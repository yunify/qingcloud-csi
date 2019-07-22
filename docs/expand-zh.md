# 扩容
扩容功能将扩大存储卷可用容量。由于云平台限制，本存储插件仅支持离线扩容硬盘。示例 YAML 文件在 https://github.com/yunify/qingcloud-csi/tree/master/deploy/disk/example/volume 内。

## 准备工作
- Kubernetes 1.14+ 集群
- Kubernetes 组件的 `feature-gate` 增加 `ExpandCSIVolumes=true`
- 配置了 QingCloud CSI storageclass，并将其 `allowVolumeExpansion` 字段值设置为 `true`
- 创建一个存储卷并挂载至 Pod，参考[存储卷管理](volume-zh.md)

## 卸载存储卷
```
$ kubectl scale deploy nginx --replicas=0
```

## 扩容存储卷
- 修改存储卷容量
```
$ kubectl patch pvc pvc-test -p '{"spec":{"resources":{"requests":{"storage": "40Gi"}}}}'
persistentvolumeclaim/pvc-test patched
```
- 挂载存储卷
```
$ kubectl scale deploy nginx --replicas=1
```
- 完成扩容
```
$ kubectl get pvc pvc-test
NAME       STATUS   VOLUME                                     CAPACITY   ACCESS MODES   STORAGECLASS    AGE
pvc-test   Bound    pvc-906f5760-a935-11e9-9a6a-5254ef68c8c1   40Gi       RWO            csi-qingcloud   6m7s
$ kubectl get po
NAME                     READY   STATUS    RESTARTS   AGE
nginx-6c444c9b7f-d6n29   1/1     Running   0          3m38s
```

## 检查
- 进入 Pod 查看
```
$ kubectl exec -ti nginx-6c444c9b7f-d6n29 /bin/bash
root@nginx-6c444c9b7f-d6n29:/# s
bash: s: command not found
root@nginx-6c444c9b7f-d6n29:/# df -ah
Filesystem      Size  Used Avail Use% Mounted on
...
/dev/vdc         40G   49M   40G   1% /mnt
...
```