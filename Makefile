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

.PHONY: all blockplugin

BLOCK_IMAGE_NAME=dockerhub.qingcloud.com/csiplugin/csi-qingcloud
BLOCK_IMAGE_VERSION=v0.2.0
BLOCK_PLUGIN_NAME=blockplugin
ROOT_PATH=$(pwd)
PACKAGE_LIST=./cmd/block ./pkg/block ./pkg/server ./pkg/server/instance ./pkg/server/volume

blockplugin:
	if [ ! -d ./vendor ]; then dep ensure; fi
	CGO_ENABLED=0 GOOS=linux go build -a -ldflags '-extldflags "-static"' -o  _output/${BLOCK_PLUGIN_NAME} ./cmd/block

blockplugin-container: blockplugin
	cp _output/${BLOCK_PLUGIN_NAME} deploy/block/docker
	docker build -t $(BLOCK_IMAGE_NAME):$(BLOCK_IMAGE_VERSION) deploy/block/docker

fmt:
	go fmt ${PACKAGE_LIST}
	gofmt -s -w -l ${PACKAGE_LIST}

clean:
	go clean -r -x
	rm -rf ./_output
	rm -rf deploy/block/docker/${BLOCK_PLUGIN_NAME}
