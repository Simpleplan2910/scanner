.DEFAULT_GOAL := cmd/v1/scanner
.PHONY: clean
SHELL = /usr/bin/env bash

VPREFIX := weeio_v2/internals/version
GO=go
TESTS=.
VERSION ?= $(IMAGE_TAG)

ifeq ($(VERSION),)
	VERSION := dev
endif

GIT_REVISION := $(shell git rev-parse --short HEAD)
GIT_BRANCH := $(shell git rev-parse --abbrev-ref HEAD)
GO_LDFLAGS   := -s -w -X $(VPREFIX).Branch=$(GIT_BRANCH) -X $(VPREFIX).Version=$(VERSION) -X $(VPREFIX).Revision=$(GIT_REVISION) -X $(VPREFIX).BuildUser=$(shell whoami)@$(shell hostname) -X $(VPREFIX).BuildDate=$(shell date -u +"%Y-%m-%dT%H:%M:%SZ")
GO_FLAGS     := -ldflags "-extldflags \"-static\" $(GO_LDFLAGS)" 
# Packages lists
TE_PACKAGES=$(shell $(GO) list ./...)

start-docker:  ## Starts the docker containers for local development.
	@echo Starting docker containers
	docker-compose -f docker-compose-local.yaml up -d 

clean-docker:
	@echo Removing docker containers

	docker-compose down -f docker-compose-local.yaml  -v
	docker-compose rm  -f docker-compose-local.yaml  -v

stop-docker: ## Stops the docker containers for local development.
	@echo Stopping docker containers

	docker-compose stop

test: start-docker 
	$(GO) test -run=$(TESTS) $(TE_PACKAGES)

test-ci: 
	$(GO) test -run=$(TESTS) $(TE_PACKAGES) 

cmd/v1/scanner: cmd/v1/main.go
ifeq ($(o),)
	CGO_ENABLED=0 go build $(GO_FLAGS) -o $@ ./$(@D)
else
	CGO_ENABLED=0 go build $(GO_FLAGS) -o $(o) ./$(@D)
endif

clean:
	rm -rf  cmd/v1/scanner