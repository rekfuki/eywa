IMAGE = of-watchdog
DEBUG ?= false
BUILD_CONTEXT := .
DOCKER_REGISTRY = registry.eywa.rekfuki.dev
VERSION = $(shell git tag --sort version:refname | grep $(IMAGE) | tail -n1 | rev | cut -d '-' -f1 | rev)
#VERSION = 1.0.0

ifeq ("$(VERSION)","")
	$(error Missing tag for $(IMAGE))
endif

.PHONY: build-watchdog
build-watchdog:
	@echo "Build watchdog binary..."
	GO111MODULE=on CGO_ENABLED=0 go build -o $(BUILD_CONTEXT)/$(IMAGE) -v main.go

.PHONY: build
build: build-watchdog 
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
