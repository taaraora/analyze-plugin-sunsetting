SHELL := /bin/sh

MAKEFILE_PATH := $(abspath $(lastword $(MAKEFILE_LIST)))
CURRENT_DIR := $(patsubst %/,%,$(dir $(MAKEFILE_PATH)))

DOCKER_IMAGE_NAME := $(if ${TRAVIS_REPO_SLUG},${TRAVIS_REPO_SLUG},supergiant/analyze-plugin-sunsetting)
DOCKER_IMAGE_TAG := $(if ${TAG},${TAG},$(shell git describe --tags --always | tr -d v || echo 'latest'))


define LINT
	@echo "Running code linters..."
	revive
	@echo "Running code linters finished."
endef

define GOIMPORTS
	goimports -v -w -local github.com/supergiant/analyze-plugin-sunsetting ${CURRENT_DIR}
endef

define TOOLS
		if [ ! -x "`which revive 2>/dev/null`" ]; \
        then \
        	echo "revive linter not found."; \
        	echo "Installing linter... into ${GOPATH}/bin"; \
        	GO111MODULE=off go get -u github.com/mgechev/revive ; \
        fi
endef


.PHONY: default
default: lint


.PHONY: lint
lint: tools
	@$(call LINT)


.PHONY: test
test:
	go test -race ./...

.PHONY: tools
tools:
	@$(call TOOLS)

.PHONY: goimports
goimports:
	@$(call GOIMPORTS)

.PHONY: build-image
build-image: build-ui gen-assets build
	docker build -t $(DOCKER_IMAGE_NAME):$(DOCKER_IMAGE_TAG) ./dist -f ./Dockerfile
	docker tag $(DOCKER_IMAGE_NAME):$(DOCKER_IMAGE_TAG) $(DOCKER_IMAGE_NAME):latest

.PHONY: build
build:
	./scripts/build.sh

.PHONY: gen-assets
gen-assets:
	./scripts/gen-assets.sh

.PHONY: build-ui
build-ui:
	./scripts/build-ui.sh

.PHONY: push
push:
	docker push $(DOCKER_IMAGE_NAME):$(DOCKER_IMAGE_TAG)

.PHONY: gofmt
gofmt:
	go fmt ./...

.PHONY: fmt
fmt: gofmt goimports
