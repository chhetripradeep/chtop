.DEFAULT_GOAL := run

.PHONY: run
run:
	go run main.go

.PHONY: build
build:
	go build -o chtop main.go

.PHONY: fmt
fmt:
	go fmt ./...

.PHONY: lint
lint: install-tools
	golangci-lint -v run --allow-parallel-runners ./...

.PHONY: install
install:
	go install
	goimports -w -local github.com/chhetripradeep/chtop ./

.PHONY: install-tools
install-tools:
	go install golang.org/x/tools/cmd/goimports@latest
	go install github.com/golangci/golangci-lint/cmd/golangci-lint@latest
