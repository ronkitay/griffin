# SHELL          := bash
# GO             ?= go
# GOOS           ?= $(word 1, $(subst /, " ", $(word 4, $(shell go version))))

# MAKEFILE       := $(realpath $(lastword $(MAKEFILE_LIST)))
# ROOT_DIR       := $(shell dirname $(MAKEFILE))
# SOURCES        := $(wildcard *.go src/*.go src/*/*.go) $(MAKEFILE)


build:
	goreleaser build --clean --snapshot --skip=post-hooks

release: build
	goreleaser release --skip-publish


.PHONY: build release
