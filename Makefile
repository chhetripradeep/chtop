GOOS ?= $(shell go env GOOS)
GOARCH ?= $(shell go env GOARCH)
GOBUILD=CGO_ENABLED=0 installsuffix=cgo go build -trimpath

.DEFAULT_GOAL := run

.PHONY: run
run:
	go run main.go --metrics-url http://localhost:9363/metrics --queries-url localhost:9000 --config ./chtop.yaml

.PHONY: build
build:
	${GOBUILD} -o chtop-$(GOOS)-$(GOARCH) main.go

.PHONY: build-linux-amd64
build-linux-amd64:
	GOOS=linux GOARCH=amd64 $(MAKE) build

.PHONY: build-linux-arm64
build-linux-arm64:
	GOOS=linux GOARCH=arm64 $(MAKE) build

.PHONY: build-darwin-amd64
build-darwin-amd64:
	GOOS=darwin GOARCH=amd64 $(MAKE) build

.PHONY: build-darwin-arm64
build-darwin-arm64:
	GOOS=darwin GOARCH=arm64 $(MAKE) build

.PHONY: build-all-platforms
build-all-platforms: build-linux-amd64 build-linux-arm64 build-darwin-amd64 build-darwin-arm64

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

.PHONY: tar
tar:
	tar -czvf chtop-$(GOOS)-$(GOARCH).tar.gz chtop-$(GOOS)-$(GOARCH) chtop.yaml

.PHONY: tar-linux-amd64
tar-linux-amd64:
	GOOS=linux GOARCH=amd64 $(MAKE) tar

.PHONY: tar-linux-arm64
tar-linux-arm64:
	GOOS=linux GOARCH=arm64 $(MAKE) tar

.PHONY: tar-darwin-amd64
tar-darwin-amd64:
	GOOS=darwin GOARCH=amd64 $(MAKE) tar

.PHONY: tar-darwin-arm64
tar-darwin-arm64:
	GOOS=darwin GOARCH=arm64 $(MAKE) tar

.PHONY: tar-all-platforms
tar-all-platforms: tar-linux-amd64 tar-linux-arm64 tar-darwin-amd64 tar-darwin-arm64

.PHONY: upgrade-deps
upgrade-deps:
	go get -u ./...
	go mod tidy
