BINARY = checks-qe
GOBIN ?= $(shell go env GOPATH)/bin
REGISTRY ?= quay.io
IMAGE_NAME ?= redhat-best-practices-for-k8s/checks-qe
IMAGE_TAG ?= unstable
VERSION ?= dev

.PHONY: build lint vet test clean build-image

build:
	go build -ldflags "-X main.version=$(VERSION)" -o $(BINARY) ./cmd/checks-qe/

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

build-image:
	docker build \
		--build-arg VERSION=$(VERSION) \
		-t $(REGISTRY)/$(IMAGE_NAME):$(IMAGE_TAG) \
		-f Dockerfile .
