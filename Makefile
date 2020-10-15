
.PHONY: lint build test

GO_FILES=$(find . -regex ".*\.go")

lint:
	go vet .
	test $(shell goimports -l . | wc -l) -eq 0

build: dist/glmr

dist:
	mkdir dist

dist/glmr: $(GO_FILES) dist
	go build -o dist ./...


