# 安装指南
此安装指南将 CSI 插件安装在 *kube-system* namespace 内。

## 准备材料
- 基于 QingCloud 云平台主机的 Kubernetes 1.14+ 集群，并且按照要求设置 Kubernetes 参数
  - 在 Kubernetes 控制平面内将 `--allow-privileged` 项设置为 `true` 。
  - 启用（默认开启）[Mount Propagation](https://kubernetes.io/docs/concepts/storage/volumes/#mount-propagation) 特性。
  - 设置 `--feature-gates=CSINodeInfo=true,CSIDriverRegistry=true,KubeletPluginsWatcher=true`
- QingCloud CSI 安装包
下载安装包并解压
```
$ wget $(curl --silent "https://api.github.com/repos/yunify/qingcloud-csi/releases/latest" | \
  grep browser_download_url | grep install|cut -d '"' -f 4)
$ tar -xvf csi-qingcloud-install.tar.gz
$ cd csi-qingcloud-install
```
- [安装 Kustomize](https://github.com/kubernetes-sigs/kustomize/blob/master/docs/zh/INSTALL.md)

## 修改配置文件
修改安装包内 config.yaml 配置文件（config.yaml）
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

## 执行安装命令
```
kustomize build overlays/prod|kubectl apply -f -
```

## 检查 CSI 插件状态
```
$ kubectl get pods -n kube-system --selector=app=csi-qingcloud
NAME                                        READY   STATUS    RESTARTS   AGE
csi-qingcloud-controller-7bf97b7f5f-5mq4p   5/5     Running   0          3h2m
csi-qingcloud-node-9wp2p                    2/2     Running   0          3h2m
csi-qingcloud-node-vtkdk                    2/2     Running   0          3h2m
csi-qingcloud-node-zzm9t                    2/2     Running   0          3h2m
```

## 卸载方法
删除 QingCloud CSI PVC 后删除 QingCloud CSI 存储插件

```
kustomize build overlays/prod|kubectl delete -f -
```