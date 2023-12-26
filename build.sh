#!/bin/sh

go build -trimpath -ldflags "-s -w" -o out/griffin cmd/griffin/main.go && cp out/griffin ${HOME}/tools/
