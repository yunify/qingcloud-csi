# 存储卷管理
存储卷（PVC，PersistentVolumeClaim）管理功能包括动态分配存储卷，删除存储卷，挂载存储卷到 Pod，从 Pod 卸载存储卷。用户可参考[示例 YAML 文件](https://github.com/yunify/qingcloud-csi/tree/master/deploy/disk/example/volume)。

## 准备工作
- Kubernetes 1.14+ 集群
- 安装了 QingCloud CSI 存储插件
- 安装了 QingCloud CSI 存储类型

### 安装 QingCloud CSI 存储类型
- 安装
```
$ kubectl create -f sc.yaml
```
- 检查
```
$ kubectl get sc
NAME            PROVISIONER              AGE
csi-qingcloud   disk.csi.qingcloud.com   14m
```

## 创建存储卷
- 创建存储卷
```
$ kubectl create -f pvc.yaml 
persistentvolumeclaim/pvc-example created
```
- 检查存储卷
```
$ kubectl get pvc pvc-example
NAME          STATUS   VOLUME                                     CAPACITY   ACCESS MODES   STORAGECLASS    AGE
pvc-example   Bound    pvc-76429525-a930-11e9-9a6a-5254ef68c8c1   20Gi       RWO            csi-qingcloud   25m
```

## 挂载存储卷
- 创建 Deployment 挂载存储卷
```
$ kubectl create -f deploy-nginx.yaml 
deployment.apps/deploy-nginx created
```
- 访问容器内挂载存储卷的目录
```
$ kubectl exec -ti deploy-nginx-84474cf674-zfhbs /bin/bash
# cd /mnt
# ls
lost+found
```

## 卸载存储卷
- 删除 deployment 卸载存储卷
```
$ kubectl delete deploy deploy-nginx
deployment.extensions "deploy-nginx" deleted
```

## 删除存储卷
- 删除存储卷
```
$ kubectl delete pvc pvc-example
persistentvolumeclaim "pvc-example" deleted
```
- 检查存储卷
```
$ kubectl get pvc pvc-example
Error from server (NotFound): persistentvolumeclaims "pvc-example" not found
```