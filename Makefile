
.PHONY: lint build test fmt

GO_FILES=$(find . -regex ".*\.go")

lint:
	go vet ./...
	test $(shell goimports -l . | wc -l) -eq 0

fmt:
	goimports -w .

build: dist/glmr

test:
	go test -v ./...

dist:
	mkdir dist

dist/glmr: $(GO_FILES) dist
	go build -o dist ./...


