ifndef WATCHDOG
    $(error WATCHDOG is not set)
endif

ifeq ("$(wildcard ./templates/$(WATCHDOG))","")
    $(error Directory ./templates/$(WATCHDOG)/ does not exist)
endif

IMAGE = $(WATCHDOG)-watchdog
VERSION = $(shell git describe --abbrev=0 --tags | grep $(IMAGE) | rev | cut -d '-' -f1 | rev)

ifeq ("$(VERSION)","")
	$(error Missing tag for $(IMAGE))
endif

BUILD_CONTEXT := ./templates/$(WATCHDOG)
DOCKER_REGISTRY = registry.eywa.rekfuki.dev

.PHONY: build-watchdog
build-watchdog:
	@echo "Build watchdog binary..."
	GO111MODULE=on CGO_ENABLED=0 go build -o $(BUILD_CONTEXT)/$(IMAGE) -v main.go

.PHONY: build-handler

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