.PHONY: all blockplugin

BLOCK_IMAGE_NAME=dockerhub.qingcloud.com/wiley/csi-qingcloud
BLOCK_IMAGE_VERSION=latest
BLOCK_PLUGIN_NAME=blockplugin

blockplugin:
	if [ ! -d ./vendor ]; then dep ensure; fi
	CGO_ENABLED=0 GOOS=linux go build -a -ldflags '-extldflags "-static"' -o  _output/${BLOCK_PLUGIN_NAME} ./cmd/block

blockplugin-container: blockplugin
	cp _output/${BLOCK_PLUGIN_NAME} deploy/block/docker
	docker build -t $(BLOCK_IMAGE_NAME):$(BLOCK_IMAGE_VERSION) deploy/block/docker

clean:
	go clean -r -x
	rm -rf ./_output
	rm -rf deploy/block/docker/${BLOCK_PLUGIN_NAME}
