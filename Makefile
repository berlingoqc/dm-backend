PROJECT_NAME := dm-backend
VERSION := $(shell git describe --abbrev=0 --tags)
PKG := github.com/berlingoqc/dm-backend

RELEASE := $(PROJECT_NAME).tar.gz

PKG_LIST := $(shell go list ${PKG}/...)
TEST_FILES := $(shell find . -name '*.go' | grep -v _test.go)

LDFLAGS := -ldflags "-X ${PKG}/api.Version=$(VERSION)"

GOBUILD := go build -v --tags="fts5" $(LDFLAGS)

.PHONY: all dep build clean test install release

all: clean test build

testall: lint test race msan

install:
	@install dm-backend /usr/bin/
	@cp dm-backend.service /etc/systemd/system/
	@mkdir -p /etc/dm
	@cp config.json /etc/dm/

configureLocalEnv:
	@mkdir -p ${HOME}/.dm/{pipeline,config,script}
	@cp config/aria2.conf ${HOME}/.dm/{config}

release: build
	@mkdir -p ./release/
	@tar -zcvf ./release/$(RELEASE) $(PROJECT_NAME) config.json dm-backend.service VERSION

build: dep
	$(GOBUILD) ./cmd/dm-backend

clean:
	@rm -rf ./release ./test *.exe

lint:
	@golint -set_exit_status ${PKG_LIST}

test: dep
	@go test -v -short ${PKG_LIST}

race: dep
	@go test -v -race -short ${PKG_LIST}

msan: dep
	@go test -v -msan -short ${PKG_LIST}

dep:
	@go get -v -d ./...

genkeys:
	openssl ecparam -genkey -name secp384r1 -out server.key
	openssl req -new -x509 -sha256 -key server.key -out server.crt -days 3650
