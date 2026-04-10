BINARY = checks-qe
GOBIN ?= $(shell go env GOPATH)/bin

.PHONY: build lint vet test clean

build:
	go build -o $(BINARY) ./cmd/checks-qe/

lint:
	golangci-lint run ./...

vet:
	go vet ./...

test:
	go test ./...

clean:
	rm -f $(BINARY)

list: build
	./$(BINARY) list
