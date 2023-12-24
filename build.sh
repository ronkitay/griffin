#!/bin/sh

go build -trimpath -ldflags "-s -w" -o griffin && cp griffin ${HOME}/tools/
