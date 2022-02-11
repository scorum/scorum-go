# supress output, run `make XXX V=` to be verbose
MKFILE_PATH := $(abspath $(lastword $(MAKEFILE_LIST)))
BUILD_PATH := $(dir $(MKFILE_PATH))
GOBIN ?= $(BUILD_PATH)tools/bin
ENV_PATH = PATH=$(GOBIN):$(PATH)

V := @

default: build

.PHONY: build
build:
	$(V) go build ./...

.PHONY: test
test:
	$(V)go test -mod=readonly -v ./...

.PHONY: vendor
vendor:
	$(V)go mod tidy
	$(V)go mod vendor

.PHONY: tools
tools:
	@if [ ! -f $(GOBIN)/mockgen ]; then\
		echo "Installing mockgen";\
		GOBIN=$(GOBIN) go install github.com/golang/mock/mockgen;\
	fi

.PHONY: generate
generate: tools
	$(ENV_PATH) go generate ./...