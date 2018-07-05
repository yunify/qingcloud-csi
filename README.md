# QingCloud-CSI

[![Build Status](https://travis-ci.org/yunify/qingcloud-csi.svg?branch=master)](https://travis-ci.org/yunify/qingcloud-csi)

Kubernetes volume plugin based on CSI specification which support block storage of qingcloud

## Description
QingCloud CSI plugin implements an interface between Container Storage Interface([CSI](https://github.com/container-storage-interface/)) enabled Container Orchestrator(CO) and the storage of QingCloud. Currently, QingCloud CSI plugin is tested in Kubernetes v1.10.0+ environment and should be able to work in any CSI enabled CO.

## Block Plugin

### Configuration
- StorageClass: set the name and other parameters of your block storage server

### Compiling
QingCloud CSI plugin can be complied as a binary file or a container.  We can get a binary file in _output folder. When compiled as a container, the image is stored in a local Docker's image store.

To compile a binary file:
```
$ make blockplugin
```

To compile a Docker image:
```
$ make blockplugin-container
```

You can find image in your local image store
```
$ docker images | grep csi-qingcloud
csi-qingcloud		v0.2.0		640a9519e59b		55 minutes ago		332MB
```

### Deploying
- Deploy helper containers that the Kubernetes team provides.
```
$ kubectl create -f deploy/block/kubernetes/csi-provisioner.yaml
$ kubectl create -f deploy/block/kubernetes/csi-attacher.yaml
```

- Deploy CSI plugin that storage vendor provides.
```
$ kubectl create -f deploy/block/kubernetes/csi-qingcloud.yaml
```

- Create storage class by Kubernetes cluster administrator
```
$ kubectl create -f deploy/block/kubernetes/sc.yaml
```

- Check status of CSI plugin
```
$ kubectl get pods | grep csi
csi-attacher-0        1/1       Running       0          3d
csi-provisioner-0     1/1       Running       0          3d
csi-qingcloud-pgsbn   2/2       Running       0          1h
```

### Verification
- Create PVC
```
$ kubectl create -f deploy/block/kubernetes/pvc.yaml
```

- Check PVC and PV
```
$ kubectl get pvc
NAME            STATUS    VOLUME                 CAPACITY   ACCESS MODES   STORAGECLASS      AGE
qingcloud-pvc   Bound     pvc-77a1e29168ab11e8   10Gi       RWO            csi-qingcloud     22s

$ kubectl describe pv pvc-77a1e29168ab11e8
Name:            pvc-77a1e29168ab11e8
Labels:          <none>
Annotations:     pv.kubernetes.io/provisioned-by=csi-qingcloud
Finalizers:      [kubernetes.io/pv-protection]
StorageClass:    csi-qingcloud
Status:          Bound
Claim:           default/qingcloud-pvc
Reclaim Policy:  Delete
Access Modes:    RWO
Capacity:        10Gi
Node Affinity:   <none>
Message:         
Source:
    Type:          CSI (a Container Storage Interface (CSI) volume source)
    Driver:        csi-qingcloud
    VolumeHandle:  vol-8o3x8lvh
    ReadOnly:      false
Events:            <none>
```

- Create deployment mounting PVC
```
$ kubectl create -f deploy/block/kubernetes/deploy.yaml
```

- Check deploy
```
$ kubectl get po | grep deploy
deploy-nginx-qingcloud-84474cf674-zfhbs   1/1       Running   0          1m
```

```
$ kubectl exec -ti deploy-nginx-qingcloud-84474cf674-zfhbs /bin/bash
// We can access the directoriy mounting persistent volume in container
# cd /mnt
# ls
lost+found
```

## Support
If you have any qustions or suggestions, please submit an issue at [qingcloud-csi](https://github.com/yunify/qingcloud-csi/issues)
