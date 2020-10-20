-include $(shell [ -e .build-harness ] || curl -sSL -o .build-harness "https://git.io/mintel-build-harness"; echo .build-harness)

.PHONY: init lint clean build test
init: bh/init

build: go/build

lint: go/lint go/vet

test: go/test

clean: go/clean
	rm -f glmr
