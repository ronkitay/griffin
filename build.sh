#!/bin/sh

go build -trimpath -ldflags "-s -w" -o griffin cmd/griffin/main.go && cp griffin ${HOME}/tools/
