PROJECT_NAME := dm-backend
VERSION := $(shell git describe --abbrev=0 --tags)
PKG := github.com/berlingoqc/dm-backend

RELEASE := $(PROJECT_NAME)_$(VERSION).tar.gz

PKG_LIST := $(shell go list ${PKG}/...)
TEST_FILES := $(shell find . -name '*.go' | grep -v _test.go)

LDFLAGS := -ldflags "-X ${PKG}/api.Version=$(VERSION)"

GOBUILD := go build -v $(LDFLAGS)

.PHONY: all dep build clean test install release

all: clean test build

testall: lint test race msan

install:
	@install fxbe /usr/bin/
	@cp fxbe.service /etc/systemd/system/

release: build
	@mkdir -p ./release/
	@tar -zcvf ./release/$(RELEASE) $(PROJECT_NAME)

build: dep
	$(GOBUILD)

clean:
	@rm -rf ./release ./test

lint:
	@echo "Start linting"
	@golint -set_exit_status ${PKG_LIST}
	@echo "End linting"

test: dep
	@echo "Start test"
	@go test -v -short ${PKG_LIST}
	@echo "End test"

race: dep
	@echo "Start race"
	@go test -v -race -short ${PKG_LIST}
	@echo "End race"

msan: dep
	@echo "Start msan"
	@go test -v -msan -short ${PKG_LIST} 
	@echo "End msan"

dep:
	@go get -v -d ./...

genkeys:
	openssl ecparam -genkey -name secp384r1 -out server.key
	openssl req -new -x509 -sha256 -key server.key -out server.crt -days 3650