# +-------------------------------------------------------------------------
# | Copyright (C) 2018 Yunify, Inc.
# +-------------------------------------------------------------------------
# | Licensed under the Apache License, Version 2.0 (the "License");
# | you may not use this work except in compliance with the License.
# | You may obtain a copy of the License in the LICENSE file, or at:
# |
# | http://www.apache.org/licenses/LICENSE-2.0
# |
# | Unless required by applicable law or agreed to in writing, software
# | distributed under the License is distributed on an "AS IS" BASIS,
# | WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
# | See the License for the specific language governing permissions and
# | limitations under the License.
# +-------------------------------------------------------------------------

.PHONY: all disk

DISK_IMAGE_NAME=dockerhub.qingcloud.com/csiplugin/csi-qingcloud
DISK_IMAGE_VERSION=canary
DISK_PLUGIN_NAME=qingcloud-disk-csi-driver
ROOT_PATH=$(pwd)
PACKAGE_LIST=./cmd/disk ./pkg/disk ./pkg/server ./pkg/server/instance ./pkg/server/volume ./pkg/server/zone

disk:
	if [ ! -d ./vendor ]; then dep ensure; fi
	CGO_ENABLED=0 GOOS=linux go build -a -ldflags '-extldflags "-static"' -o  _output/${DISK_PLUGIN_NAME} ./cmd/disk

disk-container: disk
	cp _output/${DISK_PLUGIN_NAME} deploy/disk/docker
	docker build -t $(DISK_IMAGE_NAME):$(DISK_IMAGE_VERSION) deploy/disk/docker

deploy:
	kubectl create configmap csi-qingcloud --from-file=config.yaml=./deploy/disk/kubernetes/config.yaml --namespace=kube-system
	kubectl create -f ./deploy/disk/kubernetes/csi-secret.yaml
	kubectl create -f ./deploy/disk/kubernetes/csi-controller-rbac.yaml
	kubectl create -f ./deploy/disk/kubernetes/csi-node-rbac.yaml
	kubectl create -f ./deploy/disk/kubernetes/csi-controller-sts.yaml
	kubectl create -f ./deploy/disk/kubernetes/csi-node-ds.yaml

undeploy:
	kubectl delete -f ./deploy/disk/kubernetes/csi-node-ds.yaml
	kubectl delete -f ./deploy/disk/kubernetes/csi-controller-sts.yaml
	kubectl delete -f ./deploy/disk/kubernetes/csi-node-rbac.yaml
	kubectl delete -f ./deploy/disk/kubernetes/csi-controller-rbac.yaml
	kubectl delete -f ./deploy/disk/kubernetes/csi-secret.yaml
	kubectl delete cm csi-qingcloud -n kube-system

fmt:
	go fmt ${PACKAGE_LIST}

fmt-deep: fmt
	gofmt -s -w -l ${PACKAGE_LIST}

sanity-test:
	${ROOT_PATH}/csi-sanity --csi.endpoint /var/lib/kubelet/plugins/disk.csi.qingcloud.com/csi.sock --csi.testvolumesize 107374182400

clean:
	go clean -r -x
	rm -rf ./_output
	rm -rf deploy/disk/docker/${DISK_PLUGIN_NAME}