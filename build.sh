#!/bin/sh

go build -trimpath -ldflags "-s -w" && cp repo-and-module-indexerr ${HOME}/tools/
