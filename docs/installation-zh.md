# 安装指南

## 准备
- **QingCloud 云平台** 和 **API 密钥**
- Kubernetes v1.13+ 集群
  - 启用 [Priviliged Pod](https://kubernetes-csi.github.io/docs/Setup.html#enable-privileged-pods)，将 Kubernetes 的 kubelet 和 kube-apiserver 组件配置项 `--allow-privileged` 设置为 `true`
  - 启用（默认开启）[Mount Propagation](https://kubernetes.io/docs/concepts/storage/volumes/#mount-propagation) 特性
  - 启用 Kubernetes 的 kubelet, controller-manager 和 kube-apiserver 组件的若干 [Feature Gate](https://kubernetes-csi.github.io/docs/Setup.html#enabling-features)
  ```
  --feature-gates=VolumeSnapshotDataSource=true,KubeletPluginsWatcher=true,CSINodeInfo=true,CSIDriverRegistry=true
  ```
- 下载安装包
```
$ wget $(curl --silent "https://api.github.com/repos/yunify/qingcloud-csi/releases/latest" | \
  grep browser_download_url | grep install|cut -d '"' -f 4)
$ tar -xvf csi-qingcloud-install.tar.gz
$ cd csi-qingcloud-install
```

## 配置
### 设置 kubelet 路径

如果Kubernetes集群的Kubelet设置了 `--root-dir` 选项（默认为 *`/var/lib/kubelet`* ），请将 `ds-node.yaml` 文件内 *`/var/lib/kubelet`* 的值替换为 `--root-dir` 选项的值。

### 创建 ConfigMap
修改安装包内配置文件（config.yaml）
> 注：在 QingCloud AppCenter 的 Kubernetes 集群内部署本插件无需配置此文件，直接使用 AppCenter 自动创建的配置文件（ `/etc/qingcloud/client.yaml` ）即可。

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

- `host`, `port`. `protocol`, `uri`: 共同构成 QingCloud IaaS 平台服务的 url.

## 部署

### 创建 ConfigMap
```
$ kubectl create configmap csi-qingcloud --from-file=config.yaml=./config.yaml --namespace=kube-system
```

### 创建 Docker 镜像仓库密钥
```
$ kubectl apply -f ./csi-secret.yaml
```

### 创建注册插件对象

```
$ kubectl create -f ./crd-csidriver.yaml
$ kubectl create -f ./crd-csinodeinfo.yaml
$ kubectl create -f ./csidriver-qingcloud.yaml
```

### 创建访问控制相关对象
```
$ kubectl apply -f ./csi-controller-rbac.yaml
$ kubectl apply -f ./csi-node-rbac.yaml
```

### 部署 CSI 插件

```
$ kubectl apply -f ./csi-controller-sts.yaml
$ kubectl apply -f ./csi-node-ds.yaml
```

## 验证

### 检查 CSI 插件状态
```
$ kubectl get pods -n kube-system --selector=app=csi-qingcloud
NAME                            READY     STATUS        RESTARTS   AGE
csi-qingcloud-controller-0      3/3       Running       0          5m
csi-qingcloud-node-kks3q        2/2       Running       0          2m
csi-qingcloud-node-pgsbn        2/2       Running       0          2m
```