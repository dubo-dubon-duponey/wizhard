.PHONY: default
default: all

.PHONY: fmt
fmt:
	go fmt ./...

.PHONY: vet
vet:
	go vet ./...

.PHONY: build
build:
	go build -v -ldflags "-s -w" -o dist/wizhard ./cmd/wizhard/main.go

.PHONY: all
all: fmt vet build
