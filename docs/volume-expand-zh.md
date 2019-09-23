# 扩容
扩容功能将扩大存储卷可用容量。由于云平台限制，本存储插件仅支持离线扩容硬盘。离线扩容硬盘流程是 1. 存储卷处于未挂载状态，2. 扩容存储卷，3. 挂载一次存储卷。示例 YAML 文件在 https://github.com/yunify/qingcloud-csi/tree/master/deploy/disk/volume 内。

## 准备工作
- Kubernetes 1.14+ 集群
- Kubernetes 组件的 `feature-gate` 增加 `ExpandCSIVolumes=true`
- 配置 QingCloud CSI 存储类型，并将其 `allowVolumeExpansion` 字段值设置为 `true`
- 创建一个存储卷并挂载至 Pod，参考[存储卷管理](volume-zh.md)

## 卸载存储卷
```
$ kubectl scale deploy deploy-nginx --replicas=0
```

## 扩容存储卷
- 修改存储卷容量
```
$ kubectl patch pvc pvc-example -p '{"spec":{"resources":{"requests":{"storage": "40Gi"}}}}'
persistentvolumeclaim/pvc-example patched
```
- 挂载存储卷
```
$ kubectl scale deploy deploy-nginx --replicas=1
```
- 完成扩容
```
$ kubectl get pvc pvc-example
NAME          STATUS   VOLUME                                     CAPACITY   ACCESS MODES   STORAGECLASS    AGE
pvc-example   Bound    pvc-906f5760-a935-11e9-9a6a-5254ef68c8c1   40Gi       RWO            csi-qingcloud   6m7s
$ kubectl get po
NAME                            READY   STATUS    RESTARTS   AGE
deploy-nginx-6c444c9b7f-d6n29   1/1     Running   0          3m38s
```

## 检查
- 进入 Pod 查看
```
$ kubectl exec -ti deploy-nginx-6c444c9b7f-d6n29 /bin/bash
root@deploy-nginx-6c444c9b7f-d6n29:/# s
bash: s: command not found
root@deploy-nginx-6c444c9b7f-d6n29:/# df -ah
Filesystem      Size  Used Avail Use% Mounted on
...
/dev/vdc         40G   49M   40G   1% /mnt
...
```