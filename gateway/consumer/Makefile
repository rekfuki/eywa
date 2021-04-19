IMAGE = gateway-consumer
DEBUG ?= false
BUILD_CONTEXT := .
DOCKER_REGISTRY = registry.eywa.rekfuki.dev
VERSION = $(shell git tag --sort version:refname | grep $(IMAGE) | tail -n1 | rev | cut -d '-' -f1 | rev)
# VERSION = 0.0.1

ifeq ("$(VERSION)","")
	$(error Missing tag for $(IMAGE))
endif

.PHONY: run
run:
	@echo "Running $(IMAGE)..."
	GO111MODULE=on go run main.go -in-cluster=false -debug=$(DEBUG)

.PHONY: go-build
go-build:
	@echo "Build $(IMAGE) binary..."
	GO111MODULE=on CGO_ENABLED=0 go build -o $(BUILD_CONTEXT)/$(IMAGE) -v main.go

.PHONY: build
build: go-build
	@echo "Build docker image..."
	docker build --tag $(IMAGE):$(VERSION) $(BUILD_CONTEXT)
	docker tag $(IMAGE):$(VERSION) $(IMAGE):latest

.PHONY: tag
tag:
	@echo "Tagging docker image (latest)"
	docker tag $(IMAGE):latest $(DOCKER_REGISTRY)/$(IMAGE):latest
	docker tag $(IMAGE):$(VERSION) $(DOCKER_REGISTRY)/$(IMAGE):$(VERSION)

.PHONY: push
push:
	@echo "Push docker image..."
	docker push $(DOCKER_REGISTRY)/$(IMAGE):latest
	docker push $(DOCKER_REGISTRY)/$(IMAGE):$(VERSION)
