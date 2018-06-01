.PHONY: all blockplugin

blockplugin:
	if [ ! -d ./vendor ]; then dep ensure; fi
	CGO_ENABLED=0 GOOS=linux go build -a -ldflags '-extldflags "-static"' -o  _output/blockplugin ./cmd/block

clean:
	go clean -r -x
	rm -rf ./_output
