# The binary to build (just the basename).
BIN := dasher

# This repo's root import path (under GOPATH).
PKG := github.com/eguevara/dasher

# Where to push the docker image.
REGISTRY ?= eguevara

# This version-strategy uses git tags to set the version string
VERSION := $(shell git describe --tags --always --dirty)

# Set an output prefix, which is the local directory if not specified
PREFIX?=$(shell pwd)
BUILDTAGS=

.PHONY: clean all fmt vet lint build test install static
.DEFAULT: default

all: clean build fmt lint test vet install

build:
	@echo "+ $@"
	@go build -tags "$(BUILDTAGS) cgo" .

static:
	@echo "+ $@"
	CGO_ENABLED=1 go build -tags "$(BUILDTAGS) cgo static_build" -ldflags "-w -extldflags -static" -o $(BIN) .

fmt:
	@echo "+ $@"
	@gofmt -s -l . | grep -v vendor | tee /dev/stderr

lint:
	@echo "+ $@"
	@go list ./... | grep -v /vendor/ | xargs -L1 golint

test: fmt lint vet
	@echo "+ $@"
	@go test -v -tags "$(BUILDTAGS) cgo" $(shell go list ./... | grep -v vendor)

vet:
	@echo "+ $@"
	@go vet $(shell go list ./... | grep -v vendor)

clean:
	@echo "+ $@"
	@rm -rf $(BIN)

install:
	@echo "+ $@"
	@go install .

version:
	@echo $(VERSION)

$(BIN):
	CGO_ENABLED=0 GOOS=linux go build -ldflags "-s" -a -installsuffix cgo -ldflags "-w" -o $(BIN) .

dbuild: $(BIN)
	docker build --rm --force-rm -t eguevara/$(BIN):1.0 .