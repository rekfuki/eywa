ifndef WATCHDOG
    $(error WATCHDOG is not set)
endif

ifeq ("$(wildcard $(WATCHDOG))","")
    $(error Directory $(WATCHDOG)/ does not exist)
endif

IMAGE = $(WATCHDOG)_watchdog
VERSION = $(shell python2 -c "from tag_versioner.version_reader import VersionReader; print VersionReader('$(IMAGE)-').get_docker_version()")
BUILD_CONTEXT := $(WATCHDOG)
DOCKER_REGISTRY = registry.eywa.rekfuki.dev
TAGGED = $(shell python2 -c "from tag_versioner.version_reader import VersionReader; print VersionReader('$(IMAGE)-').is_tagged_and_clean()")

TAGS =
CGO_ENABLED = 1
ifeq ($(WATCHDOG), nodejs)
	CGO_ENABLED = 0
else
	TAGS = $(WATCHDOG)
endif

.PHONY: build
build: ## Build the image and tag it as 'latest'
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
ifeq ($(TAGGED),True)
	docker push $(DOCKER_REGISTRY)/$(IMAGE):latest
	docker push $(DOCKER_REGISTRY)/$(IMAGE):$(VERSION)
else
	@echo "Tagged container not pushed because $(VERSION) is not clean - commit and tag please!"
endif