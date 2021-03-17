IMAGE = app
DEBUG ?= false
BUILD_CONTEXT := .
DOCKER_REGISTRY = registry.eywa.rekfuki.dev
#VERSION = $(shell git describe --abbrev=0 --tags | grep $(IMAGE) | rev | cut -d '-' -f1 | rev)
VERSION = 0.0.1

ifeq ("$(VERSION)","")
	$(error Missing tag for $(IMAGE))
endif

.PHONY: run
run:
	@echo "Running $(IMAGE)..."
	yarn build
	GO111MODULE=on go run main.go -debug=$(DEBUG)

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
