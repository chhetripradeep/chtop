.DEFAULT_GOAL := run

.PHONY: run
run:
	go run main.go

.PHONY: build
build:
	go build -o chtop main.go

.PHONY: install
install:
	go install
