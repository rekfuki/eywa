IMAGE = gateway
BUILD_CONTEXT := .
DOCKER_REGISTRY = registry.eywa.rekfuki.dev
VERSION = $(shell git describe --abbrev=0 --tags | grep $(IMAGE) | rev | cut -d '-' -f1 | rev)

ifeq ("$(VERSION)","")
	$(error Missing tag for $(IMAGE))
endif

.PHONY: run
run:
	@echo "Running gateway..."
	KUBERNETES_SERVICE_HOST=192.168.1.44 KUBERNETES_SERVICE_PORT=80k9s GO111MODULE=on go run main.go

.PHONY: go-build
go-build:
	@echo "Build gateway binary..."
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