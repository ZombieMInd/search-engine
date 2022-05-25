.PHONY: build
build:
	go build -v ./cmd/api

.PHONY: test
test:
	go test -v -race -timeout 30s ./...

.PHONY: build & run
run:
	go build -v ./cmd/api && ./api

.DEFAULT_GOAL := build
