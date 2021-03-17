IMAGE = envoy
BUILD_CONTEXT := .
DOCKER_REGISTRY = registry.eywa.rekfuki.dev
# VERSION = $(shell git describe --abbrev=0 --tags | grep $(IMAGE) | rev | cut -d '-' -f1 | rev)
VERSION = 0.0.1

ifeq ("$(VERSION)","")
	$(error Missing tag for $(IMAGE))
endif

watcher: watch/watcher.go
	go build watch/watcher.go

.PHONY: build
build: watcher
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
