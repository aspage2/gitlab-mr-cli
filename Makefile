
.PHONY: lint build test

GO_FILES=$(find . -regex ".*\.go")

lint:
	go vet .
	test $(shell goimports -l . | wc -l) -eq 0

build: dist/glmr

dist/glmr: $(GO_FILES)
	go build -o dist ./...


