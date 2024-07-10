OUT := build/khttp-connectivity-monitor
PKG := github.com/kwarunek/khttp-connectivity-monitor
SRC_FILES := $(wildcard ./cmd/*.go) $(wildcard ./fmt/**/*.go)
VERSION := $(shell  grep '^\#\# ' CHANGELOG.md | head -n1 | grep -o '\[.*\]' |tr -d '[]')-$(shell uname -i)


.PHONY: help fmt build build-dep image publish deps version release

default: help

help:  # Show help for each of the Makefile recipes.
	@tail -n +5 Makefile | grep -E '^[a-zA-Z0-9 -]+:.*#' | sort | while read -r l; do printf "\033[1;32m$$(echo $$l | cut -f 1 -d':')\033[00m:$$(echo $$l | cut -f 2- -d'#')\n"; done


deps:  # Install dependecies - Go modules
	GOPRIVATE=stash.grupa.onet go mod tidy

build:  # Make a production build (w/o debuginfo)
	CGO_ENABLED=0 go build -v -o ${OUT} -ldflags="-s -w -X main.version=${VERSION}" ./cmd

build-dev:  # Make a development build (with debuginfo, race checker)
	go build -v -race -o ${OUT}-dev -ldflags="-X main.version=dev" ./cmd

run: build-dev  # Create a development build and run it
	./${OUT}-dev

fmt:  # Apply Go standard fromatting
	gofmt -l -s -w cmd/ pkg/

clean:  # Clean
	rm -rf ${OUT}*

version:  # Print a sniffed version
	@echo ${VERSION}

release: clean fmt test-functional image  # Creates clean, production image with a tag version after completed tests 

image:
	sudo docker build --network host -t kwarunek/khttp-connectivity-monitor:${VERSION} .
