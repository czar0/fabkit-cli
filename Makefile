SRC=$(shell find . -name "*.go")

.PHONY: all fmt lint test build install_deps

default: all

all: build

fmt:
	$(info ******************** checking formatting ********************)
	@test -z $(shell gofmt -l $(SRC)) || (gofmt -d $(SRC); exit 1)

lint:
	$(info ******************** running lint tools ********************)
	golangci-lint run -v

test: install_deps lint
	$(info ******************** running tests ********************)
	richgo test -v ./...

install: install_deps
	$(info ******************** installing binary ********************)
	cd cmd/fabkit && go install

build: install_deps
	$(info ******************** building binary ********************)
	cd cmd/fabkit && go build

install_deps:
	$(info ******************** downloading dependencies ********************)
	go get -v ./...
ifeq (, $(shell which golangci-lint))
go install github.com/golangci/golangci-lint/cmd/golangci-lint@latest
endif

ifeq (, $(shell which richgo))
go install github.com/kyoh86/richgo@latest
endif
