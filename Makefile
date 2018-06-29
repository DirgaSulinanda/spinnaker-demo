#!/bin/bash

GOCMD=go
GOBUILD=${GOCMD} build
GOINSTALL=${GOCMD} install
GOTEST=${GOCMD} test
GODEP=dep

export NOW=$(shell date +'%FT%T%z')
export PKGS=$(shell go list ./... | grep -v vendor/)

APPNAME="spinnaker-demo"

all: test build run

get-deps:
	@echo "${NOW} UPDATING..."
	@${GODEP} ensure -v

test:
	@echo "${NOW} TESTING..."
	@${GOTEST} -v -cover -race ${PKGS}

build:
	@echo "${NOW} BUILDING..."
	@${GOBUILD} -race -o ${APPNAME} ./cmd/${APPNAME}

run:
	@echo "${NOW} RUNNING..."
	@./${APPNAME}

install:
	@echo "${NOW} INSTALLING..."
	@${GOINSTALL} ./cmd/${APPNAME}

docker-build:
	@CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o ${APPNAME} ./cmd/${APPNAME}
	@docker build -t ${APPNAME} -f Dockerfile .
