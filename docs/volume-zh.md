# 存储卷管理
存储卷（PVC，PersistentVolumeClaim）管理功能包括动态分配存储卷，删除存储卷，挂载存储卷到 Pod，从 Pod 卸载存储卷。用户可参考[示例 YAML 文件](https://github.com/yunify/qingcloud-csi/tree/master/deploy/disk/example/volume)。

## 准备工作
- Kubernetes 1.14+集群
- 安装了 QingCloud CSI 存储插件
- 安装了 QingCloud CSI storageclass

### 安装 QingCloud CSI storageclass
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
persistentvolumeclaim/pvc-test created
```
- 检查存储卷
```
$ kubectl get pvc pvc-test
NAME       STATUS   VOLUME                                     CAPACITY   ACCESS MODES   STORAGECLASS    AGE
pvc-test   Bound    pvc-76429525-a930-11e9-9a6a-5254ef68c8c1   20Gi       RWO            csi-qingcloud   25m
```

## 挂载存储卷
- 创建 deployment 挂载存储卷
```
$ kubectl create -f deploy.yaml 
deployment.apps/nginx created
```
- 访问容器内挂载存储卷的目录
```
$ kubectl exec -ti deploy-nginx-qingcloud-84474cf674-zfhbs /bin/bash
# cd /mnt
# ls
lost+found
```

## 卸载存储卷
- 删除 deployment 卸载存储卷
```
$ kubectl delete deploy nginx
deployment.extensions "nginx" deleted
```

## 删除存储卷
- 删除存储卷
```
$ kubectl delete pvc pvc-test
persistentvolumeclaim "pvc-test" deleted
```
- 检查存储卷
```
$ kubectl get pvc pvc-test
Error from server (NotFound): persistentvolumeclaims "pvc-test" not found
```