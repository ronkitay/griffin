#!/bin/sh

VERSION=$(git describe --tags --always 2>/dev/null || echo "dev")
go build -trimpath -ldflags "-s -w -X ronkitay.com/griffin/pkg/version.Version=${VERSION}" -o out/griffin cmd/griffin/main.go && cp out/griffin ${HOME}/tools/
