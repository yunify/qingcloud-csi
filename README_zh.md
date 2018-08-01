# QingCloud-CSI

[![Build Status](https://travis-ci.org/yunify/qingcloud-csi.svg?branch=master)](https://travis-ci.org/yunify/qingcloud-csi)
[![Go Report Card](https://goreportcard.com/badge/github.com/yunify/qingcloud-csi)](https://goreportcard.com/report/github.com/yunify/qingcloud-csi)

> [English](README.md) | 中文
## 描述
QingCloud CSI插件实现了[CSI](https://github.com/container-storage-interface/)接口，并使容器编排平台能够使用QingCloud云平台的存储资源。目前，QingCloud CSI插件已经在Kubernetes v1.10环境中通过了[CSI测试](https://github.com/kubernetes-csi/csi-test)。

## 块存储插件

插件的设计和安装使用Kubernetes社区推荐的CSI插件[架构](https://github.com/kubernetes/community/blob/master/contributors/design-proposals/storage/container-storage-interface.md#recommended-mechanism-for-deploying-csi-drivers-on-kubernetes)，插件架构包含Controller和Node两部分，在Controller部分，由有状态副本集（StatefulSet）在Kubernetes集群内创建一个Pod副本。在Node部分，每个可调度的节点由守护进程集（DaemonSet）创建一个Pod副本。

块存储插件部署后, 用户可创建访问模式（Access Mode）为单节点读写（ReadWriteOnce）的基于QingCloud的超高性能型，性能型或容量型硬盘的存储卷并挂载至工作负载。

### 编译
QingCloud CSI插件可编译为二进制文件或镜像。编译后的二进制文件存放在_output文件夹内。当编译为镜像时，镜像存储在本地的Docker镜像仓库内。

编译为二进制文件:
```
$ make blockplugin
```

构建Docker镜像:
```
$ make blockplugin-container
```

本地镜像仓库存储镜像：
```
$ docker images | grep csi-qingcloud
dockerhub.qingcloud.com/csiplugin/csi-qingcloud		v0.2.0	  c75dc27cbfd7		55 minutes ago		40MB
```

### 配置
#### 配置文件

如下所示的[配置文件](deploy/block/kubernetes/config.yaml)将会被ConfigMap所使用。
> 注: 在QingCloud AppCenter内, 请修改创建ConfigMap的[脚本](deploy/block/kubernetes/create-cm.sh)并创建引用存放在主机内的配置文件(*/etc/qingcloud/client.yaml*)的ConfigMap。 

```
qy_access_key_id: 'ACCESS_KEY_ID'
qy_secret_access_key: 'ACCESS_KEY_SECRET'
zone: 'ZONE'
host: 'api.qingcloud.com'
port: 443
protocol: 'https'
uri: '/iaas'
connection_retries: 3
connection_timeout: 30
```

- `qy_access_key_id`, `qy_secret_access_key`: 在QingCloud控制台创建Access key密钥. 此密钥需要有操作QingCloud IaaS平台资源的权限。

- `zone`: Zone字段应与Kubernetes集群所在Zone相同。CSI插件将会操作此Zone内的存储卷资源。

- `host`, `prot`. `protocol`, `uri`: 共同构成QingCloud IaaS平台服务的url.

### StorageClass

如下所示的StorageClass资源定义[文件](deploy/block/example/sc.yaml)可用来创建StorageClass对象.
```
apiVersion: storage.k8s.io/v1
kind: StorageClass
metadata:
  name: csi-qingcloud
provisioner: csi-qingcloud
parameters:
  type: "0"
  maxSize: "500"
  minSize: "10"
  stepSize: "10"
  fsType: "ext4"
reclaimPolicy: Delete 
```

- `type`: QingCloud云平台存储卷类型。总体上， `0`代表性能型硬盘。`3`代表超高性能型硬盘。`1`或`2`（根据Zone不同而参数不同）代表容量型硬盘。 详情见[QingCloud文档](https://docs.qingcloud.com/product/api/action/volume/create_volumes.html)。

- `maxSize`, `minSize`: 某种存储卷类型的存储卷容量范围。

- `stepSize`: 步长用来控制所创建存储卷的容量。

- `fsType`: 支持`ext3`, `ext4`, `xfs`. 默认为`ext4`.

### 安装
此安装指南将CSI插件部署在*kube-system* namespace内。用户也可以将插件部署在其他namespace内。Kubernetes控制平面内请勿禁用[Mount Propagation](https://kubernetes.io/docs/concepts/storage/volumes/#mount-propagation)特性。

- 创建ConfigMap
```
$ chmod +x deploy/block/kubernetes/create-cm.sh
$ ./create-cm.sh
```

- 创建Docker镜像仓库密钥
```
kubectl create -f deploy/block/kubernetes/csi-secret.yaml
```

- 创建访问控制相关对象
```
$ kubectl create -f deploy/block/kubernetes/csi-controller-rbac.yaml
$ kubectl create -f deploy/block/kubernetes/csi-node-rbac.yaml
```

- 部署CSI插件
> 注: 在QingCloud AppCenter内, 请将[DaemonSet](deploy/block/kubernetes/csi-node-ds.yaml) YAML文件的 *"/var/lib/kubelet"* 替换为 *"/data/var/lib/kubelet"*。

```
$ kubectl create -f deploy/block/kubernetes/csi-controller-sts.yaml
$ kubectl create -f deploy/block/kubernetes/csi-node-ds.yaml
```

- 检查CSI插件状态
```
$ kubectl get pods -n csi-qingcloud | grep csi
csi-qingcloud-controller-0      3/3       Running       0          5m
csi-qingcloud-node-kks3q        2/2       Running       0          2m
csi-qingcloud-node-pgsbn        2/2       Running       0          2m
```

### 验证
- 由Kubernetes集群管理员创建StorageClass
```
$ kubectl create -f deploy/block/example/sc.yaml
```

- 创建PVC
```
$ kubectl create -f deploy/block/example/pvc.yaml
```

- 创建挂载PVC的Deployment
```
$ kubectl create -f deploy/block/example/deploy.yaml
```

- 检查Pod状态
```
$ kubectl get po | grep deploy
nginx-84474cf674-zfhbs   1/1       Running   0          1m
```

- 访问容器内挂载存储卷的目录
```
$ kubectl exec -ti deploy-nginx-qingcloud-84474cf674-zfhbs /bin/bash
# cd /mnt
# ls
lost+found
```

## 支持
如果有任何问题或建议, 请提在[qingcloud-csi](https://github.com/yunify/qingcloud-csi/issues)项目提issue。
