TARGET_OS = linux windows darwin
TARGET_ARCH = amd64

VERSION_TAG ?= $(shell git describe HEAD)
LDFLAGS = -X gitlab.com/mintel/personal-dev/apage/glmr/cmd.AppVersion=$(VERSION_TAG)

.PHONY: init lint fmt clean build install build-local test

install:
	@echo "Installing to your GOPATH..."
	@go install -ldflags="$(LDFLAGS)" ./...

build: vgo/gox
	@mkdir -p build/
	./vgo/gox -os='$(TARGET_OS)' -arch='$(TARGET_ARCH)' -ldflags="$(LDFLAGS)" -output='build/glmr_{{.OS}}_{{.Arch}}' ./...

build-local:
	@mkdir -p build
	go build -ldflags="$(LDFLAGS)" -o build ./...

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
