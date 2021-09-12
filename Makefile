ROOT="${PWD}"
BIN="${ROOT}/bin"
XC_OS="linux darwin"
XC_ARCH="amd64"
XC_PARALLEL="2"
SRC=$(shell find . -name "*.go")

ifeq (, $(shell which gox))
$(warning "could not find gox in $(PATH), run: go get github.com/mitchellh/gox")
endif

.PHONY: all fmt lint test build install_deps clean

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

build: install_deps
	$(info ******************** building binary ********************)
	mkdir -p $(BIN)
	cd ${ROOT}/cmd/fabkit && gox \
		-os=$(XC_OS) \
		-arch=$(XC_ARCH) \
		-parallel=$(XC_PARALLEL) \
		-output=$(BIN)/{{.Dir}}_{{.OS}}_{{.Arch}} \
		;

install_deps:
	$(info ******************** downloading dependencies ********************)
	go get -v ./...

clean:
	rm -rf $(BIN)