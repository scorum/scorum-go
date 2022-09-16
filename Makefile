# supress output, run `make XXX V=` to be verbose
V := @

all: build test

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
