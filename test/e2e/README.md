# QingCloud-CSI E2E Test

## Description

This directory contains scripts and config templates used to run [Kubernetes external storage e2e test](https://github.com/kubernetes/kubernetes/tree/master/test/e2e/storage/external).

## Prerequisites

The test can only be run on QingCloud nodes with a installed Kubernetes cluster, as it actually creates/attaches volumes by calling the QingCloud IAAS API.

Make sure it's a clean environment with no existing `qingcloud-csi`, which will conflict with the csi driver of the e2e test.

Make sure no volume is attached to any node, which will mess up with the volume limits test. If that's impossible, add an `ClientNodeName` entry in the `testdriver_template.yaml`, set its value to a node which has no volume attached, like:

```yaml
StorageClass:
  FromExistingClassName: xxx
SnapshotClass:
  FromExistingClassName: xxx
ClientNodeName: node1
```

Notice this setting will disable the tests that run on multi nodes.

Make sure a kubeconfig file with admin access of the cluster exists under the path `${HOME}/.kube/config`.

And in order to be authorized to call that API, some configurations are needed beforehand, on the node which you will run the test:

- put your QingCloud access key id under the path `/etc/qingcloud/access_key_id`
- put your QingCloud secret access key under the path `/etc/qingcloud/secret_access_key`
- put the zone where your nodes' in under the path `/etc/qingcloud/zone`

## Optional

1. The zones that have a good network connectivity to global internet is preferred, like `ap2a`, as it will download the e2e test packages from google when it runs for the first time. If that's impossible, you can download the package somewhere else, and upload it to the test node:

```bash
# change the ${k8s_server_version} to the version of your k8s server
curl -L https://storage.googleapis.com/kubernetes-release/release/${k8s_server_version}/kubernetes-test-linux-amd64.tar.gz --output e2e-tests.tar.gz
# upload the file to the test node before execute the following 
tar -xf e2e-tests.tar.gz --directory=./qingcloud-csi/test/e2e && rm e2e-tests.tar.gz
```

2. Multiple nodes are preferred, if that's possible, as some tests will check against volume drifting from one node to another. If there is only one, those tests will be skipped by the script.

## Run

Execute `./run_e2e_test.sh` to run the e2e tests. You can run it as any user as long as the prerequisites are all met.

This script will download `ginkgo` and `e2e.test` if they are not found or incompatible with the k8s server version. And It will build a docker image from source code with a tag derived from the latest git commit hash, and push it to other nodes, if there are any.

You can skip some tests/only run some tests, by editing the following lines of the script:

```bash
ginkgo -focus='External.Storage.*' \
-skip='(.*Disruptive.*|.*stress.*|.*should resize volume when PVC is edited while pod is using it.*)' \
```

like

```bash
# skip block volume tests and only run volume expansion tests
ginkgo -focus='External.Storage.*expansion.*' \
-skip='(.*block.*.*Disruptive.*|.*stress.*|.*should resize volume when PVC is edited while pod is using it.*)' \
```

It's simple regex, refer to the [Ginkgo docs](https://onsi.github.io/ginkgo/#the-ginkgo-cli) for more detail.