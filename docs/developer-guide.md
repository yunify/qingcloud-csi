<!-- START doctoc generated TOC please keep comment here to allow auto update -->
<!-- DON'T EDIT THIS SECTION, INSTEAD RE-RUN doctoc TO UPDATE -->


- [Developer Guide](#developer-guide)
  - [Process](#process)
  - [Build Image](#build-image)
  - [E2E Test](#e2e-test)

<!-- END doctoc generated TOC please keep comment here to allow auto update -->

# Developer Guide

## Process
We use [Github workflow](https://github.com/kubernetes/community/blob/master/contributors/guide/github-workflow.md) to develop and submit code.

## Build Image
[Multi-stage](https://docs.docker.com/develop/develop-images/multistage-build/) is used to build Docker images in Docker v17.05+.

1. Example Environment
  - Ubuntu 16.04/CentOS 7.5
  - Docker 18.09.8
2. Download Repo
```cassandraql
root@dev:~# git clone https://github.com/yunify/qingcloud-csi.git
root@dev:~# cd qingcloud-csi
```
3. Build
```cassandraql
root@dev:~/qingcloud-csi# make disk-container
docker build -t csiplugin/csi-qingcloud:v1.1.0-rc.4 -f deploy/disk/docker/Dockerfile  .
Sending build context to Docker daemon   57.7MB
Step 1/11 : FROM golang:1.12.7-alpine as builder
 ---> 6b21b4c6e7a3
Step 2/11 : WORKDIR /qingcloud-csi
 ---> Using cache
 ---> d99239a7aae4
Step 3/11 : COPY . .
 ---> f1202e19b989
Step 4/11 : RUN CGO_ENABLED=0 GOOS=linux go build -a -mod=vendor  -ldflags "-s -w" -o  _output/qingcloud-disk-csi-driver ./cmd/disk
 ---> Running in 67e14ef016d2
Removing intermediate container 67e14ef016d2
 ---> d7c63e0b4bcb
Step 5/11 : FROM k8s.gcr.io/debian-base:v1.0.0
 ---> 204e96332c91
Step 6/11 : LABEL maintainers="Yunify"
 ---> Using cache
 ---> 06e7af6cb693
Step 7/11 : LABEL description="QingCloud CSI plugin"
 ---> Using cache
 ---> f5bfdbbd78bf
Step 8/11 : RUN clean-install util-linux e2fsprogs xfsprogs  mount ca-certificates udev
 ---> Using cache
 ---> cf04e131cbbb
Step 9/11 : COPY --from=builder /qingcloud-csi/_output/qingcloud-disk-csi-driver /qingcloud-disk-csi-driver
 ---> df6f0270b1db
Step 10/11 : RUN chmod +x /qingcloud-disk-csi-driver &&     mkdir -p /var/log/qingcloud-disk-csi-driver
 ---> Running in ac644e4db06e
Removing intermediate container ac644e4db06e
 ---> 9ae0b2614f7c
Step 11/11 : ENTRYPOINT ["/qingcloud-disk-csi-driver"]
 ---> Running in c7e4defedbb3
Removing intermediate container c7e4defedbb3
 ---> 3e8a3a1f45c5
Successfully built 3e8a3a1f45c5
Successfully tagged csiplugin/csi-qingcloud:canary
```

## E2E Test

1. Compile
```cassandraql
# git clone https://github.com/kubernetes-csi/csi-test.git
# cd kubernetes-csi/csi-test/cmd/csi-sanity/dist/csi-sanity
# make linux_amd64_dist
```

2. Edit Storage Class Parameters
```cassandraql
# cat parameters.yaml
type: "200"
tags: ""
fstype: "ext4"
```

3. Test
```cassandraql
# ./csi-sanity -csi.endpoint /var/lib/kubelet/plugins/disk.csi.qingcloud.com/csi.sock -csi.testvolumeparameters parameters.yaml -csi.testvolumeexpandsize 10737418240
```