TARGET_OS = linux windows darwin
TARGET_ARCH = amd64

.PHONY: init lint fmt clean build test

install:
	go install ./...

build: vgo/gox
	@mkdir -p build/
	./vgo/gox -os='$(TARGET_OS)' -arch='$(TARGET_ARCH)' -output='build/glmr_{{.OS}}_{{.Arch}}' ./...

build-local:
	@mkdir -p build
	go build -o build ./...

lint: vgo/goimports
	go vet ./...
	test $(shell ./vgo/goimports -l . | wc -l) -eq 0

fmt: vgo/goimports
	goimports -w .

test:
	go test -v ./...

clean:
	rm -rf build/
	rm -rf vgo/

vgo/goimports: 	
	@mkdir -p vgo
	./scripts/vgoget.sh golang.org/x/tools/cmd/goimports vgo

vgo/gox:
	@mkdir -p vgo
	./scripts/vgoget.sh github.com/mitchellh/gox vgo
