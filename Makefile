.PHONY: build ci lint test

build:
	go build ./...

ci: lint test build

lint:
	golangci-lint run

test:
	go test -v ./...
