# QingCloud-CSI

[![Build Status](https://travis-ci.org/yunify/qingcloud-csi.svg?branch=master)](https://travis-ci.org/yunify/qingcloud-csi)
[![Go Report Card](https://goreportcard.com/badge/github.com/yunify/qingcloud-csi)](https://goreportcard.com/report/github.com/yunify/qingcloud-csi)

> [English](README.md) | 中文
## 描述
QingCloud CSI 插件实现了 [CSI](https://github.com/container-storage-interface/) 接口，并使容器编排平台能够使用 QingCloud 云平台的存储资源。目前，QingCloud CSI 插件已经在 Kubernetes v1.14 环境中通过了 [CSI 测试](https://github.com/kubernetes-csi/csi-test)。

## 块存储插件

插件的设计和安装使用 Kubernetes 社区推荐的 CSI 插件[架构](https://github.com/kubernetes/community/blob/master/contributors/design-proposals/storage/container-storage-interface.md#recommended-mechanism-for-deploying-csi-drivers-on-kubernetes)，插件架构包含 Controller 和 Node 两部分，在 Controller 部分，由有状态副本集（StatefulSet）在 Kubernetes 集群内创建一个 Pod 副本。在 Node 部分，每个可调度的节点由守护进程集（DaemonSet）创建一个 Pod 副本。

块存储插件部署后, 用户可创建访问模式（Access Mode）为单节点读写（ReadWriteOnce）的基于 QingCloud 的企业型，基础型，NeonSAN，超高性能型，性能型或容量型硬盘的存储卷并挂载至工作负载。

### 安装
此安装指南将 CSI 插件安装在 *kube-system* namespace 内。用户也可以将插件部署在其他 namespace 内。为了 CSI 插件的正常使用，请确保在 Kubernetes 控制平面内将 `--allow-privileged` 项设置为 `true` 并且启用（默认开启）[Mount Propagation](https://kubernetes.io/docs/concepts/storage/volumes/#mount-propagation) 特性。

- 设置 Kubernetes 参数
  - 设置 `--allow-privileged=true`。
  - 启用（默认开启）[Mount Propagation](https://kubernetes.io/docs/concepts/storage/volumes/#mount-propagation) 特性。
  - 设置 `--feature-gates=CSINodeInfo=true,CSIDriverRegistry=true,KubeletPluginsWatcher=true`

- 下载安装包并解压
```
$ wget $(curl --silent "https://api.github.com/repos/yunify/qingcloud-csi/releases/latest" | \
  grep browser_download_url | grep install|cut -d '"' -f 4)
$ tar -xvf csi-qingcloud-install.tar.gz
$ cd csi-qingcloud-install
```

- 创建 ConfigMap
  * 在基于 QingCloud IaaS 平台的 Kubernetes 集群内
    1. 修改安装包内配置文件（config.yaml）
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
    - `qy_access_key_id`, `qy_secret_access_key`: 在 QingCloud 控制台创建 Access key 密钥. 此密钥需要有操作 QingCloud IaaS 平台资源的权限。

    - `zone`: `zone` 应与 Kubernetes 集群所在区相同。CSI 插件将会操作此区的存储卷资源。例如：`zone` 可以设置为 `sh1a` 和 `ap2a`。
    
    - `host`, `port`. `protocol`, `uri`: 共同构成 QingCloud IaaS 平台服务的 url。

    2. 创建 ConfigMap
    ```
    $ kubectl create configmap csi-qingcloud --from-file=config.yaml=./config.yaml --namespace=kube-system
    ```

  * 在基于 QingCloud Appcenter 的 Kubernetes 集群内

    1. 创建 ConfigMap
    ```
    $ kubectl create configmap csi-qingcloud --from-file=config.yaml=/etc/qingcloud/client.yaml --namespace=kube-system
    ```

- 创建 Docker 镜像仓库密钥
```
$ kubectl apply -f ./csi-secret.yaml
```

- 创建访问控制相关对象
```
$ kubectl apply -f ./csi-controller-rbac.yaml
$ kubectl apply -f ./csi-node-rbac.yaml
```

- 创建 CSIdriver 对象
```
$ kubectl apply -f ./csi-driver.yaml
```

- 部署 CSI 插件
> 注:  如果 Kubernetes 集群的 [kubelet](https://kubernetes.io/docs/reference/command-line-tools-reference/kubelet/) 设置了 `--root-dir` 选项（默认值为 *"/var/lib/kubelet"*），请将 [DaemonSet](deploy/block/kubernetes/csi-node-ds.yaml) YAML 文件 `spec.template.spec.containers[name=csi-qingcloud].volumeMounts[name=mount-dir].mountPath` 和 `spec.template.spec.volumes[name=mount-dir].hostPath.path` 的值 *"/var/lib/kubelet"* 替换为 `--root-dir` 选项的值。例如：在通过 QingCloud AppCenter 创建的 Kubernetes 集群内, 需要将 [DaemonSet](deploy/block/kubernetes/csi-node-ds.yaml) YAML 文件的 *"/var/lib/kubelet"* 字段替换为 *"/data/var/lib/kubelet"*。

```
$ kubectl apply -f ./csi-controller-deploy.yaml
$ kubectl apply -f ./csi-node-ds.yaml
```

- 检查 CSI 插件状态
```
$ kubectl get pods -n kube-system --selector=app=csi-qingcloud
NAME                            READY     STATUS        RESTARTS   AGE
csi-qingcloud-controller-0      3/3       Running       0          5m
csi-qingcloud-node-kks3q        2/2       Running       0          2m
csi-qingcloud-node-pgsbn        2/2       Running       0          2m
```

### 验证
- 由 Kubernetes 集群管理员创建 StorageClass
> 注：示例将创建 `type` 值为 `0` 的 StorageClass，用户可按照后续部分的说明设置 StorageClass 的参数。
```
$ kubectl apply -f https://raw.githubusercontent.com/yunify/qingcloud-csi/master/deploy/block/example/sc.yaml
```

- 创建 PVC
```
$ kubectl apply -f https://raw.githubusercontent.com/yunify/qingcloud-csi/master/deploy/block/example/pvc.yaml
```

- 创建挂载 PVC 的 Deployment
```
$ kubectl apply -f https://raw.githubusercontent.com/yunify/qingcloud-csi/master/deploy/block/example/deploy.yaml
```

- 检查 Pod 状态
```
$ kubectl get po | grep nginx
nginx-84474cf674-zfhbs   1/1       Running   0          1m
```

- 访问容器内挂载存储卷的目录
```
$ kubectl exec -ti deploy-nginx-qingcloud-84474cf674-zfhbs /bin/bash
# cd /mnt
# ls
lost+found
```

### StorageClass参数

如下所示的 StorageClass 资源定义[文件](deploy/block/example/sc.yaml)可用来创建 StorageClass 对象。
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
  replica: "2"
reclaimPolicy: Delete 
```

- `type`: 青云云平台存储卷类型。在青云公有云中， `0` 代表性能型硬盘, `3` 代表超高性能型硬盘, `1` 或 `2`（根据集群所在区不同而参数不同）代表容量型硬盘, `5` 代表企业级分布式 SAN (NeonSAN) 硬盘, `100` 代表基础型硬盘， `200` 代表 SSD 企业型硬盘。 详情见 [QingCloud 文档](https://docs.qingcloud.com/product/api/action/volume/create_volumes.html)。

- `maxSize`, `minSize`: 限制存储卷类型的存储卷容量范围，单位为GiB。青云公有云用户可参考[文档](docs/block-volume-parameter-zh.md)设置。

- `stepSize`: 设置用户所创建存储卷容量的增量，单位为GiB。青云公有云用户可参考[文档](docs/block-volume-parameter-zh.md)设置。

- `fsType`: 支持 `ext3`, `ext4`, `xfs`. 默认为 `ext4`。

- `replica`: `1` 代表单副本硬盘，`2` 代表多副本硬盘。 默认为 `2`。

## 支持
如果有任何问题或建议, 请在 [qingcloud-csi](https://github.com/yunify/qingcloud-csi/issues) 项目提 issue。
