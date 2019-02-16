.PHONY: mocks

$GOPATH/bin/mockgen:
	 GO111MODULE=off go get github.com/golang/mock/gomock
	 GO111MODULE=off go install github.com/golang/mock/mockgen

$GOPATH/src/periph.io/x/periph/conn/gpio:
	 GO111MODULE=off go get periph.io/x/periph/conn/gpio

tools: $GOPATH/bin/mockgen

mocks: tools $GOPATH/src/periph.io/x/periph/conn/gpio
	 mockgen -source $(GOPATH)/src/periph.io/x/periph/conn/gpio/gpio.go -destination internal/mocks/pin_io.go -package mocks
	 mockgen -source pkg/scanner/libnfc.go -destination internal/mocks/libnfc.go -package mocks
	 mockgen -source pkg/scanner/scanner.go -destination internal/mocks/scanner.go -package mocks
	 mockgen -source pkg/datastore/datastore.go -destination internal/mocks/datastore.go -package mocks

test:
	go test -v -cover -short ./...

test-full:
	go test -v -cover ./...

clean:
	rm -rf build/out

build: clean
	CC="arm-linux-gnueabihf-gcc" CGO_CFLAGS="-I /usr/include/nfc" GOOS=linux GOARCH=arm GOARM=6 CGO_ENABLED=1 go build -o build/out/linux/arm/open-keyless-controller --ldflags '-linkmode external -extldflags "-static"' github.com/lodge93/open-keyless/cmd/open-keyless-controller

release: build
	docker run --rm -v $(PWD)/build:/build -w /build -e PLUGIN_DEB_SYSTEMD=/build/package/systemd/open-keyless-controller.service -e PLUGIN_NAME=open-keyless-controller -e PLUGIN_VERSION=snapshot-$(shell git log -n 1 --pretty=format:"%H") -e PLUGIN_INPUT_TYPE=dir -e PLUGIN_OUTPUT_TYPE=deb -e PLUGIN_PACKAGE=/build/out/open-keyless-controller-snapshot-$(shell git log -n 1 --pretty=format:"%H").deb -e PLUGIN_COMMAND_ARGUMENTS=/build/out/linux/arm/open-keyless-controller=/usr/local/bin/ lodge93/drone-fpm:latest
