.DEFAULT_GOAL := cmd/v1/scanner
.PHONY: clean start-compse
SHELL = /usr/bin/env bash

VPREFIX := scanner/internals/version
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

cmd/v1/scanner: cmd/v1/main.go
ifeq ($(o),)
	CGO_ENABLED=0 go build $(GO_FLAGS) -o $@ ./$(@D)
else
	CGO_ENABLED=0 go build $(GO_FLAGS) -o $(o) ./$(@D)
endif

clean:
	rm -rf cmd/v1/scanner

start-compse:
	docker-compose up